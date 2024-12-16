package accounting

import (
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	mdls "inventory_money_tracking_software/cmd/models"
	"strings"
)

func isOnlyAllowedDigits(s string) bool {
	// we only allow arabic and english numerics. that's it
	// anything else is not whitelisted and we don't accepted
	// TODO test normalization and make sure everyhting conforms to this
	for _, d := range s {
		if !strings.Contains(cmn.AllowedNumerics, string(d)) {
			return false
		}
	}
	return true
}

// makes sure transactions have account id, journal id and if they're not new have their own
// transaction id
func validateTransactionBeforeJournalEntryUpdate(transactions []mdls.FinancialTransaction,
	journalEntryId int, isNewTransaction bool) error {
	for _, t := range transactions {
		if !t.IsTransactionTypeValid() {
			cmn.Log("Error while updating journal entry, transaction type is invalid %+v",
				cmn.ErrorLevels.Error, t)
			return fmt.Errorf("")
		}
		// skips the transaction id if we're creating a new transaction
		idInvalid := t.Id == 0
		if isNewTransaction {
			idInvalid = false
		}
		if t.FinancialAccountId == 0 {
			cmn.Log("Error transaction has zero account id %+v",
				cmn.ErrorLevels.Error, t)
			return fmt.Errorf("transaction has zero for the account id")
		}
		if idInvalid ||
			t.FinancialJournalEntryId == 0 || t.FinancialJournalEntryId != journalEntryId {
			cmn.Log("Error transaction doesnt have proper id or journal id %+v",
				cmn.ErrorLevels.Error, t)
			return fmt.Errorf("")
		}
	}
	return nil
}

// we dont trust the clinet input data, so we query the transaction by id ourselves,
// and make sure that its journalEntry id indeed matches our expected journal entry id
func doTransactionsBelongToJournalEntry(txKey string, tIds []int, inputtedJid int) bool {
	jIds, err := dataaccess.FindTransactionJournalEntriesByIds(txKey, tIds)
	if err != nil {
		return false
	}
	for _, jId := range jIds {
		if jId != inputtedJid {
			return false
		}
	}
	return true
}

func validateTransactionBeforeJournalEntryCreation(transactions []mdls.FinancialTransaction) cmn.ErrorCode {
	for _, t := range transactions {
		if !t.IsTransactionTypeValid() {
			cmn.Log("Error transaction type is invalid%+v",
				cmn.ErrorLevels.Error, t)
			return cmn.ErrorCodes.InvalidFinancialTransactionType
		}
		if t.FinancialAccountId == 0 {
			cmn.Log("Error transaction doesnt have proper account id %+v",
				cmn.ErrorLevels.Error, t)
			return cmn.ErrorCodes.InvalidFinancialAccountId
		}
	}
	return cmn.ErrorCodes.NoError
}

// if credits don't equal debtis for all transactions, create an imbalance transaction.
// this function will also validate and correct all transactions.
// This function returns the boolean if we indeed created an imbalance transaction, the imbalance transaction
// and an errorcode
// Second return will always be the "imbalance" FinanacialTransaction, but it will be empty if the first return boolean
// is false. And errorCode will be noError if there's no errors
func makeTransactionsBalance(
	ts []mdls.FinancialTransaction) (bool, mdls.FinancialTransaction, cmn.ErrorCode) {
	imbalanceTransaction := mdls.FinancialTransaction{}
	debits := 0
	credits := 0
	// validate transactions and count debits and credits
	for i := 0; i < len(ts); i++ {
		t := ts[i]
		if !t.IsTransactionTypeValid() {
			return false, imbalanceTransaction,
				cmn.AddMsgToErrorCode(cmn.ErrorCodes.InvalidFinancialTransactionType,
					fmt.Sprintf("transaction %+v has invalid transaction type", t))
		}
		t.EnsurePositive()
		if t.TransactionType == cmn.TransactionTypes.Debit {
			debits = debits + t.Amount
		} else if t.TransactionType == cmn.TransactionTypes.Credit {
			credits = credits + t.Amount
		}
	}
	// create imbalance transaction. Assume it is increased debit type
	if debits != credits {
		// ensurePositive below will take care of the other case
		if debits > credits { // imbalance needs to counter whoever is more
			imbalanceTransaction.TransactionType = cmn.TransactionTypes.Credit
		}
		imbalanceTransaction.Amount = debits - credits
		// ensure positive will take care of flipping in case credit was more
		imbalanceTransaction.EnsurePositive()
		return true, imbalanceTransaction, cmn.ErrorCodes.NoError
	}
	return false, imbalanceTransaction, cmn.ErrorCodes.NoError
}

// returns true and the imbalance transaction if it exits. Otherwise Returns false, and an empty transaction
func findImbalanceTransaction(imbalanceAccountId int,
	ts []mdls.FinancialTransaction) (bool, mdls.FinancialTransaction) {
	if imbalanceAccountId == 0 {
		cmn.Log("WhichTransactionIsTheImbalance was called with an imbalance account id of zero, was this intentional? maybe you have a bug somewhere", cmn.LogLevels.Info)
		return false, mdls.FinancialTransaction{}
	}
	for _, t := range ts {
		if t.FinancialAccountId == imbalanceAccountId {
			return true, t
		}
	}
	return false, mdls.FinancialTransaction{}
}

func readTransactionsByAccount(accountId int) ([]mdls.FinancialTransaction, cmn.ErrorCode) {
	if accountId == 0 {
		return nil, cmn.ErrorCodes.InvalidId
	}
	tx, err := dataaccess.ObtainTransactionKey()
	defer dataaccess.RollbackTransaction(tx)
	if err != nil {
		cmn.Log("failed to get database transaction key in readTransactionsByAccount with accountId %d",
			cmn.LogLevels.Error, accountId)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, cmn.ErrorCodes.ReadFailed
	}
	ts, err := dataaccess.ReadFinancialTransactionsByAccount(tx, accountId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, cmn.ErrorCodes.ReadFailed
	}
	return ts, cmn.ErrorCodes.NoError
}

func readTransactionsByJournalEntryId(jId int) ([]mdls.FinancialTransaction, cmn.ErrorCode) {
	if jId == 0 {
		return nil, cmn.ErrorCodes.InvalidId
	}
	tx, err := dataaccess.ObtainTransactionKey()
	defer dataaccess.RollbackTransaction(tx)
	if err != nil {
		cmn.Log("failed to get database transaction key in readTransactionsByJournalEntryId with jId %d",
			cmn.LogLevels.Error, jId)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, cmn.ErrorCodes.ReadFailed
	}
	ts, err := dataaccess.ReadFinancialTransactionsByJournalEntryId(tx, jId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, cmn.ErrorCodes.ReadFailed
	}
	return ts, cmn.ErrorCodes.NoError
}

func readTransactionsByJournalEntries(jIds []int) (map[int][]mdls.FinancialTransaction, cmn.ErrorCode) {
	if len(jIds) == 0 {
		return nil, cmn.ErrorCodes.InvalidId
	}
	tx, err := dataaccess.ObtainTransactionKey()
	defer dataaccess.RollbackTransaction(tx)
	if err != nil {
		cmn.Log("failed to get database transaction key in readTransactionsByJournalEntries with jId %d",
			cmn.LogLevels.Error, jIds)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, cmn.ErrorCodes.ReadFailed
	}
	ts, err := dataaccess.ReadFinancialTransactionsByJournalEntries(tx, jIds)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, cmn.ErrorCodes.ReadFailed
	}
	return ts, cmn.ErrorCodes.NoError
}
