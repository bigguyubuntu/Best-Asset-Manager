package dataaccess

import (
	"database/sql"
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
	"strconv"
	"strings"
)

func insertFinancialGroup(g mdls.FinancialGroup) (string, error) {
	cmn.Log("Inserting new Financial Group to database", cmn.LogLevels.Operation)
	var stmt string
	var grpId string
	var err error
	stmt = insertSql.insertFinancialGroup
	err = db.QueryRow(stmt, g.Name, g.ParentId, g.Description, g.GroupType).Scan(&grpId)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return grpId, err
}

func insertFinancialAccount(a mdls.FinancialAccount) (string, error) {
	cmn.Log("Inserting new Financial Account to database", cmn.LogLevels.Operation)
	stmt := insertSql.insertFinancialAccount
	var acntId string
	err := db.QueryRow(stmt, a.Name, a.Description, a.AccountType, a.GroupId).Scan(&acntId)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return acntId, err
}

func selectAllFinancialGroup() ([]mdls.FinancialGroup, error) {
	cmn.Log("Selecting all finanicalGroups", cmn.LogLevels.Info)
	stmt := readSql.readAllFinancialGroups
	rows, err := db.Query(stmt)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, err
	}
	defer rows.Close()

	var groups []mdls.FinancialGroup
	for rows.Next() {
		var grp mdls.FinancialGroup
		if err := rows.Scan(&grp.Id, &grp.Name, &grp.ParentId, &grp.Description, &grp.GroupType); err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return groups, err
		}
		groups = append(groups, grp)
	}
	if err = rows.Err(); err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
	}
	return groups, nil
}
func selectAllFinancialJournalEntries() ([]mdls.FinancialJournalEntry, error) {
	cmn.Log("Selecting all finanicalJournalEntries", cmn.LogLevels.Info)
	stmt := readSql.readAllFinancialJournalEntries
	rows, err := db.Query(stmt)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, err
	}
	defer rows.Close()
	var js []mdls.FinancialJournalEntry
	for rows.Next() {
		var e mdls.FinancialJournalEntry
		if err := rows.Scan(&e.Id, &e.EntryDate, &e.Description); err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return js, err
		}
		js = append(js, e)
	}
	if err = rows.Err(); err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
	}
	return js, nil
}
func selectFinancialJournalEntriesByIds(ids []int) ([]mdls.FinancialJournalEntry, error) {
	cmn.Log("Selecting all selectFinancialJournalEntriesByIds by ids %v", cmn.LogLevels.Info, ids)
	stmt := readSql.readFinancialJournalEntriesByIds(len(ids))
	s := make([]interface{}, len(ids))
	for i, v := range ids {
		s[i] = v
	}
	rows, err := db.Query(stmt, s...)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, err
	}
	defer rows.Close()
	var js []mdls.FinancialJournalEntry
	for rows.Next() {
		var e mdls.FinancialJournalEntry
		if err := rows.Scan(&e.Id, &e.EntryDate, &e.Description); err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return js, err
		}
		js = append(js, e)
	}
	if err = rows.Err(); err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
	}
	return js, nil
}

func processJournalEntriesCompositeRows(rows *sql.Rows) ([]mdls.FinancialJournalEntryComposite, error) {
	var js []mdls.FinancialJournalEntryComposite

	// we will get multiple of the same journal entry, we don't want to reqrite the je everytime
	// return true if we should create a new composite entry, and false if we should just update the previous one
	putInNewJournalEntry := func(theNewJid int) bool {
		if len(js) == 0 {
			return true
		}
		if js[len(js)-1].Id != theNewJid {
			return true
		}
		return false
	}

	for rows.Next() {
		var e mdls.FinancialJournalEntryComposite
		var t mdls.FinancialTransaction
		var a mdls.FinancialAccount

		if err := rows.Scan(&e.Id, &e.EntryDate, &e.Description,
			&t.Id, &t.Amount, &t.TransactionType,
			&a.Id, &a.Name, &a.Description, &a.Balance, &a.GroupId, &a.AccountType); err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialJournalEntryComposite{}, err
		}
		t.FinancialAccountId = a.Id
		t.FinancialJournalEntryId = e.Id
		if !t.AllRequiredFiledsAreThere() {
			err := fmt.Errorf("transaction has zeros for ids or a zero id somewhere %+v", t)
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialJournalEntryComposite{}, err
		}
		// update the maps
		if putInNewJournalEntry(e.Id) {
			ts := make(map[int]mdls.FinancialTransaction)
			as := make(map[int]mdls.FinancialAccount)
			ts[t.Id] = t
			as[a.Id] = a
			e.Transactions = ts
			e.Accounts = as
			js = append(js, e)
		} else {
			lastEntry := js[len(js)-1]
			_, ok := lastEntry.Transactions[t.Id]
			if !ok {
				lastEntry.Transactions[t.Id] = t
			}
			_, ok = lastEntry.Accounts[a.Id]
			if !ok {
				lastEntry.Accounts[a.Id] = a
			}
		}
	}
	if err := rows.Err(); err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
	}
	return js, nil
}
func selectAllFinancialJournalEntriesComposite() ([]mdls.FinancialJournalEntryComposite, error) {
	cmn.Log("Selecting all finanicalJournalEntriesComposite", cmn.LogLevels.Info)
	stmt := readSql.readAllFinancialJournalEntriesComposite
	rows, err := db.Query(stmt)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return []mdls.FinancialJournalEntryComposite{}, err
	}
	defer rows.Close()
	return processJournalEntriesCompositeRows(rows)
}

