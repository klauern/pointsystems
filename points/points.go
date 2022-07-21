package points

import (
	"context"
	"sort"
	"time"

	"encore.dev/rlog"
)

type PayerTransaction struct {
	TransactionID int       `json:"id"`
	CustomerID    int       `json:"customer_id"`
	Payer         string    `json:"payer"`
	Points        int64     `json:"points"`
	Entered       time.Time `json:"timestamp"`
}

type CustomerTransactions struct {
	CustomerID   int                 `json:"id"`
	Name         string              `json:"name"`
	Transactions []*PayerTransaction `json:"transactions"`
}

type Customer struct {
	Name         string `json:"name"`
	ID           int    `json:"id"`
	TotalBalance int    `json:"total_balance"`
}

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

//encore:api public method=POST path=/points/:customer
func AddTransaction(ctx context.Context, customer int, transaction *PayerTransaction) error {
	err := insertTransaction(ctx, customer, transaction.Payer, transaction.Points, transaction.Entered)
	return err
}

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

//encore:api public method=POST path=/spend/:id/:points
func SpendCustomerPoints(ctx context.Context, id int, points int64) (*CustomerTransactions, error) {
	transactions, err := selectCustomerBalances(ctx, id)
	if err != nil {
		return nil, err
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Entered.Before(transactions[j].Entered)
	})

	adjustments := make([]*PayerTransaction, 0)
	total := int64(0)
	for total < points {
		for _, adj := range transactions {
			if adj.Points <= (points - total) {
				adjustments = append(adjustments, &PayerTransaction{
					CustomerID: adj.CustomerID,
					Payer:      adj.Payer,
					Points:     0,
					Entered:    adj.Entered,
				})
				total += adj.Points
			} else {
				adjustments = append(adjustments, &PayerTransaction{
					CustomerID: adj.CustomerID,
					Payer:      adj.Payer,
					Points:     total - points,
					Entered:    adj.Entered,
				})
				total += points - total
			}
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
