package pgxotel

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

func (t *Tracer) TraceCopyFromStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(DBSQLTable.String(data.TableName.Sanitize())),
	}

	if conn != nil {
		opts = append(opts, connAttrFromCfgPgx(conn.Config())...)
		opts = append(opts, trace.WithAttributes(
			CopyColumns.StringSlice(data.ColumnNames),
		))

	}

	ctx, _ = t.tracer.Start(ctx, "COPY FROM "+data.TableName.Sanitize(), opts...)

	return ctx
}

func (t *Tracer) TraceCopyFromEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceCopyFromEndData) {
	span := trace.SpanFromContext(ctx)
	recordError(span, data.Err)

	if data.Err != nil {
		span.SetAttributes(RowsAffectedKey.Int64(data.CommandTag.RowsAffected()))
	}

	span.End()
}