func selectFinancialJournalEntryCompositeById(jId int) (mdls.FinancialJournalEntryComposite, error) {
	cmn.Log("Selecting finanicalJournalEntryComposite by id %d", cmn.LogLevels.Info, jId)
	stmt := readSql.readFinancialJournalEntryCompositeById
	rows, err := db.Query(stmt, jId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return mdls.FinancialJournalEntryComposite{}, err
	}
	defer rows.Close()
	compArr, err := processJournalEntriesCompositeRows(rows)
	if len(compArr) == 0 {
		return mdls.FinancialJournalEntryComposite{}, fmt.Errorf("journal entry with id %d doesnt exist", jId)
	}
	return compArr[0], err
}

func selectFinancialJournalEntryCompositeByAccountId(aId int) ([]mdls.FinancialJournalEntryComposite, error) {
	cmn.Log("Selecting finanicalJournalEntryComposite by financial account id %d", cmn.LogLevels.Info, aId)
	stmt := readSql.readFinancialJournalEntryCompositeByAccountId
	rows, err := db.Query(stmt, aId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return []mdls.FinancialJournalEntryComposite{}, err
	}
	defer rows.Close()
	js, err := processJournalEntriesCompositeRows(rows)
	if len(js) == 0 { // if there are no journal entries, we want to grab the account only
		l := []int{aId}
		as, err := selectAccountsByIds(l)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialJournalEntryComposite{}, err
		}
		if len(as) == 0 {
			cmn.HandleError(fmt.Errorf("could not find an account with id %d", aId), cmn.ErrorLevels.Error)
			return []mdls.FinancialJournalEntryComposite{}, err
		}
		var jComposite mdls.FinancialJournalEntryComposite
		jComposite.Accounts = make(map[int]mdls.FinancialAccount)
		jComposite.Accounts[aId] = as[0]
		jComposite.Transactions = make(map[int]mdls.FinancialTransaction)
		js = append(js, jComposite)
	}
	return js, err
}

func selectFinancialGroupWithType(t string) ([]mdls.FinancialGroup, error) {
	cmn.Log(fmt.Sprintf("Selecting finanicalGroupsWithType %s", t), cmn.LogLevels.Info)
	stmt := readSql.readFinanicalGroupsWithType
	rows, err := db.Query(stmt, t)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, err
	}
	defer rows.Close()

	var groups []mdls.FinancialGroup
	for rows.Next() {
		var grp mdls.FinancialGroup
		if err := rows.Scan(&grp.Id, &grp.Name, &grp.ParentId, &grp.Description, &grp.GroupType); err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return groups, err
		}
		groups = append(groups, grp)
	}
	if err = rows.Err(); err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
	}
	return groups, nil
}

func selectFinancialGroupsParentIdAndType(groupId int) (mdls.FinancialGroup, error) {
	cmn.Log(fmt.Sprintf("selectFinancialGroupsParentIdAndType with id %d", groupId), cmn.LogLevels.Info)
	stmt := readSql.readParentIdAndTypeOfFinancialGroup
	var g mdls.FinancialGroup
	err := db.QueryRow(stmt, groupId).Scan(&g.Id, &g.ParentId, &g.GroupType)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return g, err
	}
	return g, nil
}

