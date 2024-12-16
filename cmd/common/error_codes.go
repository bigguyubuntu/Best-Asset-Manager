package cmn

import "fmt"

// the error codes, the UI should decide on the message to display

type ErrorCode string
type errorCodes struct {
	NoError                          ErrorCode
	BadRequest                       ErrorCode
	NameAlreadyExists                ErrorCode // account names must be unique
	NameWasEmpty                     ErrorCode // name can't be empty
	NeedAtLeastOneGroupForEveryType  ErrorCode // we need at least one group for every financial type
	DeleteFailed                     ErrorCode
	UpdateFailed                     ErrorCode
	ReadFailed                       ErrorCode
	UpdatedAndDeletedSameTransaction ErrorCode
	CreateFailed                     ErrorCode
	ParentCantBeSameAsId             ErrorCode
	InvalidParentGroupId             ErrorCode
	GroupTypeCantBeChanged           ErrorCode
	InvalidFinancialTransactionType  ErrorCode
	InvalidFinancialAccountId        ErrorCode
	InvalidFinancialAccountType      ErrorCode
	InvalidFinancialGroupType        ErrorCode
	NotEnoughTransactions            ErrorCode
	JournalEntryFailed               ErrorCode
	JournalEntryDateCantBeEmpty      ErrorCode
	InvalidId                        ErrorCode
	InvalidStorehouseType            ErrorCode
}

// these error codes will be shared with the UI
var ErrorCodes = errorCodes{
	NoError:                          "no_error",
	BadRequest:                       "bad_request",
	NameAlreadyExists:                "name_clash",
	NameWasEmpty:                     "empty_name",
	NeedAtLeastOneGroupForEveryType:  "need_atleast_one_group_for_every_type",
	DeleteFailed:                     "delete_failed",
	UpdateFailed:                     "update_failed",
	CreateFailed:                     "create_failed",
	ReadFailed:                       "read_failed",
	ParentCantBeSameAsId:             "parent_id_cant_be_same_as_group_id",
	InvalidParentGroupId:             "invalid_parent_group_id",
	GroupTypeCantBeChanged:           "cant_update_group_type_please_create_a_new_group_instead",
	NotEnoughTransactions:            "transactions_cant_be_less_than_two",
	UpdatedAndDeletedSameTransaction: "cant_update_and_delete_same_transaction",
	InvalidFinancialTransactionType:  "transaction_type_must_be_either_debit_or_credit",
	InvalidFinancialAccountId:        "invalid_financial_account_id",
	InvalidFinancialAccountType:      "invalid_financial_account_type",
	InvalidFinancialGroupType:        "invalid_financial_group_type",
	JournalEntryFailed:               "journal_entry_failed",
	JournalEntryDateCantBeEmpty:      "journal_entry_date_cant_be_empty",
	InvalidId:                        "id_is_invalid",
	InvalidStorehouseType:            "invalid_storehouse_type",
}

func AddMsgToErrorCode(e ErrorCode, s string) ErrorCode {
	// takes a an error code and appends a message to it
	// this is useful when we want to include more information about
	// the error
	return ErrorCode(fmt.Sprintf("%s\t%s", s, e))
}
