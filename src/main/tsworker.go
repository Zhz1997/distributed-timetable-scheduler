package main

//
// start a worker process, which is implemented
// in ../ts/worker.go. typically there will be
// multiple worker processes, talking to one master.
//
// go run tsworker.go wc.so
//
// Please do not change this file.
//

import (
	"fmt"
	"os"

	"../ts"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: Unable to start worker\n")
		os.Exit(1)
	}

	ts.Worker()
}
