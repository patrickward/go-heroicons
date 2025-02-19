package heroicons

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// DefaultMissingIconSVG is the default SVG content for the missing icon
var DefaultMissingIconSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#fb2c36"><path d="M17.5 2.5L23 12L17.5 21.5H6.5L1 12L6.5 2.5H17.5ZM11 15V17H13V15H11ZM11 7V13H13V7H11Z"></path></svg>`

// IconType represents the different types of Heroicons
type IconType string

const (
	IconOutline IconType = "outline" // 24px outline icons
	IconSolid   IconType = "solid"   // 24px solid icons
	IconMini    IconType = "mini"    // 20px solid icons
	IconMicro   IconType = "micro"   // 16px solid icons
)

// IconSet defines an icon to be included in the project
type IconSet struct {
	Name string
	Type IconType
}

// Generator handles the icon generation process
type Generator struct {
	// HeroiconsPath is the path to the heroicons repository
	HeroiconsPath string
	// OutputPath is where the generated files will be written
	OutputPath string
	// PackageName is the name of the generated package
	PackageName string
	// Icons is the list of icons to include
	Icons []IconSet
	// FailOnError if true, missing icons will cause an error; otherwise, the missing icon will be used
	FailOnError bool
	// MissingIconSVG is the SVG content to use for missing icons. This overrides the default.
	MissingIconSVG string
}

// Generate creates the icon manifest and copies the required icons
func (g *Generator) Generate() error {
	if g.MissingIconSVG == "" {
		g.MissingIconSVG = DefaultMissingIconSVG
	}

	// Create output directory
	iconsPath := filepath.Join(g.OutputPath, "icons")
	if err := os.MkdirAll(iconsPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write our missing icon SVG
	missingIconPath := filepath.Join(iconsPath, "missing.svg")
	if err := os.WriteFile(missingIconPath, []byte(g.MissingIconSVG), 0644); err != nil {
		return fmt.Errorf("failed to write missing icon: %w", err)
	}

	// Copy icons and build manifest
	var missingIcons []string
	iconPaths := make(map[string]string)
	for _, icon := range g.Icons {
		srcPath := g.getIconPath(icon)
		filename := fmt.Sprintf("%s_%s.svg", icon.Type, icon.Name)
		destPath := filepath.Join(iconsPath, filename)

		if err := g.copyIcon(srcPath, destPath); err != nil {
			missingIcons = append(missingIcons, fmt.Sprintf("%s/%s", icon.Type, icon.Name))
			continue
		}

		key := fmt.Sprintf("%s/%s", icon.Type, icon.Name)
		iconPaths[key] = filename
	}

	// Generate provider.go
	if err := g.generateProvider(iconPaths); err != nil {
		return fmt.Errorf("failed to generate provider: %w", err)
	}

	// Log which icons are missing
	if len(missingIcons) > 0 {
		fmt.Printf("The following icons were not found and could not be copied:\n%s\n",
			strings.Join(missingIcons, "\n"))
	}

	return nil
}

func (g *Generator) getIconPath(icon IconSet) string {
	var dir string
	switch icon.Type {
	case IconOutline:
		dir = "24/outline"
	case IconSolid:
		dir = "24/solid"
	case IconMini:
		dir = "20/solid"
	case IconMicro:
		dir = "16/solid"
	}
	return filepath.Join(g.HeroiconsPath, "optimized", dir, icon.Name+".svg")
}

func (g *Generator) copyIcon(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer func(srcFile *os.File) {
		_ = srcFile.Close()
	}(srcFile)

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer func(destFile *os.File) {
		_ = destFile.Close()
	}(destFile)

	_, err = io.Copy(destFile, srcFile)
	return err
}

const providerTemplate = `// Code generated by heroicons generator; DO NOT EDIT.
package icons

import (
	"fmt"
	"embed"
	"html/template"
	"strings"

	"github.com/patrickward/go-heroicons"
)

//go:embed icons/*.svg
var iconFS embed.FS

// FailOnError determines whether to use a generic missing icon when an icon is not found
var FailOnError = {{ if .FailOnError }}true{{ else }}false{{ end }} 

// IconType represents the different types of Heroicons
type IconType string

const (
	IconOutline IconType = "outline" // 24px outline icons
	IconSolid   IconType = "solid"   // 24px solid icons
	IconMini    IconType = "mini"    // 20px solid icons
	IconMicro   IconType = "micro"   // 16px solid icons
)

var iconPaths = map[string]string{
{{- range $key, $path := .IconPaths }}
	"{{ $key }}": "{{ $path }}",
{{- end }}
}

// RenderIcon returns the SVG content for the specified icon with added classes
func RenderIcon(name string, iconType heroicons.IconType, class string) (template.HTML, error) {
	svg, err := getIcon(name, iconType)
	if err != nil {
		return "", err
	}

	// If class is provided, insert it into the SVG
	if class != "" {
		if strings.Contains(svg, "class=\"") {
			svg = strings.Replace(svg, "class=\"", fmt.Sprintf("class=\"%s ", class), 1)
		} else {
			svg = strings.Replace(svg, "<svg ", fmt.Sprintf("<svg class=\"%s\" ", class), 1)
		}
	}

	return template.HTML(svg), nil
}

func getMissingIcon() string {
	content, err := iconFS.ReadFile("icons/missing.svg")
	if err != nil {
		return ""
	}
	return string(content)
}

func getIcon(name string, iconType heroicons.IconType) (string, error) {
	key := fmt.Sprintf("%s/%s", iconType, name)
	filename, ok := iconPaths[key]
	if !ok {
		if FailOnError {
			return "", fmt.Errorf("icon not found: %s", key)
		}
		return getMissingIcon(), nil
	}

	filename = fmt.Sprintf("icons/%s", filename)
	content, err := iconFS.ReadFile(filename)
	if err != nil {
		if FailOnError {
			return "", fmt.Errorf("failed to read icon %s: %w", filename, err)
		}
		return getMissingIcon(), nil
	}

	return string(content), nil
}`

func (g *Generator) generateProvider(iconPaths map[string]string) error {
	tmpl, err := template.New("provider").Parse(providerTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(g.OutputPath, "provider.go"))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	data := struct {
		PackageName string
		IconPaths   map[string]string
		FailOnError bool
	}{
		PackageName: g.PackageName,
		IconPaths:   iconPaths,
		FailOnError: g.FailOnError,
	}

	return tmpl.Execute(f, data)
}
