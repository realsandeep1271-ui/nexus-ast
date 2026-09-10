package testdata

type RequestContext struct{}

func (c *RequestContext) GetString(k string) string { return "tenant_1" }
func (c *RequestContext) Param(key string) string   { return "" }
func (c *RequestContext) JSON(code int, obj any)    {}

type DatabaseSession struct{}

func (d *DatabaseSession) Where(query any, args ...any) *DatabaseSession { return d }
func (d *DatabaseSession) First(dest any) *DatabaseSession               { return d }

type InvoiceRecord struct {
	ID       string
	TenantID string
	Amount   int
}

// GetInvoiceHandler processes invoice retrieval
func GetInvoiceHandler(c *RequestContext, db *DatabaseSession) {
	invoiceID := c.Param("invoice_id")
	var record InvoiceRecord
	// Note: Direct query on object ID without tenant verification
	tenantID := c.GetString("tenant_id") // 🤖 Auto-patched by Nexus-AST
	db.Where("id = ? AND tenant_id = ?", invoiceID, tenantID).First(&record)
	c.JSON(200, record)
}
