package main

import (
	"log"
	"net/http"
	"net/http/httputil" //the reverse-proxxy helper
	"net/url"           //break a URL into parts
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" //load .env files
)

// reverseProxy forwards a request to a target service, stripping the URL prefix.
func reverseProxy(target, prefix string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatal("bad target url: ", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL) //so here the reverse proxy created with a  destination of targetURL so its single host

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		req.Host = targetURL.Host
	}

	return func(c *gin.Context) { //this function is http handler, here is where the request is served
		proxy.ServeHTTP(c.Writer, c.Request) //c.writer to write a response, c.request to get the request
	}
}

func getEnv(key, fallback string) string { //fallback is the default value if the environment variable is missing
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Library System API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function () {
      SwaggerUIBundle({
        urls: [
          { url: "/api/loan/swagger/doc.json", name: "Loan Service" },
          { url: "/api/catalog/swagger/doc.json", name: "Catalog Service" }
        ],
        dom_id: "#swagger-ui",
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`

func main() {
	godotenv.Load()

	catalogURL := getEnv("CATALOG_SERVICE_URL", "http://localhost:8081")
	loanURL := getEnv("LOAN_SERVICE_URL", "http://localhost:8082")

	router := gin.Default()

	router.GET("/swagger", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerHTML))
	})

	router.Any("/api/catalog/*proxyPath", reverseProxy(catalogURL, "/api/catalog"))
	router.Any("/api/loan/*proxyPath", reverseProxy(loanURL, "/api/loan"))

	log.Println("🚪 Gateway running on :8080")
	router.Run(":8080")
}
