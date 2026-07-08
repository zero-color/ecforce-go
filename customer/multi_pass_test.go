package customer

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGenerateMultiPassToken(t *testing.T) {
	const secretKey = "dummy-secret-key"
	createdAt := time.Date(2023, 4, 20, 10, 30, 0, 0, time.UTC)

	t.Run("generated token can be verified and decrypted back to the original customer data", func(t *testing.T) {
		data := MultiPassCustomerData{
			Email:      "sample@example.jp",
			Identifier: "123456789",
			CreatedAt:  createdAt,
		}

		token, err := GenerateMultiPassToken(secretKey, data)
		if err != nil {
			t.Fatal(err)
		}
		decrypted := decryptMultiPassToken(t, secretKey, token)
		want := map[string]any{
			"email":      "sample@example.jp",
			"identifier": "123456789",
			"created_at": "2023-04-20T10:30:00+00:00",
		}
		if !reflect.DeepEqual(decrypted, want) {
			t.Errorf("decrypted = %v, want %v", decrypted, want)
		}
	})

	t.Run("same input produces a different token each time because the IV is random", func(t *testing.T) {
		data := MultiPassCustomerData{
			Email:      "sample@example.jp",
			Identifier: "123456789",
			CreatedAt:  createdAt,
		}

		token1, err := GenerateMultiPassToken(secretKey, data)
		if err != nil {
			t.Fatal(err)
		}
		token2, err := GenerateMultiPassToken(secretKey, data)
		if err != nil {
			t.Fatal(err)
		}
		if token1 == token2 {
			t.Errorf("tokens must differ, both are %q", token1)
		}
	})

	t.Run("a fixed IV produces a deterministic token", func(t *testing.T) {
		data := MultiPassCustomerData{
			Email:      "sample@example.jp",
			Identifier: "123456789",
			CreatedAt:  createdAt,
		}
		fixedIV := bytes.Repeat([]byte{0x01}, aes.BlockSize)

		token1, err := generateMultiPassToken(bytes.NewReader(fixedIV), secretKey, data)
		if err != nil {
			t.Fatal(err)
		}
		token2, err := generateMultiPassToken(bytes.NewReader(fixedIV), secretKey, data)
		if err != nil {
			t.Fatal(err)
		}
		if token1 != token2 {
			t.Errorf("tokens must match: %q != %q", token1, token2)
		}
	})

	t.Run("the JST timezone is reflected in created_at", func(t *testing.T) {
		data := MultiPassCustomerData{
			Email:      "sample@example.jp",
			Identifier: "123456789",
			CreatedAt:  time.Date(2023, 4, 20, 19, 30, 0, 0, time.FixedZone("JST", 9*60*60)),
		}

		token, err := GenerateMultiPassToken(secretKey, data)
		if err != nil {
			t.Fatal(err)
		}
		decrypted := decryptMultiPassToken(t, secretKey, token)
		if got := decrypted["created_at"]; got != "2023-04-20T19:30:00+09:00" {
			t.Errorf("created_at = %v, want 2023-04-20T19:30:00+09:00", got)
		}
	})

	t.Run("validation errors", func(t *testing.T) {
		tests := []struct {
			name      string
			secretKey string
			data      MultiPassCustomerData
		}{
			{
				name:      "empty secret key",
				secretKey: "",
				data:      MultiPassCustomerData{Email: "a@example.com", Identifier: "1", CreatedAt: createdAt},
			},
			{
				name:      "empty email",
				secretKey: secretKey,
				data:      MultiPassCustomerData{Identifier: "1", CreatedAt: createdAt},
			},
			{
				name:      "empty identifier",
				secretKey: secretKey,
				data:      MultiPassCustomerData{Email: "a@example.com", CreatedAt: createdAt},
			},
			{
				name:      "unset created_at",
				secretKey: secretKey,
				data:      MultiPassCustomerData{Email: "a@example.com", Identifier: "1"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				token, err := GenerateMultiPassToken(tt.secretKey, tt.data)
				if err == nil {
					t.Error("expected an error")
				}
				if token != "" {
					t.Errorf("token = %q, want empty", token)
				}
			})
		}
	})
}

