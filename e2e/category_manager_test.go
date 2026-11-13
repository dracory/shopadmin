package e2e

import (
	"strings"
	"testing"
)

func TestCategoryManagerLoadsSeededCategories(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("categories"))
	if !strings.Contains(body, "Category 01") {
		t.Fatalf("expected seeded 'Category 01' in category manager")
	}
}
