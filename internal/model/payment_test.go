package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status PaymentStatus
		want   bool
	}{
		{
			name:   "valid status - ready",
			status: PaymentStatusReady,
			want:   true,
		},
		{
			name:   "valid status - paid",
			status: PaymentStatusPaid,
			want:   true,
		},
		{
			name:   "valid status - failed",
			status: PaymentStatusFailed,
			want:   true,
		},
		{
			name:   "valid status - cancelled",
			status: PaymentStatusCancelled,
			want:   true,
		},
		{
			name:   "invalid status",
			status: PaymentStatus("invalid"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}
