package dataaccess

import (
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
)

// Financial Accounts related below
func CreateAccount(a mdls.FinancialAccount) (string, error) {
	return insertFinancialAccount(a)
}

func CreateNewFinancialGroup(g mdls.FinancialGroup) (string, error) {
	return insertFinancialGroup(g)
}

func ReadAllFinancialGroups() ([]mdls.FinancialGroup, error) {
	return selectAllFinancialGroup()
}

func ReadFinancialGroupsWithType(t string) ([]mdls.FinancialGroup, error) {
	cmn.Log(fmt.Sprintf("Reading financialGroups with type %s ", t), cmn.LogLevels.Info)
	return selectFinancialGroupWithType(t)
}

func ReadFinancialGroupsParentIdAndType(groupId int) (mdls.FinancialGroup, error) {
	cmn.Log(fmt.Sprintf("ReadGroupParentIdAndType with id %d ", groupId), cmn.LogLevels.Info)
	return selectFinancialGroupsParentIdAndType(groupId)
}

func IsFinancialGroupsWithTypeLessThan2(groupType string) bool {
	cmn.Log(fmt.Sprintf("isFinancialGroupsWithTypeLessThan2 for type %s", groupType), cmn.LogLevels.Info)
	return isFinancialGroupsWithTypeLessThan2(groupType)
}

func DeleteFinancialAccoutsGroup(groupId int) bool {
	cmn.Log(fmt.Sprintf("DeleteFinancialAccoutsGroup for id %d", groupId), cmn.LogLevels.Info)
	return deleteFinancialAccoutsGroup(groupId)
}

func ReadFinancialGroupById(groupId int) mdls.FinancialGroup {
	cmn.Log(fmt.Sprintf("ReadFinancialGroupById with id %d ", groupId), cmn.LogLevels.Info)
	return selectFinancialGroupById(groupId)
}

func UpdateFinancialAccoutsGroup(groupId int, grp mdls.FinancialGroup) bool {
	cmn.Log(fmt.Sprintf("UpdateFinancialAccoutsGroup for id %d", groupId), cmn.LogLevels.Info)
	return updateFinancialAccoutsGroup(groupId, grp)
}

func ReadFinancialAccountsWithGroup(groupId int) ([]mdls.FinancialAccount, error) {
	cmn.Log(fmt.Sprintf("ReadFinancialAccountsWithGroup with groupId %d ", groupId), cmn.LogLevels.Info)
	return selectFinancialAccountsWithGroup(groupId)
}

func ReadFinancialGroupsByParentGroup(groupId int) ([]mdls.FinancialGroup, error) {
	cmn.Log(fmt.Sprintf("ReadFinancialAccountsWithGroup with groupId %d ", groupId), cmn.LogLevels.Info)
	return selectFinancialGroupsByParentGroup(groupId)
}

func InsertJournalEntry(txKey string, j mdls.FinancialJournalEntry) (int, error) {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return 0, fmt.Errorf("")
	}
	return insertJournalEntry(tx, j)
}

