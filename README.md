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
	"github.com/patrickward/go-heroicons"
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
		// Each name must match a file in the Heroicons repository.
		Icons: []heroicons.IconSet{
			{Name: "academic-cap", Type: heroicons.IconOutline},
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
- Copy the requested icons from the Heroicons repository into the internal/icons/icons directory
- Generate the internal/icons/provider.go file with the icons embedded
- Include a "missing icon" SVG for any icons not found during runtime

### 3. Use the Icons in Your Templates

In your project, you can now use the generated icons in your HTML template. 

The `RenderIcon` method can be used to render an icon in your template. It takes the icon name, type, and any additional classes you want to apply to the SVG:

```go

package main

import (
   "html/template"
   "yourproject/internal/icons"
)

func main() {
	
    // Add to an existing functions map or create a new one
    funcs := template.FuncMap{
        "icon": icons.RenderIcon,
    }
	
   // Use in your templates
   tmpl := template.Must(template.New("example").Funcs(funcs).Parse(`
           <button>
               {{icon "home" "outline" "w-6 h-6"}}
               Home
           </button>
       `))
   }
```

## Icon Types

The package supports v3 Heroicon types: 
- `outline` - 24px outline icons
- `solid` - 24px solid icons
- `mini` - 20px solid icons
- `micro` - 16px solid icons

## Missing Icons

If an icon specified in your generator configuration isn't found in the Heroicons repository during generation, the package will:

1. Log a warning during generation
2. Continue with the build process

During runtime, if an icon cannot be found, the package will render a "missing icon" SVG in place of the missing icon. The default missing icon is a red hexagon with an exclamation mark.

Alternatively, you can return an error if a missing icon is encountered by setting `FailOnError` to `true` in your generator configuration.

You can provide your own "missing icon" SVG by overriding the `MissingIconSVG` for the package:

```go
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
   },

   // Fail if any icons are missing
   FailOnError: false,

   // Custom missing icon SVG
   MissingIconSVG: `<svg xmlns="http://www.w3.org/2000/svg" ... </svg>`,
}

```

## Important Notes

- This package does not include the Heroicons SVGs. You need to provide the path to the Heroicons repository during build time.
- Only the icons you specify will be embedded in your binary.
- The generator needs to be run whenever you add or remove icons.
- Icons are embedded as SVGs and can be styled with CSS classes.
- The package uses `go:embed` to include the icons in your binary - no runtime filesystem access is needed.

## License

This package is licensed under MIT. Note that Heroicons has its own license - please check the [Heroicons repository](https://github.com/tailwindlabs/heroicons) for details.
