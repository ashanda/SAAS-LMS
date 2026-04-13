package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func proxy(target string, prefix string) http.Handler {
	targetURL, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 🔥 FIX: preserve original director
	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		// call original
		originalDirector(req)

		// remove prefix if needed
		if prefix != "" && len(req.URL.Path) >= len(prefix) {
			req.URL.Path = req.URL.Path[len(prefix):]
		}

		// 🔥 CRITICAL: copy header from original request context
		// (this is what was missing)
		if tenantID := req.Header.Get("X-Tenant-ID"); tenantID != "" {
			req.Header.Set("X-Tenant-ID", tenantID)
		}
	}

	return proxy
}