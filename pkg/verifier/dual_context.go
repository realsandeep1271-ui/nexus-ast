package verifier

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sandeep-yadav/nexus-ast/pkg/models"
)

type DualContextVerifier struct {
	client  *http.Client
	baseURL string
	tokenA  string
	tokenB  string
}

func NewDualContextVerifier(baseURL, tokenA, tokenB string) *DualContextVerifier {
	return &DualContextVerifier{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
		tokenA:  tokenA,
		tokenB:  tokenB,
	}
}

func (v *DualContextVerifier) VerifyEndpoint(ep models.Endpoint, testObjectID string) *models.Finding {
	if len(ep.Parameters) == 0 {
		return nil
	}

	targetPath := ep.Path
	for _, param := range ep.Parameters {
		targetPath = strings.Replace(targetPath, "{"+param+"}", testObjectID, 1)
	}

	targetURL := v.baseURL + targetPath
	req, err := http.NewRequest(ep.Method, targetURL, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+v.tokenB)
	req.Header.Set("User-Agent", "Nexus-AST-Verifier/1.0")

	resp, err := v.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return &models.Finding{
			Type:        models.FindingTypeCrossTenantAccess,
			Severity:    models.SeverityCritical,
			Title:       fmt.Sprintf("Cross-Tenant Boundary Violation on %s %s", ep.Method, ep.Path),
			Description: fmt.Sprintf("Context B accessed Context A resource identifier '%s' with HTTP status %d.", testObjectID, resp.StatusCode),
			Endpoint:    ep.Path,
			Method:      ep.Method,
			Remediation: "Enforce strict Object-Level Access Control (OLAC) against tenant session claims.",
		}
	}

	return nil
}
