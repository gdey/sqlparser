package main

import (
	"fmt"
	"os"
	"sort"
)

func ExampleFindTablesFromSources() {
	names, err := FindTablesFromSources(sampleSQL, sourceTables)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("Tables created using data from parcels_data and parcels_csv:")
	for _, name := range names {
		fmt.Println(" ", name)
	}

	// Output:
	// Tables created using data from parcels_data and parcels_csv:
	//   addresses
}

func ExampleFindTablesFromSources_sanDiego() {
	sql, err := os.ReadFile("testdata/san_diego_parcel_transactions.sql")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	sourceSet := map[string]struct{}{"parcel_transactions_data": {}}
	names, err := FindTablesFromSourcesMulti(string(sql), sourceSet)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	sort.Strings(names)
	fmt.Println("Tables created using data from parcel_transactions_data:")
	for _, name := range names {
		fmt.Println(" ", name)
	}

	// Output:
	// Tables created using data from parcel_transactions_data:
	//   parcel_transaction_owners
	//   parcel_transactions
}