func TestMultiPassService_Authenticate(t *testing.T) {
	t.Run("generates a token, sends it in the path, and returns the login URL", func(t *testing.T) {
		const secretKey = "dummy-secret-key"
		data := MultiPassCustomerData{
			Email:      "sample@example.jp",
			Identifier: "123456789",
			CreatedAt:  time.Date(2023, 4, 20, 10, 30, 0, 0, time.UTC),
		}
		client, mux := setup(t)
		var gotToken string
		mux.HandleFunc("GET /api/v2/customers/multi_pass/{token}", func(w http.ResponseWriter, r *http.Request) {
			gotToken = strings.TrimSuffix(r.PathValue("token"), ".json")
			fmt.Fprint(w, `{"id": 1, "login_url": "https://x.com/login/ONE_TIME_TOKEN"}`)
		})

		mp, _, err := client.MultiPass.Authenticate(context.Background(), secretKey, data)
		if err != nil {
			t.Fatal(err)
		}
		if mp.LoginURL != "https://x.com/login/ONE_TIME_TOKEN" {
			t.Errorf("LoginURL = %q, want https://x.com/login/ONE_TIME_TOKEN", mp.LoginURL)
		}
		// The token sent to ecforce must decrypt back to the original customer data.
		decrypted := decryptMultiPassToken(t, secretKey, gotToken)
		if decrypted["email"] != "sample@example.jp" || decrypted["identifier"] != "123456789" {
			t.Errorf("decrypted token = %v", decrypted)
		}
	})

	t.Run("returns an error when token generation fails", func(t *testing.T) {
		client, _ := setup(t)

		mp, _, err := client.MultiPass.Authenticate(context.Background(), "", MultiPassCustomerData{})
		if err == nil {
			t.Error("expected an error")
		}
		if mp != nil {
			t.Errorf("mp = %+v, want nil", mp)
		}
	})
}

// decryptMultiPassToken verifies and decrypts a MultiPass token following the
// same steps as ecforce, then decodes and returns the embedded customer data
// JSON. A verification helper independent of the generation logic.
func decryptMultiPassToken(t *testing.T, secretKey, token string) map[string]any {
	t.Helper()

	raw, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.Sum256([]byte(secretKey))
	encryptionKey := hash[:16]
	signatureKey := hash[16:]

	if len(raw) <= sha256.Size+aes.BlockSize {
		t.Fatalf("token too short: %d bytes", len(raw))
	}
	encrypted := raw[:len(raw)-sha256.Size]
	signature := raw[len(raw)-sha256.Size:]

	mac := hmac.New(sha256.New, signatureKey)
	mac.Write(encrypted)
	if !hmac.Equal(signature, mac.Sum(nil)) {
		t.Fatal("signature must match")
	}

	iv := encrypted[:aes.BlockSize]
	ciphertext := encrypted[aes.BlockSize:]
	if len(ciphertext)%aes.BlockSize != 0 {
		t.Fatalf("ciphertext length %d is not a multiple of the block size", len(ciphertext))
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		t.Fatal(err)
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)

	// PKCS#7 unpadding
	padLen := int(plaintext[len(plaintext)-1])
	if padLen < 1 || padLen > aes.BlockSize {
		t.Fatalf("invalid padding length %d", padLen)
	}
	plaintext = plaintext[:len(plaintext)-padLen]

	if !strings.HasPrefix(string(plaintext), "{") {
		t.Fatalf("plaintext is not JSON: %q", plaintext)
	}

	var result map[string]any
	if err := json.Unmarshal(plaintext, &result); err != nil {
		t.Fatal(err)
	}
	return result
}
