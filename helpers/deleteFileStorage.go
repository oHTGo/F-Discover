package helpers

import (
	"f-discover/env"
	"f-discover/instance"
	"f-discover/services"
	"strings"
)

func DeleteFileStorage(path string) {
	var object string
	if ok := strings.Contains(path, env.Get().STORAGE_BUCKET); ok {
		runes := []rune(path)
		position := strings.Index(path, env.Get().STORAGE_BUCKET) + len(env.Get().STORAGE_BUCKET) + 1
		object = string(runes[position:])
	} else {
		object = path
	}
	bucket := services.GetInstance().StorageBucket
	bucket.Object(object).Delete(instance.CtxBackground)
}
