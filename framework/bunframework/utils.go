package bunframework

import (
	"context"
	"github.com/uptrace/bun"
)

// Generate extra options from configuration object.
func (config Config) computeInternalOptions() []bun.DBOption {
	options := config.Options

	if config.DiscardUnknownColumns {
		options = append(options, bun.WithDiscardUnknownColumns())
	}

	return options
}

// PingDatabase is a preconfigured ping function for communication.WaitForPing. It waits for the database to become
// available.
//
//	err := communication.WaitForPingAuto(bunframework.PingDatabase(ctx, db))
func PingDatabase(ctx context.Context, database *bun.DB) func() error {
	return func() error {
		return database.PingContext(ctx)
	}
}
