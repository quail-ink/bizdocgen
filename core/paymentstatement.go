package core

import (
	"os"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type (
	PaymentStatementPayer struct {
		Name      string `yaml:"name"`
		Address   string `yaml:"addr"`
		TaxNumber string `yaml:"tax_number"`
		Contact   string `yaml:"contact"`
	}

	PaymentStatementPayee struct {
		Name      string `yaml:"name"`
		Address   string `yaml:"addr"`
		TaxNumber string `yaml:"tax_number"`
		Contact   string `yaml:"contact"`
	}

	PaymentStatementDetailItem struct {
		Title              string          `yaml:"title"`
		Desc               string          `yaml:"desc"`
		Amount             decimal.Decimal `yaml:"amount"`
		WithholdingTaxRate decimal.Decimal `yaml:"withholding_tax_rate"`
	}

	PaymentStatementParams struct {
		ID          string    `yaml:"id"`
		Date        time.Time `yaml:"date" time_format:"2006/01/02"`
		Currency    string    `yaml:"currency"`
		CompanySeal string    `yaml:"company_seal"`
		PeriodStart time.Time `yaml:"period_start" time_format:"2006/01/02"`
		PeriodEnd   time.Time `yaml:"period_end" time_format:"2006/01/02"`

		PaymentChannel string `yaml:"payment_channel"`
		PaymentTxID    string `yaml:"payment_tx_id"`

		Payer       PaymentStatementPayer        `yaml:"payer"`
		Payee       PaymentStatementPayee        `yaml:"payee"`
		DetailItems []PaymentStatementDetailItem `yaml:"detail_items"`
	}
)

func (params *PaymentStatementParams) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		logrus.WithError(err).Fatalf("failed to read YAML file")
		return err
	}

	if err := yaml.Unmarshal(data, params); err != nil {
		logrus.WithError(err).Fatalf("failed to unmarshal YAML")
		return err
	}

	return nil
}
