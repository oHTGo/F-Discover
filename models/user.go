package models

type User struct {
	ID        string          `firestore:"id" json:"id"`
	Name      string          `firestore:"name" json:"name"`
	CoverUrl  string          `firestore:"coverUrl" json:"coverUrl"`
	AvatarUrl string          `firestore:"avatarUrl" json:"avatarUrl"`
	Job       string          `firestore:"job" json:"job"`
	Quote     string          `firestore:"quote" json:"quote"`
	Following map[string]bool `firestore:"following" json:"following"`
	Followers map[string]bool `firestore:"followers" json:"followers"`
}
