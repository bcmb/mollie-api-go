package mollie

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Refund describe a refund for a certain payment.
type Refund struct {
	Resource         string       `json:"resource,omitempty"`
	ID               string       `json:"id,omitempty"`
	Amount           *Amount      `json:"amount,omitempty"`
	SettlementID     string       `json:"settlementId,omitempty"`
	SettlementAmount *Amount      `json:"settlementAmount,omitempty"`
	Description      string       `json:"description,omitempty"`
	Metadata         interface{}  `json:"metadata,omitempty"`
	Status           RefundStatus `json:"status,omitempty"`
	Lines            []*OrderLine `json:"lines,omitempty"`
	PaymentID        string       `json:"paymentId,omitempty"`
	OrderID          string       `json:"orderId,omitempty"`
	CreatedAt        *time.Time   `json:"createdAt,omitempty"`
	TestMode         bool         `json:"testmode,omitempty"`
	Links            RefundLinks  `json:"_links,omitempty"`
}

// RefundList describes how a list of refunds will be retrieved by Mollie.
type RefundList struct {
	Count    int `json:"count,omitempty"`
	Embedded struct {
		Refunds []*Refund
	} `json:"_embedded,omitempty"`
	Links PaginationLinks `json:"_links,omitempty"`
}

// RefundStatus describes the status of the refund.
type RefundStatus string

// Valid refund status.
const (
	Queued     RefundStatus = "queued"
	Pending    RefundStatus = "pending"
	Processing RefundStatus = "processing"
	Refunded   RefundStatus = "refunded"
	Failed     RefundStatus = "failed"
)

// RefundLinks describes all the possible links to be returned with
// a Refund object.
type RefundLinks struct {
	Self          *URL `json:"self,omitempty"`
	Payment       *URL `json:"payment,omitempty"`
	Settlement    *URL `json:"settlement,omitempty"`
	Order         *URL `json:"order,omitempty"`
	Documentation *URL `json:"documentation,omitempty"`
}

// RefundOptions describes refund endpoint valid query string parameters.
//
// See: https://docs.mollie.com/reference/v2/refunds-api/get-refund.
type RefundOptions struct {
	Embed EmbedValue `url:"embed,omitempty"`
}

// ListRefundOptions describes list refund endpoint valid query string parameters.
//
// See: https://docs.mollie.com/reference/v2/refunds-api/list-refunds.
type ListRefundOptions struct {
	From      string     `url:"from,omitempty"`
	Limit     int        `url:"limit,omitempty"`
	ProfileID string     `url:"profileId,omitempty"`
	Embed     EmbedValue `url:"embed,omitempty"`
}

// RefundsService instance operates over refund resources.
type RefundsService service

// Get retrieve a single refund by its ID.
//
// If you do not know the original payment’s ID, you can use the List payment refunds endpoint.
func (rs *RefundsService) Get(ctx context.Context, paymentID, refundID string, opts *RefundOptions) (res *Response, refund *Refund, err error) {
	u := fmt.Sprintf("v2/payments/%s/refunds/%s", paymentID, refundID)

	res, err = rs.client.get(ctx, u, opts)
	if err != nil {
		return
	}

	if err = json.Unmarshal(res.content, &refund); err != nil {
		return
	}

	return
}

// Create a refund payment request.
//
// See https://docs.mollie.com/reference/v2/refunds-api/create-refund.
func (rs *RefundsService) Create(ctx context.Context, paymentID string, re Refund, options *RefundOptions) (res *Response, rf *Refund, err error) {
	uri := fmt.Sprintf("v2/payments/%s/refunds", paymentID)

	if rs.client.HasAccessToken() && rs.client.config.testing {
		re.TestMode = true
	}

	res, err = rs.client.post(ctx, uri, re, options)
	if err != nil {
		return
	}

	if err = json.Unmarshal(res.content, &rf); err != nil {
		return
	}

	return
}

// Cancel try to cancel the refund request.
// The refund can only be canceled while the refund’s status is either queued or pending.
// See https://docs.mollie.com/reference/v2/refunds-api/cancel-refund
func (rs *RefundsService) Cancel(ctx context.Context, paymentID, refundID string) (res *Response, err error) {
	u := fmt.Sprintf("v2/payments/%s/refunds/%s", paymentID, refundID)

	res, err = rs.client.delete(ctx, u, nil)
	if err != nil {
		return
	}

	return
}

// ListRefund calls the top level https://api.mollie.com/v2/refunds.
//
// See https://docs.mollie.com/reference/v2/refunds-api/list-refunds.
func (rs *RefundsService) ListRefund(ctx context.Context, opts *ListRefundOptions) (res *Response, rl *RefundList, err error) {
	u := "v2/refunds"

	return rs.list(ctx, u, opts)
}

// ListRefundPayment calls the payment-specific
// https://api.mollie.com/v2/payments/*paymentId*/refunds.
//
// Only refunds for that specific payment are returned.
// See: https://docs.mollie.com/reference/v2/refunds-api/list-refunds
func (rs *RefundsService) ListRefundPayment(ctx context.Context, paymentID string, opts *ListRefundOptions) (res *Response, rl *RefundList, err error) {
	u := fmt.Sprintf("v2/payments/%s/refunds", paymentID)

	return rs.list(ctx, u, opts)
}

func (rs *RefundsService) list(ctx context.Context, uri string, opts interface{}) (res *Response, rl *RefundList, err error) {
	res, err = rs.client.get(ctx, uri, opts)
	if err != nil {
		return
	}

	if err = json.Unmarshal(res.content, &rl); err != nil {
		return
	}
	return
}
