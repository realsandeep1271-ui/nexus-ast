package report

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sandeep-yadav/nexus-ast/pkg/models"
)

var (
	cyanBold   = color.New(color.FgCyan, color.Bold)
	redBold    = color.New(color.FgRed, color.Bold)
	yellowBold = color.New(color.FgYellow, color.Bold)
	greenBold  = color.New(color.FgGreen, color.Bold)
	whiteBold  = color.New(color.FgWhite, color.Bold)
)

func PrintBanner() {
	banner := `
  _  _                     _   ___ _____ 
 | \| |_____ ___  _ ___   /_\ / __|_   _|
 | .` + "`" + ` / -_) \ / || (_-<  / _ \\__ \ | |  
 |_|\_\___/_\_\\_,_/__/ /_/ \_\___/ |_|  
 [ Autonomous Authorization & Policy Linter v1.0 ]
`
	cyanBold.Println(banner)
}

func PrintReport(report models.AuditReport) {
	whiteBold.Printf("\nAudit Target: %s\n", report.Target)
	whiteBold.Printf("Endpoints Evaluated: %d\n", report.EndpointsScanned)
	whiteBold.Printf("Policy Violations Identified: %d\n\n", len(report.Findings))

	if len(report.Findings) == 0 {
		greenBold.Println("✔ Zero authorization flaws detected. System satisfies object isolation policies.")
		return
	}

	for i, f := range report.Findings {
		fmt.Printf("--------------------------------------------------------------------------------\n")
		switch f.Severity {
		case models.SeverityCritical:
			redBold.Printf("[%d] %s: %s\n", i+1, f.Severity, f.Title)
		case models.SeverityHigh:
			yellowBold.Printf("[%d] %s: %s\n", i+1, f.Severity, f.Title)
		default:
			cyanBold.Printf("[%d] %s: %s\n", i+1, f.Severity, f.Title)
		}

		if f.SourceRef != "" {
			color.White("  Source:      %s\n", f.SourceRef)
		}
		color.White("  Endpoint:    [%s] %s\n", f.Method, f.Endpoint)
		color.White("  Description: %s\n", f.Description)
		color.Green("  Remediation: %s\n", f.Remediation)
	}
	fmt.Printf("--------------------------------------------------------------------------------\n\n")
}

func SaveJSON(report models.AuditReport, outputPath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, data, 0644)
}
