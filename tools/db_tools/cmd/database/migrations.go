package database

import (
	cmn "inventory_money_tracking_software/cmd/common"
	"os"
	"path/filepath"
	"strings"
)

var combinedUpSql = ""
var combinedDownSql = ""

func buildMigrationQuries(p string, fi os.FileInfo, err error) error {
	// visits sql files and appends them into a huge sql file
	if err != nil && err != os.ErrNotExist {
		cmn.HandleError(err, cmn.ErrorLevels.Critical)
		return err
	}

	if fi.Mode().IsRegular() && fi.Size() > 0 && (strings.Contains(fi.Name(), ".sql")) {
		file, err := os.ReadFile(p)
		cmn.HandleError(err, cmn.ErrorLevels.Critical)
		s := string(file)
		s = strings.ToLower(s)
		// remove the transaction commands as the db package has its own way
		// of doing transactions
		s = strings.ReplaceAll(s, "begin;", " ")
		s = strings.ReplaceAll(s, "commit;", " ")
		s = strings.ReplaceAll(s, "\nbegin", " ")
		s = strings.ReplaceAll(s, "\ncommit", " ")
		s = strings.ReplaceAll(s, "commit\n", " ")
		s = strings.ReplaceAll(s, "begin\n", " ")

		if strings.Contains(fi.Name(), ".up.sql") {
			combinedUpSql = combinedUpSql + "\n" + s
		} else if strings.Contains(fi.Name(), ".down.sql") {
			combinedDownSql = combinedDownSql + "\n" + s
		}
	}
	return nil
}
func UpMigration() {
	cmn.Log("Starting database up migration", cmn.LogLevels.Info)
	tx, err := db.Begin()
	defer tx.Rollback()
	cmn.HandleError(err, cmn.ErrorLevels.Critical)
	_, err = tx.Exec(combinedUpSql)
	if err != nil && !(strings.Contains(err.Error(), "already exists")) {
		cmn.HandleError(err, cmn.ErrorLevels.Critical)
	}
	err = tx.Commit()
	cmn.HandleError(err, cmn.ErrorLevels.Warning)
	cmn.Log("Database migrated up successfully", cmn.LogLevels.Info)
}
func DownMigration() {
	if cmn.GetEnvironment() == cmn.Envs.Production {
		return
	}
	cmn.Log("Starting database down migration", cmn.LogLevels.Info)
	tx, err := db.Begin()
	defer tx.Rollback()
	cmn.HandleError(err, cmn.ErrorLevels.Critical)

	_, err = tx.Exec(combinedDownSql)
	cmn.HandleError(err, cmn.ErrorLevels.Critical)
	err = tx.Commit()
	cmn.HandleError(err, cmn.ErrorLevels.Critical)

	cmn.Log("Database migrated down successfully", cmn.LogLevels.Info)
}

func InitDatabaseMigrationFiles() {
	if cmn.Envs.Production == cmn.GetEnvironment() {
		cmn.Log("migratinos shouldnt be done by the app in production... Use docker image instead to create the tables", cmn.LogLevels.Warning)
		return
	}
	migrationsFile := os.Getenv("MIGRATIONS")
	filepath.Walk(migrationsFile, buildMigrationQuries)
}
