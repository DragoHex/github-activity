package github

type Event int

const (
	_ Event = iota
	CommitCommentEvent
	CreateEvent
	DeleteEvent
	ForkEvent
	IssueCommentEvent
	IssuesEvent
	PullRequestEvent
	PullRequestReviewEvent
	PullRequestReviewCommentEvent
	PullRequestReviewThreadEvent
	PushEvent
	ReleaseEvent
)

func (e Event) String() string {
	return [...]string{
		"CommitCommentEvent",
		"CreateEvent",
		"DeleteEvent",
		"ForkEvent",
		"IssueCommentEvent",
		"IssuesEvent",
		"PullRequestEvent",
		"PullRequestReviewEvent",
		"PullRequestReviewCommentEvent",
		"PullRequestReviewThreadEvent",
		"PushEvent",
		"ReleaseEvent",
	}[e-1]
}

func (e Event) EnumIndex() int {
	return int(e)
}
