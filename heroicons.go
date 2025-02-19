// Package heroicons provides template functions for using Heroicons SVGs
package heroicons

import (
	"fmt"
	"html/template"
	"strings"
)

// IconType represents the different types of Heroicons
type IconType string

const (
	IconOutline IconType = "outline" // 24px outline icons
	IconSolid   IconType = "solid"   // 24px solid icons
	IconMini    IconType = "mini"    // 20px solid icons
	IconMicro   IconType = "micro"   // 16px solid icons
)

// IconProvider is an interface that must be implemented by the consuming project
// to provide icon content to the heroicons package
type IconProvider interface {
	// GetIcon returns the SVG content for a given icon name and type
	GetIcon(name string, iconType IconType) (string, error)
}

var provider IconProvider

// Initialize sets up the heroicons package with the given provider
func Initialize(p IconProvider) {
	provider = p
}

// RenderIcon returns the SVG content for the specified icon with added classes
func RenderIcon(name string, iconType IconType, class string) (template.HTML, error) {
	if provider == nil {
		return "", fmt.Errorf("heroicons package not initialized with an IconProvider")
	}

	svg, err := provider.GetIcon(name, iconType)
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
