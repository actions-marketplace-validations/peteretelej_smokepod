package smokepod

import (
	"context"

	"github.com/peteretelej/smokepod/pkg/smokepod/runners"
)

type Target interface {
	Exec(ctx context.Context, command string) (runners.ExecResult, error)
	Close() error
}
