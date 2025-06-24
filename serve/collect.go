package serve

import (
	"sort"
)

type LanguageCommit struct {
    Language string
    Commits  int
}

type Stat struct {
    Language       string
    CommitsCount   int
    CommitsPercent float64
}

func collect(query Query) []Stat {
    commitCounts := make(map[string]int)
	totalCommits := 0

	for _, edge := range query.User.Repositories.Edges {
		commitCount := edge.Node.DefaultBranchRef.Target.History.TotalCount
		language := edge.Node.PrimaryLanguage.Name

		// Skip repositories with an empty string for the programming language.
		if language != "" {
			// Update the commit count for the language.
			commitCounts[language] += int(commitCount)
			totalCommits += int(commitCount)
		}
	}

    // Collect language-commit pairs in a slice for sorting.
    var languageCommits []LanguageCommit
    for language, count := range commitCounts {
        languageCommits = append(languageCommits, LanguageCommit{Language: language, Commits: count})
    }

    // Sort the language-commit pairs by commit count in descending order.
    sort.SliceStable(languageCommits, func(i, j int) bool {
        return languageCommits[i].Commits > languageCommits[j].Commits
    })

	var stats []Stat
	topN := 7
	othersCommits := 0

	for i, lc := range languageCommits {
		if i < topN {
			percent := float64(lc.Commits) / float64(totalCommits) * 100
			stats = append(stats, Stat{
				Language:       lc.Language,
				CommitsCount:   lc.Commits,
				CommitsPercent: percent,
			})
		} else {
			othersCommits += lc.Commits
		}
	}

	if othersCommits > 0 {
		percent := float64(othersCommits) / float64(totalCommits) * 100
		stats = append(stats, Stat{
			Language:       "Others",
			CommitsCount:   othersCommits,
			CommitsPercent: percent,
		})
	}

	return stats
}
