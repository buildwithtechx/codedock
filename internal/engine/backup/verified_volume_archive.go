package backup

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"
)

type temporaryArchive struct{ *os.File }

func (a *temporaryArchive) Close() error {
	return errors.Join(a.File.Close(), os.Remove(a.Name()))
}

func verifiedVolumeArchive(ctx context.Context, source io.ReadCloser, checksum string) (io.ReadCloser, error) {
	file, err := os.CreateTemp("", "codedock-restore-*")
	if err != nil {
		return nil, errors.Join(err, source.Close())
	}
	archive := &temporaryArchive{file}
	failed := true
	defer func() {
		if failed {
			if err := archive.Close(); err != nil {
				slog.Warn("remove temporary restore archive", "error", err)
			}
		}
	}()
	digest := sha256.New()
	count, copyErr := io.Copy(io.MultiWriter(file, digest), io.LimitReader(source, 1024*1024*1024+1))
	closeErr := source.Close()
	if err := errors.Join(copyErr, closeErr, ctx.Err()); err != nil {
		return nil, err
	}
	if count == 0 || count > 1024*1024*1024 || hex.EncodeToString(digest.Sum(nil)) != checksum {
		return nil, fmt.Errorf("volume archive changed after review")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	compressed, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("invalid compressed volume archive: %w", err)
	}
	reader := tar.NewReader(compressed)
	var total int64
	for {
		if err := ctx.Err(); err != nil {
			compressed.Close()
			return nil, err
		}
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			compressed.Close()
			return nil, fmt.Errorf("invalid volume archive: %w", err)
		}
		name := path.Clean(header.Name)
		if path.IsAbs(name) || name == ".." || strings.HasPrefix(name, "../") || strings.Contains(name, "\\") {
			compressed.Close()
			return nil, fmt.Errorf("archive contains an unsafe path")
		}
		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeRegA && header.Typeflag != tar.TypeDir {
			compressed.Close()
			return nil, fmt.Errorf("archive contains unsupported links or special files")
		}
		total += header.Size
		if total > 10*1024*1024*1024 {
			compressed.Close()
			return nil, fmt.Errorf("expanded archive exceeds 10 GiB")
		}
	}
	if _, err := io.Copy(io.Discard, compressed); err != nil {
		return nil, errors.Join(err, compressed.Close())
	}
	if err := compressed.Close(); err != nil {
		return nil, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	failed = false
	return archive, nil
}
