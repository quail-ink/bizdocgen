package invoice

import (
	"os"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type (
	InvoiceDetailItem struct {
		Date            time.Time       `yaml:"date" time_format:"2006/01/02"`
		Title           string          `yaml:"title"`
		Desc            string          `yaml:"desc"`
		URL             string          `yaml:"url"`
		TotalExcludeTax decimal.Decimal `yaml:"total_exclude_tax"`
	}

	InvoiceSummary struct {
		PeriodStart     time.Time       `yaml:"period_start" time_format:"2006/01/02"`
		PeriodEnd       time.Time       `yaml:"period_end" time_format:"2006/01/02"`
		Title           string          `yaml:"title"`
		TotalExcludeTax decimal.Decimal `yaml:"total_exclude_tax"`
		TaxRate         decimal.Decimal `yaml:"tax_rate"`
	}

	InvoicePayment struct {
		ReceiveAccountBank   string `yaml:"receive_account_bank"`
		ReceiveAccountBranch string `yaml:"receive_account_branch"`
		ReceiveAccountNumber string `yaml:"receive_account_number"`
	}

	InvoiceParams struct {
		ID           string    `yaml:"id"`
		TaxNumber    string    `yaml:"tax_number"`
		Date         time.Time `yaml:"date" time_format:"2006/01/02"`
		Currency     string    `yaml:"currency"`
		CompanyName  string    `yaml:"company_name"`
		CompanyAddr  string    `yaml:"company_addr"`
		CompanyEmail string    `yaml:"company_email"`

		BillTo string `yaml:"bill_to"`

		// Summary
		Summary InvoiceSummary `yaml:"summary"`

		// Details
		DetailItems []InvoiceDetailItem `yaml:"detail_items"`

		// Payment Instructions
		Payment InvoicePayment `yaml:"payment"`
	}
)

func (pa *InvoiceParams) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		logrus.WithError(err).Fatalf("failed to read YAML file")
		return err
	}

	if err := yaml.Unmarshal(data, pa); err != nil {
		logrus.WithError(err).Fatalf("failed to unmarshal YAML")
		return err
	}

	return nil
}
