package models

import "cloud.google.com/go/firestore"

type Post struct {
	ID      string                 `firestore:"id" json:"id"`
	Content string                 `firestore:"content" json:"content"`
	Images  []string               `firestore:"images" json:"images"`
	Videos  []string               `firestore:"videos" json:"videos"`
	Likes   []map[string]bool      `firestore:"likes" json:"likes"`
	Author  *firestore.DocumentRef `firestore:"author" json:"author"`
}
