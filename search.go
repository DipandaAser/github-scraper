package githubscraper

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math"
	"net/http"
	"strings"
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

type Commit struct {
	RepositoryName string
	RepositoryLink string
	CommitLink     string
	CommitMessage  string
	Author         string
	CommitDate     string
	Error          error
}

type Issue struct {
	RepositoryName string
	RepositoryLink string
	Link           string
	Title          string
	Description    string
	Status         string
	IsPullRequest  bool
	Author         string
	Date           string
	Error          error
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
		return githubBaseUrl + fmt.Sprintf("/search?q=%s&type=%s", query, typee)
	default:
		return githubBaseUrl + fmt.Sprintf("/search?o=%s&q=%s&s=%s&type=%s", opt.getOrder(), query, opt.name, typee)
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

// SearchRepositories returns channel with Repository for a given search query
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

// SearchCommits returns channel with Commit for a given search query
func (s *Scraper) SearchCommits(opt sortOptions, query string, maxResult int) <-chan *Commit {
	channel := make(chan *Commit)

	go func() {
		defer close(channel)
		url := buildSearchUrl(query, searchModeCommits, opt)
		maxpage := getMaxPage(maxResult)
		for page := 1; page <= maxpage; page++ {
			res, err := s.client.Get(url + fmt.Sprintf("&p=%v", page))
			if err != nil {
				channel <- &Commit{Error: err}
				return
			}

			if !successfulLoaded(res) {
				channel <- &Commit{Error: errors.New("Failed to load page")}
				return
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				channel <- &Commit{Error: errors.New("Failed to read page")}
				return
			}

			doc.Find("div.commits-list-item > div.mt-n1").Each(func(i int, selection *goquery.Selection) {
				var commit Commit
				commit.RepositoryName, _ = selection.Find("a.link-gray").Attr("href")
				commit.RepositoryLink = githubBaseUrl + commit.RepositoryName
				commit.CommitMessage = selection.Find("div.f4 > a.message").Text()
				clink, _ := selection.Find("div.f4 > a.message").Attr("href")
				commit.CommitLink = githubBaseUrl + clink
				commit.Author = selection.Find("a.commit-author").Text()
				commit.CommitDate, _ = selection.Find("relative-time").Attr("datetime")

				channel <- &commit
			})
			res.Body.Close()
		}
	}()
	return channel
}

//SearchCommits wrapper for default Scraper
func SearchCommits(opt sortOptions, query string, maxResult int) <-chan *Commit {
	return defaultScraper.SearchCommits(opt, query, maxResult)
}

// SearchIssues returns channel with Issue for a given search query
func (s *Scraper) SearchIssues(opt sortOptions, query string, maxResult int) <-chan *Issue {
	channel := make(chan *Issue)

	go func() {
		defer close(channel)
		url := buildSearchUrl(query, searchModeIssues, opt)
		maxpage := getMaxPage(maxResult)
		for page := 1; page <= maxpage; page++ {
			res, err := s.client.Get(url + fmt.Sprintf("&p=%v", page))
			if err != nil {
				channel <- &Issue{Error: err}
				return
			}

			if !successfulLoaded(res) {
				channel <- &Issue{Error: errors.New("Failed to load page")}
				return
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				channel <- &Issue{Error: errors.New("Failed to read page")}
				return
			}

			doc.Find("div.issue-list-item").Each(func(i int, selection *goquery.Selection) {
				var issue Issue
				issue.RepositoryName, _ = selection.Find("a.muted-link.text-bold").Attr("data-hovercard-url")
				issue.RepositoryLink = githubBaseUrl + issue.RepositoryName
				iLink, _ := selection.Find("div.f4 > a").Attr("href")
				issue.Link = githubBaseUrl + iLink
				issue.Title, _ = selection.Find("div.f4 > a").Attr("title")
				issue.Description = selection.Find("p.mb-0").Text()

				iconClasses, _ := selection.Find("svg.octicon").Attr("class")
				classes := strings.Split(iconClasses, " ")
				if len(classes) >= 3 {
					issue.Status = classes[2]
				}

				if strings.Contains(iconClasses, "pull-request") {
					issue.IsPullRequest = true
				}

				issue.Author = selection.Find("div.d-inline > a.text-bold.muted-link").Text()
				issue.Date, _ = selection.Find("relative-time").Attr("datetime")

				channel <- &issue
			})
			res.Body.Close()
		}
	}()
	return channel
}
