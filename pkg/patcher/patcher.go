package patcher

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"

	"github.com/sandeep-yadav/nexus-ast/pkg/models"
)

// ApplyFixes automatically rewrites Go source files to eliminate BOLA vulnerabilities
func ApplyFixes(targetDir string, findings []models.Finding) (int, error) {
	patchedCount := 0
	filesToPatch := make(map[string]bool)

	for _, f := range findings {
		if f.SourceRef != "" {
			var fileName string
			fmt.Sscanf(f.SourceRef, "%s", &fileName)
			for i, c := range fileName {
				if c == ':' {
					fileName = fileName[:i]
					break
				}
			}
			filesToPatch[fileName] = true
		}
	}

	reWhere := regexp.MustCompile(`(?m)^(\t+)(db\.Where\("id = \?",\s*([a-zA-Z0-9_]+)\))`)

	for fileName := range filesToPatch {
		var fullPath string
		_ = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && filepath.Base(path) == fileName {
				fullPath = path
			}
			return nil
		})

		if fullPath == "" {
			continue
		}

		contentBytes, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		content := string(contentBytes)
		matches := reWhere.FindAllStringSubmatch(content, -1)
		if len(matches) > 0 {
			content = reWhere.ReplaceAllString(content, "${1}tenantID := c.GetString(\"tenant_id\") // 🤖 Auto-patched by Nexus-AST\n${1}db.Where(\"id = ? AND tenant_id = ?\", ${3}, tenantID)")
			patchedCount += len(matches)

			formatted, err := format.Source([]byte(content))
			if err == nil {
				_ = os.WriteFile(fullPath, formatted, 0644)
			} else {
				_ = os.WriteFile(fullPath, []byte(content), 0644)
			}
		}
	}

	return patchedCount, nil
}
