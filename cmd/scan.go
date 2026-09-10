package cmd

import (
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sandeep-yadav/nexus-ast/pkg/ast"
	"github.com/sandeep-yadav/nexus-ast/pkg/models"
	"github.com/sandeep-yadav/nexus-ast/pkg/parser"
	"github.com/sandeep-yadav/nexus-ast/pkg/patcher"
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
	autoFix    bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Execute static AST analysis or dynamic verification with automated remediation",
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

		if autoFix && len(auditReport.Findings) > 0 && repoPath != "" {
			color.Cyan("\n[⚡] --fix flag detected: Initiating Nexus-AST Autonomous Auto-Patcher...")
			count, err := patcher.ApplyFixes(repoPath, auditReport.Findings)
			if err == nil && count > 0 {
				color.Green("[✔] Successfully applied %d automated security patches to source code!", count)
				color.Cyan("[+] Re-verifying codebase post-remediation...")

				postFindings, _ := ast.InspectDirectory(repoPath)
				if len(postFindings) == 0 {
					color.Green("[✔] VERIFICATION PASSED: All BOLA vulnerabilities successfully mitigated! 0 flaws remain.\n")
				}
			}
		}

		if outputPath != "" {
			_ = report.SaveJSON(auditReport, outputPath)
		}
	},
}

func init() {
	scanCmd.Flags().StringVarP(&repoPath, "repo", "r", "", "Path to source code repository to analyze AST")
	scanCmd.Flags().StringVarP(&specPath, "spec", "s", "", "Path to OpenAPI/Swagger JSON specification")
	scanCmd.Flags().StringVarP(&targetURL, "url", "u", "", "Live API base URL")
	scanCmd.Flags().StringVar(&tokenA, "token-a", "", "Context A Bearer Token")
	scanCmd.Flags().StringVar(&tokenB, "token-b", "", "Context B Bearer Token")
	scanCmd.Flags().StringVar(&testObjID, "obj-id", "1001", "Sample object identifier")
	scanCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Save JSON audit findings")
	scanCmd.Flags().BoolVar(&autoFix, "fix", false, "Autonomously patch identified BOLA code vulnerabilities")
}
