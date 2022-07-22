package points

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func Test_createCustomer(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "marley",
			args: args{
				ctx:  context.Background(),
				name: "marley",
			},
			// there's already one user pre-loaded as part of migrations
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createCustomer(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("createCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createCustomer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_insertTransaction(t *testing.T) {
	type args struct {
		ctx       context.Context
		customer  int
		payer     string
		points    int64
		timestamp time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "dummy transaction",
			args: args{
				ctx:       context.Background(),
				customer:  1,
				payer:     "DELOITTE",
				points:    500,
				timestamp: time.Time{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := insertTransaction(tt.args.ctx, tt.args.customer, tt.args.payer, tt.args.points, tt.args.timestamp); (err != nil) != tt.wantErr {
				t.Errorf("insertTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_selectCustomerBalances(t *testing.T) {
	type args struct {
		ctx      context.Context
		customer int
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "get balance",
			args: args{
				ctx:      context.Background(),
				customer: 1,
			},
			wantCount: 6,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectCustomerBalances(tt.args.ctx, tt.args.customer)
			if (err != nil) != tt.wantErr {
				t.Errorf("selectCustomerBalances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("selectCustomerBalances() got = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}

func Test_spendPoints(t *testing.T) {
	type args struct {
		ctx         context.Context
		id          int
		adjustments []*PayerTransaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "spend a little",
			args: args{
				ctx: context.Background(),
				id:  1,
				adjustments: []*PayerTransaction{
					&PayerTransaction{
						TransactionID: 1,
						CustomerID:    1,
						Payer:         "DANNON",
						Points:        500,
						Entered:       time.Time{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "spend a lot",
			args: args{
				ctx: context.Background(),
				id:  1,
				adjustments: []*PayerTransaction{
					&PayerTransaction{
						TransactionID: 2,
						CustomerID:    1,
						Points:        5000,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := spendPoints(tt.args.ctx, tt.args.id, tt.args.adjustments); (err != nil) != tt.wantErr {
				t.Errorf("spendPoints() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
