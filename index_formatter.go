package stag

import (
	"strings"
)

type IndexFormatter interface {
	PrepareIndex(key string) string
}

type AllCharsIndex struct {
}

func (f AllCharsIndex) PrepareIndex(key string) string {
	return key
}

type NoSpacesIndex struct {
}

func (f NoSpacesIndex) PrepareIndex(key string) string {
	return strings.Replace(key, " ", "`", -1)
}
