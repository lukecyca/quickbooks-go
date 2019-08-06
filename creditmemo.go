package quickbooks

import (
	"encoding/json"
	"strconv"
)

// CreditMemo represents a QuickBooks Credit Memo object.
type CreditMemo struct {
	ID        string   `json:"Id,omitempty"`
	SyncToken string   `json:",omitempty"`
	MetaData  MetaData `json:",omitempty"`
	//CustomField
	DocNumber string `json:",omitempty"`
	TxnDate   Date   `json:",omitempty"`
	//DepartmentRef
	PrivateNote  string `json:",omitempty"`
	Line         []SalesItemLine
	TxnTaxDetail TxnTaxDetail `json:",omitempty"`
	CustomerRef  ReferenceType
	CustomerMemo MemoRef         `json:",omitempty"`
	BillAddr     PhysicalAddress `json:",omitempty"`
	ShipAddr     PhysicalAddress `json:",omitempty"`
	ClassRef     ReferenceType   `json:",omitempty"`
	SalesTermRef ReferenceType   `json:",omitempty"`
	//GlobalTaxCalculation
	TotalAmt json.Number `json:",omitempty"`
	//CurrencyRef
	ExchangeRate          json.Number  `json:",omitempty"`
	HomeBalance           json.Number  `json:",omitempty"`
	ApplyTaxAfterDiscount bool         `json:",omitempty"`
	PrintStatus           string       `json:",omitempty"`
	EmailStatus           string       `json:",omitempty"`
	BillEmail             EmailAddress `json:",omitempty"`
	Balance               json.Number  `json:",omitempty"`
}

// FetchCreditMemos gets the full list of CreditMemos in the QuickBooks account.
func (c *Client) FetchCreditMemos() ([]CreditMemo, error) {

	// See how many invoices there are.
	var r struct {
		QueryResponse struct {
			TotalCount int
		}
	}
	err := c.query("SELECT COUNT(*) FROM CreditMemo", &r)
	if err != nil {
		return nil, err
	}

	if r.QueryResponse.TotalCount == 0 {
		return make([]CreditMemo, 0), nil
	}

	var creditmemos = make([]CreditMemo, 0, r.QueryResponse.TotalCount)
	for i := 0; i < r.QueryResponse.TotalCount; i += queryPageSize {
		var page, err = c.fetchCreditMemoPage(i + 1)
		if err != nil {
			return nil, err
		}
		creditmemos = append(creditmemos, page...)
	}
	return creditmemos, nil
}

// Fetch one page of results, because we can't get them all in one query.
func (c *Client) fetchCreditMemoPage(startpos int) ([]CreditMemo, error) {

	var r struct {
		QueryResponse struct {
			CreditMemo    []CreditMemo
			StartPosition int
			MaxResults    int
		}
	}
	q := "SELECT * FROM CreditMemo ORDERBY Id STARTPOSITION " +
		strconv.Itoa(startpos) + " MAXRESULTS " + strconv.Itoa(queryPageSize)
	err := c.query(q, &r)
	if err != nil {
		return nil, err
	}

	// Make sure we don't return nil if there are no invoices.
	if r.QueryResponse.CreditMemo == nil {
		r.QueryResponse.CreditMemo = make([]CreditMemo, 0)
	}
	return r.QueryResponse.CreditMemo, nil
}
