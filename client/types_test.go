package client

import (
	"reflect"
	"testing"
)

func TestTransactionQuerry_ToMap(t *testing.T) {
	type fields struct {
		Type       TransactionType
		Source     TransactionSource
		Before     string
		After      string
		Commitment Commitment
		Limit      string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "Test 1",
			fields: fields{
				Type:       "transfer",
				Source:     "system_program",
				Before:     "",
				After:      "",
				Commitment: "confirmed",
				Limit:      "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := TransactionQuerry{
				Type:       tt.fields.Type,
				Source:     tt.fields.Source,
				Before:     tt.fields.Before,
				After:      tt.fields.After,
				Commitment: tt.fields.Commitment,
				Limit:      tt.fields.Limit,
			}
			if got := q.ToMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransactionQuerry.ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
