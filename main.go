package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/phillip-england/gtml/pkg/gtml"
)

const (
	DirComponents   = "components"
	DirRoutes       = "routes"
	DirDist         = "dist"
	DirStatic       = "static"
	FileStyleCSS    = "styles.css"
	DirPreinstalled = "spec/components/preinstalled_components"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		force := false
		path := ""
		for _, arg := range os.Args[2:] {
			if arg == "--force" {
				force = true
			} else if !strings.HasPrefix(arg, "-") && path == "" {
				path = arg
			}
		}
		if path == "" {
			fmt.Println("Error: Missing path argument for init.")
			fmt.Println("Usage: gtml init <PATH> [--force]")
			os.Exit(1)
		}
		runInit(path, force)

	case "compile":
		watch := false
		path := ""
		for _, arg := range os.Args[2:] {
			if arg == "--watch" {
				watch = true
			} else if !strings.HasPrefix(arg, "-") && path == "" {
				path = arg
			}
		}
		if path == "" {
			fmt.Println("Error: Missing path argument for compile.")
			fmt.Println("Usage: gtml compile <PATH> [--watch]")
			os.Exit(1)
		}

		if watch {
			gtml.WatchProject(path, gtml.CompileOptions{
				ComponentsDir: DirComponents,
				RoutesDir:     DirRoutes,
				DistDir:       DirDist,
				StaticDir:     DirStatic,
			})
		} else {
			err := gtml.CompileProject(path, gtml.CompileOptions{
				ComponentsDir: DirComponents,
				RoutesDir:     DirRoutes,
				DistDir:       DirDist,
				StaticDir:     DirStatic,
			})
			if err != nil {
				fmt.Printf("\n❌ Compilation failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("\n✅ Compilation successful!")
		}

	case "test":
		path := ""
		for _, arg := range os.Args[2:] {
			if !strings.HasPrefix(arg, "-") && path == "" {
				path = arg
			}
		}
		if path == "" {
			path = "."
		}
		runTests(path)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("gtml - A Static Site Generator")
	fmt.Println("Usage:")
	fmt.Println("  gtml init <PATH> [--force]")
	fmt.Println("  gtml compile <PATH> [--watch]")
	fmt.Println("  gtml test [PATH]")
}

func runInit(basePath string, force bool) {
	if _, err := os.Stat(basePath); err == nil && !force {
		fmt.Printf("Error: Directory '%s' already exists. Use --force to overwrite.\n", basePath)
		os.Exit(1)
	}

	if force {
		os.RemoveAll(basePath)
	}

	dirs := []string{
		filepath.Join(basePath, DirComponents),
		filepath.Join(basePath, DirRoutes),
		filepath.Join(basePath, DirStatic),
		filepath.Join(basePath, DirDist),
		filepath.Join(basePath, DirDist, DirStatic),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			panic(fmt.Sprintf("Failed to create directory %s: %v", d, err))
		}
	}

	files := map[string]string{
		filepath.Join(basePath, DirComponents, "BasicButton.html"): `<button props='text string'>{text}</button>`,
		filepath.Join(basePath, DirComponents, "GuestLayout.html"): `<html props='title string'>
  <head>
    <title>{title}</title>
  </head>
  <body>
    <BasicButton text={title} />
    <slot name='content' />
  </body>
 </html>`,
		filepath.Join(basePath, DirRoutes, "index.html"): `<GuestLayout title="Some Title">
  <slot name='content' tag='div'>
    <p>Some Content</p>
  </slot>
 </GuestLayout>`,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			panic(fmt.Sprintf("Failed to write file %s: %v", path, err))
		}
	}

	if err := copyPreinstalledComponents(basePath); err != nil {
		fmt.Printf("Warning: Failed to copy preinstalled components: %v\n", err)
	}

	if err := gtml.CompileProject(basePath, gtml.CompileOptions{
		ComponentsDir: DirComponents,
		RoutesDir:     DirRoutes,
		DistDir:       DirDist,
		StaticDir:     DirStatic,
	}); err != nil {
		fmt.Printf("Warning: Initial compilation failed: %v\n", err)
	} else {
		fmt.Printf("Initialized gtml project at %s\n", basePath)
	}
}

func copyPreinstalledComponents(basePath string) error {
	srcDir := DirPreinstalled
	dstDir := filepath.Join(basePath, DirComponents)

	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return fmt.Errorf("preinstalled components directory not found: %s", srcDir)
	}

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".html" {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(dstPath, content, 0644)
	})
}

func runTests(basePath string) {
	fmt.Printf("Running tests in %s...\n", basePath)

	failed := 0
	passed := 0

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if path == filepath.Join(basePath, "dist") || path == filepath.Join(basePath, "static") {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".html" {
			return nil
		}

		if !strings.HasSuffix(path, "-test.html") {
			return nil
		}

		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}

		fmt.Printf("  Testing %s... ", relPath)

		if err := gtml.CompileProject(filepath.Dir(path), gtml.CompileOptions{
			ComponentsDir: DirComponents,
			RoutesDir:     "routes",
			DistDir:       "dist",
			StaticDir:     "static",
		}); err != nil {
			fmt.Printf("FAILED\n    Error: %v\n", err)
			failed++
		} else {
			fmt.Printf("PASSED\n")
			passed++
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error running tests: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n--- Test Results ---\n")
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)

	if failed > 0 {
		os.Exit(1)
	}

	fmt.Println("All tests passed!")
}