func FindOrCreateImbalanceAccount(txKey string) (mdls.FinancialAccount, error) {
	a := mdls.FinancialAccount{}
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return a, fmt.Errorf("")
	}
	return findOrCreateImbalanceAccount(tx)
}
func ReadFinancialAccountBalanceAndType(txKey string, accountId int) (mdls.FinancialAccount, error) {
	a := mdls.FinancialAccount{}
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return a, fmt.Errorf("")
	}
	return selectAccountBalanceAndType(tx, accountId)
}
func ReadFinancialTransactionsByJournalEntryId(txKey string, jId int) ([]mdls.FinancialTransaction, error) {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return []mdls.FinancialTransaction{}, fmt.Errorf("")
	}
	return selectFinancialTransactionsByJournalEntryId(tx, jId)
}
func ReadFinancialTransactionsByJournalEntries(txKey string, jIds []int) (map[int][]mdls.FinancialTransaction, error) {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return nil, fmt.Errorf("")
	}
	return selectFinancialTransactionsByJournalEntries(tx, jIds)
}
func ReadFinancialTransactionsByAccountAndJournalEntry(txKey string, accountId int,
	journalEntryId int) (mdls.FinancialTransaction, error) {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return mdls.FinancialTransaction{}, fmt.Errorf("")
	}
	return selectTransactionByAccountAndJournalEntry(tx, accountId, journalEntryId)
}
func ReadFinancialTransactionsByAccount(txKey string, accountId int) ([]mdls.FinancialTransaction, error) {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return []mdls.FinancialTransaction{}, fmt.Errorf("")
	}
	return selectTransactionsByAccount(tx, accountId)
}
func UpdateFinancialAccountBalance(txKey string, account mdls.FinancialAccount) error {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return fmt.Errorf("")
	}
	return updateFinancialAccountBalance(tx, account)
}
func CreateNewFinancialTransaction(txKey string, t mdls.FinancialTransaction) (int, error) {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return 0, fmt.Errorf("")
	}
	return insertNewTransaction(tx, t)
}
func UpdateFinancialTransaction(txKey string, t mdls.FinancialTransaction) error {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return fmt.Errorf("")
	}
	return updateFinancialTransaction(tx, t)
}
func DeleteFinancialTransaction(txKey string, t mdls.FinancialTransaction) error {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return fmt.Errorf("")
	}
	return deleteFinancialTransaction(tx, t)
}

func UpdateJournalEntryInformation(txKey string, j mdls.FinancialJournalEntry) error {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return fmt.Errorf("")
	}
	return updateJournalEntryInformation(tx, j)
}

func FindTransactionJournalEntriesByIds(txKey string, ids []int) ([]int, error) {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return nil, fmt.Errorf("")
	}
	return selectTransactionJournalEntriesByIds(tx, ids)
}
func FindAccountsByIds(ids []int) ([]mdls.FinancialAccount, error) {
	return selectAccountsByIds(ids)
}
func FindGroupsByIds(ids []int) ([]mdls.FinancialGroup, error) {
	return selectGroupsByIds(ids)
}
func FindAllAccounts() ([]mdls.FinancialAccount, error) {
	return selectAllAccounts()
}
func FindAllJournalEntries() ([]mdls.FinancialJournalEntry, error) {
	return selectAllFinancialJournalEntries()
}
func FindJournalEntriesByIds(ids []int) ([]mdls.FinancialJournalEntry, error) {
	return selectFinancialJournalEntriesByIds(ids)
}
func FindAllJournalEntriesComposite() ([]mdls.FinancialJournalEntryComposite, error) {
	return selectAllFinancialJournalEntriesComposite()
}
func FindJournalEntryCompositeById(jId int) (mdls.FinancialJournalEntryComposite, error) {
	return selectFinancialJournalEntryCompositeById(jId)
}
func FindJournalEntryCompositeByAccountId(jId int) ([]mdls.FinancialJournalEntryComposite, error) {
	return selectFinancialJournalEntryCompositeByAccountId(jId)
}

func ReadFinancialGroupByNameAndParentId(txKey string, groupName string, parentGroupId int) (mdls.FinancialGroup, error) {
	cmn.Log(fmt.Sprintf("ReadFinancialGroupByNameAndParentId with name %s and parentId %d ",
		groupName, parentGroupId), cmn.LogLevels.Info)
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return mdls.FinancialGroup{}, fmt.Errorf("")
	}
	return selectFinancialGroupByNameAndParentId(tx, groupName, parentGroupId)
}

func ReadFinancialAccountByNameAndGroupId(txKey string, accountName string, groupId int) (mdls.FinancialAccount, error) {
	cmn.Log(fmt.Sprintf("ReadFinancialAccountByNameAndGroupId with name %s and groupId %d ",
		accountName, groupId), cmn.LogLevels.Info)
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return mdls.FinancialAccount{}, fmt.Errorf("")
	}
	return selectAccountByNameAndGroupId(tx, accountName, groupId)
}
