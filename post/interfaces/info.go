package IPost

import "time"

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
