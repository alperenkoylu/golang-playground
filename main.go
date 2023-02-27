package main

import (
	"encoding/json"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"

	"learning-golang/cache"
	"learning-golang/utils/database"
)

func main() {
	cache.Initialize()

	// create a new DatabaseProcess instance
	dbp := databaseprocess.New().AddQuery("SELECT * FROM kullanicilar WITH (NOLOCK) WHERE ad = @ad").AddParameter("ad", "Alperen")

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
