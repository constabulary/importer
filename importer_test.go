package importer

import (
	"go/build"
	"reflect"
	"runtime"
	"testing"

	"github.com/pkg/errors"
)

func TestSrcDirImporterImport(t *testing.T) {
	i := &srcdirImporter{
		Context: &build.Default,
		root:    runtime.GOROOT(),
	}
	tests := []struct {
		Importer
		setup func(*testing.T) func(*testing.T)
		path  string
		err   error
	}{{
		Importer: i,
		path:     "",
		err: &importErr{
			msg: "invalid import path",
		},
	}, {
		Importer: i,
		path:     ".",
		err: &importErr{
			path: ".",
			msg:  "relative import not supported",
		},
	}, {
		Importer: i,
		path:     "..",
		err: &importErr{
			path: "..",
			msg:  "relative import not supported",
		},
	}, {
		Importer: i,
		path:     "/math",
		err: &importErr{
			path: "/math",
			msg:  "cannot import absolute path",
		},
	}}

	for i, tt := range tests {
		if tt.setup != nil {
			teardown := tt.setup(t)
			defer teardown(t)
		}
		_, err := tt.Import(tt.path)
		if tt.err == nil && err != nil {
			t.Errorf("%d: srcdirImporter.Import(%q): err, got %v, want %v", i, tt.path, err, nil)
			continue
		}
		err = errors.Cause(err)
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("%d: srcdirImporter.Import(%q): err, got %v, want %v", i, tt.path, err, tt.err)
			continue
		}
	}
}

func TestGOROOT(t *testing.T) {
	i := GOROOT(&build.Default)
	tests := []struct {
		Importer
		path string
		err  error
	}{{
		Importer: i,
		path:     "math",
	}, {
		Importer: i,
		path:     "math/abs.go",
		err: &importErr{
			path: "math/abs.go",
			msg:  "not a directory",
		},
	}}

	for i, tt := range tests {
		_, err := tt.Import(tt.path)
		if tt.err == nil && err != nil {
			t.Errorf("%d: GOROOT().Import(%q): err, got %v, want %v", i, tt.path, err, nil)
			continue
		}
		err = errors.Cause(err)
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("%d: GOROOT().Import(%q): err, got %v, want %v", i, tt.path, err, tt.err)
			continue
		}

	}
}
