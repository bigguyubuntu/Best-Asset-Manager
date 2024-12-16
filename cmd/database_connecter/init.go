package database_connecter

import (
	"database/sql"
	"errors"
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
	"os"

	// _ "github.com/lib/pq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// this function configures pgx driver
func pgxConfig(connString string) string {
	config, err := pgx.ParseConfig(connString)
	cmn.HandleError(err, cmn.ErrorLevels.Critical)
	config.DefaultQueryExecMode = pgx.QueryExecModeExec
	connStr := stdlib.RegisterConnConfig(config)
	return connStr
}

// gets the database info from env vars
func GetDbInfo() mdls.DbInfo {
	dbname := ""
	env := cmn.GetEnvironment()
	if env == cmn.Envs.Test {
		dbname = os.Getenv("TEST_DB_NAME")
	} else if env == cmn.Envs.Development || env == cmn.Envs.Production {
		dbname = os.Getenv("DB_NAME")
	} else {
		cmn.HandleError(errors.New("Environment variables are not set correctly.... can't run application"), cmn.ErrorLevels.Critical)
	}
	return mdls.DbInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   dbname,
	}
}

// connects to database and returns the connection
func ConnectToDatabase(d mdls.DbInfo) *sql.DB {
	cmn.Log("Initalizing database connection", cmn.LogLevels.Info)
	// mdls.DbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s "+
	// 	"sslmode=disable",
	// 	d.Host, d.Port, d.User, d.Password, d.DbName)

	// TODO right now access doesn't use ssl make sure to make it secure later
	// const conString = "postgres://YourUserName:YourPassword@YourHostname:Port/YourDatabaseName?sslmode=disable";
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable", d.User, d.Password, d.Host, d.Port, d.DbName)

	connectionString = pgxConfig(connectionString)

	mydb, err := sql.Open("pgx", connectionString)
	cmn.HandleError(err, cmn.ErrorLevels.Critical)
	if mydb == nil {
		m := fmt.Errorf("Can't reach database, database object is nil. critical issue")
		cmn.HandleError(m, cmn.ErrorLevels.Critical)
	}
	cmn.Log("Pinging the database", cmn.LogLevels.Info)
	err = mydb.Ping()
	cmn.HandleError(err, cmn.ErrorLevels.Critical)
	cmn.Log("Database initalized successfully, connected to %s", cmn.LogLevels.Info, d.DbName)
	return mydb
}
