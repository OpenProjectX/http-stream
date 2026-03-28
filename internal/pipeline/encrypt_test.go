package pipeline

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io"
	"strings"
	"testing"
)

func TestAESCTRStageWrap(t *testing.T) {
	stage := AESCTRStage{}
	key := []byte("0123456789abcdef0123456789abcdef")
	iv := []byte("0123456789abcdef")
	input := "hello stream"

	reader, err := stage.Wrap(context.Background(), io.NopCloser(strings.NewReader(input)), StageConfig{
		"key_b64": base64.StdEncoding.EncodeToString(key),
		"iv_b64":  base64.StdEncoding.EncodeToString(iv),
	})
	if err != nil {
		t.Fatalf("Wrap() error = %v", err)
	}
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("NewCipher() error = %v", err)
	}
	expected := make([]byte, len(input))
	cipher.NewCTR(block, iv).XORKeyStream(expected, []byte(input))

	if string(got) != string(expected) {
		t.Fatalf("ciphertext mismatch: got %x want %x", got, expected)
	}
}
