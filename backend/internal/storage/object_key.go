package storage

import (
	"errors"
	"path"
	"strings"
)

func BuildObjectKey(cognitoSub, shotNumber string) string {
	return path.Join(cognitoSub, shotNumber+".jpg")
}

func ExtractCognitoSub(objectKey string) (string, error) {
	cleanedKey := strings.TrimPrefix(path.Clean(objectKey), "/")

	parts := strings.Split(cleanedKey, "/")
	if len(parts) != 2 {
		return "", errors.New(
			"object key must have the format {cognito_sub}/{file_name}",
		)
	}

	cognitoSub := strings.TrimSpace(parts[0])
	fileName := strings.TrimSpace(parts[1])

	if cognitoSub == "" {
		return "", errors.New("object key Cognito sub is empty")
	}

	if fileName == "" || fileName == "." {
		return "", errors.New("object key file name is empty")
	}

	return cognitoSub, nil
}
