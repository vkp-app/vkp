package promutil

type Response struct {
	Status    string   `json:"status"`
	Data      any      `json:"data"`
	ErrorType string   `json:"errorType"`
	Error     string   `json:"error"`
	Warnings  []string `json:"warnings"`
}
