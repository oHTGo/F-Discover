package helpers

import (
	"f-discover/env"
	"strings"
)

func DeleteFileLocal(path string) {
	if ok := strings.Contains(path, env.Get().URL); !ok {
		return
	}

	runes := []rune(path)
	position := strings.Index(path, env.Get().URL) + len(env.Get().URL) + 1
	object := string(runes[position:])

	RemoveFile(object)
}
