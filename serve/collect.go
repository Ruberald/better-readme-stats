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

	topN := 7
	othersCommits := 0

	for i, lc := range languageCommits {
		if i < topN {
			percentage := float64(lc.Commits) / float64(totalCommits) * 100
			stats = append(stats, map[string]interface{}{
				"language":        lc.Language,
				"commits_count":   lc.Commits,
				"commits_percent": fmt.Sprintf("%.2f", percentage),
			})
		} else {
			othersCommits += lc.Commits
		}
	}

	if othersCommits > 0 {
		percentage := float64(othersCommits) / float64(totalCommits) * 100
		stats = append(stats, map[string]interface{}{
			"language":        "Others",
			"commits_count":   othersCommits,
			"commits_percent": fmt.Sprintf("%.2f", percentage),
		})
	}

	response["stats"] = stats
	return response
}
