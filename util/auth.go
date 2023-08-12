package util

import (
	"bread-clock/configs"
	e "bread-clock/error"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"strings"
	"time"
)

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func GetUserInfo(ctx context.Context, provider string, accessToken string) (*UserInfo, error) {
	// currently only supports Google OAuth
	if provider != "google" {
		return nil, e.ErrAuthInvalidProvider
	}

	const googleOAuthURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleOAuthURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	bearerToken := "Bearer " + accessToken
	req.Header.Add("Authorization", bearerToken)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request user info: %w", err)
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, e.ErrAuthInvalidToken
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info: %w", err)
	}

	var user UserInfo
	if err = json.Unmarshal(respBody, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &user, nil
}

func GenerateToken(userID int, emailAddress string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["uid"] = userID
	claims["email"] = emailAddress
	claims["exp"] = time.Now().Add(30 * 24 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.Conf.AuthKey))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return []byte(configs.Conf.AuthKey), nil
	})
}

func GetToken(token string, bearerToken string) (string, error) {
	if token != "" {
		return token, nil
	}

	split := strings.Split(bearerToken, " ")
	if len(split) == 2 {
		return split[1], nil
	}

	return "", e.ErrAuthInvalidToken
}

func ExtractUserID(token *jwt.Token) (int, error) {
	if !token.Valid {
		return 0, e.ErrAuthInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, e.ErrAuthInvalidToken
	}

	userID, ok := claims["uid"].(float64)
	if !ok {
		return 0, e.ErrAuthInvalidToken
	}

	return int(userID), nil
}
