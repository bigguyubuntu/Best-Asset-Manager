package accounting

import (
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	mdls "inventory_money_tracking_software/cmd/models"
	"strings"
)

// creates a new journal entry in the database
// inserts journal entry j and inserts the transactions attached. if it's imbalanced
// then we append imbalance transsaction to the list
// we need to create the imbalance account unless it already exists.
// Note that the journal entry is balanced when put in the db,
// because of makeTransactionsBalance
// will balance debits and credits and put any imbalance in the imbalance transaction
func CreateJournalEntry(j mdls.FinancialJournalEntry,
	ts []mdls.FinancialTransaction) (int, []mdls.FinancialTransaction, cmn.ErrorCode) {
	if len(ts) < 2 {
		return 0, nil, cmn.ErrorCodes.NotEnoughTransactions
	}
	errCode := validateTransactionBeforeJournalEntryCreation(ts)
	if errCode != cmn.ErrorCodes.NoError {
		return 0, nil, errCode
	}
	txKey, err := dataaccess.ObtainTransactionKey()
	defer dataaccess.RollbackTransaction(txKey)
	if err != nil {
		cmn.Log("Failed to obtain transaction key while creating a new journal entry",
			cmn.LogLevels.Error)
		return 0, nil, cmn.ErrorCodes.CreateFailed
	}
	jId, err := dataaccess.InsertJournalEntry(txKey, j)
	if err != nil {
		return 0, nil, cmn.ErrorCodes.CreateFailed
	}
	j.Id = jId
	createdImbalanceTransaction, imbalanceTransaction, errCode := makeTransactionsBalance(ts)
	if errCode != cmn.ErrorCodes.NoError {
		return 0, nil, errCode
	}
	if createdImbalanceTransaction {
		ts = append(ts, imbalanceTransaction)
		err = handleImbalanceAccount(txKey, ts)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return 0, nil, cmn.ErrorCodes.CreateFailed
		}
	}
	queriedAccounts := make(map[int]mdls.FinancialAccount)
	modifiedAccountIds := []int{}
	createdTransactions := []mdls.FinancialTransaction{}
	// insert the new transactions
	for _, t := range ts {
		t.FinancialJournalEntryId = jId
		a, err := getAndAccumulateAccount(txKey, t.FinancialAccountId, queriedAccounts)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return 0, nil, cmn.ErrorCodes.CreateFailed
		}
		tId, err := dataaccess.CreateNewFinancialTransaction(txKey, t)
		t.Id = tId
		createdTransactions = append(createdTransactions, t)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return 0, nil, cmn.ErrorCodes.CreateFailed
		}
		err = a.UpdateAccountBalanceAfterTransaction(t)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return 0, nil, cmn.ErrorCodes.CreateFailed
		}
		queriedAccounts[a.Id] = a // make sure we keep track of account updates
		modifiedAccountIds = append(modifiedAccountIds, a.Id)
	}

	// updated the balances of modifiied accounts
	for _, aId := range modifiedAccountIds {
		a, ok := queriedAccounts[aId]
		if !ok {
			cmn.Log("failed to update the balance of a modified account",
				cmn.ErrorLevels.Error)
			return 0, nil, cmn.ErrorCodes.CreateFailed
		}
		err = dataaccess.UpdateFinancialAccountBalance(txKey, a)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return 0, nil, cmn.ErrorCodes.CreateFailed
		}
	}
	if dataaccess.CommitTransaction(txKey) != nil {
		return 0, nil, cmn.ErrorCodes.CreateFailed
	}
	return jId, createdTransactions, cmn.ErrorCodes.NoError
}

func UpdateJournalEntryTransactionsWrapper(j mdls.FinancialJournalEntry, updatedTransactions,
	createdTransactions []mdls.FinancialTransaction,
	deletedTransactions []int) ([]mdls.FinancialTransaction, cmn.ErrorCode) {
	tx, err := dataaccess.ObtainTransactionKey()
	defer dataaccess.RollbackTransaction(tx)
	if err != nil {
		cmn.Log("failed to get transaction key while updating journal id %d",
			cmn.LogLevels.Error, j.Id)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, cmn.ErrorCodes.UpdateFailed
	}
	created, errCode := updateJournalEntryTransactions(tx, j, updatedTransactions, createdTransactions, deletedTransactions)
	err = dataaccess.CommitTransaction(tx)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, cmn.ErrorCodes.UpdateFailed
	}
	return created, errCode
}

