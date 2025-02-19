# Go Heroicons Generator

This package provides a way to selectively embed [Heroicons](https://heroicons.com) SVGs into your Go projects and use them in your HTML templates. Instead of embedding all Heroicons (which would increase your binary size), this package lets you generate a custom subset of only the icons you need.

## Prerequisites

This package requires access to the Heroicons repository during build time. You'll need to:

1. Clone the Heroicons repository:
   ```bash
   git clone https://github.com/tailwindlabs/heroicons.git
   ```

2. Note the path to the cloned repository - you'll need this for configuration.

## Installation

```bash
go get github.com/patrickward/go-heroicons
```

## Usage

### 1. Create Your Icon Generator

Create a file at `internal/icons/generate/main.go`:

```go
package main

//go:generate go run .

import (
    "log"
    "github.com/patrickward/heroicons"
)

func main() {
    generator := &heroicons.Generator{
        // Path to your cloned Heroicons repository
        HeroiconsPath: "/path/to/heroicons",
        
        // Generated icons package will be saved to this directory. 
        // In this case, if the generator is run from internal/icons/generate,
        // the icons will be saved to internal/icons/icons and the provider 
        // will be saved to internal/icons/provider.go
        OutputPath: "../",
        
        // Name of the generated package
        PackageName: "icons",
        
        // List of icons you want to include. Each icon must have a name and type.
        Icons: []heroicons.IconSet{
            {Name: "home", Type: heroicons.IconOutline},
            {Name: "user", Type: heroicons.IconSolid},
            {Name: "cog", Type: heroicons.IconMini},
            {Name: "bell", Type: heroicons.IconMicro},
        },
    }

    if err := generator.Generate(); err != nil {
        log.Fatal(err)
    }
}
```

### 2. Generate the Icons

Run generation using either:

```bash
go generate ./internal/icons/generate
```

or simply:

```bash
go generate ./...
```

This will:
- Create the internal/icons directory if it doesn't exist
- Copy the requested icons from the Heroicons repository into the internal/icons/icons directory
- Generate a provider package with the icons embedded at internal/icons/icons.go
- Include a "missing icon" SVG for any icons that weren't found

### 3. Use the Icons in Your Templates

Initialize the icons package and use it in your templates:

```go
package main

import (
    "html/template"
    "yourproject/internal/icons"
)

func main() {
    // Initialize the icons
    icons.Initialize()

    // Use in your templates
    tmpl := template.Must(template.New("example").Funcs(icons.FuncMap()).Parse(`
        <button>
            {{icon "home" "outline" "w-6 h-6"}}
            Home
        </button>
    `))
}
```

Or, if you are adding to an existing functions map:

```go
func main() {
    // Initialize the icons
    icons.Initialize()
	
    // Add to an existing functions map
    funcs := template.FuncMap{
        "icon": icons.IconFunc,
    }
}
```

## Icon Types

The package supports all Heroicon types:
- `IconOutline` - 24px outline icons
- `IconSolid` - 24px solid icons
- `IconMini` - 20px solid icons
- `IconMicro` - 16px solid icons

## Missing Icons

If an icon specified in your generator configuration isn't found in the Heroicons repository, the package will:
1. Log a warning during generation
2. Substitute a generic "missing icon" SVG
3. Continue with the build process

This ensures your application won't break if an icon is renamed or removed in a future Heroicons update.

## Important Notes

- This package does not include the Heroicons SVGs. You need to have the Heroicons repository available during build time.
- Only the icons you specify will be embedded in your binary.
- The generator needs to be run whenever you want to add or remove icons.
- Icons are embedded as SVGs and can be styled with CSS classes.
- The package uses `go:embed` to include the icons in your binary - no runtime filesystem access is needed.

## License

This package is licensed under MIT. Note that Heroicons has its own license - please check the [Heroicons repository](https://github.com/tailwindlabs/heroicons) for details.
