package storage

func BuildObjectKey(shotNumber string) string {
	return shotNumber + ".jpg"
}
