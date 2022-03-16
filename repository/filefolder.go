package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mikkelstb/feedfetcher/feed"
)

type JsonFileFolder struct {
	path string
}

func (jsff *JsonFileFolder) String() string {
	return "json repository"
}

func NewJsonFileFolder(path string) (*JsonFileFolder, error) {

	var jsff JsonFileFolder
	jsff.path = path

	filepath, err := os.Stat(jsff.path)
	if err != nil {
		return nil, err
	}
	if !filepath.IsDir() {
		err = fmt.Errorf("error: %v is not a folder", path)
		return nil, err
	}
	return &jsff, nil
}

func (jsff JsonFileFolder) WriteSingle(ni feed.NewsItem) (string, error) {

	filename := fmt.Sprintf("%03d_%v_%v.json", ni.FeedId, ni.Docdate.Format("0601021504"), ni.Id())
	folder_path := filepath.Join(jsff.path, ni.Feed, ni.Docdate.Format("2006/01"))

	// Check if folder exists
	// If not: try create
	jsff.checkFolder(folder_path)

	// Check if file exist
	// If so ignore writing
	_, err := os.Stat(filepath.Join(folder_path, filename))
	if err == nil {
		return "", fmt.Errorf("file already exists")
	}

	//Convert NewsItem to NewsItemXML
	json_data, err := ni.ToJson()
	if err != nil {
		return "", err
	}

	//Attempt to write to folder
	err = os.WriteFile(filepath.Join(folder_path, filename), json_data, 0666)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (jsff JsonFileFolder) checkFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0766)
		if err != nil {
			return err
		}
	}
	return nil
}

func (jsff JsonFileFolder) Close() error {
	return nil
}
