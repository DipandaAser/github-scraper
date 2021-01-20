# Github Scraper
Github API is annoying to work with for a basic task like repo search you wont need to request a tokens. 

This library allows you to save time, and quickly start.
No tokens needed. No restrictions. Extremely fast.

# Installation
```
go get -u github.com/DipandaAser/github-scraper
```

# Usage
Let's start with a trivial example

```go
package main

import (
	"fmt"
	githubscraper "github.com/DipandaAser/github-scraper"
)

func main() {
	s := githubscraper.New()
	for repo := range s.SearchRepositories(githubscraper.DefaultSortOption, "go", 20) {
		fmt.Println(repo.Name)
	}
}
```

## Default Scraper (Ad hoc)

In simple cases, you can use the default scraper without creating an object instance
```go
package main

import (
	"fmt"
	githubscraper "github.com/DipandaAser/github-scraper"
)

func main() {
	for repo := range githubscraper.SearchRepositories(githubscraper.DefaultSortOption, "go", 20) {
		fmt.Println(repo.Name)
	}
}
```