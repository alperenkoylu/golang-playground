package databaseprocess

import (
	"context"
	"database/sql"
	"fmt"
	"learning-golang/cache"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type DatabaseProcess struct {
	query      string
	parameters map[string]interface{}
}

func New() *DatabaseProcess {
	return &DatabaseProcess{
		query:      "",
		parameters: make(map[string]interface{}),
	}
}

func (dbp *DatabaseProcess) AddQuery(query string) *DatabaseProcess {
	dbp.query = query
	return dbp
}

func (dbp *DatabaseProcess) AddParameter(key string, value interface{}) *DatabaseProcess {
	dbp.parameters[key] = value
	return dbp
}

func (dbp *DatabaseProcess) GetDataTable() ([]map[string]interface{}, error) {
	fmt.Println("connstr", cache.DatabaseConfig.ConnectionString)
	db, err := sql.Open("sqlserver", cache.DatabaseConfig.ConnectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cache.DatabaseConfig.TimeoutAsSeconds)*time.Second)
	defer cancel()

	// Prepare the query with named parameters
	namedParams := make([]interface{}, 0, len(dbp.parameters))

	if len(dbp.parameters) > 0 {
		for key, value := range dbp.parameters {
			namedParams = append(namedParams, sql.Named(key, value))
		}
	}

	// Prepare the query with placeholders for the parameters
	stmt, err := tx.PrepareContext(ctx, dbp.query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query with the parameters
	rows, err := stmt.QueryContext(ctx, namedParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build the result list
	var result []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val == nil {
				rowMap[col] = nil
			} else {
				rowMap[col] = val
			}
		}
		result = append(result, rowMap)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}
