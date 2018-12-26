package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	prettyjson "github.com/hokaccha/go-prettyjson"
)

func main() {
	helpFlag := flag.Bool("h", false, "print help information")

	flag.Usage = func() {
		usage := "Usage: %s [flags] file1.json file2.json ...\n\n"
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// for show `help` or empty json file list
	if *helpFlag || flag.NArg() < 1 {
		flag.Usage()
		return
	}

	// check json files is exist
	jsonFilenames := flag.Args()
	isExist := true
	for _, js := range jsonFilenames {
		if _, err := os.Stat(js); os.IsNotExist(err) {
			// json is not exists
			fmt.Fprintf(os.Stderr, "file `%s` is not exist\n", js)
			isExist = false
		}
	}
	if !isExist {
		os.Exit(1)
	}

	// show to os.Stdout
	for _, js := range jsonFilenames {
		fmt.Fprintf(os.Stdout, "Show `%s` file:\n", js)

		// read file data
		bs, err := ioutil.ReadFile(js)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Reading error : %v", err)
			os.Exit(1)
		}

		// unmarshal data
		var data interface{}
		json.Unmarshal(bs, &data)

		// print pretty result
		s, err := prettyjson.Marshal(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(s))
	}
}
