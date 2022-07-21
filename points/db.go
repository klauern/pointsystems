package points

import (
	"context"
	"time"

	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

func insertTransaction(ctx context.Context, customer int, payer string, points int64, timestamp time.Time) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO transactions (customer, payer, amount, entered_time)
		VALUES ($1, $2, $3, $4)
`, customer, payer, points, timestamp)
	return err
}

const SelectCustomer = `
SELECT id, customer, amount, payer, entered_time FROM transactions where customer = $1 LIMIT 100;
`

func selectCustomerBalances(ctx context.Context, customer int) ([]*PayerTransaction, error) {
	transactions, err := sqldb.Query(ctx, SelectCustomer, customer)

	if err != nil {
		return nil, err
	}
	t := make([]*PayerTransaction, 0)
	for transactions.Next() {
		var transaction PayerTransaction
		err = transactions.Scan(
			&transaction.TransactionID,
			&transaction.CustomerID,
			&transaction.Points,
			&transaction.Payer,
			&transaction.Entered)
		if err != nil {
			rlog.Error("not able to scan row", "err", err)
			return t, err
		}
		//rlog.Info("found row", "row", t)
		t = append(t, &transaction)
	}
	return t, nil
}

func createCustomer(ctx context.Context, name string) (int, error) {
	id := 0
	tx, err := sqldb.Begin(ctx)
	if err != nil {
		return id, err
	}
	err = tx.QueryRow(ctx, `INSERT INTO customers (name) VALUES ($1) RETURNING id;`, name).Scan(&id)
	if err != nil {
		return id, err
	}

	_, err = tx.Exec(ctx, `INSERT INTO customer_totals (customer) VALUES ('"$1" => "$2"')`, id, 0)
	if err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func updateTransaction(ctx context.Context, id int, amount int) error {
	_, err := sqldb.Exec(ctx, `UPDATE transactions SET amount = $1 WHERE id = $2;`, amount, id)
	return err
}

func deleteTransaction(ctx context.Context, id int) error {
	_, err := sqldb.Exec(ctx, `DELETE FROM transactions WHERE id = $1;`, id)
	return err
}

func spendPoints(ctx context.Context, id int, adjustments []*PayerTransaction) error {
	tx, err := sqldb.Begin(context.Background())
	if err != nil {
		rlog.Error("not able to create database transaction", "err", err)
		return err
	}
	for _, adj := range adjustments {
		rlog.Info("adjusting", "id", id, "amount", adj.Points, "transaction_id", adj.TransactionID)
		if adj.Points == 0 {
			result, err := tx.Exec(ctx, `DELETE FROM transactions WHERE id = $1;`, adj.TransactionID)
			if err != nil {
				rlog.Error("not able to delete row", "transaction_id", adj.TransactionID, "err", err)
				return err
			}
			rlog.Info("affected rows", "rowcount", result.RowsAffected())
		} else {
			result, err := tx.Exec(ctx, `UPDATE transactions SET amount = (amount - $1) WHERE id = $2;`, adj.Points, adj.TransactionID)
			if err != nil {
				rlog.Error("unable to update balance for transaction", "transation_id", adj.TransactionID, "err", err)
				return err
			}
			rlog.Info("affected rows", "rowcount", result.RowsAffected())
		}
	}
	err = tx.Commit()
	if err != nil {
		rlog.Error("error committing transaction", "err", err)
	}
	return err
}
