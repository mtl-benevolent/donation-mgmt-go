package main

import (
	"bytes"
	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/ptr"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/playwright-community/playwright-go"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const HTML_TEMPLATE = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>PDF test</title>
</head>
<body>
  <h1>{{ .Title }}</h1>
  <p>{{ .Content }}</p>
	<p>Money in fr-CA: {{ money .ReceiptAmountInCents "fr-CA" "CAD" }}
	<p>Money in en-CA: {{ money .ReceiptAmountInCents "en-CA" "CAD" }}
	<p>Generated on {{ .CreatedAt.Format "2006-01-02"  }}</p>
</body>

</html>
`

type TemplateValues struct {
	Title                string
	Content              string
	ReceiptAmountInCents int64
	CreatedAt            time.Time
}

func FormatMoney(cents int64, localeTag string, currencyCode string) (string, error) {
	locale, err := language.Parse(localeTag)
	if err != nil {
		return "", fmt.Errorf("invalid locale: %w", err)
	}

	unit, err := currency.ParseISO(currencyCode)
	if err != nil {
		return "", fmt.Errorf("invalid currency code (%s): %w", currencyCode, err)
	}

	p := message.NewPrinter(locale)

	amount := float64(cents) / 100
	result := p.Sprintf("%v", currency.NarrowSymbol(unit.Amount(amount)))

	return result, nil
}

func getHtml(value TemplateValues) (string, error) {
	tmpl, err := template.New("pdf-receipt").Funcs(template.FuncMap{
		"money": FormatMoney,
	}).Parse(HTML_TEMPLATE)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	sb := &strings.Builder{}
	if err := tmpl.Execute(sb, value); err != nil {
		return "", fmt.Errorf("error rendering template: %w", err)
	}

	return sb.String(), nil
}

func exec(l *slog.Logger) error {
	pw, err := playwright.Run(&playwright.RunOptions{
		SkipInstallBrowsers: true,
	})

	if err != nil {
		return fmt.Errorf("error launching Playwright: %w", err)
	}

	defer func() {
		err := pw.Stop()
		if err != nil {
			l.Error("error shutting down Playwright", slog.Any("err", err))
		}
	}()

	l.Info("launching Chromium")
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless:        ptr.Wrap(true),
		ChromiumSandbox: ptr.Wrap(false),
		Args: []string{
			"--no-sandbox",
			"--disable-setuid-sandbox",
			"--disable-dev-shm-usage",
			"--disable-gpu",
			"--single-process",
		},
	})
	if err != nil {
		return fmt.Errorf("error launching browser: %w", err)
	}

	defer func() {
		if err := browser.Close(); err != nil {
			l.Error("error shutting down browser", slog.Any("err", err))
		}
	}()

	l.Info("Creating new page")
	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("error creating new page: %w", err)
	}

	l.Info("Rendering the template")
	htmlContent, err := getHtml(TemplateValues{
		Title:                "My awesome PDF",
		Content:              "<script>alert('hack attempt');</script>",
		ReceiptAmountInCents: 5996,
		CreatedAt:            time.Now(),
	})
	if err != nil {
		return err
	}

	l.Info("Setting the HTML content of the page")
	if err = page.SetContent(htmlContent, playwright.PageSetContentOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		return fmt.Errorf("error setting the HTML content for the page: %w", err)
	}

	l.Info("Exporting the HTML content as a PDF")
	content, err := page.PDF(playwright.PagePdfOptions{
		Format:          ptr.Wrap("Letter"),
		PrintBackground: ptr.Wrap(true),
		HeaderTemplate:  nil,
		FooterTemplate:  nil,
		Tagged:          ptr.Wrap(true),
	})
	if err != nil {
		return fmt.Errorf("error generating PDF: %w", err)
	}

	fileName := fmt.Sprintf("/tmp/pdfs/%s.pdf", ulid.Make())
	l.Info("Writing PDF content to file", slog.String("path", fileName))

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %w", fileName, err)
	}

	size, err := io.Copy(file, bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("error writing PDF content to disk: %w", err)
	}

	l.Info("Wrote PDF file", slog.Any("path", fileName), slog.Int64("bytes", size))
	return nil
}

func main() {
	logger.BootstrapLogger(&config.AppConfiguration{
		LogLevel:     "INFO",
		LogAddSource: true,
	})
	logger := logger.ForComponent("generate_pdf")
	err := exec(logger)

	if err != nil {
		logger.Error("Error while executing program", slog.Any("err", err))
		os.Exit(255)
	}
}
