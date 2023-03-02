package openai

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	Data struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   any    `json:"param"`
		Code    any    `json:"code"`
	} `json:"error"`
}

func (c *APIError) Error() string {
	return fmt.Sprintf("chatgpt error %s: %s", c.Data.Type, c.Data.Message)
}

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
