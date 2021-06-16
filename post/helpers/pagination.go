package PostHelpers

import "cloud.google.com/go/firestore"

func Paginate(x []*firestore.DocumentSnapshot, skip int, size int) []*firestore.DocumentSnapshot {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}
