package githubscraper

import (
	"net/http"
	"time"
)

type Scraper struct {
	client *http.Client
}

var defaultScraper *Scraper

// New creates a Scraper object
func New() *Scraper {
	return &Scraper{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}
