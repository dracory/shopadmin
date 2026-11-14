package shared

import "testing"

func TestProductSectionPositionConstants(t *testing.T) {
	positions := []ProductSectionPosition{
		PositionTop, PositionBeforeForm, PositionAfterForm, PositionSidebar, PositionBottom,
	}
	for i, p := range positions {
		if p == "" {
			t.Fatalf("position %d is empty", i)
		}
	}
	if PositionTop != "top" || PositionBottom != "bottom" {
		t.Fatal("unexpected position constant values")
	}
}

func TestProductSectionFields(t *testing.T) {
	s := ProductSection{
		ID:       "sec",
		Position: PositionTop,
		HTML:     "<div>x</div>",
		CSS:      ".x{}",
		JS:       "x()",
		CSSURLs:  []string{"a.css"},
		JSURLs:   []string{"a.js"},
	}
	if s.ID != "sec" || s.Position != PositionTop || len(s.CSSURLs) != 1 || len(s.JSURLs) != 1 {
		t.Fatalf("unexpected section fields: %+v", s)
	}
}
