package IUser

type Info struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
	Following int    `json:"following"`
	Followers int    `json:"followers"`
}
