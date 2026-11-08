package networking

import (
	"archive/tar"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

func DeriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func EncryptBundle(src io.Reader, dst io.Writer, passphrase string) error {
	key := DeriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	plaintext, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	_, err = dst.Write(ciphertext)
	return err
}

func DecryptBundle(src io.Reader, dst io.Writer, passphrase string) error {
	key := DeriveKey(passphrase)
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return errors.New("bundle is too short to be valid")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("decryption failed: wrong passphrase or corrupted bundle")
	}
	_, err = dst.Write(plaintext)
	return err
}

func CreateTarGz(dst io.Writer, files map[string][]byte) error {
	gw := gzip.NewWriter(dst)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for name, data := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0644,
			Size: int64(len(data)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func ExtractTarGz(src io.Reader) (map[string][]byte, error) {
	gr, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	files := make(map[string][]byte)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		files[hdr.Name] = data
	}
	return files, nil
}
