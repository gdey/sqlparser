// FindTablesFromSources parses a SQL script and reports table names that are
// created using data from a given set of source tables (e.g. CREATE TABLE ... AS
// SELECT ... FROM source_tables).
package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/gdey/sqlparser"
)

// Source table names we care about (parcels_data and parcels_csv from the script).
var sourceTables = map[string]struct{}{
	"parcels_data": {},
	"parcels_csv":  {},
}

// Sample SQL: Los Angeles County Addresses Transform Script (excerpt).
// The script creates the "addresses" table from parcels_data and parcels_csv.
// We use a minimal parseable subset (this parser does not support CREATE INDEX with
// column list, or USING(column) in JOIN, so those are simplified).
const sampleSQL = `
DELETE FROM parcels_data WHERE ain IS NULL;
DELETE FROM parcels_csv WHERE ain IS NULL;

CREATE TABLE addresses AS
	SELECT
		trim(data.ain) reference_id,
		'06' state_fips_code,
		trim(data.ain) parcel_apn,
		nullif(trim(csv.street_name), '') street_name,
		nullif(trim(csv.zip), '') postal,
		'CA' state
	FROM
		parcels_data AS data
	LEFT JOIN
		parcels_csv AS csv ON data.ain = csv.ain
	WHERE
		data.ain IS NOT NULL
	GROUP BY
		data.ain
;

UPDATE addresses SET jurisdiction = 'Example';
`

// FindTablesFromSources parses sql and returns table names that are created using
// data from any of the given source tables (e.g. CREATE TABLE ... AS SELECT ...
// FROM source_tables). sourceSet keys are the source table names.
func FindTablesFromSources(sql string, sourceSet map[string]struct{}) ([]string, error) {
	tree, _, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	posStmts, ok := tree.(sqlparser.PositionedStatements)
	if !ok {
		return nil, fmt.Errorf("expected PositionedStatements")
	}
	var out []string
	for _, ps := range posStmts {
		if ps.Statement == nil {
			continue
		}
		if ctas, ok := ps.Statement.(*sqlparser.CreateTableAsSelect); ok {
			if sqlparser.SelectStatementReferencesAny(ctas.Select, sourceSet) {
				out = append(out, string(ctas.Table))
			}
		}
	}
	sort.Strings(out)
	return out, nil
}

func main() {
	created, err := FindTablesFromSources(sampleSQL, sourceTables)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}
	fmt.Println("Tables created using data from parcels_data and parcels_csv:")
	for _, name := range created {
		fmt.Println(" ", name)
	}
}
