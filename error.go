package main

import "fmt"

type OpenAIError struct {
	Data struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   any    `json:"param"`
		Code    any    `json:"code"`
	} `json:"error"`
}

func (c *OpenAIError) Error() string {
	return fmt.Sprintf("chatgpt error %s: %s", c.Data.Type, c.Data.Message)
}
