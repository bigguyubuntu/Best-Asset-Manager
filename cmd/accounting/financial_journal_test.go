package accounting

import (
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
	"testing"
)

func TestMakeTransactionsBalance(t *testing.T) {
	amount := 5000
	t1 := mdls.FinancialTransaction{
		Amount:          amount,
		TransactionType: cmn.TransactionTypes.Credit,
	}
	t2 := mdls.FinancialTransaction{
		Amount:          amount,
		TransactionType: cmn.TransactionTypes.Debit,
	}
	// case one, we don't create imbalance account when we give it
	// balanced transactions
	createdImbalanceT, _, errorCode :=
		makeTransactionsBalance([]mdls.FinancialTransaction{t1, t2})

	if createdImbalanceT || errorCode != cmn.ErrorCodes.NoError {
		t.Error("Expected that no imbalance transaction is created for balanced transaction")
	}

	// case two, create imbalance account when we give it imbalanced debits and credits
	// add a debit transaction, so now the debits are more than credits
	// so imbalance transaction amount must be credit type and its amount equal to amount2
	amount2 := 2000
	t3 := mdls.FinancialTransaction{
		Amount:          amount2,
		TransactionType: cmn.TransactionTypes.Debit,
	}
	createdImbalanceT, imbalancedT, _ :=
		makeTransactionsBalance([]mdls.FinancialTransaction{t1, t2, t3})

	if !createdImbalanceT || imbalancedT.Amount != amount2 ||
		imbalancedT.TransactionType != cmn.TransactionTypes.Credit {
		t.Errorf(
			"Expected to create imbalance transaction of type credit and amount 2000, but got %+v",
			imbalancedT)
	}
	// case 3 same as case two but with opposite transaction types
	t4 := mdls.FinancialTransaction{
		Amount:          amount2,
		TransactionType: cmn.TransactionTypes.Credit,
	}
	createdImbalanceT, imbalancedT, _ =
		makeTransactionsBalance([]mdls.FinancialTransaction{t1, t2, t4})

	if !createdImbalanceT || imbalancedT.Amount != amount2 ||
		imbalancedT.TransactionType != cmn.TransactionTypes.Debit {
		t.Errorf(
			"Expected to create imbalance transaction of type credit and amount 2000, but got %+v",
			imbalancedT)
	}

	// case 4 - same as case 1 but with more fields
	originalTransactions := []mdls.FinancialTransaction{{
		Id:                      0,
		Amount:                  50,
		FinancialJournalEntryId: 0,
		FinancialAccountId:      1,
		TransactionType:         cmn.TransactionTypes.Credit},
		{
			Id:                      0,
			Amount:                  50,
			FinancialJournalEntryId: 0,
			FinancialAccountId:      2,
			TransactionType:         cmn.TransactionTypes.Debit},
	}

	createdImbalanceT, _, errorCode = makeTransactionsBalance(originalTransactions)
	if createdImbalanceT || errorCode != cmn.ErrorCodes.NoError {
		t.Error("Expected that no imbalance transaction is created for balanced transaction")
	}

	// case 5 extra case where we expect imbalance
	originalTransactions = []mdls.FinancialTransaction{{
		TransactionType: cmn.TransactionTypes.Debit, Amount: 50000,
	},
		{
			TransactionType: cmn.TransactionTypes.Credit, Amount: 250000,
		},
	}
	createdImbalanceT, imbalancedT, errorCode = makeTransactionsBalance(originalTransactions)
	if !createdImbalanceT || errorCode != cmn.ErrorCodes.NoError {
		t.Error("Expected an imbalance transaction is created for imbalanced transaction")
	}
	expected := 250000 - 50000
	if imbalancedT.Amount != expected {
		t.Errorf("expected amount to be %d but it was %d", expected, imbalancedT.Amount)
	}
	if imbalancedT.TransactionType != cmn.TransactionTypes.Debit {
		t.Errorf("Expected transaction type to be debit but it was %v", imbalancedT.TransactionType)
	}
	originalTransactions = append(originalTransactions, imbalancedT)
	j := mdls.FinancialJournalEntry{}
	if !j.AreTransactionsBalanced(originalTransactions) {
		t.Error("expected that all transactions will be balanced once we include the imbalance transaction")
	}

}
