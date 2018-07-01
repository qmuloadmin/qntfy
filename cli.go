package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/qmuloadmin/qntfy/stats"
)

func main() {
	// We need to read in a list of files of arbitrary length, an output filename and a keyword file.
	// The output filename will default to output.tsv
	var outName, keyFile string
	flag.StringVar(&outName, "o", "output.tsv", "The filename of the output report")
	flag.StringVar(&keyFile, "k", "", "The name of the keyfile")
	flag.Parse()
	inputFiles := flag.Args()
	// make sure values were provided
	if keyFile == "" {
		log.Fatal("Must provide keyword file argument '-k'")
	}
	if len(inputFiles) == 0 {
		log.Fatal("Must specify at least one input file name")
	}
	err := stats.ProcessFiles(outName, keyFile, inputFiles)
	if err != nil {
		fmt.Println("Error encountered processing files: ", err)
	}
}
