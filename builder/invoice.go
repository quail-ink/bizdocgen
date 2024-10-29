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
	"github.com/shopspring/decimal"
)

func (b *Builder) BuildInvoiceHeader() ([]marotoCore.Row, error) {
	tInvoiceID := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceID", nil)
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

	leftCol.Add(text.New(b.iParams.CompanyName, props.Text{Size: 14, Top: 8, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}))
	lines := strings.Split(b.iParams.CompanyAddr, "\n")
	for ix, line := range lines {
		leftCol.Add(text.New(line, props.Text{Size: 9, Top: float64(6*ix + 16), Align: align.Left, Color: b.fgColor}))
	}
	leftCol.Add(text.New(b.iParams.CompanyEmail, props.Text{Size: 9, Top: float64(6*(len(lines)) + 16), Align: align.Left, Color: b.fgColor}))

	rs := row.New(42).WithStyle(borderBottomStyle).Add(
		leftCol,
		col.New(6).Add(
			text.New(fmt.Sprintf("%s: %s", tInvoiceID, b.iParams.ID), props.Text{Size: 9, Top: 16, Align: align.Right, Color: b.fgColor}),
			text.New(fmt.Sprintf("%s: %s", tTaxID, b.iParams.TaxNumber), props.Text{Size: 9, Top: 22, Align: align.Right, Color: b.fgColor}),
			text.New(fmt.Sprintf("%s: %s", tIssueDate, b.iParams.Date.Format("2006/01/02")), props.Text{Size: 9, Top: 28, Align: align.Right, Color: b.fgColor}),
			text.New(fmt.Sprintf("%s: %s - %s", tPeriod,
				b.iParams.Summary.PeriodStart.Format("2006/01/02"),
				b.iParams.Summary.PeriodEnd.Format("2006/01/02"),
			), props.Text{Size: 9, Top: 34, Align: align.Right, Color: b.fgColor}),
		),
	)

	rows := []marotoCore.Row{
		rs,
		row.New(6),
	}
	return rows, nil
}

func (b *Builder) BuildInvoiceBillTo() []marotoCore.Row {
	tBillTo := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceBillTo", nil)

	billTo := col.New(8)
	billTo.Add(text.New(b.iParams.BillToCompany, props.Text{Size: 9, Top: float64(0), Style: fontstyle.Bold, Color: b.fgColor}))
	billTo.Add(text.New(b.iParams.BillToAddress, props.Text{Size: 9, Top: float64(6), Color: b.fgColor}))

	return []marotoCore.Row{
		text.NewRow(8, tBillTo, props.Text{Size: 10, Top: 0, Style: fontstyle.Bold, Color: b.fgColor}),
		row.New(12).Add(billTo),
	}
}

