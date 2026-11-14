package shopadmin

import (
	"context"
	"net/http"
	"testing"

	"github.com/dracory/shopstore"
)

// Verify the re-exported aliases resolve to usable types.

func TestPositionAliases(t *testing.T) {
	positions := []ProductSectionPosition{
		PositionTop, PositionBeforeForm, PositionAfterForm, PositionSidebar, PositionBottom,
	}
	for i, p := range positions {
		if p == "" {
			t.Fatalf("position %d is empty", i)
		}
	}
}

func TestSectionAlias(t *testing.T) {
	var s ProductSection = ProductSection{ID: "x", Position: PositionTop, HTML: "h"}
	if s.ID != "x" || s.Position != PositionTop {
		t.Fatalf("unexpected section %+v", s)
	}
}

func TestRendererAlias(t *testing.T) {
	var renderer ProductSectionRenderer = func(r *http.Request, p shopstore.ProductInterface) (ProductSection, error) {
		return ProductSection{ID: "r"}, nil
	}
	sec, err := renderer(nil, nil)
	if err != nil || sec.ID != "r" {
		t.Fatalf("unexpected renderer result %+v, %v", sec, err)
	}
}

func TestCustomerResolverAlias(t *testing.T) {
	var _ CustomerResolverInterface = (*aliasResolver)(nil)
}

type aliasResolver struct{}

func (a *aliasResolver) FindByID(ctx context.Context, id string) (string, string) {
	return "", ""
}

func (a *aliasResolver) SearchIDs(ctx context.Context, name, email string) ([]string, error) {
	return nil, nil
}
