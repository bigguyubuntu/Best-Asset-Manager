package dataaccess

import (
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
)

func CreateTag(txKey string, t mdls.Tag) (int, error) {
	if txKey != "" {
		tx, ok := transactionMap.Load(txKey)
		if !ok {
			cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
			return 0, fmt.Errorf("")
		}
		return insertTag(tx, t)
	}
	return insertTag(db, t)
}

func CreateTagFinancialJournalEntry(txKey string, tagId int, jId int) error {
	if txKey != "" {
		tx, ok := transactionMap.Load(txKey)
		if !ok {
			cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
			return fmt.Errorf("")
		}
		return insertTagFinancialJournalEntry(tx, tagId, jId)
	}
	return insertTagFinancialJournalEntry(db, tagId, jId)
}

func ReadTagById(txKey string, tagId int) (mdls.Tag, error) {
	if txKey != "" {
		tx, ok := transactionMap.Load(txKey)
		if !ok {
			cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
			return mdls.Tag{}, fmt.Errorf("")
		}
		return selectTagById(tx, tagId)
	}
	return selectTagById(db, tagId)
}
func ReadTagsByFinancialJournalEntryId(txKey string, jId int) ([]mdls.Tag, error) {
	if txKey != "" {
		tx, ok := transactionMap.Load(txKey)
		if !ok {
			cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
			return nil, fmt.Errorf("")
		}
		return selectTagsByFinancialJournalEntryId(tx, jId)
	}
	return selectTagsByFinancialJournalEntryId(db, jId)
}
