package dataaccess

import (
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
)

func intLimit(id int, sqlStmnt string) {
	if cmn.MaxInt/1000 < id {
		// TODO Fire a warning signal, but for now just log
		cmn.Log(fmt.Sprintf("We're approching the limit of integers in the db, this happend after the statment %s", sqlStmnt),
			cmn.ErrorLevels.Warning)
	}
}

func bigIntLimit(id int, sqlStmnt string) {
	if cmn.MaxBigInt/1000 < id {
		// TODO Fire a warning signal, but for now just log
		cmn.Log(fmt.Sprintf("We're approching the limit of bigint, this happend after the statment %s", sqlStmnt),
			cmn.ErrorLevels.Warning)
	}
}
