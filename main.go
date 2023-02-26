package main

import (
	"encoding/json"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"

	databaseprocess "learning-golang/utils/database"
)

func main() {
	// create a new DatabaseProcess instance
	dbp := databaseprocess.New()

	// add a SELECT query to the DatabaseProcess instance
	dbp = dbp.AddQuery("SELECT * FROM kullanicilar WITH (NOLOCK)")

	// execute the query and get the results
	result, err := dbp.GetDataTable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	jsonString := string(jsonBytes)
	fmt.Println("Result:", jsonString)
}
