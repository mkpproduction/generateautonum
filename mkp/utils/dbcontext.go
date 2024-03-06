package mkputils

import (
	"context"
	"database/sql"
)

type (
	DBContext struct {
		Query string
		DB    *sql.DB
		Tx    *sql.Tx
	}
	dbcFunc func(r *DBContext)
)

func defaultDBContext() DBContext {
	return DBContext{}
}

func DB(DB *sql.DB) dbcFunc {
	return func(r *DBContext) {
		r.DB = DB
	}
}

func Tx(tx *sql.Tx) dbcFunc {
	return func(r *DBContext) {
		r.Tx = tx
	}
}

func Query(query string) dbcFunc {
	return func(r *DBContext) {
		r.Query = query
	}
}

func ExecuteRowContext(query string, dbc *sql.DB, txc *sql.Tx, args ...interface{}) (i interface{}, err error) {
	var stmt *sql.Stmt

	if txc != nil {
		stmt, err = txc.PrepareContext(context.Background(), query)
	} else {
		stmt, err = dbc.PrepareContext(context.Background(), query)
	}

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(context.Background(), args...).Scan(&i)
	if err != nil {
		return i, err
	}

	return i, nil
}

func ExecuteContext(query string, dbc *sql.DB, txc *sql.Tx, args ...interface{}) (i sql.Result, err error) {
	var stmt *sql.Stmt

	if txc != nil {
		stmt, err = txc.PrepareContext(context.Background(), query)
	} else {
		stmt, err = dbc.PrepareContext(context.Background(), query)
	}

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	i, err = stmt.ExecContext(context.Background(), args...)
	if err != nil {
		return i, err
	}

	return i, nil
}

func ParseID(result interface{}) int64 {
	return result.(int64)
}
