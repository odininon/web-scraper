package main

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"os"
)

type Scrapers struct {
	Scrapers []Scraper `json:"scrapers"`
}

type Scraper struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Query string `json:"query"`
}

type ScrapeResults struct {
	Name string
	Hits []string
	Url  string
}

func (s Scraper) Run() (ScrapeResults, error) {
	doc, err := htmlquery.LoadURL(s.Url)
	if err != nil {
		return ScrapeResults{}, err
	}

	list, err := htmlquery.QueryAll(doc, s.Query)
	if err != nil {
		return ScrapeResults{}, err
	}

	hits := make([]string, 0)
	for _, n := range list {
		hits = append(hits, fmt.Sprintf("%s", htmlquery.InnerText(n)))
	}

	return ScrapeResults{Name: s.Name, Hits: hits, Url: s.Url}, nil
}

type ScrapeRunner interface {
	Run() (ScrapeResults, error)
}

func loadScrapers() (scrapers Scrapers, err error) {
	jsonFile, err := os.Open("scrapers.json")
	if err != nil {
		return Scrapers{}, err
	}

	defer func() {
		if ferr := jsonFile.Close(); ferr != nil {
			err = ferr
		}
	}()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Scrapers{}, err
	}

	if err = json.Unmarshal(byteValue, &scrapers); err != nil {
		return Scrapers{}, err
	}

	return scrapers, nil
}
