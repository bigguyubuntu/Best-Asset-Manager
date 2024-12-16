package mdls

import (
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
)

type FinancialGroup struct {
	dateTracking
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ParentId    int    `json:"parentId"`
	Description string `json:"description"`
	GroupType   string `json:"groupType"`
	Balance     int    `json:"balance"`
}

type FinancialAccount struct {
	dateTracking
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AccountType string `json:"accountType"`
	GroupId     int    `json:"groupId"`
	Balance     int    `json:"balance"`
}

// updates the account's balance after transaction has taken place.
// the account must have at least the account type and balance set.
func (a *FinancialAccount) UpdateAccountBalanceAfterTransaction(t FinancialTransaction) error {
	if a.AccountType == "" {
		err := fmt.Errorf("Cant update an account that doesnt have its account type set %+v", a)
		return err
	}
	// in this case augment the account in case of debit, and reduce it in case of credit
	if a.AccountType == cmn.AccountTypes.DebitIncreased {
		if t.TransactionType == cmn.TransactionTypes.Debit {
			a.Balance = a.Balance + t.Amount
		} else if t.TransactionType == cmn.TransactionTypes.Credit {
			a.Balance = a.Balance - t.Amount
		}
	} else {
		if t.TransactionType == cmn.TransactionTypes.Credit {
			a.Balance = a.Balance + t.Amount
		} else if t.TransactionType == cmn.TransactionTypes.Debit {
			a.Balance = a.Balance - t.Amount
		}
	}
	return nil
}

// if we take the accounts of a balanced journal entry, their balances should be zero
// that is after summing all balances do debits equal credits?
func (a *FinancialAccount) DoAccountsBalance(as []FinancialAccount) bool {
	s := 0
	for _, a := range as {
		if a.AccountType == cmn.AccountTypes.DebitIncreased {
			s = s + a.Balance
		} else if a.AccountType == cmn.AccountTypes.CreditIncreased {
			s = s - a.Balance
		} else {
			cmn.Log("major error, we got an account that has an invalid type %+v",
				cmn.LogLevels.Critical, a)
		}
	}
	return s == 0
}

// deleting a transaction means that we remove the effects of it on the account balance
// meaning if it was a credit transaction and and it is a credit_increase account
// then we got to subtract the amount of that transaction. And vice versa.
// so credit a debit_increase aaccount or debit a credit_increased amount
func (a *FinancialAccount) UpdateAccountBalanceAfterDeletingTransaction(t FinancialTransaction) {
	t.Amount = -1 * t.Amount
	t.EnsurePositive() // negative amount will make it flip transaction type
	a.UpdateAccountBalanceAfterTransaction(t)
}

// A journal entry with all the data related to it.
// All the transactions that are in that entry
// as well as all the accounts related to those transactions.
type FinancialJournalEntryComposite struct {
	FinancialJournalEntry
	Transactions map[int]FinancialTransaction `json:"transactions"`
	Accounts     map[int]FinancialAccount     `json:"accounts"`
}
type FinancialJournalEntry struct {
	dateTracking
	Id          int    `json:"id"`
	EntryDate   string `json:"entryDate"`
	Description string `json:"description"`
}

func (j *FinancialJournalEntry) AreTransactionsBalanced(ts []FinancialTransaction) bool {
	debits := 0
	credits := 0
	// validate transactions and count debits and credits
	for i := 0; i < len(ts); i++ {
		t := ts[i]
		if t.TransactionType == cmn.TransactionTypes.Debit {
			debits = debits + t.Amount
		} else if t.TransactionType == cmn.TransactionTypes.Credit {
			credits = credits + t.Amount
		}
	}
	return debits == credits
}

type FinancialTransaction struct {
	dateTracking
	Id                      int    `json:"id"`
	Amount                  int    `json:"amount"`
	FinancialJournalEntryId int    `json:"financialJournalEntryId"`
	FinancialAccountId      int    `json:"financialAccountId"`
	TransactionType         string `json:"transactionType"` //todo update this to timestame type
}

// we only allow debits or credits financial types. if we get anything else return false
func (t *FinancialTransaction) IsTransactionTypeValid() bool {
	if t.TransactionType == cmn.TransactionTypes.Credit ||
		t.TransactionType == cmn.TransactionTypes.Debit {
		return true
	}
	return false
}

// after reading from the database we want to make
// sure that our transaction is valid for the UI to process
// this check makes sure that we dont have an null required values
func (t *FinancialTransaction) AllRequiredFiledsAreThere() bool {
	if t.IsTransactionTypeValid() && t.FinancialJournalEntryId != 0 &&
		t.Id != 0 && t.FinancialAccountId != 0 {
		return true
	}
	return false
}

// if we have negative amounts, then make them positive and flip the transaction type
// for example if we have amount = -5 and type is debit. Then make it +5 and credit
func (t *FinancialTransaction) EnsurePositive() {
	if t.Amount < 0 {
		t.Amount = -1 * t.Amount
		// flip the financnial type
		if t.TransactionType == cmn.TransactionTypes.Debit {
			t.TransactionType = cmn.TransactionTypes.Credit
		} else {
			t.TransactionType = cmn.TransactionTypes.Debit
		}
	}
}

// fuses the calling t and t2 into a single transaction that would have the final balance of both
// for example if we have two debits, one with 2$ and the other with 5$ then we will
// update the transaction to a debit with 7$.
// and if it's credit 2$ instead then we will update to debit transaction with 3$.
func (t *FinancialTransaction) CombineTransactions(t2 FinancialTransaction) {
	t.EnsurePositive()
	t2.EnsurePositive()
	if t.TransactionType == cmn.TransactionTypes.Debit {
		if t2.TransactionType == cmn.TransactionTypes.Debit {
			t.Amount = t.Amount + t2.Amount
		} else {
			t.Amount = t.Amount - t2.Amount
		}
	} else { // case when t is a credit instead
		if t2.TransactionType == cmn.TransactionTypes.Debit {
			t.Amount = t.Amount - t2.Amount
		} else {
			t.Amount = t.Amount + t2.Amount
		}
	}
	t.EnsurePositive()
}

// checks the fields of 2 transactions and tells you if they have equal values
func (t *FinancialTransaction) IsTransactionEqualTo(t2 FinancialTransaction) bool {
	if t.Amount == t2.Amount &&
		t.CreatedAt == t2.CreatedAt &&
		t.FinancialAccountId == t2.FinancialAccountId &&
		t.TransactionType == t2.TransactionType &&
		t.FinancialJournalEntryId == t2.FinancialJournalEntryId &&
		t.Id == t2.Id {
		return true
	}
	return false
}

func (t *FinancialTransaction) IsThereEqualTransactionInSlice(l []FinancialTransaction) bool {
	for _, x := range l {
		if t.IsTransactionEqualTo(x) {
			return true
		}
	}
	return false
}

// default values

var ImbalanceGroup = FinancialGroup{
	Name:        cmn.ImabalanceGroupName,
	GroupType:   cmn.GroupTypes.Expenses,
	Description: "Placeholder group for the imbalance account. You had imbalanced transactions (debits and credits were not equal) so the difference was put in an imbalance account",
	ParentId:    -1,
}

// don't forget to add Group Id before comitting to the database
var ImbalanceAccount = FinancialAccount{
	Name:        cmn.ImabalanceAccountName,
	AccountType: cmn.AccountTypes.DebitIncreased,
	Description: "All transactions that don't have a counter account go here. It's a double entry accounting right"}
