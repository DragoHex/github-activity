package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

func TestGetActivity(t *testing.T) {
	tests := []struct {
		name             string
		userName         string
		resp             string
		wantGithubEvents []GitHubEvent
		wantErr          error
	}{
		{
			name:     "CreateEvent",
			userName: "test_user1",
			resp: `[
			{
		      "type": "CreateEvent",
		      "repo": {
		          "name": "test_user1/test-repo",
		          "url": "https://api.github.com/repos/DragoHex/task-tracker"
		      }
			}
			]`,
			wantGithubEvents: []GitHubEvent{
				{
					Type: "CreateEvent",
					Repo: Repo{
						Name: "test_user1/test-repo",
						URL:  "https://api.github.com/repos/DragoHex/task-tracker",
					},
				},
			},
		},
		{
			name:             "return error if the user doesn't exists'",
			userName:         "test_user",
			wantGithubEvents: []GitHubEvent{},
			wantErr:          ErrUserNotFound,
		},
	}

	userList := []string{"test_user1", "test_user2", "test_user3"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock server
			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					for _, user := range userList {
						if tt.userName == user {
							fmt.Fprintln(w, tt.resp)
							return
						}
					}
					http.NotFound(w, r)
				}),
			)
			defer ts.Close()

			oldGitHubURL := GitHubEventsURL
			defer func() {
				GitHubEventsURL = oldGitHubURL
			}()

			GitHubEventsURL = fmt.Sprintf("%s/%%s/events", ts.URL)

			gitHubEvents := NewGitHubEvents(tt.userName, 10)
			err := gitHubEvents.GetActivity()
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
			} else if tt.wantErr != nil {
				t.Error("expecting error but didn't get any'")
			}
			assertGitHubEvents(t, tt.wantGithubEvents, gitHubEvents.Events)
		})
	}
}

func assertGitHubEvents(t testing.TB, exp, result []GitHubEvent) {
	t.Helper()
	fmt.Println("check event length")
	assert.Equal(t, len(exp), len(result))
	for i, event := range exp {
		assert.Equal(t, event.Type, result[i].Type)
		assert.Equal(t, event.Repo, result[i].Repo)
		assert.Equal(t, event.Payload, result[i].Payload)

		fmt.Println("check commits http.ResponseWriter, r *http.Request length")
		assert.Equal(t, len(event.Commits), len(result[i].Commits))
		for _, commit := range result[i].Commits {
			assert.Equal(t, event.Commits[i], commit)
		}
	}
}

func Test_processEvents(t *testing.T) {
	tests := []struct {
		name   string
		events []GitHubEvent
		want   string
	}{
		{
			name: "valid Create Event",
			events: []GitHubEvent{
				{
					Type: "CreateEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
				},
			},
			want: `User Activities:
- Created a new repo /testUser/testRepo
`,
		},
	}
	gitHubEvents := NewGitHubEvents("test_user", 10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitHubEvents.Events = tt.events
			got := gitHubEvents.ProcessEvents()
			assert.Equal(t, tt.want, got)
		})
	}
}
