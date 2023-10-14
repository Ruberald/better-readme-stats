package serve

import (
	"fmt"
	"sort"
)

type LanguageCommit struct {
	Language string
	Commits  int
}

func collect(query Query) map[string]interface{} {
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

    response := make(map[string]interface{})

    stats := []map[string]interface{}{}
    for _, lc := range languageCommits {
        percentage := float64(lc.Commits) / float64(totalCommits) * 100
        languageStats := map[string]interface{}{
            "language":        lc.Language,
            "commits_count":   lc.Commits,
            "commits_percent": fmt.Sprintf("%.2f", percentage),
        }
        stats = append(stats, languageStats)
    }

    response["stats"] = stats

    return response
}
