package notifybox

import (
	"context"

	"route256/loms/internal/loms/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

type Notifier struct {
	db DBTX
}

func New(db DBTX) *Notifier {
	return &Notifier{db: db}
}

func (n *Notifier) WithTx(tx pgx.Tx) *Notifier {
	return &Notifier{db: tx}
}

func (n *Notifier) CreateEvent(ctx context.Context, event models.Event) error {
	qry := `INSERT INTO outbox.notifier (order_id, status)
			VALUES ($1, $2);`

	if _, err := n.db.Exec(ctx, qry, event.OrderID, event.Status); err != nil {
		return err
	}

	return nil
}

func (n *Notifier) FetchNextBatch(ctx context.Context, batchSize int) ([]*models.Event, error) {
	qry := `SELECT id, order_id, status, created_at
			 FROM outbox.notifier
			 WHERE is_sent = false
			 ORDER BY id
			 LIMIT $1 FOR UPDATE SKIP LOCKED`

	rows, err := n.db.Query(ctx, qry, batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[models.Event])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Notifier) MarkAsSent(ctx context.Context, events []*models.Event) error {
	if len(events) == 0 {
		return nil
	}

	ids := make([]int64, len(events))
	for i, event := range events {
		ids[i] = event.ID
	}

	qry := `UPDATE outbox.notifier
	        SET is_sent = true
	        WHERE id = ANY($1)`

	_, err := n.db.Exec(ctx, qry, ids)
	if err != nil {
		return err
	}

	return nil
}
