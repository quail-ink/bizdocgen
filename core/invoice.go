package core

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
		URLs            []string        `yaml:"urls"`
		TotalExcludeTax decimal.Decimal `yaml:"total_exclude_tax"`
	}

	InvoiceSummary struct {
		PeriodStart     time.Time       `yaml:"period_start" time_format:"2006/01/02"`
		PeriodEnd       time.Time       `yaml:"period_end" time_format:"2006/01/02"`
		Title           string          `yaml:"title"`
		TotalExcludeTax decimal.Decimal `yaml:"total_exclude_tax"`
		TotalIncludeTax decimal.Decimal `yaml:"total_include_tax"`
		TaxRate         decimal.Decimal `yaml:"tax_rate"`
	}

	InvoicePayment struct {
		Disabled              bool   `yaml:"show"`
		PaymentID             string `yaml:"payment_id"`
		Method                string `yaml:"method"`
		ReceiveAccountBank    string `yaml:"receive_account_bank"`
		ReceiveAccountBranch  string `yaml:"receive_account_branch"`
		ReceiveDepositType    string `yaml:"receive_deposit_type"`
		ReceiveAccountNumber  string `yaml:"receive_account_number"`
		ReceiveAccountName    string `yaml:"receive_account_name"`
		ReceiveAccountRouting string `yaml:"receive_account_routing"`
		ReceiveAccountSwift   string `yaml:"receive_account_swift"`
	}

	InvoiceParams struct {
		ID           string    `yaml:"id"`
		TaxNumber    string    `yaml:"tax_number"`
		Date         time.Time `yaml:"date" time_format:"2006/01/02"`
		Currency     string    `yaml:"currency"`
		CompanyName  string    `yaml:"company_name"`
		CompanyAddr  string    `yaml:"company_address"`
		CompanyEmail string    `yaml:"company_email"`
		CompanySeal  string    `yaml:"company_seal"`

		BillToCompany string `yaml:"bill_to_company"`
		BillToAddress string `yaml:"bill_to_address"`

		// Summary
		Summary InvoiceSummary `yaml:"summary"`

		// Details
		DetailItems []InvoiceDetailItem `yaml:"detail_items"`

		// Payment Instructions
		Payment InvoicePayment `yaml:"payment"`
	}
)

func (params *InvoiceParams) Load(filename string) error {
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
