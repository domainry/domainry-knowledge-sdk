package knowledgesdk

import (
	"context"
	"time"
)

const (
	KnowledgePollStatusOpen   = "open"
	KnowledgePollStatusClosed = "closed"
)

// KnowledgePagePoll is live business state referenced by one stable document
// Block. Ballots are stored independently so voting never rewrites the page
// draft or an immutable published version.
type KnowledgePagePoll struct {
	ID              string                    `json:"id"`
	LibraryID       string                    `json:"library_id"`
	PageID          string                    `json:"page_id"`
	SourceBlockID   string                    `json:"source_block_id"`
	Question        string                    `json:"question"`
	Options         []KnowledgePagePollOption `json:"options"`
	Multiple        bool                      `json:"multiple"`
	Status          string                    `json:"status"`
	Revision        int64                     `json:"revision"`
	CreatedByUserID string                    `json:"created_by_user_id"`
	UpdatedByUserID string                    `json:"updated_by_user_id"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

type KnowledgePagePollOption struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type KnowledgePagePollOptionWrite struct {
	ID    string `json:"id,omitempty"`
	Label string `json:"label"`
}

type KnowledgePagePollResult struct {
	OptionID string `json:"option_id"`
	Label    string `json:"label"`
	Votes    int64  `json:"votes"`
}

// KnowledgePagePollView deliberately exposes only aggregate results and the
// caller's own ballot. It never discloses voter identities.
type KnowledgePagePollView struct {
	Poll              KnowledgePagePoll         `json:"poll"`
	Results           []KnowledgePagePollResult `json:"results"`
	TotalVoters       int64                     `json:"total_voters"`
	TotalSelections   int64                     `json:"total_selections"`
	SelectedOptionIDs []string                  `json:"selected_option_ids"`
	BallotRevision    int64                     `json:"ballot_revision"`
}

type KnowledgePagePollCreate struct {
	ClientID      string   `json:"client_id"`
	SourceBlockID string   `json:"source_block_id"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	Multiple      bool     `json:"multiple"`
}

type KnowledgePagePollUpdate struct {
	ClientID         string                         `json:"client_id"`
	ExpectedRevision int64                          `json:"expected_revision"`
	Question         string                         `json:"question"`
	Options          []KnowledgePagePollOptionWrite `json:"options"`
	Multiple         bool                           `json:"multiple"`
	Status           string                         `json:"status"`
}

type KnowledgePagePollVote struct {
	ClientID               string   `json:"client_id"`
	ExpectedBallotRevision int64    `json:"expected_ballot_revision"`
	OptionIDs              []string `json:"option_ids"`
}

type KnowledgePagePollPage struct {
	Items []KnowledgePagePollView `json:"items"`
}

type PagePollService interface {
	KnowledgePagePolls(context.Context, string, string, Authority) (KnowledgePagePollPage, error)
	CreateKnowledgePagePoll(context.Context, string, string, KnowledgePagePollCreate, Authority) (KnowledgePagePollView, error)
	UpdateKnowledgePagePoll(context.Context, string, string, string, KnowledgePagePollUpdate, Authority) (KnowledgePagePollView, error)
	VoteKnowledgePagePoll(context.Context, string, string, string, KnowledgePagePollVote, Authority) (KnowledgePagePollView, error)
}