func UpdateJournalEntryInformationWrapper(j mdls.FinancialJournalEntry) cmn.ErrorCode {
	tx, err := dataaccess.ObtainTransactionKey()
	defer dataaccess.RollbackTransaction(tx)
	if err != nil {
		cmn.Log("failed to get transaction key while updating journal id %d",
			cmn.LogLevels.Error, j.Id)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return cmn.ErrorCodes.UpdateFailed
	}
	code := updateJournalEntryInformation(tx, j)
	err = dataaccess.CommitTransaction(tx)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return cmn.ErrorCodes.UpdateFailed
	}
	return code
}

// updates both the journal entry info and the transactions.
// returns any created transactions, and the error code
func UpdateBothJournalEntryInformatinAndTransactions(j mdls.FinancialJournalEntry,
	updatedTransactions, createdTransactions []mdls.FinancialTransaction,
	deletedTransactions []int) ([]mdls.FinancialTransaction, cmn.ErrorCode) {
	tx, err := dataaccess.ObtainTransactionKey()
	defer dataaccess.RollbackTransaction(tx)
	if err != nil {
		cmn.Log("failed to get transaction key while updating journal id %d",
			cmn.LogLevels.Error, j.Id)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return nil, cmn.ErrorCodes.UpdateFailed
	}

	// update transactions, and rollback if not successful
	newTransactions, errCode := updateJournalEntryTransactions(tx, j, updatedTransactions, createdTransactions, deletedTransactions)
	if errCode != cmn.ErrorCodes.NoError {
		dataaccess.RollbackTransaction((tx))
		return nil, errCode
	}
	// update journal entry info, and rollback if not successful
	errCode = updateJournalEntryInformation(tx, j)
	if errCode != cmn.ErrorCodes.NoError {
		dataaccess.RollbackTransaction((tx))
		return nil, errCode
	}

	err = dataaccess.CommitTransaction(tx)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
	}
	return newTransactions, cmn.ErrorCodes.NoError
}

