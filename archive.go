package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mikkelstb/feedfetcher/feed"
)

func init() {

	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	log.SetOutput(f)
}

type Archive struct {
	path string
}

func NewArchive(path string) (*Archive, error) {

	var archive Archive
	archive.path = path

	fmt.Println("Hello World. Initializing Archive")

	filepath, err := os.Stat(archive.path)
	if err != nil {
		return nil, err
	}
	if !filepath.IsDir() {
		err = fmt.Errorf("error: %v is not a folder", path)
		return nil, err
	}
	return &archive, nil
}

func (a Archive) writeNewsItemAsJson(ni feed.NewsItem) (string, error) {
	filename := fmt.Sprintf("%03d_%v_%v.json", ni.FeedId, ni.Docdate.Format("0601021504"), ni.Id())
	//filename := strings.Join([]string{ni.FeedId, ni.Docdate.Format("_0601021504_"), ni.GetId(), ".json"}, "")
	folder_path := filepath.Join(a.path, ni.Feed, ni.Docdate.Format("2006/01"))

	// Check if folder exists
	// If not: try create
	a.checkFolder(folder_path)

	//Convert NewsItem to NewsItemXML
	json_data, err := ni.ToJson()
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(filepath.Join(folder_path, filename)); err == nil {
		return "", fmt.Errorf("skipping file %v, file already exists", filename)
	}

	err = os.WriteFile(filepath.Join(folder_path, filename), json_data, 0666)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (a Archive) checkFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0766)
		if err != nil {
			log.Default().Println(err.Error())
			return err
		}
	}
	return nil
}

// func (a Archive) getSources() ([]string, error) {
// 	list, err := os.ReadDir(a.path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var sources []string
// 	for i := range list {
// 		if list[i].IsDir() {
// 			sources = append(sources, list[i].Name())
// 		}
// 	}
// 	return sources, nil
// }

// func (a Archive) writeNewsItem(ni feed.NewsItem) (string, error) {
// 	filename := strings.Join([]string{ni.Feed, ni.Docdate.Format("_20060102_"), ni.GetId(), ".xml"}, "")
// 	folder_path := filepath.Join(a.path, ni.Feed, ni.Docdate.Format("2006/01/02"))

// 	// Check if folder exists
// 	// If not: try create
// 	a.checkFolder(folder_path)

// 	//Convert NewsItem to NewsItemXML
// 	xml_data, _ := ni.ToXml()

// 	if _, err := os.Stat(filepath.Join(folder_path, filename)); err == nil {
// 		return "", fmt.Errorf("skipping file %v, file already exists", filename)
// 	}

// 	err := os.WriteFile(filepath.Join(folder_path, filename), xml_data, 0666)
// 	if err != nil {
// 		return "", err
// 	}

// 	return filename, nil
// }
