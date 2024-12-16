package dataaccess

import (
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
)

func insertTag[T sqlConection](conn T, t mdls.Tag) (int, error) {
	cmn.Log("insert new tag to database", cmn.LogLevels.Operation)
	s := insertSql.insertTag
	id := 0
	err := conn.QueryRow(s, t.Name, t.Description).Scan(&id)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return id, err
}
func insertTagFinancialJournalEntry[T sqlConection](conn T, tid int, jId int) error {
	cmn.Log("insert new TagFinancialJournalEntry to database", cmn.LogLevels.Operation)
	s := insertSql.insertTagFinancialJournalEntry
	_, err := conn.Exec(s, tid, jId)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return err
}
func selectTagById[T sqlConection](conn T, tId int) (mdls.Tag, error) {
	cmn.Log("selecting tag by id", cmn.LogLevels.Operation)
	s := readSql.readTagById
	t := mdls.Tag{Id: tId}
	err := conn.QueryRow(s, tId).Scan(&t.Name, &t.Description)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return t, err
}

func selectTagsByFinancialJournalEntryId[T sqlConection](conn T, jId int) ([]mdls.Tag, error) {
	cmn.Log("selectTagsByFinancialJournalEntryId", cmn.LogLevels.Operation)
	s := readSql.readTagsByFinancialJournalEntryId
	tags := []mdls.Tag{}
	rows, err := conn.Query(s, jId)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return tags, err
	}
	for rows.Next() {
		t := mdls.Tag{}
		err = rows.Scan(&t.Id, &t.Name, &t.Description)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		tags = append(tags, t)
	}
	return tags, err
}

// func selectTagsByIds[T sqlConection](conn T, t mdls.Tag) (int, error) {
// 	cmn.Log("insert new tag to database", cmn.LogLevels.Operation)
// 	s := insertSql.insertTag
// 	id := 0
// 	err := conn.QueryRow(s, t.Name, t.Description).Scan(&id)
// 	cmn.HandleError(err, cmn.ErrorLevels.Error)
// 	return id, err
// }
