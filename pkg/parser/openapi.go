package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sandeep-yadav/nexus-ast/pkg/models"
)

type OpenAPISpec struct {
	Paths map[string]map[string]struct {
		Summary    string   `json:"summary"`
		Security   []any    `json:"security"`
		Parameters []struct {
			Name string `json:"name"`
			In   string `json:"in"`
		} `json:"parameters"`
	} `json:"paths"`
}

func ParseSpec(filePath string) ([]models.Endpoint, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file: %w", err)
	}

	var spec OpenAPISpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse JSON spec: %w", err)
	}

	var endpoints []models.Endpoint
	paramRegex := regexp.MustCompile(`\{([^}]+)\}`)

	for path, methods := range spec.Paths {
		for method, details := range methods {
			m := strings.ToUpper(method)
			if m != "GET" && m != "POST" && m != "PUT" && m != "DELETE" && m != "PATCH" {
				continue
			}

			matches := paramRegex.FindAllStringSubmatch(path, -1)
			var params []string
			for _, match := range matches {
				if len(match) > 1 {
					params = append(params, match[1])
				}
			}

			endpoints = append(endpoints, models.Endpoint{
				Path:       path,
				Method:     m,
				Parameters: params,
				HasAuth:    len(details.Security) > 0,
			})
		}
	}

	return endpoints, nil
}
