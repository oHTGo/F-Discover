package IUser

type Info struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
	Quote     string `json:"quote"`
	Following int    `json:"following"`
	Followers int    `json:"followers"`
}
