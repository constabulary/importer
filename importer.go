// importer provides mechanisms for loading go/build.Package
// structures from source packages on disk.
package importer

import (
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// Importer imports a package.
type Importer interface {

	// Import imports the package from importpath.
	Import(importpath string) (*build.Package, error)
}

// GOROOT returns an Importer which loads packages from the standard library.
func GOROOT(ctx *build.Context) Importer {
	return &srcdirImporter{
		Context: ctx,
		root:    filepath.Join(runtime.GOROOT()),
	}
}

type srcdirImporter struct {
	*build.Context
	root string
}

func (i *srcdirImporter) Import(importpath string) (*build.Package, error) {
	if importpath == "" {
		return nil, errors.WithStack(&importErr{path: importpath, msg: "invalid import path"})
	}

	if importpath == "." || importpath == ".." || strings.HasPrefix(importpath, "./") || strings.HasPrefix(importpath, "../") {
		return nil, errors.WithStack(&importErr{path: importpath, msg: "relative import not supported"})
	}

	if strings.HasPrefix(importpath, "/") {
		return nil, errors.WithStack(&importErr{path: importpath, msg: "cannot import absolute path"})
	}

	var p *build.Package

	loadPackage := func(importpath, dir string) error {
		pkg, err := i.ImportDir(dir, 0)
		if err != nil {
			return err
		}
		p = pkg
		p.ImportPath = importpath
		return nil
	}

	// if this is the stdlib, then search vendor first.
	// this isn't real vendor support, just enough to make net/http compile.
	if i.root == runtime.GOROOT() {
		importpath := path.Join("vendor", importpath)
		dir := filepath.Join(i.root, "src", filepath.FromSlash(importpath))
		fi, err := os.Stat(dir)
		if err == nil && fi.IsDir() {
			err := loadPackage(importpath, dir)
			return p, err
		}
	}

	dir := filepath.Join(i.root, "src", filepath.FromSlash(importpath))
	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, errors.WithStack(&importErr{path: importpath, msg: "not a directory"})
	}
	err = loadPackage(importpath, dir)
	return p, err
}

type importErr struct {
	path string
	msg  string
}

func (e *importErr) Error() string {
	return fmt.Sprintf("import %q: %v", e.path, e.msg)
}
