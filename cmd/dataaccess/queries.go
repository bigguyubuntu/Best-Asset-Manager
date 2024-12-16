package dataaccess

import "fmt"

type selectQueries struct {
	readItem                                  string
	readAllItems                              string
	readStorehouse                            string
	readStorehouseByName                      string
	readStorehouseItemQuantity                string
	readInventoryMovement                     string
	readInventoryMovementByInventoryJournalId string
	// Accounting below
	readAllFinancialGroups                                    string
	readFinanicalGroupsWithType                               string
	isNumberOfFinanicalGroupsLessThan2                        string
	readParentIdAndTypeOfFinancialGroup                       string
	readFinanicalGroupById                                    string
	readFinanicalAccountsByGroupId                            string
	readFinanicalGroupsByParentGroupId                        string
	readGroupsByIds                                           func(int) string
	readFinancialAccountBalanceAndType                        string
	readAccountByNameAndGroupId                               string
	readFinancialGroupByNameAndParentId                       string
	readFinancialJournalEntry                                 string
	readFinancialTransactionsBelongingToJournalEntry          string
	readFinancialTransactionsBelongingToJournalEntries        func(int) string
	readFinancialTransactionBelogningToAccountAndJournalEntry string
	readFinancialTransactionBelogningToAccount                string
	readAccountsByIds                                         func(int) string
	readAllAccounts                                           string
	readAllFinancialJournalEntries                            string
	readFinancialJournalEntriesByIds                          func(int) string
	readAllFinancialJournalEntriesComposite                   string
	readFinancialJournalEntryCompositeById                    string
	readFinancialJournalEntryCompositeByAccountId             string
	readTransactionsJournalEntriesByIds                       func(int) string
	// general below
	readTagById                       string
	readTagsByIds                     func(int) string
	readTagsByFinancialJournalEntryId string
}
type insertQueries struct {
	insertItem                   string
	insertStorehouse             string
	insertInventoryMovement      string
	insertStorehouseItemQuantity string
	insertInventoryJournalEntry  string
	// accounting below
	insertFinancialAccount                 string
	insertFinancialGroup                   string
	insertFinancialTransaction             string
	insertFinancialJournalEntry            string
	insertFinancialJournalEntryWithoutDate string
	// general queries below
	insertTag                      string
	insertTagFinancialJournalEntry string
}

type deleteQueries struct {
	deleteGroup       string
	deleteTransaction string
}

type updateQueries struct {
	updateItem                   string
	updateItemCount              string
	updateStorehouse             string
	updateInventoryMovement      string
	updateStorehouseItemQuantity string
	// accounting update is below
	updateFinancialGroup          string
	updateFinancialAccountBalance string
	updateJournalEntry            string
	updateTransaction             string
}

// queries
var insertSql = insertQueries{
	insertItem: `INSERT INTO Item (item_name, description, image_link, inventory_unit)
	VALUES ($1, $2, $3, $4) RETURNING ID;`,
	insertStorehouse:             `INSERT INTO storehouse (storehouse_name, description, storehouse_type) VALUES ($1, $2, $3) RETURNING id;`,
	insertInventoryMovement:      `INSERT INTO inventory_movement(storehouse, item, quantity, inventory_movement_type, inventory_journal_entry) VALUES ($1, $2, $3, $4, $5);`,
	insertStorehouseItemQuantity: `INSERT INTO storehouse_item_quantity(storehouse, item, quantity) VALUES ($1, $2,$3);`,
	insertInventoryJournalEntry: `INSERT INTO 
	inventory_journal_entry(description) VALUES ($1) RETURNING id;`,
	// accounting is below
	insertFinancialGroup: `INSERT INTO Financial_group
	(group_name, parent_id, description, financial_group_type) VALUES($1, $2, $3, $4)
	 RETURNING ID;`,
	insertFinancialAccount: `INSERT INTO Financial_account 
	(account_name, description, financial_account_type, financial_group)
	 VALUES($1, $2, $3, $4) RETURNING ID;`,
	insertFinancialJournalEntry: `INSERT INTO Financial_Journal_Entry(description, entry_date)
	VALUES ($1, $2) RETURNING id`,
	// uses the db default date
	insertFinancialJournalEntryWithoutDate: `INSERT INTO Financial_Journal_Entry(description)
	VALUES ($1) RETURNING id`,
	insertFinancialTransaction: `INSERT INTO Financial_Transaction
	(amount, financial_transaction_type, financial_journal_entry, financial_account)
	VALUES ($1, $2, $3, $4) RETURNING id`,
	insertTag: `INSERT INTO tag (tag_name, description) VALUES($1, $2) RETURNING id;`,
	// general is below
	insertTagFinancialJournalEntry: `INSERT INTO tag_financial_journal_entry(tag, financial_journal_entry)
	 VALUES($1, $2)`,
}

