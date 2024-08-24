package builder

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
	"github.com/shopspring/decimal"
)

func (b *Builder) BuildPsHeader() ([]marotoCore.Row, error) {
	tTitle := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementTitle", nil)
	tDate := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementIssueDate", nil)
	tPeriod := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPeriod", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}
	leftCol := col.New(6)

	if b.psParams.CompanySeal != "" {
		fd, err := os.Open(b.psParams.CompanySeal)
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
			Percent: 32,
			Left:    0,
			Top:     16,
		}))
	}

	leftCol.Add(text.New(tTitle, props.Text{Size: 14, Top: 8, Align: align.Left, Style: fontstyle.Bold}))

	rs := row.New(28).WithStyle(borderBottomStyle).Add(
		leftCol,
		col.New(6).Add(
			text.New(fmt.Sprintf("%s: %s", tDate, b.psParams.Date.Format("2006/01/02")),
				props.Text{Size: 10, Top: 9, Align: align.Right}),
			text.New(fmt.Sprintf("%s: %s - %s",
				tPeriod,
				b.psParams.PeriodStart.Format("2006/01/02"),
				b.psParams.PeriodEnd.Format("2006/01/02"),
			), props.Text{Size: 10, Top: 16, Align: align.Right}),
		),
	)

	rows := []marotoCore.Row{
		rs,
	}
	return rows, nil
}

func (b *Builder) BuildPsPayer() []marotoCore.Row {
	tPayee := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPayer", nil)
	tName := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserName", nil)
	tAddress := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserAddress", nil)
	tTaxID := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserTaxID", nil)
	tContact := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserContact", nil)

	return []marotoCore.Row{
		text.NewRow(14, tPayee, props.Text{Size: 12, Top: 8, Style: fontstyle.Bold}),

		row.New(6).Add(
			text.NewCol(4, tName, props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, b.psParams.Payer.Name, props.Text{Size: 10, Top: 2, Align: align.Right}),
		),
		row.New(6).Add(
			text.NewCol(4, tAddress, props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, b.psParams.Payer.Address, props.Text{Size: 10, Top: 2, Align: align.Right}),
		),
		row.New(6).Add(
			text.NewCol(4, tTaxID, props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, b.psParams.Payer.TaxNumber, props.Text{Size: 10, Top: 2, Align: align.Right}),
		),
		row.New(6).Add(
			text.NewCol(4, tContact, props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, b.psParams.Payer.Contact, props.Text{Size: 10, Top: 2, Align: align.Right}),
		),
	}
}

func (b *Builder) BuildPsPayee() []marotoCore.Row {
	tPayee := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementPayee", nil)
	tName := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserName", nil)
	tAddress := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserAddress", nil)
	tTaxID := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserTaxID", nil)
	tContact := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementUserContact", nil)

	return []marotoCore.Row{
		text.NewRow(14, tPayee, props.Text{Size: 12, Top: 8, Style: fontstyle.Bold}),
		row.New(6).Add(
			text.NewCol(3, tName, props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(9, b.psParams.Payee.Name, props.Text{Size: 10, Top: 2, Align: align.Right}),
		),
		row.New(6).Add(
			text.NewCol(3, tAddress, props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(9, b.psParams.Payee.Address, props.Text{Size: 10, Top: 2, Align: align.Right}),
		),
		row.New(6).Add(
			text.NewCol(4, tTaxID, props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, b.psParams.Payee.TaxNumber, props.Text{Size: 10, Top: 2, Align: align.Right}),
		),
		row.New(6).Add(
			text.NewCol(4, tContact, props.Text{Size: 10, Top: 2, Align: align.Left}),
			text.NewCol(8, b.psParams.Payee.Contact, props.Text{Size: 10, Top: 2, Align: align.Right}),
		),
		row.New(4),
	}
}

func (b *Builder) BuildPsChannelRows() []marotoCore.Row {
	tChannelTitle := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannelTitle", nil)
	tChannel := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannel", nil)
	tTxID := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementChannelTxID", nil)
	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}
	return []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tChannelTitle, props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
		),
		row.New(10).Add(
			text.NewCol(6, tChannel, props.Text{Size: 10, Top: 4, Align: align.Left}),
			text.NewCol(6, b.psParams.PaymentChannel, props.Text{Size: 10, Top: 4, Align: align.Right}),
		),
		row.New(12).Add(
			text.NewCol(6, tTxID, props.Text{Size: 10, Top: 4, Align: align.Left}),
			text.NewCol(6, b.psParams.PaymentTxID, props.Text{Size: 10, Top: 4, Align: align.Right}),
		),
	}
}

