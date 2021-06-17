package interfaces

import (
	"cloud.google.com/go/firestore"
)

type CurrentUser struct {
	ID        string
	Reference *firestore.DocumentRef
}
