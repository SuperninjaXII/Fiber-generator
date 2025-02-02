package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

// Embed the templates directory
//
//go:embed templates/*
var templateFs embed.FS

var baseDir string

// readTemplateFile reads the content of a template file and replaces {AppName} with the baseDir.
func readTemplateFile(filePath string) (string, error) {
	content, err := templateFs.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", filePath, err)
	}
	contentStr := string(content)
	contentStr = strings.Replace(contentStr, "{AppName}", baseDir, -1) // Ensure replacement takes effect
	return contentStr, nil
}

// runGoModInit runs the 'go mod init' command in the specified baseDir.
func runGoModInit(baseDir string) error {
	cmd := exec.Command("go", "mod", "init", baseDir)
	cmd.Dir = baseDir // Set the working directory to the baseDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error running 'go mod init %s': %s", baseDir, string(output))
	}
	fmt.Printf("Successfully initialized Go module: %s\n", baseDir)
	return nil
}

// generate handles the generation of the project structure, templates, and initializes the Go module.
func generate(ctx *cli.Context) error {
	fmt.Println("generating")

	// Define base directories
	viewDir := filepath.Join(baseDir, "views")
	publicDir := filepath.Join(baseDir, "public")
	publicCSSDir := filepath.Join(publicDir, "css")
	publicJSDir := filepath.Join(publicDir, "js")
	publicLibDir := filepath.Join(publicDir, "lib")
	routesDir := filepath.Join(baseDir, "routes")
	controllersDir := filepath.Join(baseDir, "controllers")

	// Define template paths
	templatePaths := map[string]string{
		filepath.Join(viewDir, "index.html"):                  "templates/html/index.html",
		filepath.Join(publicCSSDir, "style.css"):              "templates/css/style.css",
		filepath.Join(publicJSDir, "app.js"):                  "templates/js/index.js",
		filepath.Join(publicLibDir, "htmx.min.js"):            "templates/lib/htmx.min.js", // Adjust as needed
		filepath.Join(routesDir, "userRoutes.go"):             "templates/go/userRoutes.go",
		filepath.Join(controllersDir, "CreateUserHandler.go"): "templates/go/userHandler.go",
		filepath.Join(baseDir, "app.go"):                      "templates/go/app.go",
		filepath.Join(baseDir, "Makefile"):                    "templates/Makefile",
	}

	// Create directories
	dirs := []string{viewDir, publicCSSDir, publicJSDir, publicLibDir, routesDir, controllersDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	// Create files from templates
	for filePath, templatePath := range templatePaths {
		content, err := readTemplateFile(templatePath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("error creating file %s: %w", filePath, err)
		}
	}

	// Run go mod init
	if err := runGoModInit(baseDir); err != nil {
		return err
	}

	fmt.Println("done generating views")
	fmt.Println("done generating public")
	fmt.Println("done generating routes and controllers")

	return nil
}

func main() {
	app := &cli.App{
		Name:  "fiber-gen",
		Usage: "generate go fiber templates",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Value:       "fiberApp",
				Destination: &baseDir,
				Required:    true,
			},
		},
		Action: generate,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
