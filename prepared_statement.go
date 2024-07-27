package pgxotel

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (t *Tracer) TracePrepareStart(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
	}

	prepareStmtNameKey := attribute.Key("pgx.prepared_statement")
	if data.Name != "" {
		trace.WithAttributes(prepareStmtNameKey.String(data.Name))
	}

	if conn != nil {
		opts = append(opts, connAttrFromCfgPgx(conn.Config())...)
		opts = append(opts, trace.WithAttributes(DBStatement.String(data.SQL)))
	}

	spanName := "PREPARED STATEMENT | " + data.SQL

	ctx, _ = t.tracer.Start(ctx, spanName, opts...)

	return ctx
}

func (t *Tracer) TracePrepareEnd(ctx context.Context, _ *pgx.Conn, data pgx.TracePrepareEndData) {
	span := trace.SpanFromContext(ctx)
	recordError(span, data.Err)

	span.End()
}
