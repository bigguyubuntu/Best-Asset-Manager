package cmn

// these would ideally sit the accounting package
// but since these are needed in other packages too, like dataaccess
// we put them in common package to avoid circular dependancies

// accounting constants

const ImabalanceAccountName = "Imbalance"
const ImabalanceGroupName = "Imbalance"

type groupTypes struct {
	Assets      string
	Expenses    string
	Income      string
	Liabilities string
	Equity      string
}
type accountTypes struct {
	DebitIncreased  string
	CreditIncreased string
}
type transactionTypes struct {
	Debit  string
	Credit string
}

var GroupTypes = groupTypes{
	Assets:      "assets",
	Expenses:    "expenses",
	Income:      "income",
	Liabilities: "liabilities",
	Equity:      "equity",
}

var AccountTypes = accountTypes{
	DebitIncreased:  "debit_increased",
	CreditIncreased: "credit_increased",
}

var TransactionTypes = transactionTypes{
	Debit:  "debit",
	Credit: "credit",
}
