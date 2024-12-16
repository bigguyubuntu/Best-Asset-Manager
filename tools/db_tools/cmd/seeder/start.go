package seeder

import (
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	"net/http"
	"os"
)

var baseUrl string

// must be in dev or test environment to do so
func seedWithDefaults() {
	env := os.Getenv("ENV_TYPE")
	if env != "dev" && env != "test" {
		fmt.Printf("Environment is not set properly. Or is not dev nor test")
		os.Exit(1)
	}
	seedFinancialGroups()
	seedFinancialAccounts()
	seedFinancialJournalEntries()
	//----
	seedTags()
}

// bootstraps the database with initial test data
// seeding uses the backend api to populate the database
func SeedDatabase() {
	cmn.Log("seeding database with default values", cmn.LogLevels.Info)
	base := os.Getenv("BACKEND_HOST") + ":" + os.Getenv("BACKEND_PORT")
	healthUrl := base + "/ping"
	baseUrl = base + "/api"
	res, err := http.Get(healthUrl)
	if err != nil {
		fmt.Printf("error making http request to backend: %s\n", err)
		os.Exit(1)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Printf("error talking to db, invalid status code. client: status code: %d\n", res.StatusCode)
		os.Exit(1)
	}

	seedWithDefaults()
}
