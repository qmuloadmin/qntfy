package main

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/qmuloadmin/qntfy/stats"
)

type Filenames struct {
	KeyFile    string
	InputFiles []string
}

func main() {
	router := gin.Default()
	// e.g. /file/benchmarks.txt to get stored benchmarks
	// or /file/output.tsv to get previous results
	router.GET("/file/:name", getFile)
	router.POST("/stats/:outfile", postStats)
	router.Run()
}

// this is like... VERY poorly implemented right now. It doesn't check to see if the
// outfile already exists.
func postStats(context *gin.Context) {
	outFile := context.Param("outfile")
	filenames := new(Filenames)
	// marshall payload into filenames struct
	context.BindJSON(filenames)
	if err := stats.ProcessFiles(outFile, filenames.KeyFile, filenames.InputFiles); err == nil {
		context.Status(204)
	} else {
		context.AbortWithStatus(400)
	}
}

func getFile(context *gin.Context) {
	file := context.Param("name")
	contents, err := ioutil.ReadFile(file)
	writer := bytes.NewReader(contents)
	if err != nil {
		// assume for now that the file doesn't exist.
		context.AbortWithStatus(404)
		return
	}
	context.DataFromReader(200, int64(len(contents)), "text/plain", writer, map[string]string{})
}
