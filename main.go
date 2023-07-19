package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type Alert struct {
	Identifier string `json:"identifier"`
	Severity   string `json:"severity"`
	Summary    string `json:"summary"`
	URL        string `json:"url"`
}

func main() {
	userName, err := getUserName(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	client, err := api.DefaultGraphQLClient()
	if err != nil {
		log.Fatal(err)
	}

	var query struct {
		Search struct {
			PageInfo struct {
				EndCursor *graphql.String
			}
			Nodes []struct {
				Repository struct {
					Name                string
					VulnerabilityAlerts struct {
						Nodes []struct {
							Number           int
							SecurityAdvisory struct {
								Summary string
							}
							State                 string
							SecurityVulnerability struct {
								Severity string
								Package  struct {
									Name string
								}
							}
						}
					} `graphql:"vulnerabilityAlerts(first:100, states:OPEN)"`
				} `graphql:"... on Repository"`
			}
		} `graphql:"search(first: 100, type: REPOSITORY, query: $query, after: $after)"`
	}

	var after *graphql.String
	var alerts []Alert
	for {
		variables := map[string]interface{}{
			"query": graphql.String(fmt.Sprintf("user:%s", userName)),
			"after": after,
		}
		if err := client.Query("search", &query, variables); err != nil {
			log.Fatal(err)
		}
		for _, node := range query.Search.Nodes {
			if len(node.Repository.VulnerabilityAlerts.Nodes) == 0 {
				continue
			}
			for _, alert := range node.Repository.VulnerabilityAlerts.Nodes {
				alerts = append(alerts, Alert{
					Identifier: fmt.Sprintf("%s#%d", node.Repository.Name, alert.Number),
					Severity:   alert.SecurityVulnerability.Severity,
					Summary:    alert.SecurityAdvisory.Summary,
					URL:        fmt.Sprintf("https://github.com/%s/%s/security/dependabot/%d", userName, node.Repository.Name, alert.Number),
				})
			}
		}
		if query.Search.PageInfo.EndCursor == nil {
			break
		}
		after = query.Search.PageInfo.EndCursor
	}
	if err := json.NewEncoder(os.Stdout).Encode(alerts); err != nil {
		log.Fatal(err)
	}
}

func getUserName(ctx context.Context) (string, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return "", err
	}
	response := struct{ Login string }{}
	err = client.Get("user", &response)
	if err != nil {
		return "", err
	}
	return response.Login, nil
}
