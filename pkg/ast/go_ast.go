package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/sandeep-yadav/nexus-ast/pkg/models"
)

func InspectDirectory(targetDir string) ([]models.Finding, error) {
	var findings []models.Finding
	fset := token.NewFileSet()

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil
		}

		fileFindings := inspectASTNode(fset, node, path)
		findings = append(findings, fileFindings...)
		return nil
	})

	return findings, err
}

func inspectASTNode(fset *token.FileSet, node *ast.File, filePath string) []models.Finding {
	var findings []models.Finding

	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		var extractedParams []string
		hasTenantVerification := false
		hasDirectQuery := false

		ast.Inspect(funcDecl.Body, func(inner ast.Node) bool {
			callExpr, ok := inner.(*ast.CallExpr)
			if !ok {
				return true
			}

			if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				methodName := sel.Sel.Name
				if methodName == "Param" || methodName == "PathValue" || methodName == "Query" {
					if len(callExpr.Args) > 0 {
						if lit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
							paramName := strings.Trim(lit.Value, "\"")
							if isSensitiveIdentifier(paramName) {
								extractedParams = append(extractedParams, paramName)
							}
						}
					}
				}

				if methodName == "Where" || methodName == "First" || methodName == "Find" || methodName == "Delete" {
					hasDirectQuery = true
				}

				if methodName == "Get" || methodName == "Value" || methodName == "MustGet" {
					hasTenantVerification = true
				}
			}
			return true
		})

		if len(extractedParams) > 0 && hasDirectQuery && !hasTenantVerification {
			pos := fset.Position(funcDecl.Pos())
			for _, param := range extractedParams {
				findings = append(findings, models.Finding{
					Type:        models.FindingTypeUnboundObjectID,
					Severity:    models.SeverityCritical,
					Title:       fmt.Sprintf("Authorization Policy Violation: Unbound Identifier in %s", funcDecl.Name.Name),
					Description: fmt.Sprintf("Controller extracts sensitive object identifier '%s' from URL parameters and queries database without asserting session tenant boundaries.", param),
					Endpoint:    fmt.Sprintf("Handler: %s", funcDecl.Name.Name),
					Method:      "STATIC_AST",
					SourceRef:   fmt.Sprintf("%s:%d", filepath.Base(filePath), pos.Line),
					Remediation: fmt.Sprintf("Ensure queries include authenticated tenant context: db.Where(\"id = ? AND tenant_id = ?\", %s, session.TenantID).", param),
				})
			}
		}

		return true
	})

	return findings
}

func isSensitiveIdentifier(param string) bool {
	p := strings.ToLower(param)
	identifiers := []string{"id", "user_id", "org_id", "tenant_id", "account_id", "invoice_id", "doc_id", "uuid"}
	for _, id := range identifiers {
		if p == id || strings.HasSuffix(p, "_id") {
			return true
		}
	}
	return false
}
