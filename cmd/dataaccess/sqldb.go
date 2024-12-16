package dataaccess

import (
	"database/sql"
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/database_connecter"

	// _ "github.com/lib/pq"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

type sqlConection interface {
	*sql.DB | *sql.Tx
	// either uses sql.DB or sql.Tx depending on who's calling
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

func initDatabase() {
	db = database_connecter.ConnectToDatabase(database_connecter.GetDbInfo())
	if cmn.GetEnvironment() == cmn.Envs.Production {
		initTemplate()
	}
}

// inserts the inital default groups to the db,
// names are unique so it's okay to run this query over and over
// this should be done for new users only. Other users will have their
// own configs that they made themselves
func initTemplate() {
	cmn.Log("Inserting template to tables", cmn.LogLevels.Operation)
	stmt := insertSql.insertFinancialGroup
	// todo use batch quering instead of one by one
	_, err := db.Exec(stmt, "Assets", -1, "All asset accounts", cmn.GroupTypes.Assets)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
	}
	_, err = db.Exec(stmt, "Liabilities", -1, "All liability accounts", cmn.GroupTypes.Liabilities)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
	}
	_, err = db.Exec(stmt, "Equity", -1, "All Owner's Equity accounts", cmn.GroupTypes.Equity)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
	}
}
