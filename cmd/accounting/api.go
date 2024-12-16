package accounting

import (
	"encoding/json"
	"fmt"
	"inventory_money_tracking_software/cmd/api"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func findAccountsByIdsHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	idsString := strings.Split(strings.Trim(r.URL.Query().Get("ids"), ""), ",")
	ids := []int{}
	for i := 0; i < len(idsString); i++ {
		n, err := strconv.Atoi(idsString[i])
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
			return
		}
		ids = append(ids, n)
	}
	as, errCode := findAccountsByIds(ids)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	if len(as) == 0 {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed,
			"Didn't find any accounts with the provided ids.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(as)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findAllAccountsHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	as, errCode := findAllAccounts()
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	if len(as) == 0 {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed,
			"Didn't find any accounts.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(as)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}
func findAllJournalEntriesHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	js, errCode := findAllJournalEntries()
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	if len(js) == 0 {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed,
			"Didn't find any journal entries.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(js)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findAllJournalEntriesByIdsHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	idsString := strings.Split(strings.Trim(r.URL.Query().Get("ids"), ""), ",")
	ids := []int{}
	for i := 0; i < len(idsString); i++ {
		n, err := strconv.Atoi(idsString[i])
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
			return
		}
		ids = append(ids, n)
	}
	js, errCode := findJournalEntriesByIds(ids)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	if len(js) == 0 {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed,
			"Didn't find any journal entries.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(js)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findAllJournalEntriesCompositeHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	js, errCode := findAllJournalEntriesComposite()
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(js)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findJournalEntryCompositeByIdHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	idString := strings.Trim(r.URL.Query().Get("journal_entry_id"), "")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
		return
	}
	js, errCode := findJournalEntryCompositeById(idInt)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(js)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findJournalEntryCompositeByAccountIdHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	idString := strings.Trim(r.URL.Query().Get("account_id"), "")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
		return
	}
	js, errCode := findJournalEntryCompositeByAccountId(idInt)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(js)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findAllGroupsHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	gs, errCode := ReadAllFinancialGroups()
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}
	if len(gs) == 0 {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed,
			"Didn't find any groups")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(gs)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findGroupsByTypeHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	groupType := strings.Trim(r.URL.Query().Get("type"), "")
	if !cmn.IsGroupTypeValid(groupType) {
		api.HandleRequestError(w, cmn.ErrorCodes.InvalidFinancialGroupType,
			fmt.Sprintf("Group type can only be one of: %s, %s, %s, %s or %s",
				cmn.GroupTypes.Assets, cmn.GroupTypes.Expenses, cmn.GroupTypes.Liabilities,
				cmn.GroupTypes.Income, cmn.GroupTypes.Equity))
		return
	}
	gs, errCode := ReadFinancialGroupsWithType(groupType)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}
	if len(gs) == 0 {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed,
			"Didn't find any groups with the provided group type")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(gs)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findAccountsByGroupIdHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	idString := strings.Trim(r.URL.Query().Get("group_id"), "")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
		return
	}
	as, errCode := ReadFinancialAccountsWithGroup(idInt)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(as)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findGroupByNameAndParent(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	parentId := strings.Trim(r.URL.Query().Get("parent_group_id"), "")
	idInt, err := strconv.Atoi(parentId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
		return
	}
	name := strings.Trim(r.URL.Query().Get("group_name"), "")
	g, errCode := ReadFinancialGroupByNameAndParentId(name, idInt)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(g)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findGroupsByIdsHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	idsString := strings.Split(strings.Trim(r.URL.Query().Get("ids"), ""), ",")
	ids := []int{}
	for i := 0; i < len(idsString); i++ {
		n, err := strconv.Atoi(idsString[i])
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
			return
		}
		ids = append(ids, n)
	}
	gs, errCode := findGroupsByIds(ids)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	if len(gs) == 0 {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed,
			"Didn't find any group with the provided ids.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(gs)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findAccountByNameAndGroup(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	parentId := strings.Trim(r.URL.Query().Get("group_id"), "")
	idInt, err := strconv.Atoi(parentId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
		return
	}
	name := strings.Trim(r.URL.Query().Get("account_name"), "")
	a, errCode := ReadFinancialAccountByNameAndGroupId(name, idInt)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(a)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findGroupsByParentGroup(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	idString := strings.Trim(r.URL.Query().Get("parent_group_id"), "")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
		return
	}
	gs, errCode := ReadFinancialGroupsByParentGroup(idInt)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(gs)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, api.StdErrMsg)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var allowedMethods = []string{http.MethodPost, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "text/plain")
		var acnt mdls.FinancialAccount
		err := json.NewDecoder(r.Body).Decode(&acnt)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.CreateFailed, "")
			return
		}
		aId, errorCode := CreateAccount(acnt)
		if errorCode != cmn.ErrorCodes.NoError {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, errorCode, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, aId)
	}
	return
}
func createJournalEntryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var allowedMethods = []string{http.MethodPost, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var jsonShape struct {
			JournalEntry mdls.FinancialJournalEntry  `json:"journal_entry"`
			Transactions []mdls.FinancialTransaction `json:"transactions"`
		}
		err := json.NewDecoder(r.Body).Decode(&jsonShape)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.CreateFailed, "")
			return
		}
		journalEntryId, createdTransactions, errorCode := CreateJournalEntry(jsonShape.JournalEntry, jsonShape.Transactions)
		if errorCode != cmn.ErrorCodes.NoError {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, errorCode, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		respShape := struct {
			JournalEntryId int                         `json:"created_journal_entry_id"`
			TransactionIds []mdls.FinancialTransaction `json:"created_transactions"`
		}{
			JournalEntryId: journalEntryId,
			TransactionIds: createdTransactions,
		}
		// io.WriteString(w, journalEntryId)
		resp, err := json.Marshal(respShape)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.CreateFailed, "")
			return
		}
		w.Write(resp)
		return
	}
	return
}

