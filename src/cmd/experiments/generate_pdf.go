package main

import (
	"bytes"
	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/ptr"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/oklog/ulid/v2"
	"github.com/playwright-community/playwright-go"
)

const BAGE_PAGE_HTML = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>PDF Go generation experiment</title>
</head>
<body>
  <h1>Test PDF</h1>
  <p>This PDF was generated using Go + Playwright</p>
</body>

</html>
`

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

	l.Info("Setting the HTML content of the page")
	if err = page.SetContent(BAGE_PAGE_HTML, playwright.PageSetContentOptions{
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
