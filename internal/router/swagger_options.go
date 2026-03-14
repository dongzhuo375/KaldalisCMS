package router

import "strings"

// SwaggerOptions keeps OpenAPI exposure policy at the router boundary.
type SwaggerOptions struct {
	Enabled     bool
	Path        string
	Title       string
	Version     string
	Description string
}

func normalizeSwaggerPath(path string) string {
	clean := strings.TrimSpace(path)
	if clean == "" {
		return "/swagger"
	}
	if !strings.HasPrefix(clean, "/") {
		clean = "/" + clean
	}
	return strings.TrimRight(clean, "/")
}
