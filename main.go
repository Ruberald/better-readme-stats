package main

import (
    "context"
    "fmt"
    "os"
    "sort"

    godotenv "github.com/joho/godotenv"
    "github.com/shurcooL/githubv4"
    "golang.org/x/oauth2"
)

type LanguageCommit struct {
	Language string
	Commits  int
}

func main() {
    godotenv.Load(".env")

    accessToken := os.Getenv("TOKEN")

    // Create a GitHub GraphQL client with your personal access token.
    src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
    httpClient := oauth2.NewClient(context.Background(), src)
    client := githubv4.NewClient(httpClient)

    // Replace with the GitHub username you want to query.
	targetUser := "Ruberald"

    type CommitTarget struct {
        History struct {
            TotalCount githubv4.Int
        }
    }

	// Send the GraphQL request using the query string and variables.
	var response struct {
		User struct {
			Repositories struct {
				Edges []struct {
					Node struct {
						Name           string
						PrimaryLanguage struct {
							Name string
						}
						DefaultBranchRef struct {
							Target struct {
                                CommitTarget `graphql:"... on Commit"`
							}
						}
					}
				}
			} `graphql:"repositories(first: 100)"`
		} `graphql:"user(login: $username)"`
	}

	variables := map[string]interface{}{
		"username": githubv4.String(targetUser),
	}

	err := client.Query(context.Background(), &response, variables)

	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}

    commitCounts := make(map[string]int)
    // Process the response and calculate commit counts for each non-empty language.
	totalCommits := 0

	for _, edge := range response.User.Repositories.Edges {
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

    // Display the sorted commit counts and percentages for each language.
    fmt.Println("Commit Counts by Language (Descending Order):")
    for _, lc := range languageCommits {
        percentage := float64(lc.Commits) / float64(totalCommits) * 100
        fmt.Printf("%s =>\n", lc.Language)
        fmt.Printf("Number of total commits in all repos: %d\n", lc.Commits)
        fmt.Printf("Percentage of total commits: %.2f%%\n", percentage)
        fmt.Println()
    }
}