func (b *Builder) BuildPsSummaryRows() []marotoCore.Row {
	tSummary := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummary", nil)
	tSummaryAmount := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummaryAmount", nil)
	tRevenue := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummaryRevenue", nil)
	tWithholdingTax := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementWithholdingTax", nil)
	tNetAmount := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementSummaryNetAmount", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	total := decimal.NewFromFloat(0.0)
	totalTax := decimal.NewFromFloat(0.0)
	for _, item := range b.psParams.DetailItems {
		tax := item.Amount.Mul(item.WithholdingTaxRate)
		total = total.Add(item.Amount)
		totalTax = totalTax.Add(tax)
	}
	totalWithoutTax := total.Sub(totalTax)

	return []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(8, tSummary, props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, tSummaryAmount, props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
		row.New(14).Add(
			text.NewCol(8, tRevenue, props.Text{Size: 10, Top: 4, Align: align.Left}),
			text.NewCol(4, fmt.Sprintf("%s %s", total.Round(b.Round), b.psParams.Currency), props.Text{Size: 10, Top: 4, Align: align.Right}),
		),
		row.New(10).WithStyle(borderBottomStyle).Add(
			text.NewCol(6, tWithholdingTax, props.Text{Size: 10, Top: 0, Align: align.Left}),
			text.NewCol(6, fmt.Sprintf("-%s %s", totalTax.Round(b.Round), b.psParams.Currency), props.Text{Size: 10, Top: 0, Align: align.Right}),
		),
		row.New(16).Add(
			text.NewCol(6, tNetAmount, props.Text{Size: 12, Top: 4, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(6, fmt.Sprintf("%s %s", totalWithoutTax.Round(b.Round), b.psParams.Currency), props.Text{Size: 12, Top: 4, Align: align.Right, Style: fontstyle.Bold}),
		),
	}
}

func (b *Builder) BuildPsDetailsRows() []marotoCore.Row {
	tDetails := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementDetails", nil)
	tAmount := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementDetailsAmount", nil)
	tTax := b.i18nBundle.MusT(b.cfg.Lang, "PaymentStatementDetailsTax", nil)

	borderBottomStyle := &props.Cell{
		BorderType:  border.Bottom,
		BorderColor: &props.Color{Red: 200, Green: 200, Blue: 200},
	}

	rows := []marotoCore.Row{
		row.New(16).WithStyle(borderBottomStyle).Add(
			text.NewCol(4, tDetails, props.Text{Size: 12, Top: 8, Align: align.Left, Style: fontstyle.Bold}),
			text.NewCol(4, tAmount, props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
			text.NewCol(4, tTax, props.Text{Size: 12, Top: 8, Align: align.Right, Style: fontstyle.Bold}),
		),
	}

	for ix, item := range b.psParams.DetailItems {
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
				text.New(fmt.Sprintf("%s %s", netAmount.Round(b.Round), b.psParams.Currency), props.Text{Size: 10, Top: paddingTop, Align: align.Right}),
			),
			col.New(4).Add(
				text.New(fmt.Sprintf("%s %s", tax.Round(b.Round), b.psParams.Currency), props.Text{Size: 10, Top: paddingTop, Align: align.Right}),
			),
		))

	}
	return rows
}
