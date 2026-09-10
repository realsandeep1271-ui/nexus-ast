package patcher

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sandeep-yadav/nexus-ast/pkg/models"
)

// ApplyFixes automatically rewrites Go source files to eliminate BOLA vulnerabilities
func ApplyFixes(targetDir string, findings []models.Finding) (int, error) {
	patchedCount := 0
	filesToPatch := make(map[string][]models.Finding)

	for _, f := range findings {
		if f.SourceRef != "" {
			parts := strings.Split(f.SourceRef, ":")
			fileName := parts[0]
			filesToPatch[fileName] = append(filesToPatch[fileName], f)
		}
	}

	for fileName, fileFindings := range filesToPatch {
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
		originalContent := content

		for _, f := range fileFindings {
			reParam := regexp.MustCompile(`'([a-zA-Z0-9_]+)'`)
			matches := reParam.FindStringSubmatch(f.Description)
			paramName := "objectID"
			if len(matches) > 1 {
				paramName = matches[1]
			}

			oldWhere := fmt.Sprintf(`Where("id = ?", %s)`, paramName)
			newWhere := fmt.Sprintf(`Where("id = ? AND tenant_id = ?", %s, tenantID)`, paramName)

			if strings.Contains(content, oldWhere) {
				injection := fmt.Sprintf("\n\ttenantID := c.GetString(\"tenant_id\") // 🤖 Auto-patched by Nexus-AST\n\tdb.", paramName)
				content = strings.Replace(content, oldWhere, newWhere, 1)
				content = strings.Replace(content, "\tdb.", injection, 1)
				patchedCount++
			}
		}

		if content != originalContent {
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
