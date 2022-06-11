package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mikkelstb/feedfetcher/feed"
)

var month_pattern regexp.Regexp = *regexp.MustCompile(`(^[0][1-9]||[1][012]$)`)
var year_pattern regexp.Regexp = *regexp.MustCompile(`^2\d{3}$`)
var year_month_pattern regexp.Regexp = *regexp.MustCompile(`(^2\d{3}([0][1-9])||([1][012])$)`)

type JsonFileFolder struct {
	path     string
	root     fs.FS
	max_days int
}

func (jsff *JsonFileFolder) String() string {
	return "json repository"
}

func NewJsonFileFolder(path string) (*JsonFileFolder, error) {

	var jsff JsonFileFolder
	jsff.path = path
	jsff.root = os.DirFS(path)

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

	filename := fmt.Sprintf("%s.json", ni.GetId())
	folder_path := filepath.Join(jsff.path, ni.Feed, fmt.Sprintf(ni.GetDocdate().Format("200601")))

	// Check if folder exists
	// If not: try create
	jsff.checkFolder(folder_path)

	// Check if file exist
	// If so ignore writing
	_, err := os.Stat(filepath.Join(folder_path, filename))
	if err == nil {
		return "", fmt.Errorf("file already exists")
	}

	//Convert NewsItem to Json
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

func (jsff JsonFileFolder) GetSources() []string {
	var dirs []string
	elements, err := ioutil.ReadDir(jsff.path)
	if err != nil {
		fmt.Println("OK!")
		return nil
	}

	for i := range elements {
		if elements[i].IsDir() {
			dirs = append(dirs, elements[i].Name())
		}
	}
	return dirs
}

func (jsff JsonFileFolder) GetNewsItem(filename string) (*feed.NewsItem, error) {
	ni := new(feed.NewsItem)

	json_file, err := os.ReadFile(filepath.Join(jsff.path, filename))
	if err != nil {
		return nil, err
	}

	json.Unmarshal(json_file, &ni)

	return ni, nil
}

func (jsff JsonFileFolder) getMonths(source string) []string {
	years, err := os.ReadDir(filepath.Join(jsff.path, source))
	if err != nil {
		fmt.Println(err)
	}
	var result []string
	for y := range years {
		if year_month_pattern.Match([]byte(years[y].Name())) && years[y].IsDir() {
			result = append(result, years[y].Name())
		}
	}
	return result
}

func (jsff JsonFileFolder) ListAllFiles(source string) []string {

	var files []string
	for _, month := range jsff.getMonths(source) {

		fs.WalkDir(jsff.root, filepath.Join(source, month), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Println(err.Error())
			}
			if !d.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	}
	return files
}

func (jsff JsonFileFolder) EraseOldArticles() (int, error) {
	return 0, errors.New("Not implemented yet")
}
