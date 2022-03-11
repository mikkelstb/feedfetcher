package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mikkelstb/feedfetcher/config"
	"github.com/mikkelstb/feedfetcher/feed"
)

var config_file string
var infologger *log.Logger
var cfg *config.Config

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

	// Init Archive directive
	archive, err := NewArchive(cfg.Archive_path)
	if err != nil {
		infologger.Panic(err)
	}

	// Load database
	db, err := setupDB(cfg.DB_file_path)
	if err != nil {
		infologger.Panic(err)
	}
	defer db.Close()

	sources := getSources(cfg.Sources)

	// Loop through feeds in config
	for _, source := range sources {

		var filenames []string

		err := source.Process()
		if err != nil {
			infologger.Println(err)
			continue
		}

		newsitems := source.GetNewsitems()

		for i := range newsitems {
			filename, err := archive.writeNewsItemAsJson(newsitems[i])
			if err != nil {
				infologger.Println(err)
			} else {
				filenames = append(filenames, filename)
			}

			err = db.InsertNewsItem(newsitems[i])
			if err != nil {
				infologger.Println(err)
			}
		}
		fmt.Printf("Number of files added: %v\n", len(filenames))
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

func setupDB(filename string) (*SQLiteRepository, error) {

	db_config, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	return NewSQLiteRepository(db_config), nil
}
