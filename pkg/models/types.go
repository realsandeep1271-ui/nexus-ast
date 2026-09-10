package models

import "time"

type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
	SeverityInfo     Severity = "INFO"
)

type FindingType string

const (
	FindingTypeUnboundObjectID    FindingType = "AUTH_UNBOUND_OBJECT_ID"
	FindingTypeCrossTenantAccess  FindingType = "AUTH_CROSS_TENANT_ACCESS"
	FindingTypeMissingAuthContext FindingType = "AUTH_MISSING_CONTEXT"
)

type Endpoint struct {
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	HandlerName string   `json:"handler_name,omitempty"`
	Parameters  []string `json:"parameters"`
	HasAuth     bool     `json:"has_auth"`
	SourceFile  string   `json:"source_file,omitempty"`
	LineNumber  int      `json:"line_number,omitempty"`
}

type Finding struct {
	Type        FindingType `json:"type"`
	Severity    Severity    `json:"severity"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Endpoint    string      `json:"endpoint"`
	Method      string      `json:"method"`
	SourceRef   string      `json:"source_ref,omitempty"`
	Remediation string      `json:"remediation"`
}

type AuditReport struct {
	Target           string    `json:"target"`
	Timestamp        time.Time `json:"timestamp"`
	EndpointsScanned int       `json:"endpoints_scanned"`
	Findings         []Finding `json:"findings"`
}

type TenantContext struct {
	PrimaryToken   string
	SecondaryToken string
}
