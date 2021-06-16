package IPost

import "time"

type Author struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Job       string `json:"job"`
}

type Info struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	VideoUrl  string    `json:"videoUrl"`
	Likes     int       `json:"likes"`
	Comments  int       `json:"comments"`
	Author    Author    `json:"author"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"createdAt"`
}

type InfoWithoutAuthor struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	VideoUrl  string    `json:"videoUrl"`
	Likes     int       `json:"likes"`
	Comments  int       `json:"comments"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"createdAt"`
}
