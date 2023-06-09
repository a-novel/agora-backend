package communication

import (
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWaitForPing(t *testing.T) {
	t.Log("WARNING: this test uses timeouts and may take a certain time to complete")

	t.Run("Success", func(t *testing.T) {
		require.NoError(t, WaitForPingAuto(func() error {
			return nil
		}))
	})

	t.Run("Success/WithRetries", func(t *testing.T) {
		var count int

		require.NoError(t, WaitForPingAuto(func() error {
			if count > 3 {
				return nil
			}

			count++
			return validation.ErrNil
		}))
	})

	t.Run("Error/Timeout", func(t *testing.T) {
		require.ErrorIs(t, WaitForPingAuto(func() error {
			return validation.ErrNil
		}), validation.ErrNil)
	})
}
