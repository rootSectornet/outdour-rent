package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rentoutdoor/api/internal/usecase/auth"
)

const googleTokenInfoURL = "https://oauth2.googleapis.com/tokeninfo?id_token="

type googleVerifier struct {
	httpClient *http.Client
	clientID   string
}

// NewGoogleVerifier creates a Google ID token verifier.
func NewGoogleVerifier(clientID string) auth.GoogleTokenVerifier {
	return &googleVerifier{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		clientID:   clientID,
	}
}

func (v *googleVerifier) Verify(ctx context.Context, idToken string) (*auth.GoogleClaims, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleTokenInfoURL+idToken, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google token verification failed: %s", string(body))
	}

	var payload struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Aud           string `json:"aud"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}

	// Validate audience matches our client ID
	if v.clientID != "" && payload.Aud != v.clientID {
		return nil, fmt.Errorf("token audience mismatch: expected %s, got %s", v.clientID, payload.Aud)
	}

	return &auth.GoogleClaims{
		Sub:           payload.Sub,
		Email:         payload.Email,
		EmailVerified: payload.EmailVerified == "true",
		Name:          payload.Name,
		Picture:       payload.Picture,
	}, nil
}
