package main

//
// start the master process, which is implemented
// in ../ts/master.go
//
// go run tsmaster.go pg*.txt
//
// Please do not change this file.
//

import (
	"fmt"
	"os"
	"time"

	"../ts"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: tsmaster inputfiles...\n")
		os.Exit(1)
	}

	m := ts.MakeMaster(os.Args[1:], 10)
	for m.Done() == false {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
