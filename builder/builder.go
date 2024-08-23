package builder

import (
	"log"
	"log/slog"
	"strings"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	marotoCore "github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/quail-ink/bizdocgen/core"
	"github.com/quail-ink/bizdocgen/i18n"
	"github.com/quail-ink/bizdocgen/invoice"
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
	headers, err := invoice.BuildInvoiceHeader(b.iParams)
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
	newPage.Add(
		text.NewRow(14, "BILL TO", props.Text{Size: 12, Top: 20, Style: fontstyle.Bold}),
	)

	// bill to
	billTo := col.New(12)
	lines := strings.Split(b.iParams.BillTo, "\n")
	for ix, line := range lines {
		line = strings.TrimSpace(line)
		billTo.Add(text.New(line, props.Text{Size: 10, Top: float64(6*ix + 8)}))
	}

	receiveRow := row.New(30).Add(billTo)
	newPage.Add(receiveRow)

	summary := invoice.BuildInvoiceSummaryRows(b.iParams)

	newPage.Add(summary...)

	details := invoice.BuildInvoiceDetailsRows(b.iParams)

	newPage.Add(details...)

	payment := invoice.BuildInvoicePaymentRows(b.iParams)

	newPage.Add(payment...)

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
