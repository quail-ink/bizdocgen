package paymentstatement

import (
	"fmt"
	"io"
	"log"
	"os"

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
	"github.com/quail-ink/bizdocgen/core"
	"github.com/shopspring/decimal"
)

func BuildHeader(params *core.PaymentStatementParams) ([]marotoCore.Row, error) {
	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}
	leftCol := col.New(6)

	if params.CompanySeal != "" {
		fd, err := os.Open(params.CompanySeal)
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
			Percent: 24,
			Left:    44,
			Top:     8,
		}))
	}

	leftCol.Add(text.New("Payment Statement", props.Text{Size: 14, Top: 8, Align: align.Left, Style: fontstyle.Bold}))

	rs := row.New(28).WithStyle(borderBottomStyle).Add(
		leftCol,
		col.New(6).Add(
			text.New(fmt.Sprintf("Issue Date: %s", params.Date.Format("2006/01/02")), props.Text{Size: 10, Top: 8, Align: align.Right}),
			text.New(fmt.Sprintf("Period: %s - %s",
				params.PeriodStart.Format("2006/01/02"),
				params.PeriodEnd.Format("2006/01/02"),
			), props.Text{Size: 10, Top: 16, Align: align.Right}),
		),
	)

	rows := []marotoCore.Row{
		rs,
	}
	return rows, nil
}

func BuildPayer(params *core.PaymentStatementParams) []marotoCore.Row {
	return []marotoCore.Row{
		text.NewRow(14, "PAYMENT FROM", props.Text{Size: 12, Top: 8, Style: fontstyle.Bold}),

		row.New(6).Add(
			text.NewCol(4, "Name", props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, params.Payer.Name, props.Text{Size: 10, Top: 2, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(6).Add(
			text.NewCol(4, "Address", props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, params.Payer.Address, props.Text{Size: 10, Top: 2, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(6).Add(
			text.NewCol(4, "Tax Number", props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, params.Payer.TaxNumber, props.Text{Size: 10, Top: 2, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(6).Add(
			text.NewCol(4, "Contact", props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, params.Payer.Contact, props.Text{Size: 10, Top: 2, Align: align.Right, Style: fontstyle.Bold}),
		),
	}
}

func BuildPayee(params *core.PaymentStatementParams) []marotoCore.Row {
	return []marotoCore.Row{
		text.NewRow(14, "PAYMENT TO", props.Text{Size: 12, Top: 8, Style: fontstyle.Bold}),

		row.New(6).Add(
			text.NewCol(4, "Name", props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, params.Payee.Name, props.Text{Size: 10, Top: 2, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(6).Add(
			text.NewCol(4, "Address", props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, params.Payee.Address, props.Text{Size: 10, Top: 2, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(6).Add(
			text.NewCol(4, "Tax Number", props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, params.Payee.TaxNumber, props.Text{Size: 10, Top: 2, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(6).Add(
			text.NewCol(4, "Contact", props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, params.Payee.Contact, props.Text{Size: 10, Top: 2, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(4),
	}
}

func BuildSummaryRows(params *core.PaymentStatementParams) []marotoCore.Row {
	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	total := decimal.NewFromFloat(0.0)
	totalTax := decimal.NewFromFloat(0.0)
	for _, item := range params.DetailItems {
		tax := item.Amount.Mul(item.WithholdingTaxRate)
		total = total.Add(item.Amount)
		totalTax = totalTax.Add(tax)
	}
	totalWithoutTax := total.Sub(totalTax)

	return []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, "Summary", props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, "Amount", props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(14).Add(
			text.NewCol(8, "Revenue", props.Text{Size: 10, Top: 4, Align: align.Left}),
			text.NewCol(4, fmt.Sprintf("%s %s", total.RoundDown(2), params.Currency), props.Text{Size: 10, Top: 4, Align: align.Right}),
		),
		row.New(10).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, "Withholding Tax", props.Text{Size: 10, Top: 0, Align: align.Left}),
			text.NewCol(6, fmt.Sprintf("-%s %s", totalTax.RoundDown(2), params.Currency), props.Text{Size: 10, Top: 0, Align: align.Right}),
		),
		row.New(20).Add(
			text.NewCol(6, "Payment Amount (excluding tax)", props.Text{Size: 12, Top: 4, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(6, fmt.Sprintf("%s %s", totalWithoutTax.RoundDown(2), params.Currency), props.Text{Size: 12, Top: 4, Align: align.Right, Style: fontstyle.Bold}),
		),
	}
}

func BuildDetailsRows(params *core.PaymentStatementParams) []marotoCore.Row {
	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	rows := []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(4, "Details", props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, "Amount", props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
			text.NewCol(4, "Withholding Tax", props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
	}

	for ix, item := range params.DetailItems {
		paddingTop := float64(0)
		rowHeight := float64(8)
		if ix == 0 {
			paddingTop = float64(4)
			rowHeight = float64(12)
		}
		tax := item.Amount.Mul(item.WithholdingTaxRate)
		netAmount := item.Amount.Sub(tax)
		rows = append(rows, row.New(rowHeight).Add(
			col.New(4).Add(
				text.New(item.Title, props.Text{Size: 10, Top: paddingTop, Align: align.Left}),
			),
			col.New(4).Add(
				text.New(fmt.Sprintf("%s %s", netAmount, params.Currency), props.Text{Size: 10, Top: paddingTop, Align: align.Right}),
			),
			col.New(4).Add(
				text.New(fmt.Sprintf("%s %s", tax, params.Currency), props.Text{Size: 10, Top: paddingTop, Align: align.Right}),
			),
		))

	}
	return rows
}
