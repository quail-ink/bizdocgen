package builder

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	marotoCore "github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func (b *Builder) BuildInvoiceHeader() ([]marotoCore.Row, error) {
	tTaxID := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceTaxID", nil)
	tIssueDate := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceIssueDate", nil)
	tPeriod := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePeriod", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}
	leftCol := col.New(6)

	if b.iParams.CompanySeal != "" {
		fd, err := os.Open(b.iParams.CompanySeal)
		if err != nil {
			log.Printf("failed to open seal file: %v\n", err)
			return nil, err
		}
		defer fd.Close()
		buf, err := io.ReadAll(fd)
		if err != nil {
			log.Printf("failed to read seal file: %v\n", err)
			return nil, err
		}

		leftCol.Add(image.NewFromBytes(buf, extension.Png, props.Rect{
			Center:  false,
			Percent: 20,
			Left:    34,
			Top:     7,
		}))
	}

	leftCol.Add(text.New(b.iParams.CompanyName, props.Text{Size: 14, Top: 8, Align: align.Left, Style: fontstyle.Bold}))
	lines := strings.Split(b.iParams.CompanyAddr, "\n")
	for ix, line := range lines {
		leftCol.Add(text.New(line, props.Text{Size: 10, Top: float64(6*ix + 16), Align: align.Left}))
	}
	leftCol.Add(text.New(b.iParams.CompanyEmail, props.Text{Size: 10, Top: float64(6*(len(lines)) + 16), Align: align.Left}))

	rs := row.New(40).WithStyle(borderBottomStyle).Add(
		leftCol,
		col.New(6).Add(
			text.New(fmt.Sprintf("%s: %s", tTaxID, b.iParams.TaxNumber), props.Text{Size: 10, Top: 16, Align: align.Right}),
			text.New(fmt.Sprintf("%s: %s", tIssueDate, b.iParams.Date.Format("2006/01/02")), props.Text{Size: 10, Top: 22, Align: align.Right}),
			text.New(fmt.Sprintf("%s: %s - %s", tPeriod,
				b.iParams.Summary.PeriodStart.Format("2006/01/02"),
				b.iParams.Summary.PeriodEnd.Format("2006/01/02"),
			), props.Text{Size: 10, Top: 28, Align: align.Right}),
		),
	)

	rows := []marotoCore.Row{
		rs,
		row.New(10),
	}
	return rows, nil
}

func (b *Builder) BuildInvoiceBillTo() []marotoCore.Row {
	tBillTo := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceBillTo", nil)

	billTo := col.New(8)
	lines := strings.Split(b.iParams.BillTo, "\n")
	for ix, line := range lines {
		line = strings.TrimSpace(line)
		billTo.Add(text.New(line, props.Text{Size: 10, Top: float64(6 * ix)}))
	}

	return []marotoCore.Row{
		text.NewRow(10, tBillTo, props.Text{Size: 12, Top: 0, Style: fontstyle.Bold}),
		row.New(16).Add(billTo),
	}
}

func (b *Builder) BuildInvoicePaymentRows() []marotoCore.Row {
	tPayment := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePayment", nil)
	tBankName := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankName", nil)
	tBankBranch := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankBranch", nil)
	tBankAccount := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankAccount", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	rows := []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tPayment, props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, "", props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
	}

	rows = append(rows, row.New(12).Add(
		col.New(2).Add(
			text.New(tBankName, props.Text{Size: 10, Top: 4, Align: align.Left}),
		),
		col.New(10).Add(
			text.New(b.iParams.Payment.ReceiveAccountBank, props.Text{Size: 10, Top: 4, Align: align.Right}),
		),
	))

	if b.iParams.Payment.ReceiveAccountBranch != "" {
		rows = append(rows, row.New(8).Add(
			col.New(2).Add(
				text.New(tBankBranch, props.Text{Size: 10, Top: 0, Align: align.Left}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.ReceiveAccountBranch, props.Text{Size: 10, Top: 0, Align: align.Right}),
			),
		))
	}

	rows = append(rows, row.New(8).Add(
		col.New(2).Add(
			text.New(tBankAccount, props.Text{Size: 10, Top: 0, Align: align.Left}),
		),
		col.New(10).Add(
			text.New(b.iParams.Payment.ReceiveAccountNumber, props.Text{Size: 10, Top: 0, Align: align.Right}),
		),
	))

	if b.iParams.Payment.ReceiveAccountSwift != "" {
		rows = append(rows, row.New(8).Add(
			col.New(2).Add(
				text.New("SWIFT", props.Text{Size: 10, Top: 0, Align: align.Left}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.ReceiveAccountSwift, props.Text{Size: 10, Top: 0, Align: align.Right}),
			),
		))
	}

	if b.iParams.Payment.ReceiveAccountRouting != "" {
		rows = append(rows, row.New(8).Add(
			col.New(2).Add(
				text.New("Routing Number", props.Text{Size: 10, Top: 0, Align: align.Left}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.ReceiveAccountRouting, props.Text{Size: 10, Top: 0, Align: align.Right}),
			),
		))
	}
	return rows
}