func selectFinancialGroupByNameAndParentId(tx *sql.Tx, groupName string, parentId int) (mdls.FinancialGroup, error) {
	cmn.Log("select group by name %s and parent id %d", cmn.LogLevels.Operation, groupName, parentId)
	g := mdls.FinancialGroup{Name: groupName, ParentId: parentId}
	stmt := readSql.readFinancialGroupByNameAndParentId
	err := tx.QueryRow(stmt, groupName, parentId).Scan(
		&g.Id, &g.Description, &g.GroupType, &g.Balance)
	if err != sql.ErrNoRows {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
	} else {
		cmn.HandleError(err, cmn.ErrorLevels.Warning)
	}
	return g, err
}

func findOrCreateImbalanceGroup(tx *sql.Tx) (mdls.FinancialGroup, error) {
	cmn.Log("Try and query imbalance group", cmn.LogLevels.Operation)
	g, err := selectFinancialGroupByNameAndParentId(tx, cmn.ImabalanceGroupName, -1)
	// case if it doesn't exists
	if err != nil && (err == sql.ErrNoRows || strings.Contains(err.Error(), "not exist")) {
		cmn.Log("Imbalance group doesnt exist, will create it", cmn.LogLevels.Operation)
		g = mdls.ImbalanceGroup
		gId, err := insertFinancialGroup(g)
		if err != nil {
			cmn.Log("failed at creating imbalance group", cmn.LogLevels.Error)
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return g, err
		}
		gIdInt, err := strconv.Atoi(gId)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		g.Id = gIdInt
	} else if err != nil { // incase we failed at quring and the error is not that it doesn't exist
		cmn.Log("failed to query imbalance group",
			cmn.LogLevels.Operation)
		return g, err
	}
	return g, nil
}

func selectAccountByNameAndGroupId(tx *sql.Tx, accountName string, groupId int) (mdls.FinancialAccount, error) {
	stmt := readSql.readAccountByNameAndGroupId
	var a = mdls.FinancialAccount{Name: accountName, GroupId: groupId}
	a.Name = accountName
	cmn.Log("try and query imbalance account", cmn.LogLevels.Operation)
	err := tx.QueryRow(stmt, accountName, groupId).Scan(&a.Id, &a.Description, &a.Balance, &a.AccountType)
	cmn.HandleError(err, cmn.ErrorLevels.Warning)
	return a, err
}

func selectTransactionByAccountAndJournalEntry(tx *sql.Tx, accountId int,
	journalEntryId int) (mdls.FinancialTransaction, error) {
	cmn.Log("selectTransactionByAccountAndJournalEntry", cmn.LogLevels.Operation)
	stmt := readSql.readFinancialTransactionBelogningToAccountAndJournalEntry
	var t mdls.FinancialTransaction
	t.FinancialAccountId = accountId
	t.FinancialJournalEntryId = journalEntryId
	err := tx.QueryRow(stmt, accountId, journalEntryId).Scan(&t.Id,
		&t.Amount, &t.TransactionType)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return t, err
}
func selectTransactionsByAccount(tx *sql.Tx, accountId int) ([]mdls.FinancialTransaction, error) {
	cmn.Log("selectTransactionByAccount", cmn.LogLevels.Operation)
	ts := []mdls.FinancialTransaction{}
	stmt := readSql.readFinancialTransactionBelogningToAccount
	rows, err := tx.Query(stmt, accountId)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return ts, err
	}
	for rows.Next() {
		var t mdls.FinancialTransaction
		t.FinancialAccountId = accountId
		err = rows.Scan(&t.Id, &t.Amount, &t.TransactionType, &t.FinancialJournalEntryId)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return ts, err
		}
		ts = append(ts, t)
	}
	return ts, err
}

