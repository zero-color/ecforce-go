package customer

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/zero-color/ecforce-go"
)

// MultiPass is the result of a MultiPass authentication: the identified
// customer, an API authentication token and a one-time login URL.
type MultiPass struct {
	ID                  int64  `json:"id"`
	Email               string `json:"email"`
	AuthenticationToken string `json:"authentication_token"`
	// LoginURL logs anyone in as the customer without credentials; handle it
	// with care. It expires after 5 minutes and can be used only once.
	LoginURL string `json:"login_url"`
}

// Get exchanges a MultiPass token for an API authentication token and a
// one-time login URL. A MultiPass token that has already issued an
// authentication token cannot be used again.
//
// ecforce API docs: GET /api/v2/customers/multi_pass/:token
func (s *MultiPassService) Get(ctx context.Context, token string) (*MultiPass, *ecforce.Response, error) {
	path := fmt.Sprintf("customers/multi_pass/%s.json", url.PathEscape(token))
	mp := new(MultiPass)
	resp, err := ecforce.Do(ctx, s.client, http.MethodGet, path, nil, nil, mp)
	if err != nil {
		return nil, resp, err
	}
	return mp, resp, nil
}

// Authenticate generates a MultiPass token for the customer data using the
// shop's MultiPass secret key and exchanges it for an API authentication token
// and a one-time login URL. It is the one-call form of GenerateMultiPassToken
// followed by Get.
func (s *MultiPassService) Authenticate(ctx context.Context, secretKey string, data MultiPassCustomerData) (*MultiPass, *ecforce.Response, error) {
	token, err := GenerateMultiPassToken(secretKey, data)
	if err != nil {
		return nil, nil, err
	}
	return s.Get(ctx, token)
}

// multiPassCreatedAtLayout is the layout used for the MultiPass token's
// created_at field (e.g. "2023-04-20T10:30:00+00:00"). The numeric offset form
// is used so UTC is rendered as "+00:00" rather than "Z".
const multiPassCreatedAtLayout = "2006-01-02T15:04:05-07:00"

// MultiPassCustomerData is the customer data embedded in a MultiPass token.
type MultiPassCustomerData struct {
	// Email is the customer's email address.
	Email string
	// Identifier is the customer ID in the source system.
	Identifier string
	// CreatedAt is the issue time of the MultiPass token.
	CreatedAt time.Time
}

func (d MultiPassCustomerData) validate() error {
	if d.Email == "" {
		return errors.New("ecforce: multipass: email is required")
	}
	if d.Identifier == "" {
		return errors.New("ecforce: multipass: identifier is required")
	}
	if d.CreatedAt.IsZero() {
		return errors.New("ecforce: multipass: created_at is required")
	}
	return nil
}

// GenerateMultiPassToken generates an authentication token for ecforce's
// MultiPass API from the shop's MultiPass secret key and the customer data.
//
// Specification:
//  1. Hash the secret key with SHA256.
//  2. Use the first 128 bits of the hash as the encryption key and the last
//     128 bits as the signature key.
//  3. Marshal the customer data to JSON and encrypt it with AES-128-CBC
//     (random IV). Encrypted data = IV + ciphertext.
//  4. Sign the encrypted data with HMAC-SHA256 using the signature key.
//  5. Concatenate the encrypted data and the signature and URL-safe Base64
//     encode it to form the token.
//
// The IV is random, so the same input produces a different token each time.
// Note that a generated token can only be used once.
//
// ecforce API docs: https://apidoc.ec-force.com/apidoc/v2/customer/index.html#api-MultiPass_API-MultiPass
func GenerateMultiPassToken(secretKey string, data MultiPassCustomerData) (string, error) {
	return generateMultiPassToken(rand.Reader, secretKey, data)
}

func generateMultiPassToken(randReader io.Reader, secretKey string, data MultiPassCustomerData) (string, error) {
	if secretKey == "" {
		return "", errors.New("ecforce: multipass: secret key is required")
	}
	if err := data.validate(); err != nil {
		return "", err
	}

	payload := struct {
		Email      string `json:"email"`
		Identifier string `json:"identifier"`
		CreatedAt  string `json:"created_at"`
	}{
		Email:      data.Email,
		Identifier: data.Identifier,
		CreatedAt:  data.CreatedAt.Format(multiPassCreatedAtLayout),
	}
	customerData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("ecforce: multipass: marshal customer data: %w", err)
	}

	hash := sha256.Sum256([]byte(secretKey))
	encryptionKey := hash[:16]
	signatureKey := hash[16:]

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("ecforce: multipass: new cipher: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(randReader, iv); err != nil {
		return "", fmt.Errorf("ecforce: multipass: generate iv: %w", err)
	}

	padded := pkcs7Pad(customerData, aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)

	// Encrypted data = IV + ciphertext
	encrypted := make([]byte, 0, len(iv)+len(ciphertext))
	encrypted = append(encrypted, iv...)
	encrypted = append(encrypted, ciphertext...)

	mac := hmac.New(sha256.New, signatureKey)
	mac.Write(encrypted)
	signature := mac.Sum(nil)

	token := make([]byte, 0, len(encrypted)+len(signature))
	token = append(token, encrypted...)
	token = append(token, signature...)

	return base64.URLEncoding.EncodeToString(token), nil
}

// pkcs7Pad applies PKCS#7 padding. AES-CBC requires the input length to be a
// multiple of the block size.
func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	return append(data, bytes.Repeat([]byte{byte(padLen)}, padLen)...)
}
