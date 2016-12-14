// +build go1.7

package importer

import (
	"go/build"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestGOROOTVendor(t *testing.T) {
	i := GOROOT(&build.Default)
	tests := []struct {
		Importer
		path string
		err  error
	}{{
		Importer: i,
		path:     "golang_org/x/net/http2/hpack",
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
