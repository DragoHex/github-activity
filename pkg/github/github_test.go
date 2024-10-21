package github

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestGetActivity(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		resp     string
		wantErr  error
	}{
		{
			name:     "return error if the user doesn't exists'",
			userName: "56984654*",
			wantErr:  ErrUserNotFound,
		},
		{
			name:     "return error if the user doesn't exists'",
			userName: "DragoHex",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetActivity(tt.userName)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
			} else if tt.wantErr != nil {
				t.Error("expecting error but didn't get any'")
			}

			// TODO: Remove once completed
			fmt.Println("got: ", got)
			fmt.Println("err: ", err)

			assert.Equal(t, tt.resp, got)
		})
	}
}

func TestProcessEvents(t *testing.T) {
	tests := []struct {
		name   string
		events GitHubEvents
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessEvents(tt.events)
			fmt.Println("got: ", got)
			assert.Equal(t, tt.want, got)
		})
	}
}
