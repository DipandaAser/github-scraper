package githubscraper

import (
	"testing"
)

func TestScraper_SearchRepositories(t *testing.T) {
	t.Run("Test SearchRepositories", func(t *testing.T) {
		s := New()
		repoCount := 0
		for _ = range s.SearchRepositories(DefaultSortOption, "go", 20) {
			repoCount++
		}
		if repoCount < 20 {
			t.Fatal("SearchRepositories faild. Can't get 20 items.")
		}
	})
}