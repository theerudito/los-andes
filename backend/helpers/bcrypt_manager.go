package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func DerivarClave(valor string) []byte {
	hash := sha256.Sum256([]byte(valor))
	return hash[:]
}

func EncriptarDato(textPlano string) (string, error) {

	key := os.Getenv("Secret_Key")

	keyBytes := DerivarClave(key)

	block, err := aes.NewCipher(keyBytes)

	if err != nil {
		return "", fmt.Errorf("error creando cifrador: %v", err)
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return "", fmt.Errorf("error creando GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generando nonce: %v", err)
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(textPlano), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil

}

func DesencriptarDato(textCifrado string) (string, error) {

	key := os.Getenv("Secret_Key")

	keyBytes := DerivarClave(key)

	data, err := base64.StdEncoding.DecodeString(textCifrado)

	if err != nil {
		return "", fmt.Errorf("error decodificando base64: %v", err)
	}

	block, err := aes.NewCipher(keyBytes)

	if err != nil {
		return "", fmt.Errorf("error creando cifrador: %v", err)
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return "", fmt.Errorf("error creando GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return "", fmt.Errorf("texto cifrado demasiado corto")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)

	if err != nil {
		return "", fmt.Errorf("error desencriptando: %v", err)
	}

	return string(plainText), nil

}
