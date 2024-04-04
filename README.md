# Invoice Generator

A simple invoice generator.

## Usage

```go
package main

import (
  "os"
  "log"
  "github.com/quail-ink/quail-invoice"
)

func main() {
	builder, err := invoice.NewBuilder(
		Config{
			SealImage:      "./sample-seal.png",
			FontName:       "noto-sans-cjk",
			FontNormal:     "./fonts/NotoSansCJK-JP/NotoSansCJKjp-Regular.ttf",
			FontItalic:     "./fonts/NotoSansCJK-JP/NotoSansCJKjp-Italic.ttf",
			FontBold:       "./fonts/NotoSansCJK-JP/NotoSansCJKjp-Bold.ttf",
			FontBoldItalic: "./fonts/NotoSansCJK-JP/NotoSansCJKjp-BoldItalic.ttf",
		},
		"./sample-params.yaml")
	if err != nil {
		log.Fatal("failed to create builder")
		return
	}

	buf, err := builder.GenerateInvoice()
	if buf == nil || err != nil {
		log.Fatal("failed to generate invoice")
		return
	}

	filename := "sample-invoice.pdf"
	if err := os.WriteFile(filename, buf, 0666); err != nil {
		log.Fatal("failed to write to file")
		return
	}
}
```