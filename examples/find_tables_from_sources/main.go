// FindTablesFromSources parses a SQL script and reports table names that are
// created using data from a given set of source tables (e.g. CREATE TABLE ... AS
// SELECT ... FROM source_tables).
package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

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
		parcels_csv AS csv USING(ain)
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
	posStmts, _, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
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

// FindTablesFromSourcesMulti parses a multi-statement SQL script (by splitting on ";\n")
// and returns table names that are created using data from any of the given source tables.
// Use this for scripts that contain BEGIN/COMMIT or multiple statements.
func FindTablesFromSourcesMulti(sql string, sourceSet map[string]struct{}) ([]string, error) {
	sql = strings.ReplaceAll(sql, "\r\n", "\n")
	segments := strings.Split(sql, ";\n")
	seen := make(map[string]struct{})
	var out []string
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" || seg == "BEGIN" || seg == "COMMIT" {
			continue
		}
		posStmts, _, err := sqlparser.Parse(seg)
		if err != nil {
			prefix := seg
			if len(prefix) > 60 {
				prefix = prefix[:60] + "..."
			}
			return nil, fmt.Errorf("parse segment %q: %w", prefix, err)
		}
		if len(posStmts) == 0 {
			continue
		}
		for _, ps := range posStmts {
			if ps.Statement == nil {
				continue
			}
			if ctas, ok := ps.Statement.(*sqlparser.CreateTableAsSelect); ok {
				if sqlparser.SelectStatementReferencesAny(ctas.Select, sourceSet) {
					name := string(ctas.Table)
					if _, ok := seen[name]; !ok {
						seen[name] = struct{}{}
						out = append(out, name)
					}
				}
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
