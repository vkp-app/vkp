package syncers

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"path/filepath"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

func untar(ctx context.Context, r io.Reader, dst string) error {
	log := logging.FromContext(ctx).WithValues("dst", dst)
	log.V(1).Info("unpacking tar")

	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			log.Error(err, "error reading tar")
			return err
		case header == nil:
			continue
		}
		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					log.Error(err, "failed to create unpacked directory", "path", target)
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				log.Error(err, "failed to create unpacked file", "path", target)
				return err
			}

			// copy contents
			if _, err := io.Copy(f, tr); err != nil {
				log.Error(err, "failed to copy contents", "path", target)
				_ = f.Close()
				return err
			}
			_ = f.Close()
		}
	}
}
