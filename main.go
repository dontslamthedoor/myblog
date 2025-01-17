package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// TemplateCache stores parsed templates
var TemplateCache *template.Template

func main() {
	// Load templates from the views directory
	TemplateCache = parseTemplates()

	// Initialize the router
	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Log each request
	r.Use(middleware.Recoverer) // Recover from panics

	// Serve static files
	fileServer(r, "/static", "./static")
	fileServer(r, "/static/assets", "./node_modules/@fortawesome/fontawesome-free")

	// Routes
	r.Get("/", renderPage("home"))
	r.Get("/blog", renderPage("Blog"))
	r.Get("/projects", renderPage("Projects"))

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// parseTemplates parses all templates in the views directory
func parseTemplates() *template.Template {
	// Parse all template files (layout and components)
	tmpl, err := template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
	return tmpl
}

// renderPage renders a page by injecting it into the layout template
func renderPage(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := TemplateCache.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			"Content": page + ".html",
		})
		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			log.Printf("Error rendering template %s: %v\n", page, err)
		}
	}
}

// fileServer sets up a static file server for serving files from a directory
func fileServer(r chi.Router, path string, root string) {
	fs := http.StripPrefix(path, http.FileServer(http.Dir(root)))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	r.Get(path+"*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
