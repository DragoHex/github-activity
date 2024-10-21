package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrUserNotFound = errors.New("error user not found")
	GitHubEventsURL = "https://api.github.com/users/%s/events"
)

type GitHubEvents []GitHubEvent

type GitHubEvent struct {
	Type    string   `json:"type,omitempty"`
	Repo    Repo     `json:"repo,omitempty"`
	Payload Payload  `json:"payload,omitempty"`
	Commits []Commit `json:"commits,omitempty"`
}

type Repo struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Payload struct {
	Action string `json:"action,omitempty"`
	Title  string `json:"title,omitempty"`
}

type Commit struct {
	SHA     string `json:"sha,omitempty"`
	Message string `json:"message,omitempty"`
}

func GetActivity(userName string) (string, error) {
	githubURL := fmt.Sprintf(GitHubEventsURL, userName)
	resp, err := http.Get(githubURL)
	if err != nil {
		return "", fmt.Errorf("error in fetching activity: %s", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", ErrUserNotFound
	}

	var data []byte
	_, err = resp.Body.Read(data)
	if err != nil {
		return "", fmt.Errorf("error in reading the response: %s", err)
	}

	var events GitHubEvents

	err = json.Unmarshal(data, &events)
	if err != nil {
		return "", fmt.Errorf("error in unmarshalling: %s", err)
	}

	result := ProcessEvents(events)
	return result, nil
}

func ProcessEvents(events GitHubEvents) string {
	activity := "User Activities:\n"

	for i, event := range events {
		fmt.Println("event type: ", event.Type)
		fmt.Println("Event Type: ", Event(1).String())
		switch event.Type {
		case Event(1).String():
		case Event(2).String():
			act := fmt.Sprintf("- Created a new repo %s\n", event.Repo.Name)
			activity = activity + act
		case Event(3).String():
		case Event(4).String():
		case Event(5).String():
		case Event(6).String():
		case Event(7).String():
		case Event(8).String():
		case Event(9).String():
		case Event(10).String():
		case Event(11).String():
		case Event(12).String():
		case Event(13).String():
		case Event(14).String():
		case Event(15).String():
		case Event(16).String():
		case Event(17).String():
		}
		if i == 10 {
			break
		}
	}

	return activity
}
