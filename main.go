package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jrallison/go-workers"
)

var redisFlag = flag.String("redis", "redis://localhost:6379/0/", "Redis database to read queue from")
var exercismFlag = flag.String("exercism", "http://localhost:4567", "Url of exercism api, e.g. http://exercism.io")
var analysseurFlag = flag.String("analysseur", "http://localhost:8989", "Url of analysseur api, e.g. http://analysseur.exercism.io")

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s \n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	url, err := url.Parse(*redisFlag)
	if err != nil {
		panic(err)
	}
	config := map[string]string{
		// location of redis instance
		"server": url.Host,
		// instance of the database
		"database": strings.Trim(url.Path, "/"),
		// number of connections to keep open with redis
		"pool": "30",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "1",
	}
	if url.User != nil {
		password, ok := url.User.Password()
		if ok {
			config["password"] = password
		}
	}

	workers.Configure(config)

	analyzer := NewAnalyzer(*exercismFlag, *analysseurFlag, NewAuth().Key())
	workers.Process("analyze", analyzer.process, 4)

	workers.Run()
}
