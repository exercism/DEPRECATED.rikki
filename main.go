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

	workers.Configure(redisConfig())

	analyzer := NewAnalyzer(*exercismFlag, *analysseurFlag, NewAuth().Key())
	workers.Process("analyze", analyzer.process, 4)

	workers.Run()
}

func redisConfig() map[string]string {
	url, err := url.Parse(*redisFlag)
	if err != nil {
		panic(err)
	}
	config := map[string]string{
		"server":   url.Host,
		"database": strings.Trim(url.Path, "/"),
		"pool":     "30",
		"process":  "1",
	}
	if url.User != nil {
		password, ok := url.User.Password()
		if ok {
			config["password"] = password
		}
	}
	return config
}

func commentDir() string {
	dir := os.Getenv("RIKKI_FEEDBACK_DIR")
	if dir == "" {
		dir = "comments"
	}
	return dir
}
