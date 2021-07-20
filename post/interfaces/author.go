package IPost

type Author struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AvatarUrl    string `json:"avatarUrl"`
	FollowStatus int    `json:"followStatus"`
	Job          string `json:"job"`
}
