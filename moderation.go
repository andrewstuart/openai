package openai

import "context"

const (
	ModerationModelStable = "text-moderation-stable"
	ModerationModelLatest = "text-moderation-latest"
)

// ModerationReq details a request for moderation to the OpenAI endpoint.
type ModerationReq struct {
	Input string  `json:"input"`
	Model *string `json:"model,omitempty"`
}

// ModerationRes holds the resulst of the OpenAI moderation call.
type ModerationRes struct {
	ID      string             `json:"id"`
	Model   string             `json:"model"`
	Results []ModerationResult `json:"results"`
}

// ModerationResult is the specific results of the input.
type ModerationResult struct {
	Flagged        bool               `json:"flagged"`
	Categories     map[string]bool    `json:"categories"`
	CategoryScores map[string]float64 `json:"category_scores"`
}

// Moderation returns whether or not OpenAI considers the input text to be
// rule-breaking, and exactly how rule-breaking.
func (c Client) Moderation(ctx context.Context, req ModerationReq) (*ModerationRes, error) {
	var res ModerationRes
	err := c.c.R().Post("moderations").JSON(req).Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
