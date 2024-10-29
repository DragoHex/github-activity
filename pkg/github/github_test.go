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
		assert.Equal(t, event.Payload.Action, result[i].Payload.Action)
		assert.Equal(t, event.Payload.Issue, result[i].Payload.Issue)

		fmt.Println("check commits http.ResponseWriter, r *http.Request length")
		assert.Equal(t, len(event.Payload.Commits), len(result[i].Payload.Commits))
		for _, commit := range result[i].Payload.Commits {
			assert.Equal(t, event.Payload.Commits[i], commit)
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
		{
			name: "valid Push Event",
			events: []GitHubEvent{
				{
					Type: "PushEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
					Payload: Payload{
						Commits: []Commit{},
					},
				},
			},
			want: `User Activities:
- Pushed 0 commits to /testUser/testRepo
`,
		},
		{
			name: "valid Pull Event",
			events: []GitHubEvent{
				{
					Type: "PullRequestEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
				},
			},
			want: `User Activities:
- Created a Pull Request to /testUser/testRepo
`,
		},
		{
			name: "valid Issue Event",
			events: []GitHubEvent{
				{
					Type: "IssuesEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
					Payload: Payload{
						Action: "closed",
						Issue: Issue{
							Title: "Test Issue",
						},
					},
				},
			},
			want: `User Activities:
- closed "Test Issue" issue for repo /testUser/testRepo
`,
		},
		{
			name: "valid Commit Comment Event",
			events: []GitHubEvent{
				{
					Type: "CommitCommentEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
				},
			},
			want: `User Activities:
- Left commit comment for repo /testUser/testRepo
`,
		},
		{
			name: "valide Fork Event",
			events: []GitHubEvent{
				{
					Type: "ForkEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
				},
			},
			want: `User Activities:
- Forked repo /testUser/testRepo
`,
		},
		{
			name: "valid Issue Comment Event",
			events: []GitHubEvent{
				{
					Type: "IssueCommentEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
					Payload: Payload{
						Action: "closed",
						Issue: Issue{
							Title: "Test Issue",
						},
					},
				},
			},
			want: `User Activities:
- Commented on issue "Test Issue" for repo /testUser/testRepo
`,
		},
		{
			name: "valid PR review event",
			events: []GitHubEvent{
				{
					Type: "PullRequestReviewEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
				},
			},
			want: `User Activities:
- Reviewed a PR for repo /testUser/testRepo
`,
		},
		{
			name: "valid PR review comment event",
			events: []GitHubEvent{
				{
					Type: "PullRequestReviewCommentEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
				},
			},
			want: `User Activities:
- Left a PR review comment for repo /testUser/testRepo
`,
		},
		{
			name: "valid PR review thread event",
			events: []GitHubEvent{
				{
					Type: "PullRequestReviewThreadEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
				},
			},
			want: `User Activities:
- Created a PR review thread for repo /testUser/testRepo
`,
		},
		{
			name: "valid release event",
			events: []GitHubEvent{
				{
					Type: "ReleaseEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
				},
			},
			want: `User Activities:
- Created a release for repo /testUser/testRepo
`,
		},
		{
			name: "valid delete event",
			events: []GitHubEvent{
				{
					Type: "DeleteEvent",
					Repo: Repo{
						Name: "/testUser/testRepo",
					},
					Payload: Payload{
						Action: "closed",
						Issue: Issue{
							Title: "Test Issue",
						},
						RefType: "branch",
					},
				},
			},
			want: `User Activities:
- Deleted branch for repo /testUser/testRepo
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
