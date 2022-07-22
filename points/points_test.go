package points

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestAddTransaction(t *testing.T) {
	type args struct {
		ctx         context.Context
		customer    int
		transaction *PayerTransaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "little",
			args: args{
				ctx:      context.Background(),
				customer: 1,
				transaction: &PayerTransaction{
					CustomerID: 1,
					Payer:      "PONYSHOW",
					Points:     150000,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddTransaction(tt.args.ctx, tt.args.customer, tt.args.transaction); (err != nil) != tt.wantErr {
				t.Errorf("AddTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCustomerBalances(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name      string
		args      args
		wantTotal int
		wantErr   bool
	}{
		{
			name: "get nick's balance",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantTotal: 500,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCustomerBalances(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCustomerBalances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			total := int64(0)
			for _, c := range got.Transactions {
				fmt.Println("points", c.Points)
				total += c.Points
			}
			if !reflect.DeepEqual(total, tt.wantTotal) {
				t.Errorf("GetCustomerBalances() got = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}

func TestNewCustomer(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *Customer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCustomer(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCustomer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpendCustomerPoints(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     int
		points int64
	}
	tests := []struct {
		name    string
		args    args
		want    *CustomerTransactions
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SpendCustomerPoints(tt.args.ctx, tt.args.id, tt.args.points)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpendCustomerPoints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpendCustomerPoints() got = %v, want %v", got, tt.want)
			}
		})
	}
}
