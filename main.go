package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mikkelstb/feedfetcher/config"
	"github.com/mikkelstb/feedfetcher/feed"
	"github.com/mikkelstb/feedfetcher/repository"
)

var config_file string
var infologger *log.Logger
var cfg *config.Config
var repositories []repository.Archive

var loop bool

func init() {
	flag.StringVar(&config_file, "config", "./config.json", "filepath for config file")
	flag.BoolVar(&loop, "loop", false, "set if program should run once per 2 hours")
}

func main() {

	flag.Parse()

	// Read configfile
	var err error
	cfg, err = config.Read(config_file)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("error: config file not read, aborting")
		os.Exit(1)
	}

	// Set up logfile
	logfile, err := os.OpenFile(cfg.Logfile_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	infologger = log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)

	//Set up repositories

	for _, rep := range cfg.Repositories {
		var err error
		var r repository.Archive

		if !rep.Active {
			continue
		}

		switch rep.Type {
		case "sqlite3":
			r, err = repository.NewSQLite(rep.Address)
			repositories = append(repositories, r)
		case "jsonfilefolder":
			r, err = repository.NewJsonFileFolder(rep.Address)
			repositories = append(repositories, r)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		infologger.Printf("running feedfetcher @ %s\n", time.Now().Format("2006-01-02 15:04"))
		run()
		if !loop {
			break
		}
		infologger.Printf("finished feedfetcher @ %s\n", time.Now().Format("2006-01-02 15:04"))
		time.Sleep(2 * time.Hour)
	}
}

func run() {

	sources := getSources(cfg.Sources)

	// Loop through feeds in config
	for _, source := range sources {

		err := source.Process()
		if err != nil {
			infologger.Println(err)
			continue
		}

		newsitems := source.GetNewsitems()

		for i := range newsitems {
			for rep := range repositories {
				result, err := repositories[rep].WriteSingle(newsitems[i])

				if err != nil {
					infologger.Println(err.Error())
				} else {
					infologger.Printf("added %v\n", result)
				}
			}
		}
	}
}

func getSources(conf []config.SourceConfig) []*feed.Source {

	var sources []*feed.Source
	for _, source_cfg := range conf {

		if !source_cfg.Active {
			continue
		}
		sources = append(sources, feed.NewSource(source_cfg))
	}
	return sources
}
