package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/thechiranjeevvyas/ProEmailHunter/banner"
	"github.com/thechiranjeevvyas/ProEmailHunter/extractor"
)

func main() {
	concurrent := flag.Int("c", 30, "Number of concurrent requests")
	timeout := flag.Int("t", 15, "Request timeout in seconds")
	silent := flag.Bool("silent", false, "Silent mode")
	version := flag.Bool("version", false, "Print version and exit")
	verbose := flag.Bool("verbose", false, "Enable verbose output")

	flag.Parse()

	if *version {
		banner.PrintBanner()
		banner.PrintVersion()
		return
	}

	if !*silent {
		banner.PrintBanner()
	}

	timeoutDuration := time.Duration(*timeout) * time.Second

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domainURL := strings.TrimSpace(scanner.Text())
		if domainURL != "" {
			extractor.ProcessDomain(domainURL, *concurrent, timeoutDuration, *verbose)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}
