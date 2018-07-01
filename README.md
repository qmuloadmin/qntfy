# Stats

A simple (toy) library and tools for processing basic stats about text files.

Essentially, scans provided text files concurrently, collects data on the number of characters/runes in each line, the number of words (whitespace separated, not word-boundary) in each line, and counts occurences of duplicate lines and words/tokens that match a list of keywords. 

## Installing the library:

`go get github.com/qmuloadmin/qntfy`

## Running the client

There are three ways to run the client, assuming you don't want to write your own. There is a cli, a stupidly basic web (REST) api, and a docker service. 

### CLI Usage

The cli can be built with `go build cli.go`, or in the docker image, it is already built (for Linux 64bit).

The cli expects to be provided with a `-k` (keyfile name) parameter, and a list of input files, provided after `-k`. Optionally, the default report filename (output.tsv) can be overridden with `-o` for output filename. 

e.g.

`./cli -k example_data/testkeys.txt example_data/test1.txt example_data/text2.txt example_data/test3.txt example_data/test4.txt`

To run the provided example input failed test 1-4.

### Web API Usage

The API binds to port 8080 (in docker service, binds to 80 externally). There are two routes, `/file/:filename` and `/stats/:outfile`:

`GET /file:filename` returns the contents of `filename` or 404.

E.g.
`GET localhost:8080/file/benchmarks.txt`

This can be used to retrieve the report file content after executing with POST:

`POST /stats/:outfile` will return 204 unless the payload is incorrect, in which case it will return 400.

E.g.
`POST localhost:8080/stats/output.tsv`

With payload (application/json):
```
{
    "keyFile": "example_data/testkeys.txt",
    "inputFiles": [
        "example_data/test1.txt",
        "example_data/test2.txt",
        "example_data/test3.txt",
        "example_data/test4.txt"
    ]
}
```

Will generate a new file called `output.tsv`. 

### Docker Service

The docker service builds both the cli and the web api for you, and runs the web api on port 80. 

All you need to do is `docker-compose up` in the directory with `docker-compose.yml`. 
