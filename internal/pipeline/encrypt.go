package pipeline

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
)

const EncryptAESCTR = "encrypt.aes_ctr"

type AESCTRStage struct{}

func (AESCTRStage) Name() string {
	return EncryptAESCTR
}

func (AESCTRStage) Wrap(_ context.Context, src io.ReadCloser, cfg StageConfig) (io.ReadCloser, error) {
	keyText := cfg["key_b64"]
	ivText := cfg["iv_b64"]
	if keyText == "" || ivText == "" {
		return nil, fmt.Errorf("encrypt.aes_ctr requires key_b64 and iv_b64")
	}

	key, err := base64.StdEncoding.DecodeString(keyText)
	if err != nil {
		return nil, fmt.Errorf("decode key_b64: %w", err)
	}

	iv, err := base64.StdEncoding.DecodeString(ivText)
	if err != nil {
		return nil, fmt.Errorf("decode iv_b64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("iv length must be %d bytes", block.BlockSize())
	}

	stream := cipher.NewCTR(block, iv)
	return &cipherReadCloser{
		Reader: cipher.StreamReader{
			S: stream,
			R: src,
		},
		closer: src,
	}, nil
}

type cipherReadCloser struct {
	io.Reader
	closer io.Closer
}

func (c *cipherReadCloser) Close() error {
	return c.closer.Close()
}