var deleteSql = deleteQueries{
	deleteGroup:       `DELETE FROM financial_group WHERE id=$1`,
	deleteTransaction: `DELETE FROM financial_transaction WHERE id=$1`,
}

var readSql = selectQueries{
	// general below
	readTagById: `SELECT tag_name, description FROM tag where id=$1;`,
	readTagsByFinancialJournalEntryId: `SELECT tag.id, tag.tag_name, tag.Description FROM tag
	LEFT JOIN tag_financial_journal_entry ON
	tag_financial_journal_entry.tag=tag.id
	WHERE
	tag_financial_journal_entry.financial_journal_entry=$1;
	`,
	//inventory below
	readItem:                   `SELECT item_name, description, image_link, count, inventory_unit FROM Item WHERE id=$1`,
	readAllItems:               `select * from Items`,
	readStorehouse:             `SELECT storehouse_name, description, storehouse_type from storehouse WHERE ID=$1`,
	readStorehouseByName:       `SELECT id, description, storehouse_type from storehouse WHERE storehouse_name=$1`,
	readStorehouseItemQuantity: `SELECT quantity from storehouse_item_quantity WHERE storehouse=$1 AND item=$2`,
	readInventoryMovement:      `SELECT storehouse, item, quantity, inventory_movement_type, inventory_journal_entry WHERE id=$1`,
	readInventoryMovementByInventoryJournalId: `SELECT id, storehouse, item, quantity, inventory_movement_type FROM inventory_movement WHERE inventory_journal_entry=$1`,
	// accounting below
	readAllFinancialGroups:              "select id, group_name, parent_id, description, financial_group_type from financial_group ORDER BY group_name",
	readFinanicalGroupsWithType:         `select id, group_name, parent_id, description, financial_group_type from Financial_group where financial_group_type=$1`,
	readFinanicalGroupById:              `select id, group_name, parent_id, description, financial_group_type from Financial_group where id=$1`,
	isNumberOfFinanicalGroupsLessThan2:  "select id from financial_group where financial_group_type=$1 LIMIT 2",
	readParentIdAndTypeOfFinancialGroup: `select id, parent_id, financial_group_type from financial_group where id=$1`,
	// depricated, the name is no longer unique. TODO
	readFinancialGroupByNameAndParentId:              `SELECT id, description, financial_group_type, balance from financial_group WHERE group_name=$1 AND parent_id=$2`,
	readFinanicalGroupsByParentGroupId:               `SELECT id, group_name, description, financial_group_type, balance from financial_group WHERE parent_id=$1 ORDER BY group_name`,
	readFinanicalAccountsByGroupId:                   `SELECT id, account_name, description, financial_account_type, balance FROM financial_account WHERE financial_group=$1 ORDER BY account_name;`,
	readFinancialAccountBalanceAndType:               `SELECT balance, financial_account_type FROM financial_account WHERE id=$1`,
	readAccountByNameAndGroupId:                      `select id, description, balance, financial_account_type FROM financial_account WHERE account_name=$1 AND financial_group=$2`,
	readFinancialJournalEntry:                        "select id, entry_date, description FROM financial_journal_entry WHERE id=$1",
	readFinancialTransactionsBelongingToJournalEntry: `SELECT id, amount, financial_transaction_type, financial_account FROM financial_transaction WHERE financial_journal_entry=$1`,
	readFinancialTransactionsBelongingToJournalEntries: func(numIds int) string {
		p := getEnoughPlaceHolders(numIds)
		return fmt.Sprintf(
			`SELECT id, amount, financial_transaction_type, financial_account, financial_journal_entry
			FROM financial_transaction
		 	WHERE financial_journal_entry IN (%s);`, p)
	},
	readFinancialTransactionBelogningToAccountAndJournalEntry: `SELECT id, amount, financial_transaction_type FROM financial_transaction WHERE financial_journal_entry=$1 AND financial_account=$2`,
	readFinancialTransactionBelogningToAccount:                `SELECT id, amount, financial_transaction_type, financial_journal_entry FROM financial_transaction WHERE financial_account=$1`,
	readAllAccounts: "SELECT id, account_name, description, financial_account_type, financial_group, balance FROM financial_account;",
	readAccountsByIds: func(numIds int) string {
		p := getEnoughPlaceHolders(numIds)
		return fmt.Sprintf("select id, account_name, description, financial_account_type, financial_group, balance FROM financial_account WHERE id IN (%s);", p)
	},
	readGroupsByIds: func(numIds int) string {
		p := getEnoughPlaceHolders(numIds)
		return fmt.Sprintf("select id, group_name, description, financial_group_type, parent_id, balance FROM financial_group WHERE id IN (%s);", p)
	},
	readTransactionsJournalEntriesByIds: func(numIds int) string {
		p := getEnoughPlaceHolders(numIds)
		return fmt.Sprintf("SELECT financial_journal_entry FROM financial_transaction WHERE id IN (%s);", p)
	},
	readFinancialJournalEntriesByIds: func(numIds int) string {
		p := getEnoughPlaceHolders(numIds)
		return fmt.Sprintf("SELECT id, entry_date, description FROM financial_journal_entry WHERE id IN (%s) ORDER BY entry_date ASC;", p)
	},
	readAllFinancialJournalEntries: "SELECT id, entry_date, description FROM financial_journal_entry ORDER BY entry_date ASC;",
	readAllFinancialJournalEntriesComposite: `
SELECT financial_journal_entry.id, financial_journal_entry.entry_date,
    financial_journal_entry.description,
    financial_transaction.id AS transaction_id, financial_transaction.amount,
    financial_transaction.financial_transaction_type,
    financial_account.id AS account_id, financial_account.account_name,
    financial_account.description AS account_description,
    financial_account.balance, financial_account.financial_group,
    financial_account.financial_account_type
FROM financial_journal_entry 
LEFT JOIN financial_transaction ON
    financial_journal_entry.id=financial_transaction.financial_journal_entry
LEFT JOIN financial_account ON financial_account.id=financial_transaction.financial_account
ORDER BY financial_journal_entry.id;`,
	readFinancialJournalEntryCompositeById: `
SELECT financial_journal_entry.id, financial_journal_entry.entry_date,
    financial_journal_entry.description,
    financial_transaction.id AS transaction_id, financial_transaction.amount,
    financial_transaction.financial_transaction_type,
    financial_account.id AS account_id, financial_account.account_name,
    financial_account.description AS account_description,
    financial_account.balance, financial_account.financial_group,
    financial_account.financial_account_type
FROM financial_journal_entry 
LEFT JOIN financial_transaction ON
    financial_journal_entry.id=financial_transaction.financial_journal_entry
LEFT JOIN financial_account ON financial_account.id=financial_transaction.financial_account
WHERE financial_journal_entry.id=$1;
`,
	// the order by desc makes sure that the newest journal entries are up at the top
	readFinancialJournalEntryCompositeByAccountId: `
SELECT financial_journal_entry.id, financial_journal_entry.entry_date,
    financial_journal_entry.description,
    financial_transaction.id AS transaction_id, financial_transaction.amount,
    financial_transaction.financial_transaction_type,
    financial_account.id AS account_id, financial_account.account_name,
    financial_account.description AS account_description,
    financial_account.balance, financial_account.financial_group,
    financial_account.financial_account_type
FROM financial_journal_entry 
LEFT JOIN financial_transaction ON
    financial_journal_entry.id=financial_transaction.financial_journal_entry
LEFT JOIN financial_account ON financial_account.id=financial_transaction.financial_account
WHERE financial_journal_entry.id IN (
    SELECT financial_journal_entry from financial_transaction WHERE financial_account=$1
) ORDER BY financial_journal_entry.entry_date DESC, financial_journal_entry.id;
`,
}

