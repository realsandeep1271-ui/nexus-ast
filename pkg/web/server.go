package web

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed dashboard.html
var dashboardHTML []byte

// StartServer starts the high-speed Nexus-AST web visualizer
func StartServer(port int) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(dashboardHTML)
	})

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("[+] Nexus-AST Visual Dashboard running at: http://localhost:%d\n", port)
	fmt.Printf("[+] Press Ctrl+C to terminate server.\n")
	return http.ListenAndServe(addr, nil)
}
