package dat

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v4"
)

const (
	requestLogsTable           = "request_logs"
	responseBodyTable           = "response_logs"
	requestLogsTableBodyColumn = "body"
	responseLogsTableBodyColumn = "body"
)

// AccessLogMapper provides CRUD methods for the tables request_logs and response_logs
type AccessLogMapper struct {
	ctx  context.Context
	txn  pgx.Tx
	psql sq.StatementBuilderType
}

// NewAccessLogMapper creates a new AccessLogMapper for transaction
func NewAccessLogMapper(ctx context.Context, txn pgx.Tx) *AccessLogMapper {
	return &AccessLogMapper{
		ctx:  ctx,
		txn:  txn,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (al *AccessLogMapper) InsertRequest(input string) error {
	sql, params, err := al.psql.
		Insert(requestLogsTable).
		Columns(requestLogsTableBodyColumn).
		Values(input).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := al.txn.Query(al.ctx, sql, params...)
	rows.Close()
	return err
}

func (al *AccessLogMapper) InsertResponse(input string) error {
	sql, params, err := al.psql.
		Insert(responseBodyTable).
		Columns(responseLogsTableBodyColumn).
		Values(input).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := al.txn.Query(al.ctx, sql, params...)
	rows.Close()
	return err
}
