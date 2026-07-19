package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// reverseProxy forwards a request to a target service, stripping the URL prefix.
func reverseProxy(target, prefix string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatal("bad target url: ", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		req.Host = targetURL.Host
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	godotenv.Load()

	catalogURL := getEnv("CATALOG_SERVICE_URL", "http://localhost:8081")
	loanURL := getEnv("LOAN_SERVICE_URL", "http://localhost:8082")

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "gateway"})
	})

	router.Any("/api/catalog/*proxyPath", reverseProxy(catalogURL, "/api/catalog"))
	router.Any("/api/loan/*proxyPath", reverseProxy(loanURL, "/api/loan"))

	log.Println("🚪 Gateway running on :8080")
	router.Run(":8080")
}
