package githubscraper

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math"
	"net/http"
	"net/url"
)

const (
	resultPerPage int = 10
	githubBaseUrl     = "https://github.com"
)

type Repository struct {
	Name                string
	Url                 string
	Description         string
	Topics              []string
	Stars               string
	Licence             string
	ProgrammingLanguage string
	UpdateTime          string
	Error               error
}

type searchMode string

const (
	searchModeRepositories searchMode = "repositories"
	searchModeCode         searchMode = "code"
	searchModeCommits      searchMode = "commits"
	searchModeIssues       searchMode = "issues"
	searchModeDiscussions  searchMode = "discussions"
	searchModePackages     searchMode = "registrypackages"
	searchModeMarketPlace  searchMode = "marketPlace"
	searchModeTopics       searchMode = "topics"
	searchModeWikis        searchMode = "wikis"
	searchModeUsers        searchMode = "users"
)

func buildSearchUrl(query string, typee searchMode, opt sortOptions) string {

	switch typee {
	case searchModeMarketPlace, searchModeTopics:
		return githubBaseUrl + url.QueryEscape(fmt.Sprintf("/search?q=%s&type=%s", query, typee))
	default:
		return githubBaseUrl + url.QueryEscape(fmt.Sprintf("/search?o=%s&q=%s&s=%s&type=%s", opt.getOrder(), query, opt.name, typee))
	}
}

func successfulLoaded(res *http.Response) bool {
	if res.StatusCode != 200 {
		return false
	}
	return true
}

func getMaxPage(maxResult int) int {
	return int(math.Round(float64(maxResult) / float64(resultPerPage)))
}

// SearchRepositories returns channel with repositories for a given search query
func (s *Scraper) SearchRepositories(opt sortOptions, query string, maxResult int) <-chan *Repository {
	channel := make(chan *Repository)

	go func() {
		defer close(channel)
		url := buildSearchUrl(query, searchModeRepositories, opt)
		maxpage := getMaxPage(maxResult)
		for page := 1; page <= maxpage; page++ {
			res, err := s.client.Get(url + fmt.Sprintf("&p=%v", page))
			if err != nil {
				channel <- &Repository{Error: err}
				return
			}

			if !successfulLoaded(res) {
				channel <- &Repository{Error: errors.New("Failed to load page")}
				return
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)

			if err != nil {
				channel <- &Repository{Error: errors.New("Failed to read page")}
				return
			}

			doc.Find("li.repo-list-item > div.mt-n1").Each(func(i int, selection *goquery.Selection) {
				var repo Repository
				repo.Name, _ = selection.Find("a").Attr("href") //repo name
				repo.Url = githubBaseUrl + repo.Name
				repo.Description = selection.Find("p.mb-1").Text()                                 //description
				repo.Stars = selection.Find("div.mr-3 > a.muted-link").Text()                      //stars
				repo.ProgrammingLanguage = selection.Find("[itemprop=programmingLanguage]").Text() //programmingLanguage
				repo.UpdateTime, _ = selection.Find("relative-time").Attr("datetime")              //updatedtime

				selection.Find("a.topic-tag").Each(func(itag int, selectiontag *goquery.Selection) {
					repo.Topics = append(repo.Topics, selectiontag.Text()) // each tags topics
				})

				channel <- &repo
			})
			res.Body.Close()
		}
	}()
	return channel
}

// SearchRepositories wrapper for default Scraper
func SearchRepositories(opt sortOptions, query string, maxResult int) <-chan *Repository {
	return defaultScraper.SearchRepositories(opt, query, maxResult)
}