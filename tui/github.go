package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type contribsMsg struct{ counts map[string]int } // date "2006-01-02" → count

// fetchContribs fetches the GitHub contribution calendar via the GraphQL API.
// Reads the token from the GITHUB_TOKEN environment variable.
func fetchContribs(username string) tea.Cmd {
	return func() tea.Msg {
		token := os.Getenv("GITHUB_TOKEN")
		if username == "" || token == "" {
			return contribsMsg{}
		}
		return fetchContribsGraphQL(username, token)
	}
}

const graphqlQuery = `{"query":"query($login:String!){user(login:$login){contributionsCollection{contributionCalendar{weeks{contributionDays{date contributionCount}}}}}}","variables":{"login":%q}}`

func fetchContribsGraphQL(username, token string) tea.Msg {
	body := bytes.NewBufferString(fmt.Sprintf(graphqlQuery, username))
	req, err := http.NewRequest(http.MethodPost, "https://api.github.com/graphql", body)
	if err != nil {
		return contribsMsg{}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return contribsMsg{}
	}
	defer resp.Body.Close()
	var result struct {
		Data struct {
			User struct {
				ContributionsCollection struct {
					ContributionCalendar struct {
						Weeks []struct {
							ContributionDays []struct {
								Date              string `json:"date"`
								ContributionCount int    `json:"contributionCount"`
							} `json:"contributionDays"`
						} `json:"weeks"`
					} `json:"contributionCalendar"`
				} `json:"contributionsCollection"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return contribsMsg{}
	}
	counts := make(map[string]int)
	for _, week := range result.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		for _, day := range week.ContributionDays {
			counts[day.Date] = day.ContributionCount
		}
	}
	return contribsMsg{counts: counts}
}
