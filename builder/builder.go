package builder

import (
	"log"
	"log/slog"

	"github.com/johnfercher/maroto/v2/pkg/components/page"
	marotoCore "github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/quail-ink/bizdocgen/core"
	"github.com/quail-ink/bizdocgen/i18n"
)

type (
	Config struct {
		FontName       string
		FontNormal     string
		FontItalic     string
		FontBold       string
		FontBoldItalic string

		Lang string
	}

	Builder struct {
		cfg        Config
		i18nBundle *i18n.I18nBundle
		iParams    *core.InvoiceParams
		psParams   *core.PaymentStatementParams
		Round      int32
	}
)

func NewInvoiceBuilder(cfg Config, params *core.InvoiceParams) (*Builder, error) {
	i18nBundle := i18n.New()
	if cfg.Lang == "" {
		cfg.Lang = "en"
	}
	round := 2
	if params.Currency == "JPY" || params.Currency == "円" {
		round = 0
	}
	return &Builder{
		cfg:        cfg,
		i18nBundle: i18nBundle,
		iParams:    params,
		Round:      int32(round),
	}, nil
}

func NewInvoiceBuilderFromFile(cfg Config, filename string) (*Builder, error) {
	params := &core.InvoiceParams{}
	if err := params.Load(filename); err != nil {
		return nil, err
	}
	return NewInvoiceBuilder(cfg, params)
}

func NewPaymentStatementBuilder(cfg Config, params *core.PaymentStatementParams) (*Builder, error) {
	i18nBundle := i18n.New()
	if cfg.Lang == "" {
		cfg.Lang = "en"
	}
	round := 2
	if params.Currency == "JPY" || params.Currency == "円" {
		round = 0
	}
	return &Builder{
		cfg:        cfg,
		i18nBundle: i18nBundle,
		psParams:   params,
		Round:      int32(round),
	}, nil
}

func NewPaymentStatementBuilderFromFile(cfg Config, filename string) (*Builder, error) {
	params := &core.PaymentStatementParams{}
	if err := params.Load(filename); err != nil {
		return nil, err
	}
	return NewPaymentStatementBuilder(cfg, params)
}

func (b *Builder) GenerateInvoice() ([]byte, error) {
	headers, err := b.BuildInvoiceHeader()
	if err != nil {
		log.Printf("failed to build invoice header: %v\n", err)
		return nil, err
	}

	m, err := b.CreateMetricsDecorator(headers)
	if err != nil {
		log.Printf("failed to register header: %v\n", err)
		return nil, err
	}

	newPage := page.New()

	receiveRows := b.BuildInvoiceBillTo()
	newPage.Add(receiveRows...)

	summary := b.BuildInvoiceSummaryRows()

	newPage.Add(summary...)

	details := b.BuildInvoiceDetailsRows()

	newPage.Add(details...)

	if !b.iParams.Payment.Disabled {
		payment := b.BuildInvoicePaymentRows()
		newPage.Add(payment...)
	}

	m.AddPages(newPage)

	return b.getBytesFromMaroto(m)
}

func (b *Builder) GeneratePaymentStatement() ([]byte, error) {
	headers, err := b.BuildPsHeader()
	if err != nil {
		log.Printf("failed to build header: %v\n", err)
		return nil, err
	}

	m, err := b.CreateMetricsDecorator(headers)
	if err != nil {
		log.Printf("failed to register header: %v\n", err)
		return nil, err
	}

	newPage := page.New()

	payer := b.BuildPsPayer()
	newPage.Add(payer...)

	payee := b.BuildPsPayee()
	newPage.Add(payee...)

	channel := b.BuildPsChannelRows()
	newPage.Add(channel...)

	summary := b.BuildPsSummaryRows()
	newPage.Add(summary...)

	details := b.BuildPsDetailsRows()
	newPage.Add(details...)

	m.AddPages(newPage)

	return b.getBytesFromMaroto(m)
}

func (b *Builder) getBytesFromMaroto(maroto marotoCore.Maroto) ([]byte, error) {
	document, err := maroto.Generate()
	if err != nil {
		slog.Error("failed to generate document from maroto", "error", err)
		return nil, err
	}

	bytes := document.GetBytes()
	return bytes, nil
}
