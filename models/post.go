package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Post struct {
	ID      string                 `firestore:"id" json:"id"`
	Content string                 `firestore:"content" json:"content"`
	Images  []string               `firestore:"images" json:"images"`
	Videos  []string               `firestore:"videos" json:"videos"`
	Likes   map[string]bool        `firestore:"likes" json:"likes"`
	Comment map[string]Comment     `firestore:"comments" json:"comments"`
	Author  *firestore.DocumentRef `firestore:"author" json:"author"`
}

type Comment struct {
	ID        string                 `firestore:"id" json:"id"`
	Content   string                 `firestore:"content" json:"content"`
	CreatedAt time.Time              `firestore:"createdAt" json:"createdAt"`
	UpdatedAt time.Time              `firestore:"updatedAt" json:"updatedAt"`
	Author    *firestore.DocumentRef `firestore:"author" json:"author"`
}
