package random_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ictsc/ictsc-outlands/backend/pkg/random"
	"github.com/stretchr/testify/require"
)

func TestNewString(t *testing.T) {
	t.Parallel()

	randomNum, err := rand.Int(rand.Reader, big.NewInt(100))
	require.NoError(t, err)

	type args struct {
		digit int
	}

	tests := []struct {
		name      string
		args      args
		assertion require.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				digit: int(randomNum.Int64() + 1),
			},
			assertion: require.NoError,
		},
		{
			name: "fail (digit is 0)",
			args: args{
				digit: 0,
			},
			assertion: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := random.NewString(tt.args.digit)
			tt.assertion(t, err)

			if err == nil {
				require.Len(t, got, tt.args.digit)
			}
		})
	}
}
