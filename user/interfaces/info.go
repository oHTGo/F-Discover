package IUser

type Info struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CoverUrl  string `json:"coverUrl"`
	Job       string `json:"job"`
	AvatarUrl string `json:"avatarUrl"`
	Quote     string `json:"quote"`
	Following int    `json:"following"`
	Followers int    `json:"followers"`
}