// update journal entry
// the user will commit an entire journal entry with transactions
// and will submit 3 arrays of transactions, 1. updated, 2. newly created
// 3. original transactions 4. deleted(just a list of ids)
// this function we will balance and compile a list of transactions
// to be created and others to be either deleted or updated
// and send those to approperitate dataaccess function
// We will query the database and make sure all deleted and updated
// transactions belong to the journalId
// noError if updated worked otherwise will return approeritate error
// note that when update is called we always create the imbalance account if it doesn't
// already exist
func updateJournalEntryTransactions(tx string, j mdls.FinancialJournalEntry, updatedTransactions,
	createdTransactions []mdls.FinancialTransaction, deletedTransactionIds []int) ([]mdls.FinancialTransaction, cmn.ErrorCode) {
	if j.Id == 0 {
		cmn.Log("journal id is zero, can't update transactions...", cmn.LogLevels.Error)
		return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
	}
	err := validateTransactionBeforeJournalEntryUpdate(updatedTransactions, j.Id, false)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		if strings.Contains(err.Error(), "transaction has zero for the account id") {
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.InvalidFinancialAccountId
		}
		return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
	}
	err = validateTransactionBeforeJournalEntryUpdate(createdTransactions, j.Id, true)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		if strings.Contains(err.Error(), "transaction has zero for the account id") {
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.InvalidFinancialAccountId
		}
		return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
	}
	editedTransactionIds := []int{}
	// ensures that every transaction is seen only once
	// we can't update and delete the same transaction,
	// and we store a map of ids for later use
	deletedIds := make(map[int]bool)
	updatedMap := make(map[int]mdls.FinancialTransaction)
	for _, tId := range deletedTransactionIds {
		deletedIds[tId] = true
		editedTransactionIds = append(editedTransactionIds, tId)
	}
	for _, t := range updatedTransactions {
		updatedMap[t.Id] = t
		_, ok := deletedIds[t.Id]
		if ok {
			cmn.Log("You can't have repeated transaction ids... saw id %d twice",
				cmn.ErrorLevels.Error, t.Id)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdatedAndDeletedSameTransaction
		}
		editedTransactionIds = append(editedTransactionIds, t.Id)
	}

	// make sure that all the transactions we are trying to edit actually belong the journal
	// entry for real, and that the client isn't trying to reach and update some
	// other journal entry.
	if !doTransactionsBelongToJournalEntry(tx, editedTransactionIds, j.Id) {
		cmn.Log("While updating journal id %d, the user tried editing a transaction that doesnt belong to the journal entry... ",
			cmn.LogLevels.Info, j.Id)
		return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
	}
	// figure out what transactions were unchanged, and which need an update
	originalTransactions, err := dataaccess.ReadFinancialTransactionsByJournalEntryId(tx, j.Id)
	if err != nil {
		cmn.Log("failed to read transactions while updating journal id %d",
			cmn.LogLevels.Error, j.Id)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
	}
	ia, err := dataaccess.FindOrCreateImbalanceAccount(tx)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		cmn.Log("Couldnt find or create the imbalance account", cmn.LogLevels.Info)
		return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
	}
	thereIsImbalanceTransactionAlready, oldImbalanceT := findImbalanceTransaction(ia.Id, originalTransactions)
	// if there's one already then:
	// 1. if the transactions without it are balanced, then delete the imbalance t
	// 2. if the transactions without it are imbalanced, then update it
	// If there is not one already:
	// 1. if transactions are imbalanced then create a new imbalance transaction

	unChangedTransactions := []mdls.FinancialTransaction{}
	outdatedTransactions := []mdls.FinancialTransaction{}
	deletedTransactions := []mdls.FinancialTransaction{}
	for _, t := range originalTransactions {
		if thereIsImbalanceTransactionAlready && t.Id == oldImbalanceT.Id { // dont add the imbalance transaction yet
			continue
		}
		if _, ok := updatedMap[t.Id]; ok {
			outdatedTransactions = append(outdatedTransactions, t)
		} else if _, ok := deletedIds[t.Id]; ok {
			deletedTransactions = append(deletedTransactions, t)
		} else {
			unChangedTransactions = append(unChangedTransactions, t)
		}
	}

	clientInputTs := append(updatedTransactions, createdTransactions...)
	clientInputTs = append(clientInputTs, unChangedTransactions...)
	isThereImbalance, imbalanceT, errorCode := makeTransactionsBalance(clientInputTs)
	if errorCode != cmn.ErrorCodes.NoError {
		return []mdls.FinancialTransaction{}, errorCode
	}
	imbalanceT.FinancialJournalEntryId = j.Id
	imbalanceT.FinancialAccountId = ia.Id
	_, clientUpdatedImbalanceTransaction := updatedMap[oldImbalanceT.Id]
	if isThereImbalance {
		cmn.Log("Jouranl entry transactions need balancing", cmn.LogLevels.Info)
		if thereIsImbalanceTransactionAlready { // if it's already here, then just update it
			cmn.Log("Updating existing imbalance transaction", cmn.LogLevels.Info)
			imbalanceT.Id = oldImbalanceT.Id
			outdatedTransactions = append(outdatedTransactions, oldImbalanceT)
			updatedTransactions = append(updatedTransactions, imbalanceT)
		} else { // create it if it doesn't exist
			cmn.Log("Creating an imbalance transaction", cmn.LogLevels.Info)
			createdTransactions = append(createdTransactions, imbalanceT)
		}
	} else { // transactions are already balanced
		cmn.Log("Jouranl entry transactions are already balanced", cmn.LogLevels.Info)
		// if imbalance is not needed, and there is an imbalance transaction already, then delete the
		// imbalance transaction. Unless, we are editing the imbalance transaction (for example the
		// journalEntry is already balanced because of the imbalanced transaction and the user decided
		// to not use the imbalance account but instead use a different account)
		if clientUpdatedImbalanceTransaction {
			cmn.Log("Client updated the existing imbalance transaction", cmn.LogLevels.Info)
			// the imbalance transaction wasn't added to the outdatedTransactions in the loop above
			// so we have to add it now to ensure that updatedTransactions corospond to outdatedTransactions
			outdatedTransactions = append(outdatedTransactions, oldImbalanceT)
		} else if thereIsImbalanceTransactionAlready {
			cmn.Log("Deleting the imbalance transaction", cmn.LogLevels.Info)
			deletedTransactions = append(deletedTransactions, oldImbalanceT)
		}
	}

	cmn.Log("updateJournalEntryTransactions.\nIncluding any imbalances, there are %d created, %d unchanged, %d updated, and %d deleted transacations",
		cmn.LogLevels.Info, len(createdTransactions),
		len(unChangedTransactions), len(updatedTransactions), len(deletedTransactions))

	// we might have many transactions that work on the same accounts
	// so we store those accounts instead of having to read them evertime
	// to spare writes and reads form the db. At the end we will update
	// all accounts once all transactions have processed
	queriedAccounts := make(map[int]mdls.FinancialAccount)
	modifiedAccountIds := []int{}
	// undo the effects of outdated transactions and update with new ones
	// both updated and outdated transactions have same len
	for i, t := range updatedTransactions {
		oldT := outdatedTransactions[i]
		a, err := getAndAccumulateAccount(tx, t.FinancialAccountId, queriedAccounts)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
		oldA, err := getAndAccumulateAccount(tx, oldT.FinancialAccountId, queriedAccounts) // in case the update switched accounts
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
		err = dataaccess.UpdateFinancialTransaction(tx, t)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
		// undo old transaction effects
		if oldA.Id == a.Id {
			a.UpdateAccountBalanceAfterDeletingTransaction(oldT)
		} else if a.Id != oldA.Id {
			oldA.UpdateAccountBalanceAfterDeletingTransaction(oldT)
			queriedAccounts[oldA.Id] = oldA // make sure we keep track account upates
			modifiedAccountIds = append(modifiedAccountIds, oldT.FinancialAccountId)
		}
		// update with new transaction
		err = a.UpdateAccountBalanceAfterTransaction(t)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
		queriedAccounts[a.Id] = a // make sure we keep track account upates
		modifiedAccountIds = append(modifiedAccountIds, t.FinancialAccountId)
	}
	for _, t := range deletedTransactions {
		a, err := getAndAccumulateAccount(tx, t.FinancialAccountId, queriedAccounts)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
		err = dataaccess.DeleteFinancialTransaction(tx, t)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
		a.UpdateAccountBalanceAfterDeletingTransaction(t)
		queriedAccounts[a.Id] = a // make sure we keep track of account upates
		modifiedAccountIds = append(modifiedAccountIds, a.Id)
	}

	newTransactions := []mdls.FinancialTransaction{}

	// insert the new transactions
	for _, t := range createdTransactions {
		a, err := getAndAccumulateAccount(tx, t.FinancialAccountId, queriedAccounts)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
		tId, err := dataaccess.CreateNewFinancialTransaction(tx, t)
		if err != nil {
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
		t.Id = tId
		err = a.UpdateAccountBalanceAfterTransaction(t)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
		queriedAccounts[a.Id] = a // make sure we keep track of account upates
		modifiedAccountIds = append(modifiedAccountIds, a.Id)
		newTransactions = append(newTransactions, t)
	}

	// updated the balances of all modifiied accounts
	for _, aId := range modifiedAccountIds {
		a, ok := queriedAccounts[aId]
		if !ok {
			cmn.Log("failed to update the balance of a modified account",
				cmn.ErrorLevels.Error)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
		err = dataaccess.UpdateFinancialAccountBalance(tx, a)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return []mdls.FinancialTransaction{}, cmn.ErrorCodes.UpdateFailed
		}
	}
	return newTransactions, cmn.ErrorCodes.NoError
}

