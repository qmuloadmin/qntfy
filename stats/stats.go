package stats

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/montanaflynn/stats"
)

var dupCount int64
var dupLock = sync.Mutex{}

// ProcessFiles builds an aggregate report for all files provided in inputFiles
func ProcessFiles(outName string, keyFile string, inputFiles []string) error {
	// initialize seenLines
	var seenLines set
	seenLines.init()

	// read in and parse the keywords
	if data, err := ioutil.ReadFile(keyFile); ok(err) {
		keywordSlice := strings.Split(string(data), "\n")
		keywords := *new(set)
		keywords.init()
		for _, word := range keywordSlice {
			keywords.add(strings.TrimSpace(word))
		}
		/*
		 For each file, and then each line in each file, we need to compute
		 running median and std dev for two different metrics, and count tokens
		 There is a lot of information missing that could affect the design of this.
		 for instance, if the inputFiles are relatively equal length, and there are
		 always more than 2, it's probably most efficient to parse each file in a
		 separate goroutine than it is to try to parse alternating lines.
		 Running line by line in goroutines can't be super efficient... there's just
		 not much computation involved, context switching may only involve three registers
		 but that's not nothing, and we also have to handle mutexes locking and unlocking
		 eww. I don't like lacking information.
		 so here's what we're going to do for now. May change mind later. Assume files
		 are similar in length, and that there are always more than 2. I think its
		 reasonable... and it's the simplest to implement. Lazy? Probably.
		*/
		count := len(inputFiles)
		writer := resultWriter{}
		writer.init(count)
		for _, file := range inputFiles {
			go processFile(file, &seenLines, keywords, &writer)
		}

		// for each result written to the channel, update/extend the total
		// counts for calculating stats.
		rep := report{}
		totalLengthCounts := make([]float64, 0, 1000)
		totalTokenCounts := make([]float64, 0, 1000)
		for result := range writer.c {
			extend(&totalLengthCounts, result.charCounts)
			extend(&totalTokenCounts, result.tokenCounts)
			rep.kwMerge(result.keywords)
		}

		// Build and write the report
		lm, _ := stats.Median(totalLengthCounts)
		ls, _ := stats.StandardDeviation(totalLengthCounts)
		tm, _ := stats.Median(totalTokenCounts)
		ts, _ := stats.StandardDeviation(totalTokenCounts)
		rep.dCount = dupCount
		rep.devLength = ls
		rep.devTokens = ts
		rep.medLength = lm
		rep.medTokens = tm
		if file, err := os.Create(outName); ok(err) {
			fmt.Fprint(file, rep)
		} else {
			return errors.New("Unable to write to outfile: " + err.Error())
		}
	} else {
		return errors.New("Unexpected error; failed to read keyfile: " + err.Error())
	}
	return nil
}

func extend(sl1 *[]float64, sl2 []float64) {
	for _, v := range sl2 {
		*sl1 = append(*sl1, v)
	}
}

func processLine(line string, seenLines *set, keywords set, kwCount map[string]int) (float64, float64) {
	line = strings.TrimSpace(line)
	// if the line has already been seen, increment the global counter
	if seenLines.contains(line) {
		/*
			if these two mutexes become a significant slowdown (requires profiling
			of a large dataset to really know) we would probably make another goroutine
			and (buffered) channel that is responsible solely for handling line duplication.
			It would not need a mutex and would not block. Then the main routine would
			block on that routine's completion along with the individual file routines.
		*/
		inc()
	} else {
		// otherwise, add the line to seen lines
		seenLines.safeAdd(line)
	}
	// split the line over whitespace, keep track of how many tokens we inspected
	count := 0
	for _, token := range strings.Fields(line) {
		if keywords.contains(token) {
			kwCount[token]++
		}
		count++
	}
	return float64(len(line)), float64(count)
}

func processFile(fileName string, seenLines *set, keywords set, writer *resultWriter) {
	kwCount := make(map[string]int)
	lens := make([]float64, 0, 1000)
	num := make([]float64, 0, 1000)
	if file, err := os.Open(fileName); ok(err) {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		// read the file line by line
		for scanner.Scan() {
			line := scanner.Text()
			ln, count := processLine(line, seenLines, keywords, kwCount)
			lens = append(lens, ln)
			num = append(num, float64(count))
		}
	} else {
		log.Println("Failed to open input file: ", err)
	}
	r := result{
		charCounts:  lens,
		tokenCounts: num,
		keywords:    kwCount,
	}
	writer.write(r)
}
