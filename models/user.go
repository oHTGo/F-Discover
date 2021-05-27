package models

import "cloud.google.com/go/firestore"

type User struct {
	ID        string                            `firestore:"id" json:"id"`
	Name      string                            `firestore:"name" json:"name"`
	AvatarUrl string                            `firestore:"avatarUrl" json:"avatarUrl"`
	Following map[string]*firestore.DocumentRef `firestore:"following" json:"following"`
	Followers map[string]*firestore.DocumentRef `firestore:"followers" json:"followers"`
}
