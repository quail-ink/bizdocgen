package invoice

import (
	"log"
	"strings"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type (
	Config struct {
		FontName       string
		FontNormal     string
		FontItalic     string
		FontBold       string
		FontBoldItalic string
	}

	Builder struct {
		cfg    Config
		params *InvoiceParams
	}
)

func NewBuilder(cfg Config, paramFile string) (*Builder, error) {
	params := &InvoiceParams{}
	if err := params.Load(paramFile); err != nil {
		return nil, err
	}

	return &Builder{
		cfg:    cfg,
		params: params,
	}, nil
}

func (b *Builder) GenerateInvoice() ([]byte, error) {
	headers, err := buildInvoiceHeader(b.params)
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
	lines := strings.Split(b.params.BillTo, "\n")
	for ix, line := range lines {
		line = strings.TrimSpace(line)
		billTo.Add(text.New(line, props.Text{Size: 10, Top: float64(6*ix + 8)}))
	}

	receiveRow := row.New(30).Add(billTo)
	newPage.Add(receiveRow)

	summary := buildInvoiceSummaryRows(b.params)

	newPage.Add(summary...)

	details := buildInvoiceDetailsRows(b.params)

	newPage.Add(details...)

	payment := buildInvoicePaymentRows(b.params)

	newPage.Add(payment...)

	m.AddPages(newPage)

	document, err := m.Generate()
	if err != nil {
		log.Printf("failed to generate invoice: %v\n", err)
		return nil, err
	}

	bytes := document.GetBytes()
	return bytes, nil
}