func updateJournalEntryInformation(txKey string, j mdls.FinancialJournalEntry) cmn.ErrorCode {
	if j.Id == 0 {
		cmn.Log("journal entry has invalid id", cmn.LogLevels.Error)
		return cmn.ErrorCodes.UpdateFailed
	}
	err := dataaccess.UpdateJournalEntryInformation(txKey, j)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return cmn.ErrorCodes.UpdateFailed
	}
	return cmn.ErrorCodes.NoError
}

func findAllJournalEntries() ([]mdls.FinancialJournalEntry, cmn.ErrorCode) {
	js, err := dataaccess.FindAllJournalEntries()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return []mdls.FinancialJournalEntry{}, cmn.ErrorCodes.ReadFailed
	}
	return js, cmn.ErrorCodes.NoError
}
func findJournalEntriesByIds(ids []int) ([]mdls.FinancialJournalEntry, cmn.ErrorCode) {
	js, err := dataaccess.FindJournalEntriesByIds(ids)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return []mdls.FinancialJournalEntry{}, cmn.ErrorCodes.ReadFailed
	}
	return js, cmn.ErrorCodes.NoError
}

func findAllJournalEntriesComposite() ([]mdls.FinancialJournalEntryComposite, cmn.ErrorCode) {
	js, err := dataaccess.FindAllJournalEntriesComposite()
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return []mdls.FinancialJournalEntryComposite{}, cmn.ErrorCodes.ReadFailed
	}
	return js, cmn.ErrorCodes.NoError
}
func findJournalEntryCompositeById(jId int) (mdls.FinancialJournalEntryComposite, cmn.ErrorCode) {
	js, err := dataaccess.FindJournalEntryCompositeById(jId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return mdls.FinancialJournalEntryComposite{}, cmn.ErrorCodes.ReadFailed
	}
	return js, cmn.ErrorCodes.NoError
}
func findJournalEntryCompositeByAccountId(jId int) ([]mdls.FinancialJournalEntryComposite, cmn.ErrorCode) {
	js, err := dataaccess.FindJournalEntryCompositeByAccountId(jId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return []mdls.FinancialJournalEntryComposite{}, cmn.ErrorCodes.ReadFailed
	}
	return js, cmn.ErrorCodes.NoError
}
