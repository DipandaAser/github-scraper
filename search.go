package githubscraper

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

const (
	resultPerPage int = 10
	githubBaseUrl     = "https://github.com"
)

type Repository struct {
	Name                string   `json:"name"`
	Url                 string   `json:"url"`
	Description         string   `json:"description"`
	Topics              []string `json:"topics"`
	Stars               string   `json:"stars"`
	Licence             string   `json:"licence"`
	ProgrammingLanguage string   `json:"programming_language"`
	UpdateTime          string   `json:"update_time"`
	Error               error    `json:"error"`
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

type User struct {
	Name     string
	Pseudo   string
	Bio      string
	Avatar   string
	Link     string
	Location string
	Error    error
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

const githubMaxPageResult = 100

func buildSearchUrl(query string, typee searchMode, opt sortOptions) string {

	switch typee {
	case searchModeMarketPlace, searchModeTopics:
		return githubBaseUrl + fmt.Sprintf("/search?q=%s&type=%s", query, typee)
	default:
		return githubBaseUrl + fmt.Sprintf("/search?o=%s&q=%s&s=%s&type=%s", opt.order, query, opt.name, typee)
	}
}

func successfulLoaded(res *http.Response) bool {
	if res.StatusCode != 200 {
		return false
	}
	return true
}

// SearchRepositories returns channel with Repository for a given search query
func (s *Scraper) SearchRepositories(opt sortOptions, query string, maxResult int) <-chan *Repository {
	channel := make(chan *Repository)

	go func() {
		defer close(channel)
		url := buildSearchUrl(query, searchModeRepositories, opt)
		var resultCount int
		for page := 1; page <= githubMaxPageResult; page++ {
			if resultCount == maxResult {
				channel <- &Repository{Error: errors.New("maxResult")}
				return
			}
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
				if resultCount == maxResult {
					return
				}
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
				resultCount++
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
		var resultCount int
		for page := 1; page <= githubMaxPageResult; page++ {
			if resultCount == maxResult {
				channel <- &Commit{Error: errors.New("maxResult")}
				return
			}
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
				if resultCount == maxResult {
					return
				}
				var commit Commit
				commit.RepositoryName, _ = selection.Find("a.link-gray").Attr("href")
				commit.RepositoryLink = githubBaseUrl + commit.RepositoryName
				commit.CommitMessage = selection.Find("div.f4 > a.message").Text()
				clink, _ := selection.Find("div.f4 > a.message").Attr("href")
				commit.CommitLink = githubBaseUrl + clink
				commit.Author = selection.Find("a.commit-author").Text()
				commit.CommitDate, _ = selection.Find("relative-time").Attr("datetime")

				resultCount++
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
		var resultCount int
		for page := 1; page <= githubMaxPageResult; page++ {
			if resultCount == maxResult {
				channel <- &Issue{Error: errors.New("maxResult")}
				return
			}
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
				if resultCount == maxResult {
					return
				}
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

				resultCount++
				channel <- &issue
			})
			res.Body.Close()
		}
	}()
	return channel
}

// SearchIssues wrapper for default Scraper
func SearchIssues(opt sortOptions, query string, maxResult int) <-chan *Issue {
	return defaultScraper.SearchIssues(opt, query, maxResult)
}

// SearchUsers returns channel with User for a given search query
func (s *Scraper) SearchUsers(opt sortOptions, query string, maxResult int) <-chan *User {
	channel := make(chan *User)

	go func() {
		defer close(channel)
		urll := buildSearchUrl(query, searchModeUsers, opt)
		var resultCount int
		for page := 1; page <= githubMaxPageResult; page++ {
			if resultCount == maxResult {
				channel <- &User{Error: errors.New("maxResult")}
				return
			}
			res, err := s.client.Get(urll + fmt.Sprintf("&p=%v", page))
			if err != nil {
				channel <- &User{Error: err}
				return
			}

			if !successfulLoaded(res) {
				channel <- &User{Error: errors.New("Failed to load page")}
				return
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				channel <- &User{Error: errors.New("Failed to read page")}
				return
			}

			doc.Find("div.user-list-item").Each(func(i int, selection *goquery.Selection) {
				if resultCount == maxResult {
					return
				}
				// if true is means that this selection is an organisation
				if selection.Find("span.user-following-container").Text() == "" {
					return
				}

				var user User
				nameSelection := selection.Find("div.f4 > a.mr-1")
				user.Name = nameSelection.Text()
				user.Pseudo, _ = nameSelection.Attr("href")
				user.Link = githubBaseUrl + user.Pseudo
				user.Pseudo = strings.Replace(user.Pseudo, "/", "", -1)

				user.Bio = selection.Find("p.mb-1").Text()
				user.Location = selection.Find("div.mr-3").Text()

				user.Avatar, _ = selection.Find("img.avatar-user").Attr("src")
				if user.Avatar != "" {
					u, _ := url.Parse(user.Avatar)
					values, _ := url.ParseQuery(u.RawQuery)
					values.Set("s", "500")
					u.RawQuery = values.Encode()
					user.Avatar = u.String()
				}

				resultCount++
				channel <- &user
			})
			res.Body.Close()
		}
	}()
	return channel
}

// SearchUsers wrapper for default Scraper
func SearchUsers(opt sortOptions, query string, maxResult int) <-chan *User {
	return defaultScraper.SearchUsers(opt, query, maxResult)
}
