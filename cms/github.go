package cms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type PinnedRepo struct {
	Name        string
	Description string
	URL         string
	Readme      string
	Stars       int
	Language    string
	UpdatedAt   string
}

func (c *Site) loadGitHubPinnedItems(ct ContentType) ([]ContentItem, error) {
	username := strings.TrimSpace(ct.GitHubUser)
	if username == "" {
		username = strings.TrimSpace(c.Site.GitHub)
	}
	if username == "" {
		return nil, fmt.Errorf("github username is empty for content type %q", ct.Name)
	}

	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is required to load pinned repositories")
	}

	repos, err := FetchPinnedRepos(username, token, 6)
	if err != nil {
		return nil, err
	}

	items := make([]ContentItem, 0, len(repos))
	for _, repo := range repos {
		slug := strings.ToLower(repo.Name)
		body := strings.TrimSpace(repo.Readme)
		if body == "" {
			body = buildRepoFallbackBody(repo)
		}

		excerpt := strings.TrimSpace(repo.Description)
		if excerpt == "" {
			excerpt = extractExcerpt(body, 200)
		}

		items = append(items, ContentItem{
			Title:   repo.Name,
			Slug:    slug,
			Path:    "/" + ct.Name + "/" + slug,
			Body:    body,
			Excerpt: excerpt,
			Metadata: map[string]string{
				"source":   "GitHub",
				"stars":    fmt.Sprintf("%d", repo.Stars),
				"language": repo.Language,
				"updated":  repo.UpdatedAt,
				"url":      repo.URL,
			},
			ContentDir: "",
		})
	}

	return items, nil
}

func (c *Site) loadGitHubPinnedItemBySlug(ct ContentType, slug string) (ContentItem, error) {
	items, err := c.loadGitHubPinnedItems(ct)
	if err != nil {
		return ContentItem{}, err
	}
	for _, item := range items {
		if strings.EqualFold(item.Slug, slug) {
			return item, nil
		}
	}
	return ContentItem{}, fmt.Errorf("pinned repository %q not found", slug)
}

func buildRepoFallbackBody(repo PinnedRepo) string {
	var sb strings.Builder
	sb.WriteString("# ")
	sb.WriteString(repo.Name)
	sb.WriteString("\n\n")
	if strings.TrimSpace(repo.Description) != "" {
		sb.WriteString(repo.Description)
		sb.WriteString("\n\n")
	}
	if strings.TrimSpace(repo.URL) != "" {
		sb.WriteString("[View on GitHub](")
		sb.WriteString(repo.URL)
		sb.WriteString(")\n")
	}
	return sb.String()
}

func FetchPinnedRepos(username, token string, limit int) ([]PinnedRepo, error) {
	if limit <= 0 {
		limit = 6
	}
	if limit > 20 {
		limit = 20
	}
	const query = `query($login:String!,$first:Int!){user(login:$login){pinnedItems(first:$first,types:[REPOSITORY]){nodes{... on Repository{name description url stargazerCount updatedAt primaryLanguage{name} object(expression:"HEAD:README.md"){... on Blob{text}}}}}}}`

	payload := map[string]any{
		"query": query,
		"variables": map[string]any{
			"login": username,
			"first": limit,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.github.com/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("github graphql returned %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var result struct {
		Data struct {
			User struct {
				PinnedItems struct {
					Nodes []struct {
						Name            string `json:"name"`
						Description     string `json:"description"`
						URL             string `json:"url"`
						StargazerCount  int    `json:"stargazerCount"`
						UpdatedAt       string `json:"updatedAt"`
						PrimaryLanguage *struct {
							Name string `json:"name"`
						} `json:"primaryLanguage"`
						Object *struct {
							Text string `json:"text"`
						} `json:"object"`
					} `json:"nodes"`
				} `json:"pinnedItems"`
			} `json:"user"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("github graphql error: %s", result.Errors[0].Message)
	}

	repos := make([]PinnedRepo, 0, len(result.Data.User.PinnedItems.Nodes))
	for _, n := range result.Data.User.PinnedItems.Nodes {
		if strings.TrimSpace(n.Name) == "" {
			continue
		}
		repo := PinnedRepo{
			Name:        n.Name,
			Description: n.Description,
			URL:         n.URL,
			Stars:       n.StargazerCount,
			UpdatedAt:   strings.TrimSpace(n.UpdatedAt),
		}
		if n.PrimaryLanguage != nil {
			repo.Language = strings.TrimSpace(n.PrimaryLanguage.Name)
		}
		if repo.Language == "" {
			repo.Language = "n/a"
		}
		if t, err := time.Parse(time.RFC3339, repo.UpdatedAt); err == nil {
			repo.UpdatedAt = t.Format("2006-01-02")
		}
		if repo.UpdatedAt == "" {
			repo.UpdatedAt = "n/a"
		}
		if n.Object != nil {
			repo.Readme = n.Object.Text
		}
		repos = append(repos, repo)
	}

	return repos, nil
}
