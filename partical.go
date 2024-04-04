package invoice

import (
	"fmt"
	"io"
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
	"github.com/sirupsen/logrus"
)

func buildInvoiceHeader(params *InvoiceParams, sealfile string) ([]marotoCore.Row, error) {
	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	fd, err := os.Open(sealfile)
	if err != nil {
		logrus.WithError(err).Error("failed to open seal file")
		return nil, err
	}
	defer fd.Close()
	buf, err := io.ReadAll(fd)
	if err != nil {
		logrus.WithError(err).Error("failed to read seal file")
		return nil, err
	}

	leftCol := col.New(6)
	leftCol.Add(image.NewFromBytes(buf, extension.Png, props.Rect{
		Center:  false,
		Percent: 20,
		Left:    34,
		Top:     7,
	}))
	leftCol.Add(text.New(params.CompanyName, props.Text{Size: 14, Top: 8, Align: align.Left, Style: fontstyle.Bold}))
	lines := strings.Split(params.CompanyAddr, "\n")
	for ix, line := range lines {
		leftCol.Add(text.New(line, props.Text{Size: 10, Top: float64(6*ix + 16), Align: align.Left}))
	}
	leftCol.Add(text.New(params.CompanyEmail, props.Text{Size: 10, Top: float64(6*(len(lines)) + 16), Align: align.Left}))

	rs := row.New(40).WithStyle(borderBottomStyle).Add(
		leftCol,
		col.New(6).Add(
			text.New(fmt.Sprintf("Tax ID: %s", params.TaxNumber), props.Text{Size: 10, Top: 16, Align: align.Left}),
			text.New(fmt.Sprintf("Invoice Issue Date: %s", params.Date.Format("2006/01/02")), props.Text{Size: 10, Top: 22, Align: align.Left}),
			text.New(fmt.Sprintf("Invoice Period: %s - %s",
				params.Summary.PeriodStart.Format("2006/01/02"),
				params.Summary.PeriodEnd.Format("2006/01/02"),
			), props.Text{Size: 10, Top: 28}),
		),
	)

	rows := []marotoCore.Row{
		rs,
		row.New(10),
	}
	return rows, nil
}

func buildInvoicePaymentRows(params *InvoiceParams) []marotoCore.Row {
	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	rows := []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, "Payment Instructions", props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, "", props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
	}

	rows = append(rows, row.New(12).Add(
		col.New(2).Add(
			text.New("Bank Name", props.Text{Size: 10, Top: 4, Align: align.Left}),
		),
		col.New(10).Add(
			text.New(params.Payment.ReceiveAccountBank, props.Text{Size: 10, Top: 4, Align: align.Right}),
		),
	))
	rows = append(rows, row.New(8).Add(
		col.New(2).Add(
			text.New("Bank Branch", props.Text{Size: 10, Top: 0, Align: align.Left}),
		),
		col.New(10).Add(
			text.New(params.Payment.ReceiveAccountBranch, props.Text{Size: 10, Top: 0, Align: align.Right}),
		),
	))
	rows = append(rows, row.New(8).Add(
		col.New(2).Add(
			text.New("Bank Account", props.Text{Size: 10, Top: 0, Align: align.Left}),
		),
		col.New(10).Add(
			text.New(params.Payment.ReceiveAccountNumber, props.Text{Size: 10, Top: 0, Align: align.Right}),
		),
	))
	return rows
}

func buildInvoiceDetailsRows(params *InvoiceParams) []marotoCore.Row {
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
			text.NewCol(8, "Details", props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, "", props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
	}

	for ix, item := range params.DetailItems {
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
		}
	}
	return rows
}

func buildInvoiceSummaryRows(params *InvoiceParams) []marotoCore.Row {
	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	subtotal := params.Summary.TotalExcludeTax
	tax := subtotal.Mul(params.Summary.TaxRate).RoundDown(2)
	total := subtotal.Add(tax).RoundDown(2)

	return []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, "Summary", props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, "Amount", props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(14).Add(
			text.NewCol(8, params.Summary.Title, props.Text{Size: 10, Top: 4, Align: align.Left}),
			text.NewCol(4, fmt.Sprintf("%s %s", subtotal.RoundDown(2), params.Currency), props.Text{Size: 10, Top: 4, Align: align.Right}),
		),
		row.New(10).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, "Tax", props.Text{Size: 10, Top: 0, Align: align.Left}),
			text.NewCol(6, fmt.Sprintf("%s %s", tax, params.Currency), props.Text{Size: 10, Top: 0, Align: align.Right}),
		),
		row.New(20).Add(
			text.NewCol(6, "Total (including tax)", props.Text{Size: 12, Top: 4, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(6, fmt.Sprintf("%s %s", total, params.Currency), props.Text{Size: 12, Top: 4, Align: align.Right, Style: fontstyle.Bold}),
		),
	}
}
