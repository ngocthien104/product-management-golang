package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	slugRegex   = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	searchRegex = regexp.MustCompile(`^[\p{L}\p{N}\s.-]+$`)
)

var AllowedCategories = map[string]bool{
	"electronics": true,
	"books":       true,
	"clothing":    true,
	"home":        true,
}

func Required(fieldName, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func StringLength(fieldName, value string, min, max int) error {
	length := len([]rune(value))
	if length < min || length > max {
		return fmt.Errorf("%s must be between %d and %d characters", fieldName, min, max)
	}
	return nil
}

func Slug(value string) error {
	if !slugRegex.MatchString(value) {
		return fmt.Errorf("slug must contain lowercase letters, numbers, hyphens or dots")
	}
	return nil
}

func Search(value string) error {
	if value == "" {
		return nil
	}
	if err := StringLength("search", value, 1, 50); err != nil {
		return err
	}
	if !searchRegex.MatchString(value) {
		return fmt.Errorf("search may only contain letters, numbers, spaces, dots and hyphens")
	}
	return nil
}

func PositiveInt(fieldName, value string) (int64, error) {
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", fieldName)
	}
	return n, nil
}

func Category(value string) error {
	if !AllowedCategories[value] {
		return fmt.Errorf("category must be one of: electronics, books, clothing, home")
	}
	return nil
}

func ProductInput(name, slug, category string, price float64) error {
	if err := Required("name", name); err != nil {
		return err
	}
	if err := StringLength("name", name, 2, 100); err != nil {
		return err
	}
	if err := Required("slug", slug); err != nil {
		return err
	}
	if err := StringLength("slug", slug, 2, 120); err != nil {
		return err
	}
	if err := Slug(slug); err != nil {
		return err
	}
	if err := Category(category); err != nil {
		return err
	}
	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	return nil
}
