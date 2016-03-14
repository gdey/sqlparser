/*
sqlfmt — Will take an sql file and print out a pretty version, if it can not parse the sql file. It will
then show where the error is and exit with an error code of 1
*/
package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/gdey/sqlparser"
)

const SytnaxErrorCode = 1

// FormatSQL will walk the tree and produce "pretty" SQL.
// Given a statement like: select a,b from a where a = 1
// Return
//    SELECT a, b
//    FROM a
//    WHERE a = 1;
func formatSQL(tree sqlparser.Statement) (string, error) {
	return "", nil
}
func main() {
	// We will have flags later.
	flag.Parse()
	statusCode := 0

	for i := 1; i < len(os.Args); i++ {
		file := os.Args[i]
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		buff := string(buf)
		_, err = sqlparser.Parse(buff)
		if err != nil {
			tokenError, _ := err.(*sqlparser.TokenizerError)
			FormatErrorMessage(os.Stderr, tokenError, file, buff, true)
			statusCode = SytnaxErrorCode
			continue
		}
		// fmt.Println(sqlparser.String(statement))

	}
	os.Exit(statusCode)
}
