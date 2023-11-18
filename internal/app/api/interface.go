package api

type ExecRequest struct {
	Input string `json:"input"` // hex encoded
}

type ExecResponse struct {
	Output string `json:"output"` // hex encoded
	Error  string `json:"error"`
}
