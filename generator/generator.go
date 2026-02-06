package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/core"
	"github.com/zenaxo/valibase/internal/gen"
)

// Options controls how TypeScript types are generated.
// Zero values are valid and will use defaults.
type Options struct {
	// OutPath is optional if you prefer to pass the output path to GenerateTypes directly.
	// If you set OutPath, you can pass an empty outPath to GenerateTypes.
	OutPath string
}

// GenerateTypes generates TypeScript types for all PocketBase collections and writes them to outPath.
// If opts.OutPath is set and outPath is empty, opts.OutPath will be used.
func GenerateTypes(app core.App, outPath string, opts ...Options) error {
	var o Options
	if len(opts) > 0 {
		o = opts[0]
	}
	if outPath == "" {
		outPath = o.OutPath
	}
	if outPath == "" {
		return fmt.Errorf("GenerateTypes: outPath is required")
	}

	colls, err := app.FindAllCollections()
	if err != nil {
		return fmt.Errorf("GenerateTypes: FindAllCollections: %w", err)
	}

	ts := gen.GenerateTS(gen.BuildCollections(colls))

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("GenerateTypes: mkdir %s: %w", filepath.Dir(outPath), err)
	}

	if err := os.WriteFile(outPath, []byte(ts), 0o644); err != nil {
		return fmt.Errorf("GenerateTypes: write %s: %w", outPath, err)
	}

	return nil
}
