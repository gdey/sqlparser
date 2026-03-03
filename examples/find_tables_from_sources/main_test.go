package main

import (
	"fmt"
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
