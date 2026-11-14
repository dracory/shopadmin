package shared

import "testing"

func TestBreadcrumbFields(t *testing.T) {
	b := Breadcrumb{Name: "Home", URL: "/admin"}
	if b.Name != "Home" || b.URL != "/admin" {
		t.Fatalf("unexpected breadcrumb values: %+v", b)
	}
}
