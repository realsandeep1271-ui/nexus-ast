package testdata

type ContextWrapper struct{}

func (c *ContextWrapper) GetString(k string) string { return "tenant_1" }
func (c *ContextWrapper) Param(key string) string   { return "" }
func (c *ContextWrapper) JSON(code int, obj any)    {}

type QueryEngine struct{}

func (q *QueryEngine) Where(query any, args ...any) *QueryEngine { return q }
func (q *QueryEngine) First(dest any) *QueryEngine               { return q }
func (q *QueryEngine) Delete(dest any) *QueryEngine              { return q }

type VehicleLocation struct {
	ID        string
	TenantID  string
	Latitude  float64
	Longitude float64
}

// GetVehicleLocationHandler has BOLA: extracts vehicle_id and queries without tenant binding
func GetVehicleLocationHandler(c *ContextWrapper, db *QueryEngine) {
	vehicleID := c.Param("vehicle_id")
	var loc VehicleLocation
	// UNBOUND QUERY: Developer forgot to check session tenant/owner!
	tenantID := c.GetString("tenant_id") // 🤖 Auto-patched by Nexus-AST
	db.Where("id = ? AND tenant_id = ?", vehicleID, tenantID).First(&loc)
	c.JSON(200, loc)
}

// DeleteVehicleHandler has BOLA: accepts vehicle_id directly into Delete statement
func DeleteVehicleHandler(c *ContextWrapper, db *QueryEngine) {
	vehicleID := c.Param("vehicle_id")
	var loc VehicleLocation
	// UNBOUND DELETION
	tenantID := c.GetString("tenant_id") // 🤖 Auto-patched by Nexus-AST
	db.Where("id = ? AND tenant_id = ?", vehicleID, tenantID).Delete(&loc)
	c.JSON(200, "deleted")
}
