package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/exercism/rikki/analysis/crystal"
	"github.com/exercism/rikki/analysis/ruby"
	"github.com/jrallison/go-workers"
)

var redisFlag = flag.String("redis", "redis://localhost:6379/0/", "Redis database to read queue from")
var exercismFlag = flag.String("exercism", "http://localhost:4567", "Url of exercism api, e.g. http://exercism.io")
var analysseurFlag = flag.String("analysseur", "http://localhost:8989", "Url of analysseur api, e.g. http://analysseur.exercism.io")
var crystalAnalyzerFlag = flag.String("crystal-analyzer", "http://localhost:3000", "Url of crystal-analyzer api, e.g. http://crystal-analyzer.exercism.io")

var lgr = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s \n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	workers.Configure(redisConfig())

	exercism := NewExercism(*exercismFlag, NewAuth().Key())

	ruby.Host = *analysseurFlag
	crystal.Host = *crystalAnalyzerFlag

	analyzer, err := NewAnalyzer(exercism, commentDir())
	if err != nil {
		lgr.Print(err)
		os.Exit(1)
	}
	workers.Process("analyze", analyzer.process, 4)

	hello, err := NewHello(exercism, commentDir())
	if err != nil {
		lgr.Print(err)
		os.Exit(1)
	}
	workers.Process("hello", hello.process, 4)

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

func read(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}