// when we add imbalance transactions they must go in an imbalance account.
// this function find the imbalance account if it exists. Otherwise
// will create it and puts it in an imbalance group.
func findOrCreateImbalanceAccount(tx *sql.Tx) (mdls.FinancialAccount, error) {
	cmn.Log("findOrCreateImbalanceAccount", cmn.LogLevels.Operation)
	imGroup, err := findOrCreateImbalanceGroup(tx)
	if err != nil {
		return mdls.FinancialAccount{}, err
	}
	a, err := selectAccountByNameAndGroupId(tx, cmn.ImabalanceAccountName, imGroup.Id)
	// we create the immbalance account in this case
	if err != nil && (err == sql.ErrNoRows || strings.Contains(err.Error(), "not exist")) {
		cmn.Log("Imbalance account doesn't exists, will create it", cmn.LogLevels.Info)
		a = mdls.ImbalanceAccount
		a.GroupId = imGroup.Id
		aId, err := insertFinancialAccount(a)
		if err != nil {
			cmn.Log("failed at creating imbalance account", cmn.LogLevels.Info)
			return a, err
		}
		aIdint, err := strconv.Atoi(aId)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		a.Id = aIdint
		// we successfully created the imbalance account
	} else if err != nil { // in this case we exit and report error
		return a, err
	}
	cmn.Log("Successfully obtained imbalance account with Id %d", cmn.LogLevels.Operation, a.Id)
	return a, nil
}

// make sure to update account balances after inserting new transactions
func insertNewTransaction(tx *sql.Tx, t mdls.FinancialTransaction) (int, error) {
	stmt := insertSql.insertFinancialTransaction
	if t.FinancialAccountId == 0 || t.FinancialJournalEntryId == 0 {
		m := fmt.Sprintf(
			"Can't insert financial transaction either invalid account id  or journal entry id%+v ",
			t)
		cmn.Log(m, cmn.LogLevels.Critical)
		err := fmt.Errorf(m)
		return 0, err
	}
	// insert the transaction to the db
	cmn.Log("inserting transaction under journalId %d",
		cmn.LogLevels.Operation, t.FinancialJournalEntryId)
	tId := 0
	err := tx.QueryRow(stmt, t.Amount, t.TransactionType,
		t.FinancialJournalEntryId, t.FinancialAccountId).Scan(&tId)
	if err != nil {
		m := fmt.Sprintf("Can't insert financial transaction %+v ", t)
		cmn.Log(m, cmn.LogLevels.Critical)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return 0, err
	}
	return tId, nil
}

// updates an account balance in the database with the balance in the
// passed account
func updateFinancialAccountBalance(tx *sql.Tx, account mdls.FinancialAccount) error {
	if account.Id == 0 {
		err := fmt.Errorf("Account %+v has an invalid id, can't update it", account)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return err
	}
	_, err := tx.Exec(updateSql.updateFinancialAccountBalance,
		account.Balance, account.Id)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return err
	}
	return nil
}
func selectAccountBalanceAndType(tx *sql.Tx, accountId int) (mdls.FinancialAccount, error) {
	cmn.Log("query AccountBalanceAndType for account %d ",
		cmn.LogLevels.Operation, accountId)
	currentAccount := mdls.FinancialAccount{Id: accountId}
	var accountType string
	var balance int
	readAccountBalanceAndType := readSql.readFinancialAccountBalanceAndType
	err := tx.QueryRow(readAccountBalanceAndType, accountId).Scan(&balance, &accountType)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return currentAccount, err
	}
	currentAccount.AccountType = accountType
	currentAccount.Balance = balance
	currentAccount.Id = accountId
	return currentAccount, nil
}

func selectFinancialTransactionsByJournalEntryId(tx *sql.Tx,
	jId int) ([]mdls.FinancialTransaction, error) {
	ts := []mdls.FinancialTransaction{}
	stmt := readSql.readFinancialTransactionsBelongingToJournalEntry
	rows, err := tx.Query(stmt, jId)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return ts, err
	}
	for rows.Next() {
		t := mdls.FinancialTransaction{}
		t.FinancialJournalEntryId = jId
		err = rows.Scan(&t.Id, &t.Amount, &t.TransactionType,
			&t.FinancialAccountId)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return ts, err
		}
		ts = append(ts, t)
	}
	return ts, nil
}
func selectFinancialTransactionsByJournalEntries(tx *sql.Tx,
	jIds []int) (map[int][]mdls.FinancialTransaction, error) {
	//id is journal entry id, value is all transactions under that journal entry
	var ts = make(map[int][]mdls.FinancialTransaction)
	stmt := readSql.readFinancialTransactionsBelongingToJournalEntries(len(jIds))

	s := make([]interface{}, len(jIds))
	for i, v := range jIds {
		s[i] = v
	}
	rows, err := tx.Query(stmt, s...)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return ts, err
	}
	for rows.Next() {
		t := mdls.FinancialTransaction{}
		err = rows.Scan(&t.Id, &t.Amount, &t.TransactionType,
			&t.FinancialAccountId, &t.FinancialJournalEntryId)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return ts, err
		}
		ts[t.FinancialJournalEntryId] = append(ts[t.FinancialJournalEntryId], t)
	}
	return ts, nil
}