func (b *Builder) BuildInvoiceDetailsRows() []marotoCore.Row {
	tDetails := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceDetails", nil)

	colorSecondary := &props.Color{
		Red:   100,
		Green: 100,
		Blue:  100,
	}
	colorLink := &props.Color{
		Red:   0,
		Green: 0,
		Blue:  255,
	}

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	rows := []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tDetails, props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, "", props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
	}

	for ix, item := range b.iParams.DetailItems {
		paddingTop := float64(0)
		rowHeight := float64(8)
		if ix == 0 {
			paddingTop = float64(4)
			rowHeight = float64(12)
		}
		rows = append(rows, row.New(rowHeight).Add(
			col.New(2).Add(
				text.New(item.Date.Format("2006/01/02"), props.Text{Size: 10, Top: paddingTop, Align: align.Left}),
			),
			col.New(10).Add(
				text.New(item.Title, props.Text{Size: 10, Top: paddingTop, Align: align.Left}),
			),
		))
		rows = append(rows, row.New(8).Add(
			col.New(2),
			col.New(10).Add(
				text.New(item.Desc, props.Text{Size: 8, Top: 0, Align: align.Left, Color: colorSecondary}),
			),
		))
		if item.URL != "" {
			url := item.URL
			rows = append(rows, row.New(8).Add(
				col.New(2),
				col.New(10).Add(
					text.New(item.URL, props.Text{Size: 8, Top: 0, Align: align.Left, Hyperlink: &url, Color: colorLink}),
				),
			))
		} else if len(item.URLs) > 0 {
			for _, url := range item.URLs {
				rows = append(rows, row.New(8).Add(
					col.New(2),
					col.New(10).Add(
						text.New(url, props.Text{Size: 8, Top: 0, Align: align.Left, Hyperlink: &url, Color: colorLink}),
					),
				))
			}
		}
	}
	return rows
}

func (b *Builder) BuildInvoiceSummaryRows() []marotoCore.Row {
	tSummary := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummary", nil)
	tAmount := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryAmount", nil)
	tJct := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryJct", nil)
	tTotal := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryTotalWithTax", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	subtotal := b.iParams.Summary.TotalExcludeTax
	tax := subtotal.Mul(b.iParams.Summary.TaxRate).RoundDown(2)
	total := subtotal.Add(tax).RoundDown(2)

	return []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tSummary, props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, tAmount, props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(14).Add(
			text.NewCol(8, b.iParams.Summary.Title, props.Text{Size: 10, Top: 4, Align: align.Left}),
			text.NewCol(4, fmt.Sprintf("%s %s", subtotal.RoundDown(2), b.iParams.Currency), props.Text{Size: 10, Top: 4, Align: align.Right}),
		),
		row.New(10).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, tJct, props.Text{Size: 10, Top: 0, Align: align.Left}),
			text.NewCol(6, fmt.Sprintf("%s %s", tax, b.iParams.Currency), props.Text{Size: 10, Top: 0, Align: align.Right}),
		),
		row.New(20).Add(
			text.NewCol(6, tTotal, props.Text{Size: 12, Top: 4, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(6, fmt.Sprintf("%s %s", total, b.iParams.Currency), props.Text{Size: 12, Top: 4, Align: align.Right, Style: fontstyle.Bold}),
		),
	}
}
