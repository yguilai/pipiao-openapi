package id

import (
    "github.com/google/uuid"
    "strconv"
    "strings"
)

func New(prefix string) string {
    return prefix + "_" + strconv.Itoa(int(uuid.New().ID()))
}

func SimpleUUID() string {
    return strings.ReplaceAll(uuid.NewString(), "-", "")
}