// finds all the journal entry ids from the transaction table for the transaction
// ids provided
func selectTransactionJournalEntriesByIds(tx *sql.Tx, ids []int) ([]int, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	stmt := readSql.readTransactionsJournalEntriesByIds(len(ids))
	// need to change the type from int to any
	s := make([]interface{}, len(ids))
	for i, v := range ids {
		s[i] = v
	}
	rows, err := tx.Query(stmt, s...)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, err
	}
	jIds := []int{}
	for rows.Next() {
		jId := 0
		err = rows.Scan(&jId)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return nil, err
		}
		jIds = append(jIds, jId)
	}
	return jIds, nil
}
func selectAccountsByIds(ids []int) ([]mdls.FinancialAccount, error) {
	as := []mdls.FinancialAccount{}
	stmt := readSql.readAccountsByIds(len(ids))
	// need to change the type from int to any
	s := make([]interface{}, len(ids))
	for i, v := range ids {
		s[i] = v
	}
	rows, err := db.Query(stmt, s...)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return as, err
	}
	for rows.Next() {
		a := mdls.FinancialAccount{}
		err = rows.Scan(&a.Id, &a.Name, &a.Description, &a.AccountType, &a.GroupId, &a.Balance)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return as, err
		}
		as = append(as, a)
	}
	return as, nil
}
func selectGroupsByIds(ids []int) ([]mdls.FinancialGroup, error) {
	gs := []mdls.FinancialGroup{}
	stmt := readSql.readGroupsByIds(len(ids))
	// need to change the type from int to any
	s := make([]interface{}, len(ids))
	for i, v := range ids {
		s[i] = v
	}
	rows, err := db.Query(stmt, s...)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return gs, err
	}
	for rows.Next() {
		g := mdls.FinancialGroup{}
		err = rows.Scan(&g.Id, &g.Name, &g.Description, &g.GroupType, &g.ParentId, &g.Balance)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return gs, err
		}
		gs = append(gs, g)
	}
	return gs, nil
}
func selectAllAccounts() ([]mdls.FinancialAccount, error) {
	as := []mdls.FinancialAccount{}
	stmt := readSql.readAllAccounts
	rows, err := db.Query(stmt)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return as, err
	}
	for rows.Next() {
		a := mdls.FinancialAccount{}
		err = rows.Scan(&a.Id, &a.Name, &a.Description, &a.AccountType, &a.GroupId, &a.Balance)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return as, err
		}
		as = append(as, a)
	}
	return as, nil
}

func updateFinancialTransaction(tx *sql.Tx, t mdls.FinancialTransaction) error {
	if !t.AllRequiredFiledsAreThere() {
		return fmt.Errorf("Transaction is invalid, it doesnt have all required fields %+v\n", t)
	}
	_, err := tx.Exec(updateSql.updateTransaction,
		t.Amount, t.FinancialAccountId, t.TransactionType, t.Id)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return err
}
func deleteFinancialTransaction(tx *sql.Tx, t mdls.FinancialTransaction) error {
	_, err := tx.Exec(deleteSql.deleteTransaction, t.Id)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return err
}

// updates the journal entry's information (not transactions) if passed j
// is different from the quried journal entry
func updateJournalEntryInformation(tx *sql.Tx, j mdls.FinancialJournalEntry) error {
	cmn.Log("updateJournalEntryInformation with Id %d",
		cmn.LogLevels.Operation, j.Id)
	if j.Id == 0 {
		err := fmt.Errorf("journal entry has invalid id")
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return err
	}
	if j.EntryDate == "" {
		err := fmt.Errorf("journal entry has invalid entry date")
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return err
	}

	stmt := updateSql.updateJournalEntry
	newEntryDate := j.EntryDate
	newDescription := j.Description
	_, err := tx.Exec(stmt, newEntryDate, newDescription, j.Id)
	if err != nil {
		cmn.Log("failed at updating journal entry id:%d entry date",
			cmn.LogLevels.Error, j.Id)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return err
	}
	return nil
}

