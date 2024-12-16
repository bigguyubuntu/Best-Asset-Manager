package accounting

import (
	"errors"
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	mdls "inventory_money_tracking_software/cmd/models"
)

func ValidateAccount(acnt mdls.FinancialAccount) cmn.ErrorCode {
	if !cmn.IsAccountTypeValid(acnt.AccountType) {
		cmn.Log(
			fmt.Sprintf(
				"user request error at account creation\n, account type isn't right he inputed %+v the UI doesn't allow this, invesigate this further",
				acnt), cmn.LogLevels.Info)
		return cmn.ErrorCodes.InvalidFinancialAccountType
	}
	if cmn.IsStringEmpty(acnt.Name) {
		return cmn.ErrorCodes.NameWasEmpty
	}
	if acnt.GroupId <= 0 {
		return cmn.ErrorCodes.InvalidParentGroupId
	}
	return cmn.ErrorCodes.NoError
}

func CreateAccount(acnt mdls.FinancialAccount) (string, cmn.ErrorCode) {
	// creates a new finaincial account, if there's error it will return an error code

	errorCode := ValidateAccount(acnt)
	if errorCode != cmn.ErrorCodes.NoError {
		return "", errorCode
	}
	// todo check if account type matches that of the parent group
	acntId, err := dataaccess.CreateAccount(acnt)
	if err != nil {
		errorCode = cmn.ErrorCodes.CreateFailed
		cmn.Log("error creating financial account", cmn.ErrorLevels.Error)
		e := err.Error()
		if cmn.IsDuplicateKeyConstraintViolated(e) {
			errorCode = cmn.ErrorCodes.NameAlreadyExists
		}
	} else {
		m := fmt.Sprintf("Created financial account %s", acntId)
		cmn.Log(m, cmn.LogLevels.Info)
	}
	return acntId, errorCode
}

func ReadFinancialAccountsWithGroup(grpId int) ([]mdls.FinancialAccount, cmn.ErrorCode) {
	acnts, err := dataaccess.ReadFinancialAccountsWithGroup(grpId)
	if err != nil {
		return nil, cmn.ErrorCodes.ReadFailed
	}
	return acnts, cmn.ErrorCodes.NoError
}

// if there's an imbalance transaction this function will take care of creating the
// imbalance account and assigning it to the imbalance Transaction
// you MUST make sure that the imbalance transaction is indeed created
// the last element of transactions argument.
// This func always assumes last elem to be the imbalance transaction.
func handleImbalanceAccount(txKey string, transactions []mdls.FinancialTransaction) error {
	// if we have an imbalance transaction, then obtain imbalance account
	// and set the imblanace transaction to belong to the imbalance account
	imbalanceAccount, err := dataaccess.FindOrCreateImbalanceAccount(txKey)
	if err != nil {
		return err
	}
	//  the last transaction is imbalane and it belongs to the imbalance account
	imbalanceTransaction := transactions[len(transactions)-1]
	imbalanceTransaction.FinancialAccountId = imbalanceAccount.Id
	imbalanceAccount.UpdateAccountBalanceAfterTransaction(imbalanceTransaction)
	transactions[len(transactions)-1] = imbalanceTransaction
	return nil
}

func findAccountsByIds(ids []int) ([]mdls.FinancialAccount, cmn.ErrorCode) {
	if len(ids) == 0 {
		return []mdls.FinancialAccount{}, cmn.ErrorCodes.NoError
	}
	as, err := dataaccess.FindAccountsByIds(ids)
	if err != nil {
		return []mdls.FinancialAccount{}, cmn.ErrorCodes.ReadFailed
	}
	return as, cmn.ErrorCodes.NoError
}

func findAllAccounts() ([]mdls.FinancialAccount, cmn.ErrorCode) {
	as, err := dataaccess.FindAllAccounts()
	if err != nil {
		return []mdls.FinancialAccount{}, cmn.ErrorCodes.ReadFailed
	}
	return as, cmn.ErrorCodes.NoError
}

// queries for accountId and updates the quriedAccounts map.
// if the account already exists in the map it will return
// the value from the map instead of qurying the database again
// basically it builds up a cache using the quriedAccounts param
// since it is pass by pointer
// note it doesn't query the entire account. Just the balance and type
func getAndAccumulateAccount(txKey string, accountId int, queriedAccounts map[int]mdls.FinancialAccount) (mdls.FinancialAccount, error) {
	if accountId == 0 {
		err := errors.New("Tried to query for account but its id is 0, thats wrong... is there a bug somewhere?")
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return mdls.FinancialAccount{}, err
	}
	_, ok := queriedAccounts[accountId]
	if !ok {
		// query for the current account type and balance
		currentAccount, err := dataaccess.ReadFinancialAccountBalanceAndType(txKey, accountId)
		if err != nil {
			return mdls.FinancialAccount{}, err
		}
		currentAccount.Id = accountId
		queriedAccounts[accountId] = currentAccount
	}
	return queriedAccounts[accountId], nil
}
