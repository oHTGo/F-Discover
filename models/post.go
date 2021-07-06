package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Post struct {
	ID           string                 `firestore:"id" json:"id"`
	Content      string                 `firestore:"content" json:"content"`
	ThumbnailUrl string                 `firestore:"thumbnailUrl" json:"thumbnailUrl"`
	VideoUrl     string                 `firestore:"videoUrl" json:"videoUrl"`
	Likes        map[string]bool        `firestore:"likes" json:"likes"`
	Comments     map[string]Comment     `firestore:"comments" json:"comments"`
	Location     string                 `firestore:"location" json:"location"`
	Author       *firestore.DocumentRef `firestore:"author" json:"author"`
	CreatedAt    time.Time              `firestore:"createdAt" json:"createdAt"`
}

type Comment struct {
	ID        string                 `firestore:"id" json:"id"`
	Content   string                 `firestore:"content" json:"content"`
	CreatedAt time.Time              `firestore:"createdAt" json:"createdAt"`
	Author    *firestore.DocumentRef `firestore:"author" json:"author"`
}
