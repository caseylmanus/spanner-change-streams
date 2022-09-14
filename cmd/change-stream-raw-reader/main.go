package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"cloud.google.com/go/spanner"
	"github.com/caseylmanus/spanner-change-streams/changestreams"
)

func main() {
	var db, instance, project, name string
	flag.StringVar(&db, "d", "", "spanner database name")
	flag.StringVar(&instance, "i", "", "spanner instance name")
	flag.StringVar(&project, "p", "", "gcp project name")
	flag.StringVar(&name, "n", "", "change stream name")
	flag.Parse()
	if missingArg(db, instance, project, name) {
		flag.Usage()
		os.Exit(1)
	}
	ctx := context.Background()
	connString := fmt.Sprintf("projects/%s/instances/%s/databases/%s", project, instance, db)

	client, err := spanner.NewClient(ctx, connString)
	if err != nil {
		fmt.Printf("Cannot connect to %s using application default credentials", connString)
		os.Exit(1)
	}
	response, errors := changestreams.Subscribe(ctx, client, name)
	for {
		select {
		case dc := <-response:
			fmt.Println(dc)
		case err := <-errors:
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
func missingArg(args ...string) bool {
	for _, a := range args {
		if a == "" {
			return true
		}
	}
	return false
}
