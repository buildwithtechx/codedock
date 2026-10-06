package backup

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"testing"
)

func TestVerifiedVolumeArchiveRejectsChangedAndUnsafeInput(t *testing.T) {
	for _, test := range []struct {
		name, path string
		kind       byte
		changed    bool
		valid      bool
	}{
		{"valid", "./data/file", tar.TypeReg, false, true},
		{"changed", "data", tar.TypeReg, true, false},
		{"traversal", "../../escape", tar.TypeReg, false, false},
		{"absolute", "/escape", tar.TypeReg, false, false},
		{"link", "data", tar.TypeSymlink, false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer
			compressed := gzip.NewWriter(&buffer)
			writer := tar.NewWriter(compressed)
			if err := writer.WriteHeader(&tar.Header{Name: test.path, Mode: 0600, Typeflag: test.kind, Linkname: "../../escape"}); err != nil {
				t.Fatal(err)
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			if err := compressed.Close(); err != nil {
				t.Fatal(err)
			}
			digest := sha256.Sum256(buffer.Bytes())
			checksum := hex.EncodeToString(digest[:])
			if test.changed {
				checksum = "wrong"
			}
			archive, err := verifiedVolumeArchive(context.Background(), io.NopCloser(bytes.NewReader(buffer.Bytes())), checksum)
			if !test.valid {
				if err == nil {
					archive.Close()
					t.Fatal("unsafe archive accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			actual, err := io.ReadAll(archive)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(actual, buffer.Bytes()) {
				t.Fatal("verified input changed")
			}
			if err := archive.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
