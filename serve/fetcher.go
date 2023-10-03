package serve

import (
	"context"
	"fmt"
	"os"

	godotenv "github.com/joho/godotenv"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type CommitTarget struct {
    History struct {
        TotalCount githubv4.Int
    }
}

type Query struct {
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

func fetch(targetUser string) Query {

    var query Query

    godotenv.Load(".env")

    accessToken := os.Getenv("TOKEN")

    src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
    httpClient := oauth2.NewClient(context.Background(), src)
    client := githubv4.NewClient(httpClient)

	variables := map[string]interface{}{
		"username": githubv4.String(targetUser),
	}

	err := client.Query(context.Background(), &query, variables)

	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}

    return query
}
