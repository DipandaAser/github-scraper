package githubscraper

type sortOptions struct {
	name  string
	order string
}

var (
	DefaultSortOption              = sortOptions{name: "", order: ""}
	RepositoryMostStars            = sortOptions{name: "stars", order: "asc"}
	RepositoryFewestStars          = sortOptions{name: "stars", order: "desc"}
	RepositoryMostForks            = sortOptions{name: "fork", order: "asc"}
	RepositoryFewestForks          = sortOptions{name: "fork", order: "desc"}
	RepositoryRecentlyUpdated      = sortOptions{name: "updated", order: "asc"}
	RepositoryLeastRecentlyUpdated = sortOptions{name: "updated", order: "desc"}
	CommitRecentlyCommitted        = sortOptions{name: "committer-date", order: "asc"}
	CommitLeastRecentlyCommitted   = sortOptions{name: "committer-date", order: "desc"}
	CommitRecentlyAuthored         = sortOptions{name: "author-date", order: "asc"}
	CommitLeastRecentlyAuthored    = sortOptions{name: "author-date", order: "desc"}
	IssuesMostCommented            = sortOptions{name: "comments", order: "desc"}
	IssuesNewest                   = sortOptions{name: "created", order: "desc"}
	IssuesOldest                   = sortOptions{name: "created", order: "asc"}
	IssuesRecentlyUpdated          = sortOptions{name: "updated", order: "desc"}
	IssuesLeastRecentlyUpdated     = sortOptions{name: "updated", order: "asc"}
)
