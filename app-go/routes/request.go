package routes

type ReportRequestBody struct {
	Name string `json:"name"`
}

type SearchRequestBody struct {
	Per string `json:"per"`
}
