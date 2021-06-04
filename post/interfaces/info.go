package IPost

type Author struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

type Info struct {
	ID      string   `json:"id"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
	Videos  []string `json:"videos"`
	Likes   int      `json:"likes"`
	Author  Author   `json:"author"`
}

type InfoWithoutAuthor struct {
	ID      string   `json:"id"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
	Videos  []string `json:"videos"`
	Likes   int      `json:"likes"`
}
