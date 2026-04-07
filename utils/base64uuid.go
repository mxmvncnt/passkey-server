package utils

import (
	"encoding/base64"

	"github.com/google/uuid"
)

func GetUuidFromBase64(b64UUID string) uuid.UUID {
	bytes, err := base64.RawURLEncoding.DecodeString(b64UUID)
	if err != nil {
		panic(err)
	}
	result, err := uuid.FromBytes(bytes)
	if err != nil {
		panic(err)
	}
	return result
}

func GetByteArrayFromBase64(b64UUID string) []byte {
	bytes, err := base64.RawURLEncoding.DecodeString(b64UUID)
	if err != nil {
		panic(err)
	}
	return bytes
}
