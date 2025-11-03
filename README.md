# ProEmailHunter

<p align="center">
<a href="#"><img src="https://madewithlove.org.in/badge.svg"></a>
<a href="https://x.com/the_cv_xo"><img src="https://img.shields.io/badge/twitter-%40rix4uni-blue.svg"></a>
<a href="https://github.com/thechiranjeevvyas/ProEmailHunter/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg"></a>
<a href="#"><img src="https://img.shields.io/badge/Made%20with-Bash-1f425f.svg"></a>
<a href="https://github.com/thechiranjeevvyas/ProEmailHunter"><img src="https://img.shields.io/badge/github-%40thechiranjeev-orange"></a>
</p>

> High-speed Go email scraper that crawls sites and internal links concurrently to collect email addresses for reconnaissance, research, or sales intelligence.

## Features

- 🚀 **High-performance concurrent crawling** - Process multiple domains simultaneously
- 🔗 **Deep link traversal** - Automatically follows internal links to discover more emails
- ⚡ **Blazing fast** - Built with Go for maximum speed and efficiency
- 🎯 **Flexible configuration** - Customize concurrency, timeouts, and output verbosity
- 📊 **Multiple input methods** - Single domain, file input, or piped input

## Installation

```bash
git clone --depth 1 https://github.com/thechiranjeevvyas/ProEmailHunter.git
cd emailextractor
go install
```

## Usage

### Command-line Options

```
Usage of emailextractor:
  -c int
        Number of concurrent requests (default 30)
  -silent
        Silent mode
  -t int
        Request timeout in seconds (default 15)
  -verbose
        Enable verbose output
  -version
        Print the version of the tool and exit
```

### Examples

#### Single Domain

```bash
echo "https://www.shopify.com" | proemailhunter
```

#### Multiple Domains from File

**Create a domains file:**

```bash
cat domains.txt
https://www.shopify.com
http://testphp.vulnweb.com
```

**Run the scraper:**

```bash
cat domains.txt | emailextractor
```

#### Advanced: Custom Concurrency & Timeout

```bash
cat domains.txt | emailextractor -c 50 -t 30 --verbose
```

This command processes domains with:

- 50 concurrent requests
- 30-second timeout per request
- Verbose output enabled

## Use Cases

- 🔍 **Reconnaissance** - Security research and penetration testing
- 📈 **Sales Intelligence** - Lead generation and contact discovery
- 📚 **Research** - Academic or market research data collection
- 🎯 **OSINT** - Open-source intelligence gathering

## Notes

- Always ensure you have permission to scrape target websites
- Respect robots.txt and website terms of service
- Use appropriate concurrency levels to avoid overwhelming target servers
- Consider legal and ethical implications of email collection in your jurisdiction

## Author

Created by **[@thechiranjeevvyas](https://github.com/thechiranjeevvyas)**

---

⭐ If you find this tool useful, please consider giving it a star on GitHub!
