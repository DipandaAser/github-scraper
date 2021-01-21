package githubscraper

type sortOptions struct {
	name  string
	order int
}

var (
	DefaultSortOption              = sortOptions{name: "", order: -1}
	RepositoryMostStars            = sortOptions{name: "stars", order: 1}
	RepositoryFewestStars          = sortOptions{name: "stars", order: 0}
	RepositoryMostForks            = sortOptions{name: "fork", order: 1}
	RepositoryFewestForks          = sortOptions{name: "fork", order: 0}
	RepositoryRecentlyUpdated      = sortOptions{name: "updated", order: 1}
	RepositoryLeastRecentlyUpdated = sortOptions{name: "updated", order: 0}
	CommitRecentlyCommitted        = sortOptions{name: "committer-date", order: 1}
	CommitLeastRecentlyCommitted   = sortOptions{name: "committer-date", order: 0}
	CommitRecentlyAuthored         = sortOptions{name: "author-date", order: 1}
	CommitLeastRecentlyAuthored    = sortOptions{name: "author-date", order: 0}
)

func (so *sortOptions) getOrder() string {
	switch so.order {
	case 0:
		return "desc"
	case 1:
		return "asc"
	default:
		return ""
	}
}
