package points

import (
	"context"
	"errors"
	"sort"
	"time"

	"encore.dev/rlog"
)

// PayerTransaction represents a single transaction that was granted to a customer from a specific Payer.
type PayerTransaction struct {
	TransactionID int       `json:"id"`
	CustomerID    int       `json:"customer_id"`
	Payer         string    `json:"payer"`
	Points        int64     `json:"points"`
	Entered       time.Time `json:"timestamp"`
}

// CustomerTransactions represents the list of transactions for a given customer.
type CustomerTransactions struct {
	CustomerID   int                 `json:"id"`
	Name         string              `json:"name"`
	Transactions []*PayerTransaction `json:"transactions"`
}

// Customer represents a customer in the database.
type Customer struct {
	Name         string `json:"name"`
	ID           int    `json:"id"`
	TotalBalance int    `json:"total_balance"`
}

// NewCustomer creates a new customer in the database. The return value is the unique customer id.
//encore:api public method=POST path=/customer/:name
func NewCustomer(ctx context.Context, name string) (*Customer, error) {
	id, err := createCustomer(ctx, name)
	if err != nil {
		return nil, err
	}
	return &Customer{
		Name: name,
		ID:   id,
	}, nil
}

// AddTransaction will add a transaction to a customer's balance sheet.
//encore:api public method=POST path=/points/:customer
func AddTransaction(ctx context.Context, customer int, transaction *PayerTransaction) error {
	err := insertTransaction(ctx, customer, transaction.Payer, transaction.Points, transaction.Entered)
	return err
}

// GetCustomerBalances returns the list of balances across various transactions from payers.
//encore:api public method=GET path=/points/:id
func GetCustomerBalances(ctx context.Context, id int) (*CustomerTransactions, error) {
	transactions, err := selectCustomerBalances(ctx, id)
	if err != nil {
		return nil, err
	}
	return &CustomerTransactions{
		CustomerID:   id,
		Transactions: transactions,
	}, nil
}

// SpendCustomerPoints will walk through a customer's list of balances from various payers and spend those points.
// What is returned is the list of adjustments to existing transactions.  If a transaction is marked as having 0,
// it can be assumed to have been deleted from the database.
//encore:api public method=POST path=/spend/:id/:points
func SpendCustomerPoints(ctx context.Context, id int, points int64) (*CustomerTransactions, error) {
	transactions, err := selectCustomerBalances(ctx, id)
	if err != nil {
		return nil, err
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Entered.Before(transactions[j].Entered)
	})
	customerTotal := int64(0)
	for _, t := range transactions {
		customerTotal += t.Points
	}
	if customerTotal < points {
		return nil, errors.New("customer doesn't have enough points to spend")
	}

	adjustments := make([]*PayerTransaction, 0)
	total := int64(0)
	for _, adj := range transactions {
		if total < points {
			originalPoints := adj.Points
			leftToRemove := points - total
			if adj.Points <= leftToRemove {
				adjustments = append(adjustments, &PayerTransaction{
					CustomerID:    adj.CustomerID,
					Payer:         adj.Payer,
					Points:        0,
					Entered:       adj.Entered,
					TransactionID: adj.TransactionID,
				})
				total += originalPoints
			} else {
				adjustments = append(adjustments, &PayerTransaction{
					CustomerID:    adj.CustomerID,
					Payer:         adj.Payer,
					Points:        leftToRemove,
					Entered:       adj.Entered,
					TransactionID: adj.TransactionID,
				})
				total += originalPoints - leftToRemove
			}
		} else {
			break
		}
	}
	err = spendPoints(ctx, id, adjustments)
	if err != nil {
		rlog.Error("unable to spend points on customer", "err", err)
		return nil, err
	}

	return &CustomerTransactions{
		CustomerID:   id,
		Name:         transactions[0].Payer,
		Transactions: adjustments,
	}, nil
}
