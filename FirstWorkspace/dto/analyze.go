package dto

type AnalyzeRequest struct {
	URL string `json:"url"`
}

type AnalyzeResult struct {
	Title      string  `json:"title"`
	Folder     string  `json:"folder"`
	Confidence float64 `json:"confidence"`
}

type AnalyzeResponse struct {
	Title           string  `json:"title"`
	URL             string  `json:"url"`
	SuggestedFolder *string `json:"suggested_folder"`
	IsNewFolder     bool    `json:"is_new_folder"`
	Confidence      float64 `json:"confidence"`
	Message         string  `json:"message,omitempty"`
}
