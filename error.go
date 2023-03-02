package openai

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError is a representation of the error body returned by OpenAI apis, and
// will be returned by Client calls if encountered.
type APIError struct {
	Data struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   any    `json:"param"`
		Code    any    `json:"code"`
	} `json:"error"`
}

// Error implements error
func (c *APIError) Error() string {
	return fmt.Sprintf("chatgpt error %s: %s", c.Data.Type, c.Data.Message)
}

// Parse the error as a roundtripper requestfunc and return the APIError if one
// was encountered.
func parseOpenAIError(r *http.Response) (*http.Response, error) {
	if r.StatusCode > 399 {
		defer r.Body.Close()
		var e APIError
		err := json.NewDecoder(r.Body).Decode(&e)
		if err != nil {
			return r, err
		}
		return r, &e
	}
	return r, nil
}