// updates just the description, and entrydate. To update transactions use the other update function
func updateJournalEntryInformationHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var allowedMethods = []string{http.MethodPost, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	if r.Method == http.MethodPost {
		j := mdls.FinancialJournalEntry{}
		err := json.NewDecoder(r.Body).Decode(&j)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.UpdateFailed, "")
			return
		}
		errorCode := UpdateJournalEntryInformationWrapper(j)
		if errorCode != cmn.ErrorCodes.NoError {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, errorCode, "")
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
	return
}

// updates the transactions of a journal entry
func updateJournalEntryTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var allowedMethods = []string{http.MethodPost, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// request shape
		var jsonShape struct {
			JournalEntry        mdls.FinancialJournalEntry  `json:"journal_entry"`
			CreatedTransactions []mdls.FinancialTransaction `json:"created_transactions"`
			UpdatedTransactions []mdls.FinancialTransaction `json:"updated_transactions"`
			DeletedTransactions []int                       `json:"deleted_transaction_ids"`
		}
		err := json.NewDecoder(r.Body).Decode(&jsonShape)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.UpdateFailed, "")
			return
		}
		created, errorCode := UpdateJournalEntryTransactionsWrapper(jsonShape.JournalEntry,
			jsonShape.UpdatedTransactions, jsonShape.CreatedTransactions, jsonShape.DeletedTransactions)
		if errorCode != cmn.ErrorCodes.NoError {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, errorCode, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		respShape := struct {
			Transactions []mdls.FinancialTransaction `json:"created_transactions"`
		}{
			Transactions: created,
		}
		// io.WriteString(w, journalEntryId)
		resp, err := json.Marshal(respShape)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.UpdateFailed, "")
			return
		}
		w.Write(resp)
		return
	}
	return
}

// updates both the transactions of a journal entry and the entries information
func updateBothJournalEntryInformatinAndTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var allowedMethods = []string{http.MethodPost, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// request shape
		var jsonShape struct {
			JournalEntry        mdls.FinancialJournalEntry  `json:"journal_entry"`
			CreatedTransactions []mdls.FinancialTransaction `json:"created_transactions"`
			UpdatedTransactions []mdls.FinancialTransaction `json:"updated_transactions"`
			DeletedTransactions []int                       `json:"deleted_transaction_ids"`
		}
		err := json.NewDecoder(r.Body).Decode(&jsonShape)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.UpdateFailed, "")
			return
		}
		createdTs, errorCode := UpdateBothJournalEntryInformatinAndTransactions(jsonShape.JournalEntry,
			jsonShape.UpdatedTransactions, jsonShape.CreatedTransactions, jsonShape.DeletedTransactions)
		// both updates have failed case
		if errorCode != cmn.ErrorCodes.NoError {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, errorCode, "both updates failed, the transactions and the entry information updates")
			return
		}
		// success case
		w.WriteHeader(http.StatusOK)
		respBody := struct {
			TransactionIds []mdls.FinancialTransaction `json:"created_transactions"`
		}{
			TransactionIds: createdTs,
		}
		resp, err := json.Marshal(respBody)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.UpdateFailed, "")
			return
		}
		w.Write(resp)
		return
	}
	return
}

func createGroupHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var allowedMethods = []string{http.MethodPost, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "text/plain")
		var grp mdls.FinancialGroup
		err := json.NewDecoder(r.Body).Decode(&grp)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.CreateFailed, "")
			return
		}
		grpId, errorCode := CreateNewFinancialGroup(grp)
		if errorCode != cmn.ErrorCodes.NoError {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, errorCode, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, grpId)
	}
	return
}

func updateGroupHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var allowedMethods = []string{http.MethodPost, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "text/plain")
		var grp mdls.FinancialGroup
		err := json.NewDecoder(r.Body).Decode(&grp)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.CreateFailed, "")
			return
		}
		if grp.Id == 0 {
			api.HandleRequestError(w, cmn.ErrorCodes.CreateFailed, "Please pass a valid group Id")
			return
		}
		errorCode := UpdateFinancialGroup(grp.Id, grp)
		if errorCode != cmn.ErrorCodes.NoError {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, errorCode, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, strconv.Itoa(grp.Id))
	}
	return
}

func deleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var allowedMethods = []string{http.MethodDelete, http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	if r.Method != http.MethodOptions {
		idString := strings.Trim(r.URL.Query().Get("group_id"), "")
		idInt, err := strconv.Atoi(idString)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		errorCode := DeleteFinancialAccoutsGroup(idInt)
		if errorCode != cmn.ErrorCodes.NoError {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, errorCode, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, idString)
	}
	return
}

func findTransactionsByAccountIdHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	idString := strings.Trim(r.URL.Query().Get("id"), "")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
		return
	}
	ts, errCode := readTransactionsByAccount(idInt)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(ts)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func findTransactionsByJournalEntryIdHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	idString := strings.Trim(r.URL.Query().Get("id"), "")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
		return
	}
	ts, errCode := readTransactionsByJournalEntryId(idInt)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(ts)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}
func findTransactionsByJournalEntriesHandler(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	if !api.ValidateHttpRequestSecurityRules(w, r,
		allowedMethods, allowedHeaders) {
		api.HandleBadRequest(w)
		return
	}
	idsString := strings.Split(strings.Trim(r.URL.Query().Get("ids"), ""), ",")
	ids := []int{}
	for i := 0; i < len(idsString); i++ {
		n, err := strconv.Atoi(idsString[i])
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "malformed request")
			return
		}
		ids = append(ids, n)
	}
	ts, errCode := readTransactionsByJournalEntries(ids)
	if errCode != cmn.ErrorCodes.NoError {
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonResp, err := json.Marshal(ts)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		api.HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "Something went wrong, try again please.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	return
}

func RegisterHandlers() {
	mux := api.GetHttpMux()
	mux.HandleFunc(accountRoutes.findByGroupId, findAccountsByGroupIdHandler)
	mux.HandleFunc(accountRoutes.findByIds, findAccountsByIdsHandler)
	mux.HandleFunc(accountRoutes.create, createAccountHandler)
	mux.HandleFunc(accountRoutes.findAll, findAllAccountsHandler)
	mux.HandleFunc(accountRoutes.findByNameAndGroupId, findAccountByNameAndGroup)

	mux.HandleFunc(groupRoutes.create, createGroupHandler)
	mux.HandleFunc(groupRoutes.findByType, findGroupsByTypeHandler)
	mux.HandleFunc(groupRoutes.findAll, findAllGroupsHandler)
	mux.HandleFunc(groupRoutes.delete, deleteGroupHandler)
	mux.HandleFunc(groupRoutes.update, updateGroupHandler)
	mux.HandleFunc(groupRoutes.findByParentId, findGroupsByParentGroup)
	mux.HandleFunc(groupRoutes.findByNameAndParentId, findGroupByNameAndParent)
	mux.HandleFunc(groupRoutes.findByIds, findGroupsByIdsHandler)

	mux.HandleFunc(journalEntryRoutes.create, createJournalEntryHandler)
	mux.HandleFunc(journalEntryRoutes.findAll, findAllJournalEntriesHandler)
	mux.HandleFunc(journalEntryRoutes.findByIds, findAllJournalEntriesByIdsHandler)
	mux.HandleFunc(journalEntryRoutes.compositeFindAll, findAllJournalEntriesCompositeHandler)
	mux.HandleFunc(journalEntryRoutes.compositeFindById, findJournalEntryCompositeByIdHandler)
	mux.HandleFunc(journalEntryRoutes.compositeFindByAccountId, findJournalEntryCompositeByAccountIdHandler)
	mux.HandleFunc(journalEntryRoutes.update, updateJournalEntryInformationHandler)
	mux.HandleFunc(journalEntryRoutes.updateTransactions, updateJournalEntryTransactionsHandler)
	mux.HandleFunc(journalEntryRoutes.updateBothEntryAndTransactions, updateBothJournalEntryInformatinAndTransactionsHandler)

	mux.HandleFunc(transactionRoutes.findByAccountId, findTransactionsByAccountIdHandler)
	mux.HandleFunc(transactionRoutes.findByJournalEntryId, findTransactionsByJournalEntryIdHandler)
	mux.HandleFunc(transactionRoutes.findByJournalEntries, findTransactionsByJournalEntriesHandler)
}
