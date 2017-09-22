package stopwatchapp

import (
	"bytes"
	"io"
	"path"
)

// GetAssetFunc defines an interface function for retrieving data of a named file
type GetAssetFunc func(path string) (io.ReadCloser, error)

// Asset is an embedded asset that implements io.ReadCloser
type Asset struct {
	*bytes.Reader
}

// Close method for completing io.Closer
func (a *Asset) Close() error { return nil }

// NewAsset returns an asset from given bytes
func NewAsset(b []byte) *Asset {
	return &Asset{Reader: bytes.NewReader(b)}
}

// AssetRestoreFn defines an accessor function type for restoring embedded assets to
// file system.
type AssetRestoreFn func(dir, name string) error

// UnpackEmbeddedAssets unloads certain required files to given directory from the
// inmemory copy.
func UnpackEmbeddedAssets(toDir string, assetFn AssetRestoreFn) error {
	for _, asset := range EmbeddedResources() {
		if err := assetFn(toDir, asset); err != nil {
			return err
		}
	}

	return nil
}

// AssetReader returns an function for reading in-memory firmware asset as a io.Reader
func (a *App) AssetReader() GetAssetFunc {
	return func(fname string) (io.ReadCloser, error) {
		b, err := a.assetData(path.Join(resourcesDir, fname))
		if err != nil {
			return nil, err
		}

		return NewAsset(b), nil
	}
}
