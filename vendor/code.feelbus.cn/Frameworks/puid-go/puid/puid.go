package puid

import (
	"encoding/hex"
	"fmt"

	"github.com/pborman/uuid"
)

// NewPUID returns a new PUID
func NewPUID() int64 {
	t, _, _ := uuid.GetTime()
	return int64(t)
}

// NewBinaryOrderedUUID returns a binary UUID
func NewBinaryOrderedUUID() []byte {
	ordered := NewOrderedUUID()
	h, _ := hex.DecodeString(ordered)
	return h
}

func NewOrderedUUID() string {
	uid := uuid.NewUUID().String()
	return fmt.Sprintf("%s%s%s%s%s", uid[14:18], uid[9:13], uid[0:8], uid[19:23], uid[24:])
}
