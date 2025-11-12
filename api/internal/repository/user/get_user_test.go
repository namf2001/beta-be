package user

import (
	"context"
	"errors"
	"testing"

	"beta-be/internal/pkg/testent"
	"beta-be/internal/repository/ent"

	"github.com/stretchr/testify/require"
)

func TestGetByEmail(t *testing.T) {
	type args struct {
		givenEmail string
		expErr     error
	}

	tcs := map[string]args{
		"success": {
			givenEmail: "admintest@gmail.com",
		},
		"err - user not found": {
			givenEmail: "admintest10@gmail.com",
			expErr:     errors.New("ent: user not found"),
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			testent.WithEntTx(t, func(tx *ent.Tx) {
				testent.LoadTestSQLFile(t, tx, "testdata/insert_user.sql")
				repo := New(tx.Client())
				_, err := repo.GetByEmail(context.Background(), tc.givenEmail)

				// then
				if tc.expErr != nil {
					require.EqualError(t, err, tc.expErr.Error())
				} else {
					require.NoError(t, err)
				}
			})
		})
	}
}
