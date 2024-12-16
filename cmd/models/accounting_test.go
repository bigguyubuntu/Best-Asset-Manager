package mdls

import (
	cmn "inventory_money_tracking_software/cmd/common"
	"testing"
)

func TestIsTransactionTypeValid(t *testing.T) {
	tests := [5]FinancialTransaction{
		{TransactionType: "potato"},
		{TransactionType: ""},
		{TransactionType: "\""},
		{TransactionType: cmn.TransactionTypes.Debit},
		{TransactionType: cmn.TransactionTypes.Credit},
	}
	expected := [5]bool{false, false, false, true, true}
	for i := 0; i < len(tests); i++ {
		if expected[i] != tests[i].IsTransactionTypeValid() {
			t.Errorf(
				"Transaction type validation gave wrong results. Expected: %v, received (%v)\n",
				expected[i], tests[i].IsTransactionTypeValid())
		}
	}
}
func TestEnsurePositive(test *testing.T) {
	ts := [6]FinancialTransaction{
		{Amount: -5, TransactionType: cmn.TransactionTypes.Credit},
		{Amount: -5, TransactionType: cmn.TransactionTypes.Debit},
		{Amount: 5, TransactionType: cmn.TransactionTypes.Debit},
		{Amount: 5, TransactionType: cmn.TransactionTypes.Credit},
		{Amount: 0, TransactionType: cmn.TransactionTypes.Debit},
		{Amount: 0, TransactionType: cmn.TransactionTypes.Credit},
	}
	es := [6]FinancialTransaction{
		{Amount: 5, TransactionType: cmn.TransactionTypes.Debit},
		{Amount: 5, TransactionType: cmn.TransactionTypes.Credit},
		{Amount: 5, TransactionType: cmn.TransactionTypes.Debit},
		{Amount: 5, TransactionType: cmn.TransactionTypes.Credit},
		{Amount: 0, TransactionType: cmn.TransactionTypes.Debit},
		{Amount: 0, TransactionType: cmn.TransactionTypes.Credit},
	}
	for i := 0; i < len(es); i++ {
		ts[i].EnsurePositive()
		if ts[i].Amount != es[i].Amount ||
			ts[i].TransactionType != es[i].TransactionType {
			test.Errorf(
				"Amount and type were not proper. Expected: %v, received (%v)\n",
				es[i], ts[i])
		}
	}
}

func TestIsThereEqualTransactionInSlice(t *testing.T) {
	es := []FinancialTransaction{
		{Id: 1, Amount: 5, TransactionType: cmn.TransactionTypes.Debit},
		{Id: 2, Amount: 5, TransactionType: cmn.TransactionTypes.Credit},
		{Id: 3, Amount: 7, TransactionType: cmn.TransactionTypes.Debit},
		{Id: 4, Amount: 7, TransactionType: cmn.TransactionTypes.Credit},
	}
	ts := FinancialTransaction{Id: 1, Amount: 5, TransactionType: cmn.TransactionTypes.Debit}

	if !ts.IsThereEqualTransactionInSlice(es) {
		t.Errorf("Expected that there is an equal for %+v\n", ts)
	}
	ts = FinancialTransaction{Id: 1, Amount: 10000, TransactionType: cmn.TransactionTypes.Debit}
	if ts.IsThereEqualTransactionInSlice(es) {
		t.Errorf("Expected that there is no equal for %+v\n", ts)
	}
	ts = FinancialTransaction{Id: 1, Amount: 5, TransactionType: cmn.TransactionTypes.Credit}
	if ts.IsThereEqualTransactionInSlice(es) {
		t.Errorf("Expected that there is no equal for %+v\n", ts)
	}
}