func (b *Builder) BuildInvoicePaymentRows() []marotoCore.Row {
	tPayment := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePayment", nil)
	tMethod := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentMethod", nil)
	tPaymentID := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentID", nil)
	tBankName := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankName", nil)
	tBankBranch := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankBranch", nil)
	tBankDepositType := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankDepositType", nil)
	tBankAccount := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankAccount", nil)
	tBankAccountName := b.i18nBundle.MusT(b.cfg.Lang, "InvoicePaymentBankAccountName", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	rows := []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tPayment, props.Text{Size: 10, Top: 8, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}),
			text.NewCol(4, "", props.Text{Size: 10, Top: 8, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		),
	}

	if b.iParams.Payment.Method == "" {
		b.iParams.Payment.Method = "Bank"
	}
	rows = append(rows, row.New(10).Add(
		col.New(2).Add(
			text.New(tMethod, props.Text{Size: 9, Top: 4, Align: align.Left, Color: b.fgColor}),
		),
		col.New(10).Add(
			text.New(b.iParams.Payment.Method, props.Text{Size: 9, Top: 4, Align: align.Right, Color: b.fgColor}),
		),
	))

	if b.iParams.Payment.PaymentID != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tPaymentID, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.PaymentID, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if b.iParams.Payment.ReceiveAccountBank != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tBankName, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.ReceiveAccountBank, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if b.iParams.Payment.ReceiveAccountBranch != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tBankBranch, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.ReceiveAccountBranch, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if b.iParams.Payment.ReceiveAccountNumber != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tBankAccount, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.ReceiveAccountNumber, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if b.iParams.Payment.ReceiveDepositType != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tBankDepositType, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.ReceiveDepositType, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if b.iParams.Payment.ReceiveAccountName != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New(tBankAccountName, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.ReceiveAccountName, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if b.iParams.Payment.ReceiveAccountSwift != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New("SWIFT", props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.ReceiveAccountSwift, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}

	if b.iParams.Payment.ReceiveAccountRouting != "" {
		rows = append(rows, row.New(6).Add(
			col.New(2).Add(
				text.New("Routing Number", props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			),
			col.New(10).Add(
				text.New(b.iParams.Payment.ReceiveAccountRouting, props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
			),
		))
	}
	return rows
}

func (b *Builder) BuildInvoiceDetailsRows() []marotoCore.Row {
	tDetails := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceDetails", nil)

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
			text.NewCol(8, tDetails, props.Text{Size: 10, Top: 8, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}),
			text.NewCol(4, "", props.Text{Size: 10, Top: 8, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		),
	}

	for ix, item := range b.iParams.DetailItems {
		paddingTop := float64(0)
		rowHeight := float64(6)
		if ix == 0 {
			paddingTop = float64(4)
			rowHeight = float64(10)
		}
		r := row.New(rowHeight)
		r.Add(
			col.New(2).Add(
				text.New(item.Date.Format("2006/01/02"), props.Text{Size: 9, Top: paddingTop, Align: align.Left, Color: b.fgColor}),
			),
			col.New(6).Add(
				text.New(item.Title, props.Text{Size: 9, Top: paddingTop, Align: align.Left, Color: b.fgColor}),
			),
		)
		if item.TotalExcludeTax.IsPositive() || item.TotalIncludeTax.IsPositive() {
			if item.TotalIncludeTax.IsPositive() {
				r.Add(
					col.New(4).Add(
						text.New(fmt.Sprintf("%s %s", item.TotalIncludeTax.RoundDown(2), b.iParams.Currency), props.Text{Size: 9, Top: paddingTop, Align: align.Right, Color: b.fgColor}),
					),
				)
			} else {
				r.Add(
					col.New(4).Add(
						text.New(fmt.Sprintf("%s %s", item.TotalExcludeTax.RoundDown(2), b.iParams.Currency), props.Text{Size: 9, Top: paddingTop, Align: align.Right, Color: b.fgColor}),
					),
				)
			}
		}
		rows = append(rows, r)

		if item.Desc != "" {
			r := row.New(6)
			r.Add(
				col.New(2),
			)
			if item.TotalExcludeTax.IsPositive() && item.Tax.IsPositive() {
				r.Add(
					col.New(6).Add(
						text.New(item.Desc, props.Text{Size: 8, Top: 0, Align: align.Left, Color: b.fgSecondaryColor}),
					),
					col.New(4).Add(
						text.New(fmt.Sprintf("VAT: %s %s", item.Tax.RoundDown(2), b.iParams.Currency), props.Text{Size: 8, Top: 0, Align: align.Right, Color: b.fgSecondaryColor}),
					),
				)
			} else {
				r.Add(
					col.New(10).Add(
						text.New(item.Desc, props.Text{Size: 8, Top: 0, Align: align.Left, Color: b.fgSecondaryColor}),
					),
				)
			}
			rows = append(rows, r)
		}
		if item.URL != "" {
			url := item.URL
			rows = append(rows, row.New(6).Add(
				col.New(2),
				col.New(10).Add(
					text.New(item.URL, props.Text{Size: 8, Top: 0, Align: align.Left, Hyperlink: &url, Color: colorLink}),
				),
			))
		} else if len(item.URLs) > 0 {
			for _, url := range item.URLs {
				rows = append(rows, row.New(6).Add(
					col.New(2),
					col.New(10).Add(
						text.New(url, props.Text{Size: 8, Top: 0, Align: align.Left, Hyperlink: &url, Color: colorLink}),
					),
				))
			}
		}
		rows = append(rows, row.New(2))
	}
	return rows
}

func (b *Builder) BuildInvoiceSummaryRows() []marotoCore.Row {
	tSummary := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummary", nil)
	tAmount := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryAmount", nil)
	tVAT := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryVAT", nil)
	tTotal := b.i18nBundle.MusT(b.cfg.Lang, "InvoiceSummaryTotalWithTax", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	var total, tax, subtotal decimal.Decimal
	if b.iParams.Summary.TotalExcludeTax.IsPositive() {
		// tax excluded?
		subtotal = b.iParams.Summary.TotalExcludeTax
		if b.iParams.Summary.Tax.IsPositive() {
			tax = b.iParams.Summary.Tax.Round(2)
		} else if b.iParams.Summary.TaxRate.IsPositive() {
			tax = subtotal.Mul(b.iParams.Summary.TaxRate).Round(2)
		}
		total = subtotal.Add(tax).Round(2)
	} else {
		// tax included?
		total = b.iParams.Summary.TotalIncludeTax
		subtotal = total.Div(decimal.NewFromFloat(1).Add(b.iParams.Summary.TaxRate)).Round(2)
		tax = total.Sub(subtotal).Round(2)
	}

	return []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tSummary, props.Text{Size: 10, Top: 8, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}),
			text.NewCol(4, tAmount, props.Text{Size: 10, Top: 8, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		),
		row.New(12).Add(
			text.NewCol(8, b.iParams.Summary.Title, props.Text{Size: 9, Top: 4, Align: align.Left, Color: b.fgColor}),
			text.NewCol(4, fmt.Sprintf("%s %s", subtotal.RoundDown(2), b.iParams.Currency), props.Text{Size: 9, Top: 4, Align: align.Right, Color: b.fgColor}),
		),
		row.New(8).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, tVAT, props.Text{Size: 9, Top: 0, Align: align.Left, Color: b.fgColor}),
			text.NewCol(6, fmt.Sprintf("%s %s", tax, b.iParams.Currency), props.Text{Size: 9, Top: 0, Align: align.Right, Color: b.fgColor}),
		),
		row.New(10).Add(
			text.NewCol(6, tTotal, props.Text{Size: 10, Top: 4, Align: align.Left, Style: fontstyle.Bold, Color: b.fgColor}),
			text.NewCol(6, fmt.Sprintf("%s %s", total, b.iParams.Currency), props.Text{Size: 10, Top: 4, Align: align.Right, Style: fontstyle.Bold, Color: b.fgColor}),
		),
	}
}
