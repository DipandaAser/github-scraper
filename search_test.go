package githubscraper

import (
	"testing"
)

func TestScraper_SearchRepositories(t *testing.T) {
	t.Run("Test SearchRepositories", func(t *testing.T) {
		s := New()
		count := 0
		for _ = range s.SearchRepositories(DefaultSortOption, "go", 20) {
			count++
		}
		if count < 20 {
			t.Fatal("SearchRepositories faild. Can't get 20 items.")
		}
	})
}

func TestScraper_SearchCommits(t *testing.T) {
	t.Run("Test SearchCommits", func(t *testing.T) {
		s := New()
		count := 0
		for _ = range s.SearchCommits(DefaultSortOption, "go", 20) {
			count++
		}
		if count < 20 {
			t.Fatal("SearchCommits faild. Can't get 20 items.")
		}
	})
}

func TestScraper_SearchIssues(t *testing.T) {
	t.Run("Test SearchIssues", func(t *testing.T) {
		s := New()
		count := 0
		for _ = range s.SearchIssues(DefaultSortOption, "go", 20) {
			count++
		}
		if count < 20 {
			t.Fatal("SearchIssues faild. Can't get 20 items.")
		}
	})
}