var updateSql = updateQueries{
	updateItem:                   `UPDATE item SET item_name=$1, description=$2, image_link=$3, inventory_unit=$4 WHERE id=$5`,
	updateItemCount:              `UPDATE item SET count=$1 WHERE id=$2`,
	updateStorehouseItemQuantity: `UPDATE storehouse_item_quantity SET quantity=$1 WHERE storehouse=$2 AND item=$3;`,
	// accounting below
	updateFinancialGroup: `UPDATE financial_group
	SET  parent_id=$1, group_name=$2, description=$3 WHERE id=$4`,
	updateFinancialAccountBalance: `UPDATE financial_account SET balance=$1 WHERE id=$2`,
	updateJournalEntry:            "UPDATE financial_journal_entry SET entry_date=$1, description=$2 WHERE id=$3;",
	updateTransaction:             "UPDATE financial_transaction SET amount=$1, financial_account=$2, financial_transaction_type=$3 WHERE id=$4;",
}

// generates enough place holders to use in the query. this should be
// safe as I am not taking any value from the client directly.
func getEnoughPlaceHolders(n int) string {
	if n <= 0 {
		return ""
	}
	s := ""
	for i := 1; i < n+1; i++ {
		s = s + fmt.Sprintf("$%d,", i)
	}
	s = s[0 : len(s)-1] // remove last comma
	return s
}
