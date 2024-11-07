package entities

import "fmt"

func GenerateCategoriesKey() string {
	return "categories"
}

func GenerateNumberOfPagesKey(key string) string {
	return fmt.Sprintf("category:%s:pages-quantity", key)
}