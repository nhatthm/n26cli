package transaction

import (
	"context"

	"github.com/bool64/ctxd"
	"github.com/nhatthm/go-clock"
	"github.com/nhatthm/n26api/pkg/transaction"
	"github.com/spf13/cobra"

	"github.com/nhatthm/n26cli/internal/io"
	"github.com/nhatthm/n26cli/internal/service"
	"github.com/nhatthm/n26cli/internal/time"
)

// TransactionsDeps is dependencies for finding transactions.
type TransactionsDeps interface {
	Clock() clock.Clock
	TransactionsFinder() transaction.Finder
	DataWriter() io.DataWriter
	CtxdLogger() ctxd.Logger
}

// NewTransactions creates a new `transactions` command.
func NewTransactions(l *service.Locator) *cobra.Command {
	var (
		from string
		to   string
	)

	cmd := &cobra.Command{
		Use:   "transactions",
		Short: "show all transactions in a time period",
		Long:  "show all transactions in a time period",
		RunE: func(_ *cobra.Command, _ []string) error {
			return findTransactions(context.Background(), l, from, to)
		},
	}

	cmd.Flags().StringVarP(&from, "from", "f", "", "start date, format: 2006-01-02T15:04:05Z or 2006-01-02")
	cmd.Flags().StringVarP(&to, "to", "t", "", "end date, format: 2006-01-02T15:04:05Z or 2006-01-02")

	return cmd
}

func findTransactions(ctx context.Context, deps TransactionsDeps, from, to string) error {
	start, end, err := time.Period(deps.Clock().Now(), from, to)
	if err != nil {
		return err
	}

	trans, err := deps.TransactionsFinder().FindAllTransactionsInRange(ctx, start, end)
	if err != nil {
		return err
	}

	return deps.DataWriter().WriteData(trans)
}
