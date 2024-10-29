package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrUserNotFound = errors.New("error user not found")
	GitHubEventsURL = "https://api.github.com/users/%s/events"
)

type GitHubEvents struct {
	Events []GitHubEvent
	user   string
	limit  int
}

type GitHubEvent struct {
	Type    string  `json:"type,omitempty"`
	Repo    Repo    `json:"repo,omitempty"`
	Payload Payload `json:"payload,omitempty"`
}

type Repo struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Payload struct {
	Action  string   `json:"action,omitempty"`
	Issue   Issue    `json:"issue,omitempty"`
	RefType string   `json:"ref_type,omitempty"`
	Commits []Commit `json:"commits,omitempty"`
}

type Issue struct {
	Title string `json:"title,omitempty"`
	State string `json:"state,omitempty"`
}

type Commit struct {
	SHA     string `json:"sha,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewGitHubEvents(user string, limit int) *GitHubEvents {
	return &GitHubEvents{
		user:  user,
		limit: limit,
	}
}

// GetActivity makes API call to fetch events from GitHub
func (g *GitHubEvents) GetActivity() error {
	githubURL := fmt.Sprintf(GitHubEventsURL, g.user)
	resp, err := http.Get(githubURL)
	if err != nil {
		return fmt.Errorf("error in fetching activity: %s", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return ErrUserNotFound
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("call to github api failed: %s\n", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error in reading the response: %s", err)
	}

	err = json.Unmarshal(data, &g.Events)
	if err != nil {
		return fmt.Errorf("error in unmarshalling: %s", err)
	}

	return nil
}

// ProcessEvents process the events data to create activity response
func (g *GitHubEvents) ProcessEvents() string {
	activity := "User Activities:\n"

	for i, event := range g.Events {
		switch event.Type {
		case Event(1).String():
			act := fmt.Sprintf("- Left commit comment for repo %s\n", event.Repo.Name)
			activity = activity + act
		case Event(2).String():
			act := fmt.Sprintf("- Created a new repo %s\n", event.Repo.Name)
			activity = activity + act
		case Event(3).String():
			act := fmt.Sprintf("- Deleted %s for repo %s\n", event.Payload.RefType, event.Repo.Name)
			activity = activity + act
		case Event(4).String():
			act := fmt.Sprintf("- Forked repo %s\n", event.Repo.Name)
			activity = activity + act
		case Event(5).String():
			act := fmt.Sprintf(
				"- Commented on issue %q for repo %s\n",
				event.Payload.Issue.Title,
				event.Repo.Name,
			)
			activity = activity + act
		case Event(6).String():
			act := fmt.Sprintf(
				"- %s %q issue for repo %s\n",
				event.Payload.Action,
				event.Payload.Issue.Title,
				event.Repo.Name,
			)
			activity = activity + act
		case Event(7).String():
			act := fmt.Sprintf("- Created a Pull Request to %s\n", event.Repo.Name)
			activity = activity + act
		case Event(8).String():
			act := fmt.Sprintf("- Reviewed a PR for repo %s\n", event.Repo.Name)
			activity = activity + act
		case Event(9).String():
			act := fmt.Sprintf("- Left a PR review comment for repo %s\n", event.Repo.Name)
			activity = activity + act
		case Event(10).String():
			act := fmt.Sprintf("- Created a PR review thread for repo %s\n", event.Repo.Name)
			activity = activity + act
		case Event(11).String():
			commitCount := len(event.Payload.Commits)
			act := fmt.Sprintf("- Pushed %d commits to %s\n", commitCount, event.Repo.Name)
			activity = activity + act
		case Event(12).String():
			act := fmt.Sprintf("- Created a release for repo %s\n", event.Repo.Name)
			activity = activity + act
		}
		if i == g.limit {
			break
		}
	}

	return activity
}
