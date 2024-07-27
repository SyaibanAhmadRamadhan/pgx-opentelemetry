package pgxotel

import "go.opentelemetry.io/otel/attribute"

const (
	RowsAffectedKey    = attribute.Key("pgx.rows_affected")
	QueryParametersKey = attribute.Key("pgx.query.parameters")
	BatchSizeKey       = attribute.Key("pgx.batch.size")
	DBStatement        = attribute.Key("db.statement")
	DBSQLTable         = attribute.Key("db.table")
	SQLStateKey        = attribute.Key("pgx.sql_state")
	CopyColumns        = attribute.Key("pgx.copy.columns")
)
