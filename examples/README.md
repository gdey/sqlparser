# Examples

## find_tables_from_sources

Parses a SQL script and reports **table names that are created using data from** a given set of source tables.

- Detects `CREATE TABLE ... AS SELECT ... FROM source_tables` (CTAS).
- Source tables are configured in the example as `parcels_data` and `parcels_csv`.
- The sample SQL (simplified from the LA County Addresses Transform Script) creates the `addresses` table from those two sources; the example outputs `addresses`.

Run from the repository root:

```bash
go run ./examples/find_tables_from_sources/
```

Or run the testable example (output is verified by `go test`):

```bash
go test -v -run ExampleFindTablesFromSources ./examples/find_tables_from_sources/
```

Expected output:

```
Tables created using data from parcels_data and parcels_csv:
  addresses
```

The example uses the sqlparser’s `CreateTableAsSelect` AST node and `TableNamesFromTableExprs` / `SelectStatementReferencesAny` from the analyzer to walk the FROM clause (including JOINs and subqueries) and check for references to the source table names.
