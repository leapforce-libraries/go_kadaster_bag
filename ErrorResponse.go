package go_kadaster_bag

type ErrorResponse struct {
	Status        int    `json:"status"`
	Title         string `json:"title"`
	Type          string `json:"type"`
	Detail        string `json:"detail"`
	Instance      string `json:"instance"`
	Code          string `json:"code"`
	InvalidParams []struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		Reason string `json:"reason"`
	} `json:"invalidParams"`
}
