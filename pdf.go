package invoice

import (
	marotoCore "github.com/johnfercher/maroto/v2/pkg/core"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/johnfercher/maroto/v2/pkg/repository"
	"github.com/sirupsen/logrus"
)

func (b *Builder) CreateMetricsDecorator(head []marotoCore.Row) (marotoCore.Maroto, error) {
	customFonts, err := repository.New().
		AddUTF8Font(b.cfg.FontName, fontstyle.Normal, b.cfg.FontNormal).
		AddUTF8Font(b.cfg.FontName, fontstyle.Italic, b.cfg.FontItalic).
		AddUTF8Font(b.cfg.FontName, fontstyle.Bold, b.cfg.FontBold).
		AddUTF8Font(b.cfg.FontName, fontstyle.BoldItalic, b.cfg.FontBoldItalic).
		Load()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load custom fonts")
		return nil, err
	}

	cfg := config.NewBuilder().
		WithPageNumber("Page {current} of {total}", props.Bottom).
		WithCustomFonts(customFonts).
		WithDefaultFont(&props.Font{Family: b.cfg.FontName}).
		Build()

	mrt := maroto.New(cfg)

	m := maroto.NewMetricsDecorator(mrt)

	err = m.RegisterHeader(head...)
	if err != nil {
		logrus.WithError(err).Error("failed to register header")
		return nil, err
	}
	return m, nil
}
