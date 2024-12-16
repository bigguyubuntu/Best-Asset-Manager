package database

import (
	"database/sql"
	"inventory_money_tracking_software/cmd/database_connecter"
	"os"
)

var db *sql.DB

func InitDatabaseConnection() {
	dbInfo := database_connecter.GetDbInfo()
	dbInfo.User = os.Getenv("DB_TOOL_USERNAME")
	db = database_connecter.ConnectToDatabase(dbInfo)
}
