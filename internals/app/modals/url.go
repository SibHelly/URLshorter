package modals

type Url struct {
	Id           int64  `json:"id"`
	Original_url string `json:"url"`
	Alias        string `json:"alias"`
	Created_at   string `json:"created_at"`
	Expires_at   string `json:"expires_at"`
	Is_active    bool   `json:"is_active"`
	Visit_count  int    `json:"visit_count"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}