func TestUpdateAccountBalanceAfterTransaction(t *testing.T) {
	a1 := FinancialAccount{
		AccountType: cmn.AccountTypes.CreditIncreased,
		Balance:     0,
	}
	a2 := FinancialAccount{
		AccountType: cmn.AccountTypes.DebitIncreased,
		Balance:     0,
	}
	t1 := FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Credit,
		Amount:          10,
	}
	t2 := FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Debit,
		Amount:          10,
	}
	a1.UpdateAccountBalanceAfterTransaction(t1)
	if a1.Balance != t1.Amount {
		t.Errorf("Expected that credit_increased balance will increase by credit transaction")
	}
	a1.UpdateAccountBalanceAfterTransaction(t2)
	if a1.Balance != 0 {
		t.Errorf("Expected that credit_increased balance will decrease by debit transaction")
	}
	a2.UpdateAccountBalanceAfterTransaction(t1)
	if a2.Balance != -1*t1.Amount {
		t.Errorf("Expected that debit_increased balance will decrease by credit transaction")
	}
	a2.UpdateAccountBalanceAfterTransaction(t2)
	if a2.Balance != 0 {
		t.Errorf("Expected that debit_increased balance will increase by debit transaction")
	}

}
func TestUpdateAccountBalanceAfterDeletingTransaction(t *testing.T) {
	a1 := FinancialAccount{
		AccountType: cmn.AccountTypes.CreditIncreased,
		Balance:     0,
	}
	a2 := FinancialAccount{
		AccountType: cmn.AccountTypes.DebitIncreased,
		Balance:     0,
	}
	t1 := FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Credit,
		Amount:          10,
	}
	t2 := FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Debit,
		Amount:          10,
	}
	// balance is zero and we delete t1. so we need to subtract t1 from balance
	// making it negative t1.Amount
	a1.UpdateAccountBalanceAfterDeletingTransaction(t1)
	if a1.Balance != -1*t1.Amount {
		t.Errorf("Expected that credit_increased balance will decrease by deleting credit transaction")
	}
	// t2 is a debit transaction and a1 is a credit account. So to delete t2
	// we need to credit a1 with t2's amount.
	a1.UpdateAccountBalanceAfterDeletingTransaction(t2)
	if a1.Balance != 0 {
		t.Errorf("Expected that credit_increased balance will increase by deleting debit transaction")
	}
	a2.UpdateAccountBalanceAfterDeletingTransaction(t1)
	if a2.Balance != t1.Amount {
		t.Errorf("Expected that debit_increased balance will increase by deleting credit transaction")
	}
	a2.UpdateAccountBalanceAfterDeletingTransaction(t2)
	if a2.Balance != 0 {
		t.Errorf("Expected that debit_increased balance will decrease by deleting debit transaction")
	}
}

func TestCombineTransactions(t *testing.T) {
	t1 := FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Debit,
		Amount:          10,
	}
	t2 := FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Credit,
		Amount:          10,
	}
	//-- start by making them equal and see if they cancel each other out
	// while maintaining the original transaction type
	t1.CombineTransactions(t2)
	if t1.Amount != 0 || t1.TransactionType != cmn.TransactionTypes.Debit {
		t.Errorf("Expected amount of 0 and type to be debit, but go: %+v\n", t1)
	}
	// repeat but for t2
	t1.Amount = 10
	t2.CombineTransactions(t1)
	if t2.Amount != 0 || t2.TransactionType != cmn.TransactionTypes.Credit {
		t.Errorf("Expected amount of 0 and type to be credit, but go: %+v\n", t2)
	}

	//-- now we make credit more than debit
	// debit is 10 and credit is 30
	t2.Amount = 30
	t1.CombineTransactions(t2)
	if t1.Amount != 20 || t1.TransactionType != cmn.TransactionTypes.Credit {
		t.Errorf("Expected amount of 20 and type to be credit, but go: %+v\n", t1)
	}
	// repeat for t2
	t1 = FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Debit,
		Amount:          10,
	}
	t2.CombineTransactions(t1)
	if t2.Amount != 20 || t2.TransactionType != cmn.TransactionTypes.Credit {
		t.Errorf("Expected amount of 20 and type to be credit, but go: %+v\n", t2)
	}

	//-- we make the debit more than the credit now

	t1 = FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Debit,
		Amount:          100,
	}
	t2 = FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Credit,
		Amount:          10,
	}

	t1.CombineTransactions(t2)
	if t1.Amount != 90 || t1.TransactionType != cmn.TransactionTypes.Debit {
		t.Errorf("Expected amount of 90 and type to be debit, but go: %+v\n", t1)
	}
	t1.Amount = 100
	t2.CombineTransactions(t1)
	if t2.Amount != 90 || t2.TransactionType != cmn.TransactionTypes.Debit {
		t.Errorf("Expected amount of 90 and type to be debit, but go: %+v\n", t2)
	}

	//-- now we make them both debit
	t1 = FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Debit,
		Amount:          100,
	}
	t2 = FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Debit,
		Amount:          10,
	}

	t1.CombineTransactions(t2)
	if t1.Amount != 110 || t1.TransactionType != cmn.TransactionTypes.Debit {
		t.Errorf("Expected amount of 110 and type to be debit, but go: %+v\n", t1)
	}

	//-- finally both are credit
	t1 = FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Credit,
		Amount:          100,
	}
	t2 = FinancialTransaction{
		TransactionType: cmn.TransactionTypes.Credit,
		Amount:          10,
	}

	t1.CombineTransactions(t2)
	if t1.Amount != 110 || t1.TransactionType != cmn.TransactionTypes.Credit {
		t.Errorf("Expected amount of 110 and type to be credit, but go: %+v\n", t1)
	}

}
