package portifolio

import (
	"encoding/json"
	"strings"
)

func buildPrompt(template string, snapshot portfolioSnapshot) (string, error) {
	data, err := json.Marshal(snapshot)
	if err != nil {
		return "", err
	}
	return strings.Replace(template, promptDataPlaceholder, string(data), 1), nil
}
