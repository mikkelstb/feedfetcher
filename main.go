package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/mikkelstb/feedfetcher/config"
	"github.com/mikkelstb/feedfetcher/feed"
)

var config_file string
var loop bool

func init() {
	flag.StringVar(&config_file, "config", "./config.json", "filepath for config file")
	flag.BoolVar(&loop, "loop", false, "set if program should run once per 2 hours")
}

func main() {

	for {
		fmt.Printf("Running feedfetcher @ %s", time.Now().Format("2006-01-02 15:04"))
		run()
		if !loop {
			break
		}
		time.Sleep(2 * time.Hour)
	}
}

func run() {

	// Read configfile
	cfg, err := config.Read(config_file)
	if err != nil {
		panic(err)
	}

	// Init Archive directive
	archive, err := NewArchive(cfg.Archive_path)
	if err != nil {
		panic(err)
	}

	// Load database
	db, err := setupDB(cfg.DB_file_path)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sources := getSources(cfg.Sources)

	// Loop through feeds in config
	for _, source := range sources {

		var filenames []string

		err := source.Process()
		if err != nil {
			fmt.Println(err)
			continue
		}

		newsitems := source.GetNewsitems()

		for i := range newsitems {
			filename, err := archive.writeNewsItemAsJson(newsitems[i])
			if err != nil {
				fmt.Println(err)
			} else {
				filenames = append(filenames, filename)
			}

			err = db.InsertNewsItem(newsitems[i])
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println("The following files were added:")
		fmt.Print(strings.Join(filenames, ", "))
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
