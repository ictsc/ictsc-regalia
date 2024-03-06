package random_test

import (
	"math/rand/v2"
	"testing"

	"github.com/ictsc/ictsc-outlands/backend/pkg/random"
	"github.com/stretchr/testify/require"
)

func TestNewString(t *testing.T) {
	t.Parallel()

	type args struct {
		digit uint32
	}

	tests := []struct {
		name      string
		args      args
		assertion require.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				digit: uint32(rand.IntN(100) + 1),
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

	for _, tt := range tests { // nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := random.NewString(tt.args.digit)
			tt.assertion(t, err)

			if err == nil {
				require.Len(t, got, int(tt.args.digit))
			}
		})
	}
}
