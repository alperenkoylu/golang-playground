package databaseprocess

import (
	"database/sql"
	"strconv"
)

const connectionString = "server=localhost;user id=sa;password=1234;database=test"

type DatabaseProcess struct {
	query      string
	parameters map[string]interface{}
	namedArgs  []sql.NamedArg
}

func New() *DatabaseProcess {
	return &DatabaseProcess{
		query:      "",
		parameters: make(map[string]interface{}),
		namedArgs:  make([]sql.NamedArg, 0),
	}
}

func (dbp *DatabaseProcess) AddQuery(query string) *DatabaseProcess {
	if dbp.parameters == nil {
		dbp.parameters = make(map[string]interface{})
	}
	return &DatabaseProcess{
		query:      query,
		parameters: dbp.parameters,
	}
}

func (dbp *DatabaseProcess) AddParameter(key string, value interface{}) *DatabaseProcess {
	if dbp.parameters == nil {
		dbp.parameters = make(map[string]interface{})
	}
	dbp.parameters[key] = value
	namedArgs := make([]sql.NamedArg, 0, len(dbp.parameters))
	for k, v := range dbp.parameters {
		namedArgs = append(namedArgs, sql.Named(k, v))
	}
	return &DatabaseProcess{
		query:      dbp.query,
		parameters: dbp.parameters,
		namedArgs:  namedArgs,
	}
}

func (dbp *DatabaseProcess) GetDataTable() (map[string]interface{}, error) {
	db, err := sql.Open("sqlserver", connectionString)
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

	// Prepare the query with named parameters
	namedArgs := make([]sql.NamedArg, len(dbp.parameters))
	i := 0
	for k, v := range dbp.parameters {
		namedArgs[i] = sql.Named(k, v)
		i++
	}
	// Prepare the query with placeholders for the parameters
	stmt, err := tx.Prepare(dbp.query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Convert named arguments to interface slice
	var args []interface{}
	for _, arg := range namedArgs {
		args = append(args, arg.Value)
	}

	// Execute the query with the parameters
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build the result map
	result := make(map[string]interface{})
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
		result[strconv.FormatInt(rowMap["id"].(int64), 10)] = rowMap
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
