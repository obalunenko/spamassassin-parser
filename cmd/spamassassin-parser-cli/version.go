// spamassassin-parser-cli is a command line tool that shows how processing of reports works.
package main

import "fmt"

const unset = "unset"

var ( // build info
	version = unset
	date    = unset
	commit  = unset
)

func printVersion() {
	fmt.Printf("Version info: %s \n", version)
	fmt.Printf("Build date: %s \n", date)
	fmt.Printf("commit: %s \n", commit)
	fmt.Println()
}
