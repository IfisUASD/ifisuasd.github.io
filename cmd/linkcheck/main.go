package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	rootDir := "./output"
	brokenLinks := 0

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if err := checkLinks(path, rootDir); err != nil {
				brokenLinks++
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	if brokenLinks > 0 {
		fmt.Printf("\n❌ Found broken links in %d files.\n", brokenLinks)
		os.Exit(1)
	} else {
		fmt.Println("\n✅ No broken links found!")
	}
}

func checkLinks(filePath string, rootDir string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	doc, err := html.Parse(f)
	if err != nil {
		return err
	}

	hasBroken := false
	var fVisit func(*html.Node)
	fVisit = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link := a.Val
					// Ignore external links, mailto, tel, #
					if strings.HasPrefix(link, "http") || strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") || strings.HasPrefix(link, "#") {
						continue
					}

					// Resolve absolute paths relative to rootDir
					var targetPath string
					if strings.HasPrefix(link, "/") {
						targetPath = filepath.Join(rootDir, link)
					} else {
						// Relative path
						baseDir := filepath.Dir(filePath)
						targetPath = filepath.Join(baseDir, link)
					}

					// Remove query params and fragments
					if idx := strings.IndexAny(targetPath, "?#"); idx != -1 {
						targetPath = targetPath[:idx]
					}

					// Check if file or directory exists
					info, err := os.Stat(targetPath)
					if err != nil {
						// Try adding index.html if it's a directory path but missing trailing slash or index.html
						targetPathIndex := filepath.Join(targetPath, "index.html")
						if _, errIndex := os.Stat(targetPathIndex); errIndex == nil {
							continue // Found index.html
						}
						
						fmt.Printf("❌ Broken link in %s: %s (Target: %s)\n", filePath, link, targetPath)
						hasBroken = true
					} else if info.IsDir() {
						// Check for index.html inside
						if _, err := os.Stat(filepath.Join(targetPath, "index.html")); err != nil {
							fmt.Printf("❌ Broken link in %s: %s (Target directory missing index.html: %s)\n", filePath, link, targetPath)
							hasBroken = true
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fVisit(c)
		}
	}
	fVisit(doc)

	if hasBroken {
		return fmt.Errorf("broken links found")
	}
	return nil
}
