package main

import (
	"fmt"
	"log"
	"strings"
	"web-scraper/database"
	"web-scraper/email"
)

func main() {
	db := database.NewDatabase()

	err := db.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scrapers, err := loadScrapers()
	if err != nil {
		log.Fatal(err)
	}

	changedHits := make([]email.Hit, 0)
	for _, scraper := range scrapers.Scrapers {
		scrape, err := scraper.Run()
		if err != nil {
			log.Fatal(err)
		}

		changed, err := db.Update(scrape.Name, strings.Join(scrape.Hits[:], ","))

		fmt.Printf("Scrape: %s Changed: %v\n", scrape.Name, changed)

		if changed {
			changedHits = append(changedHits, email.Hit{Name: scrape.Name, Url: scrape.Url})
		}
	}

	if len(changedHits) > 0 {
		email := email.NewEmail()
		_, err := email.Send(changedHits)

		if err != nil {
			fmt.Println(err)
		}
	}
}
