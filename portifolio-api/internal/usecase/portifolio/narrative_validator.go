package portifolio

import (
	"regexp"
	"strings"
)

var tickerPattern = regexp.MustCompile(`\b[A-Z]{3,5}[0-9]?\b`)

var knownBenchmarks = []string{"IBOV", "BVSP", "SPX", "S&P", "BOVESPA", "NYSE", "NASDAQ", "B3"}

func validateNarrative(text string, allowedTickers []string) bool {
	allowed := make(map[string]bool, len(allowedTickers)+len(knownBenchmarks))
	for _, t := range allowedTickers {
		allowed[strings.ToUpper(t)] = true
	}
	for _, b := range knownBenchmarks {
		allowed[b] = true
	}

	for _, found := range tickerPattern.FindAllString(text, -1) {
		if !allowed[found] {
			return false
		}
	}
	return true
}
