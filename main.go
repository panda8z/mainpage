package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed web/*
var content embed.FS

func main() {

	// Use the Sub method to create a sub filesystem for the "web" directory
	webFS, err := fs.Sub(content, "web/static")
	if err != nil {
		fmt.Println("Error accessing embedded filesystem:", err)
		return
	}

	// Create a file server handler using the sub filesystem
	staticFileServer := http.FileServer(http.FS(webFS))

	// Handle requests to "/" by serving the index.html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexHTML, err := content.ReadFile("web/index.html")
		if err != nil {
			http.Error(w, "could not read index.html", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(indexHTML)
		if err != nil {
			http.Error(w, "could not write response", http.StatusInternalServerError)
		}
	})

	// Handle requests to "/static/" by serving static files
	http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	// Start the server
	port := 8080
	fmt.Printf("Server is listening on port %d...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

	e := gin.Default()
	e.GET("/", func(c *gin.Context) {
		indexHTML, err := content.ReadFile("web/index.html")
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Errors = append(c.Errors, c.Error(errors.New("could not read index.html")))
			c.Abort()
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(indexHTML))

	})
	e.StaticFS("/static", http.FS(webFS))
	e.Run(":9090")
}
