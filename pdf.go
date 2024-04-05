package invoice

import (
	marotoCore "github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/core/entity"

	"log"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/johnfercher/maroto/v2/pkg/repository"
)

func (b *Builder) CreateMetricsDecorator(head []marotoCore.Row) (marotoCore.Maroto, error) {
	repo := repository.New()
	useCustomFonts := false
	customFontsAdded := false
	var customFonts []*entity.CustomFont
	if b.cfg.FontName == "" {
		b.cfg.FontName = "default-font"
	}
	if b.cfg.FontNormal != "" || b.cfg.FontItalic != "" || b.cfg.FontBold != "" || b.cfg.FontBoldItalic != "" {
		if b.cfg.FontNormal != "" {
			repo = repo.AddUTF8Font(b.cfg.FontName, fontstyle.Normal, b.cfg.FontNormal)
		}
		if b.cfg.FontItalic != "" {
			repo = repo.AddUTF8Font(b.cfg.FontName, fontstyle.Italic, b.cfg.FontItalic)
		}
		if b.cfg.FontBold != "" {
			repo = repo.AddUTF8Font(b.cfg.FontName, fontstyle.Bold, b.cfg.FontBold)
		}
		if b.cfg.FontBoldItalic != "" {
			repo = repo.AddUTF8Font(b.cfg.FontName, fontstyle.BoldItalic, b.cfg.FontBoldItalic)
		}
		customFontsAdded = true
	}
	if customFontsAdded {
		var err error
		customFonts, err = repo.Load()
		if err != nil {
			log.Printf("failed to load custom fonts: %v\n", err)
		} else {
			useCustomFonts = true
		}
	}

	bu := config.NewBuilder()
	bu = bu.WithPageNumber("Page {current} of {total}", props.Bottom)
	if useCustomFonts {
		bu = bu.WithCustomFonts(customFonts)
		bu = bu.WithDefaultFont(&props.Font{Family: b.cfg.FontName})
	}

	cfg := bu.Build()

	mrt := maroto.New(cfg)

	m := maroto.NewMetricsDecorator(mrt)

	if err := m.RegisterHeader(head...); err != nil {
		log.Printf("failed to register header: %v\n", err)
		return nil, err
	}
	return m, nil
}
