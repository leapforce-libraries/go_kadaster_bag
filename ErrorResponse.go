package go_kadaster_bag

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"errorDescription"`
	ErrorDetail      string `json:"errorDetail"`
}
