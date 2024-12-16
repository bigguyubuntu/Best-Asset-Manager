package main

import (
	"fmt"
	"inventory_money_tracking_software/cmd/accounting"
	"inventory_money_tracking_software/cmd/api"
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	"inventory_money_tracking_software/cmd/tag"
	"os"
	"os/exec"
	"path/filepath"
)

func runStartupScript() {
	if cmn.Envs.Development == cmn.GetEnvironment() {
		var err error
		var stdout []byte
		d, e := os.Getwd()
		cmn.HandleError(e, cmn.ErrorLevels.Critical)
		d = filepath.Join(d, "scripts")
		// run dev script
		d = filepath.Join(d, "development_environment", "dev_env_init.sh")
		stdout, err = exec.Command("/bin/bash", d).Output()
		fmt.Println(string(stdout))
		if err != nil || string(stdout) == "" {
			cmn.Log("Error in the start up script, is the script the same as the environment?", cmn.LogLevels.Error)
			cmn.HandleError(err, cmn.ErrorLevels.Critical)
			panic("environment setup startup script failed") // we're also panicking in case the err is nil
		}
		cmn.Log("start up script was successful", cmn.LogLevels.Info)
	}
}

func initEnvironment() {
	env := os.Getenv("ENV_TYPE")
	cmn.InitCommon(env, "BAM")
}

func registerAllHandlers() {
	accounting.RegisterHandlers()
	tag.RegisterHandlers()
	api.RegisterInvalidUrlHandler() // must be the last one
}

func initProgram() {
	initEnvironment() // must be the first thing called
	cmn.Log("Program initializing...", cmn.LogLevels.Info)
	runStartupScript()
	dataaccess.InitDataAccess()
	api.InitApi()
	registerAllHandlers()
	api.StartApi()
}

func main() {
	initProgram()
}