func insertJournalEntry(tx *sql.Tx, j mdls.FinancialJournalEntry) (int, error) {
	cmn.Log("InsertJournalEntry", cmn.LogLevels.Operation)
	var stmt string
	var jId int
	var err error

	if j.EntryDate != "" { // if the date was provided.
		stmt = insertSql.insertFinancialJournalEntry
		err = tx.QueryRow(stmt, j.Description, j.EntryDate).Scan(&jId)
	} else { // if it's not provided we pass nothing to use the db default value
		stmt = insertSql.insertFinancialJournalEntryWithoutDate
		err = tx.QueryRow(stmt, j.Description).Scan(&jId)
	}
	if err != nil {
		cmn.Log(fmt.Sprintf("Can't insert journal entry id %+v ", j), cmn.LogLevels.Critical)
		cmn.HandleError(err, cmn.ErrorLevels.Error) // this is major we shouldn't have this error
		return 0, err
	}
	return jId, nil
}

func selectFinancialAccountsWithGroup(groupId int) ([]mdls.FinancialAccount, error) {
	cmn.Log(fmt.Sprintf("selectFinancialAccountsWithGroup with group id %d",
		groupId), cmn.LogLevels.Operation)
	stmt := readSql.readFinanicalAccountsByGroupId
	rows, err := db.Query(stmt, groupId)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, err
	}
	as := []mdls.FinancialAccount{}
	for rows.Next() {
		var a mdls.FinancialAccount
		err = rows.Scan(&a.Id, &a.Name, &a.Description, &a.AccountType, &a.Balance)
		cmn.HandleError(err, cmn.ErrorLevels.Info)
		a.GroupId = groupId
		as = append(as, a)
	}
	return as, nil
}

func selectFinancialGroupsByParentGroup(parentGroupId int) ([]mdls.FinancialGroup, error) {
	cmn.Log(fmt.Sprintf("selectFinancialGroupsByParentGroup with group id %d",
		parentGroupId), cmn.LogLevels.Operation)
	stmt := readSql.readFinanicalGroupsByParentGroupId
	rows, err := db.Query(stmt, parentGroupId)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, err
	}
	gs := []mdls.FinancialGroup{}
	for rows.Next() {
		var g mdls.FinancialGroup
		err = rows.Scan(&g.Id, &g.Name, &g.Description, &g.GroupType, &g.Balance)
		cmn.HandleError(err, cmn.ErrorLevels.Info)
		g.ParentId = parentGroupId
		gs = append(gs, g)
	}
	return gs, nil
}

func isFinancialGroupsWithTypeLessThan2(groupType string) bool {
	// return true if the number of rows with groupType is less than 2
	cmn.Log(fmt.Sprintf("isFinancialGroupsWithTypeLessThan2 for type %s",
		groupType), cmn.LogLevels.Operation)
	stmt := readSql.isNumberOfFinanicalGroupsLessThan2
	numIds := 0
	rows, err := db.Query(stmt, groupType)
	defer rows.Close()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return true // as far as we are concerned we couldn't locate at least 2 rows
	}
	for rows.Next() {
		numIds++
	}
	return numIds < 2
}

func deleteFinancialAccoutsGroup(groupId int) bool {
	// returns true if delete was successful
	cmn.Log(fmt.Sprintf("deleteFinancialAccoutsGroup for id %d", groupId),
		cmn.LogLevels.Operation)
	stmt := deleteSql.deleteGroup
	_, err := db.Exec(stmt, groupId)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return err == nil
}

func selectFinancialGroupById(groupId int) mdls.FinancialGroup {
	cmn.Log(fmt.Sprintf("selectFinancialGroupById with id %d", groupId),
		cmn.LogLevels.Operation)
	stmt := readSql.readFinanicalGroupById
	var grp mdls.FinancialGroup
	err := db.QueryRow(stmt, groupId).Scan(
		&grp.Id, &grp.Name, &grp.ParentId, &grp.Description, &grp.GroupType)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return grp
	}
	return grp
}

func updateFinancialAccoutsGroup(idToUpdate int, g mdls.FinancialGroup) bool {
	cmn.Log(fmt.Sprintf("updateFinancialAccoutsGroup of id with fields %d %+v",
		idToUpdate, g), cmn.LogLevels.Operation)
	stmt := updateSql.updateFinancialGroup
	_, err := db.Exec(stmt, g.ParentId, g.Name, g.Description, idToUpdate)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return false
	}
	return true
}
