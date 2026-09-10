package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/sandeep-yadav/nexus-ast/pkg/ast"
	"github.com/sandeep-yadav/nexus-ast/pkg/models"
	"github.com/sandeep-yadav/nexus-ast/pkg/parser"
	"github.com/sandeep-yadav/nexus-ast/pkg/report"
	"github.com/sandeep-yadav/nexus-ast/pkg/verifier"
)

var (
	repoPath   string
	specPath   string
	targetURL  string
	tokenA     string
	tokenB     string
	testObjID  string
	outputPath string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Execute static AST analysis or dynamic verification",
	Run: func(cmd *cobra.Command, args []string) {
		report.PrintBanner()

		auditReport := models.AuditReport{
			Target:    repoPath,
			Timestamp: time.Now(),
		}

		if repoPath != "" {
			auditReport.Target = repoPath
			findings, err := ast.InspectDirectory(repoPath)
			if err == nil {
				auditReport.Findings = append(auditReport.Findings, findings...)
				auditReport.EndpointsScanned = len(findings)
			}
		}

		if specPath != "" && targetURL != "" && tokenB != "" {
			auditReport.Target = targetURL
			endpoints, err := parser.ParseSpec(specPath)
			if err == nil {
				auditReport.EndpointsScanned += len(endpoints)
				v := verifier.NewDualContextVerifier(targetURL, tokenA, tokenB)
				for _, ep := range endpoints {
					finding := v.VerifyEndpoint(ep, testObjID)
					if finding != nil {
						auditReport.Findings = append(auditReport.Findings, *finding)
					}
				}
			}
		}

		report.PrintReport(auditReport)

		if outputPath != "" {
			_ = report.SaveJSON(auditReport, outputPath)
		}
	},
}

func init() {
	scanCmd.Flags().StringVarP(&repoPath, "repo", "r", "", "Path to source code repository to analyze AST")
	scanCmd.Flags().StringVarP(&specPath, "spec", "s", "", "Path to OpenAPI/Swagger JSON specification")
	scanCmd.Flags().StringVarP(&targetURL, "url", "u", "", "Live API base URL (e.g. https://api.target.com)")
	scanCmd.Flags().StringVar(&tokenA, "token-a", "", "Context A Bearer Token")
	scanCmd.Flags().StringVar(&tokenB, "token-b", "", "Context B Bearer Token")
	scanCmd.Flags().StringVar(&testObjID, "obj-id", "1001", "Sample object identifier")
	scanCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Save JSON audit findings")
}
