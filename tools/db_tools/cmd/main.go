package main

import (
	"db_tools/cmd/database"
	"db_tools/cmd/seeder"
	"db_tools/cmd/server"
	"flag"
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	"os"
	"strings"
)

func InitDbTools() {
	env := os.Getenv("ENV_TYPE")
	cmn.InitCommon(env, "db_tools")
	database.InitDatabaseMigrationFiles()
	database.InitDatabaseConnection()
}
func printHelpMsg() {
	fmt.Printf(`
db_tools can help you perform some database tasks that are needed for development and testing.
It can:
 1. Do up or down migrations of the database. Doing down then up migrations is as good as deleting all the data in the db.
 2. Seed the database with default values
This tool can be used in two modes, either as a cli tool or as a server.


To use in server mode pass the -server flag. It will run on the port defined in the env variable DB_TOOLS_PORT which defaults to 4444,
then you can use on of the following routes:
1. ping to check if the server is still there
2. up to do up migration
3. down to do down migration
4. seed to seed the db
5. fresh to get a fresh seeded database, it will do a down then up migrations, then seed the database.
5. exit to exit the db
For example: curl localhost:4444/down will trigger a down migration

To use the tool in cli mode:
pass the desired action after the tool invocation command.
For example
db_tools seed --> this will seed the database
db_tools down --> this will do a down migration
db_tools up --> this will do a down migration
db_tools fresh --> this will do a down migration, then up migration, then a seed

You can also use the first letter as shortcut.
For example:
	db_tools s --> same as db_tools seed

You can also run the direct go file like this:
go run cmd/main.go  down


To control whether you are editing the dev database or the test database you can control the 
env variable ENV_TYPE. You can also prefix the call to this tool with the variable you want this way it is not going to
change the env var globally.
For example:
ENV_TYPE=test go run cmd/main.go  -server
OR
ENV_TYPE=test go run cmd/main.go  down 

The default is to edit the dev server
`)
	os.Exit(0)
}
func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if strings.Contains(arg, "help") || arg == "h" {
			printHelpMsg()
		}
	}
	useServer := flag.Bool("server", false, "Whether to use the server mode instead of cli mode")
	showHelp := flag.Bool("help", false, "Whether to show help")
	showHelp2 := flag.Bool("h", false, "Whether to show help")
	flag.Parse()
	if *showHelp || *showHelp2 {
		printHelpMsg()
	}
	if !*useServer && len(os.Args) == 1 {
		printHelpMsg()
	}

	InitDbTools()

	if *useServer {
		cmn.Log("Starting in server mode", cmn.LogLevels.Info)
		server.StartServer()
		return
	}

	// command line mode instead
	if len(os.Args) > 1 {
		arg := os.Args[1]
		arg = strings.ToLower(arg)
		if strings.Contains(arg, "up") || arg == "u" {
			database.UpMigration()
		} else if strings.Contains(arg, "down") || arg == "d" {
			database.DownMigration()
		} else if strings.Contains(arg, "seed") || arg == "s" {
			seeder.SeedDatabase()
		} else if strings.Contains(arg, "fresh") || arg == "f" {
			database.DownMigration()
			database.UpMigration()
			seeder.SeedDatabase()
		} else {
			printHelpMsg()
		}
	}
}
