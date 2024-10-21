package github

type Event int

const (
	_ Event = iota
	CommitCommentEvent
	CreateEvent
	DeleteEvent
	ForkEvent
	GollumEvent
	IssueCommentEvent
	IssuesEvent
	MemberEvent
	PublicEvent
	PullRequestEvent
	PullRequestReviewEvent
	PullRequestReviewCommentEvent
	PullRequestReviewThreadEvent
	PushEvent
	ReleaseEvent
	SponsorshipEvent
	WatchEvent
)

func (e Event) String() string {
	return [...]string{
		"CommitCommentEvent",
		"CreateEvent",
		"DeleteEvent",
		"ForkEvent",
		"GollumEvent",
		"IssueCommentEvent",
		"IssuesEvent",
		"MemberEvent",
		"PublicEvent",
		"PullRequestEvent",
		"PullRequestReviewEvent",
		"PullRequestReviewCommentEvent",
		"PullRequestReviewThreadEvent",
		"PushEvent",
		"ReleaseEvent",
		"SponsorshipEvent",
		"WatchEvent",
	}[e-1]
}

func (e Event) EnumIndex() int {
	return int(e)
}
