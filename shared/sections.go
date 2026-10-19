package shared

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/shopstore"
)

// ProcessedSections contains rendered custom sections grouped by position
// and aggregated CSS/JS options for layout rendering.
type ProcessedSections struct {
	Positions map[ProductSectionPosition][]hb.TagInterface
	Options   struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}
}

// RenderPosition returns a container element holding all sections for a given position.
// Returns nil if there are no sections for that position.
func (ps ProcessedSections) RenderPosition(pos ProductSectionPosition) hb.TagInterface {
	tags, ok := ps.Positions[pos]
	if !ok || len(tags) == 0 {
		return nil
	}
	container := hb.Div().Class("shopadmin-section-container shopadmin-section-container-" + string(pos))
	for _, tag := range tags {
		container.Child(tag)
	}
	return container
}

// ProcessProductSections executes the given section renderers in order,
// handling errors gracefully and aggregating section HTML, CSS, JS, and external URLs.
func ProcessProductSections(
	renderers []ProductSectionRenderer,
	r *http.Request,
	product shopstore.ProductInterface,
	logger *slog.Logger,
) ProcessedSections {
	res := ProcessedSections{
		Positions: make(map[ProductSectionPosition][]hb.TagInterface),
	}

	for _, renderer := range renderers {
		if renderer == nil {
			continue
		}

		sec, err := renderer(r, product)
		if err != nil {
			if logger != nil {
				logger.Error("Error rendering product custom section", "error", err, "id", sec.ID)
			}
			continue
		}

		pos := sec.Position
		if pos == "" {
			pos = PositionBottom
		}

		secDiv := hb.Div().Class("shopadmin-custom-section")
		if sec.ID != "" {
			secDiv.Attr("id", "section-"+sec.ID)
		}
		if sec.HTML != "" {
			secDiv.Child(hb.Raw(sec.HTML))
		}

		res.Positions[pos] = append(res.Positions[pos], secDiv)

		if sec.CSS != "" {
			res.Options.Styles = append(res.Options.Styles, sec.CSS)
		}
		if sec.JS != "" {
			res.Options.Scripts = append(res.Options.Scripts, sec.JS)
		}
		for _, url := range sec.CSSURLs {
			url = strings.TrimSpace(url)
			if url != "" {
				res.Options.StyleURLs = appendUnique(res.Options.StyleURLs, url)
			}
		}
		for _, url := range sec.JSURLs {
			url = strings.TrimSpace(url)
			if url != "" {
				res.Options.ScriptURLs = appendUnique(res.Options.ScriptURLs, url)
			}
		}
	}

	return res
}

func appendUnique(slice []string, val string) []string {
	for _, item := range slice {
		if item == val {
			return slice
		}
	}
	return append(slice, val)
}
