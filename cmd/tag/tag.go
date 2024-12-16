package tag

import (
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	mdls "inventory_money_tracking_software/cmd/models"
)

func CreateTag(t mdls.Tag) (int, cmn.ErrorCode) {
	if cmn.IsStringEmpty(t.Name) {
		return 0, cmn.ErrorCodes.NameWasEmpty
	}
	errorCode := cmn.ErrorCodes.NoError
	tagId, err := dataaccess.CreateTag("", t)

	if err != nil {
		errorCode = cmn.ErrorCodes.CreateFailed
		cmn.Log("error creating tag", cmn.ErrorLevels.Error)
		e := err.Error()
		if cmn.IsDuplicateKeyConstraintViolated(e) {
			errorCode = cmn.ErrorCodes.NameAlreadyExists
		}
	} else {
		m := fmt.Sprintf("Created tag successfully %d", tagId)
		cmn.Log(m, cmn.LogLevels.Info)
	}
	return tagId, errorCode
}

func CreateTagFinancialJournalEntry(tagId int, jId int) cmn.ErrorCode {
	if tagId == 0 || jId == 0 {
		cmn.Log("Expected that both ids are non-zero, received tagId %d and journal entry id %d",
			cmn.LogLevels.Error, tagId, jId)
		return cmn.ErrorCodes.InvalidId
	}

	err := dataaccess.CreateTagFinancialJournalEntry("", tagId, jId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return cmn.ErrorCodes.CreateFailed
	}

	return cmn.ErrorCodes.NoError
}

func ReadTagById(tagId int) (mdls.Tag, cmn.ErrorCode) {
	t, err := dataaccess.ReadTagById("", tagId)
	if err != nil {
		return t, cmn.ErrorCodes.ReadFailed // TODO handle error reporting, we want to give a more specific error msg to the client
	}
	return t, cmn.ErrorCodes.NoError
}

func ReadTagsBelongingToFinancialJournalId(jId int) ([]mdls.Tag, cmn.ErrorCode) {
	t, err := dataaccess.ReadTagsByFinancialJournalEntryId("", jId)
	if err != nil {
		return t, cmn.ErrorCodes.ReadFailed // TODO handle error reporting, we want to give a more specific error msg to the client
	}
	return t, cmn.ErrorCodes.NoError

}
