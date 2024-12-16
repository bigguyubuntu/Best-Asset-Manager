package accounting

import (
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	mdls "inventory_money_tracking_software/cmd/models"
)

func ValidateGroup(g mdls.FinancialGroup) cmn.ErrorCode {
	if !cmn.IsGroupTypeValid(g.GroupType) {
		cmn.Log(fmt.Sprintf("user request error at group creation \n, group type isn't right he inputed %+v the UI doesn't allow this, invesigate this further",
			g), cmn.LogLevels.Info)
		return cmn.ErrorCodes.InvalidFinancialGroupType
	}
	if g.ParentId == 0 {
		return cmn.ErrorCodes.InvalidParentGroupId
	}
	if cmn.IsStringEmpty(g.Name) {
		return cmn.ErrorCodes.NameWasEmpty
	}
	return cmn.ErrorCodes.NoError
}

// creates a new finaincial account group, if there's error it will return an error code
func CreateNewFinancialGroup(grp mdls.FinancialGroup) (string, cmn.ErrorCode) {
	errorCode := ValidateGroup(grp)
	if errorCode != cmn.ErrorCodes.NoError {
		return "", errorCode
	}
	groupId, err := dataaccess.CreateNewFinancialGroup(grp)
	if err != nil {
		cmn.Log("error creating Financial Group",
			cmn.ErrorLevels.Error)
		e := err.Error()
		if cmn.IsDuplicateKeyConstraintViolated(e) {
			errorCode = cmn.ErrorCodes.NameAlreadyExists
		}
	} else {
		m := fmt.Sprintf("Created Financial Group %s", groupId)
		cmn.Log(m, cmn.LogLevels.Operation)
	}
	return groupId, errorCode
}

func ReadAllFinancialGroups() ([]mdls.FinancialGroup, cmn.ErrorCode) {
	l, err := dataaccess.ReadAllFinancialGroups()
	if err != nil {
		return []mdls.FinancialGroup{}, cmn.ErrorCodes.ReadFailed
	}
	return l, cmn.ErrorCodes.NoError
}
func ReadFinancialGroupsWithType(t string) ([]mdls.FinancialGroup, cmn.ErrorCode) {
	l, err := dataaccess.ReadFinancialGroupsWithType(t)
	if err != nil {
		return []mdls.FinancialGroup{}, cmn.ErrorCodes.ReadFailed
	}
	return l, cmn.ErrorCodes.NoError
}

func ReadFinancialGroupById(id int) mdls.FinancialGroup {
	return dataaccess.ReadFinancialGroupById(id)
}

func DeleteFinancialAccoutsGroup(id int) cmn.ErrorCode {
	//there has to be at least one group for each of three account types
	// assets, liabilities and SE. Perform the check and if we're at the
	// last one surface an error to the user that he must retain at least
	// one group for each type
	// meaning that at least one group with parentId=-1 in each of the categories

	//step1 read financial_group_type and parentId from db
	// step2 if parentId is -1, make sure it's not the last one, if it is then don't allow delete
	// step3 if step2 is clear then we allow the delete.

	g, err := dataaccess.ReadFinancialGroupsParentIdAndType(id)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return cmn.ErrorCodes.DeleteFailed // delete failed can't find account
	}
	if g.ParentId == -1 {
		if dataaccess.IsFinancialGroupsWithTypeLessThan2(g.GroupType) {
			return cmn.ErrorCodes.NeedAtLeastOneGroupForEveryType
		}
	}

	if dataaccess.DeleteFinancialAccoutsGroup(id) {
		return cmn.ErrorCodes.NoError
	}

	return cmn.ErrorCodes.DeleteFailed
}

func UpdateFinancialGroup(id int, newGroup mdls.FinancialGroup) cmn.ErrorCode {
	// todo don't allow the update if it wants to make the parent group itself

	if cmn.IsStringEmpty(newGroup.Name) {
		return cmn.ErrorCodes.NameWasEmpty
	}
	if !cmn.IsGroupTypeValid(newGroup.GroupType) {
		return cmn.ErrorCodes.InvalidFinancialGroupType
	}

	if newGroup.ParentId == id {
		return cmn.ErrorCodes.ParentCantBeSameAsId
	}

	success := dataaccess.UpdateFinancialAccoutsGroup(id, newGroup)
	if success {
		return cmn.ErrorCodes.NoError
	}
	return cmn.ErrorCodes.UpdateFailed
}

func ReadFinancialGroupsByParentGroup(grpId int) ([]mdls.FinancialGroup, cmn.ErrorCode) {
	grps, err := dataaccess.ReadFinancialGroupsByParentGroup(grpId)
	if err != nil {
		return nil, cmn.ErrorCodes.ReadFailed
	}
	return grps, cmn.ErrorCodes.NoError
}

func ReadFinancialGroupByNameAndParentId(groupName string, parentGroupId int) (mdls.FinancialGroup, cmn.ErrorCode) {
	txKey, err := dataaccess.ObtainTransactionKey()
	defer dataaccess.RollbackTransaction(txKey)
	if err != nil {
		cmn.Log("Failed to obtain transaction key while creating a new journal entry",
			cmn.LogLevels.Error)
		return mdls.FinancialGroup{}, cmn.ErrorCodes.ReadFailed
	}
	g, err := dataaccess.ReadFinancialGroupByNameAndParentId(txKey, groupName, parentGroupId)
	if err != nil {
		return mdls.FinancialGroup{}, cmn.ErrorCodes.ReadFailed
	}
	return g, cmn.ErrorCodes.NoError
}

func ReadFinancialAccountByNameAndGroupId(accountName string, groupId int) (mdls.FinancialAccount, cmn.ErrorCode) {
	txKey, err := dataaccess.ObtainTransactionKey()
	defer dataaccess.RollbackTransaction(txKey)
	if err != nil {
		cmn.Log("Failed to obtain transaction key while creating a new journal entry",
			cmn.LogLevels.Error)
		return mdls.FinancialAccount{}, cmn.ErrorCodes.ReadFailed
	}
	a, err := dataaccess.ReadFinancialAccountByNameAndGroupId(txKey, accountName, groupId)
	if err != nil {
		return mdls.FinancialAccount{}, cmn.ErrorCodes.ReadFailed
	}
	return a, cmn.ErrorCodes.NoError
}

func findGroupsByIds(ids []int) ([]mdls.FinancialGroup, cmn.ErrorCode) {
	if len(ids) == 0 {
		return []mdls.FinancialGroup{}, cmn.ErrorCodes.NoError
	}
	gs, err := dataaccess.FindGroupsByIds(ids)
	if err != nil {
		return []mdls.FinancialGroup{}, cmn.ErrorCodes.ReadFailed
	}
	return gs, cmn.ErrorCodes.NoError
}
