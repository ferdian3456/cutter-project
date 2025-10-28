package util

import (
	"cutterproject/internal/constant"
	"cutterproject/internal/model"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"time"
)

var (
	BearerPrefix            = "Bearer "
	TokenIssuer             = "cutterproject"
	AccessTokenDuration     = 15 * time.Minute
	RefreshTokenDuration    = 7 * 24 * time.Hour
	ErrInvalidSigningMethod = errors.New("invalid token signing method")
)

func GenerateAccessToken(userId int, jwtSecretKey string) (string, error) {
	if jwtSecretKey == "" {
		return "", errors.New("jwt secret key is not configured")
	}

	now := time.Now()
	claims := &model.Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    TokenIssuer,
			Subject:   fmt.Sprintf("user:%d", userId),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// GenerateRefreshToken creates a unique refresh token
// Note: This should be stored server-side with expiration time and user association
func GenerateRefreshToken() string {
	return uuid.New().String()
}

// GenerateTokenPair creates both access and refresh tokens for a user
func GenerateTokenPair(userId int, jwtSecretKey string) (model.TokenResponse, error) {
	accessToken, err := GenerateAccessToken(userId, jwtSecretKey)
	if err != nil {
		return model.TokenResponse{}, err
	}

	refreshToken := GenerateRefreshToken()

	return model.TokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresIn:  int(AccessTokenDuration.Seconds()),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: int(RefreshTokenDuration.Seconds()),
		TokenType:             "Bearer",
	}, nil
}

// ValidateAccessToken validates a JWT access token and returns the user ID
func ValidateAccessToken(accessToken string, log *zap.Logger, jwtSecretKey string) (int, error) {
	// Don't log the full token - security risk
	log.Debug("Validating access token", zap.String("accessToken", accessToken))

	if jwtSecretKey == "" {
		return 0, errors.New("jwt secret key is not configured")
	}

	// Extract token from Authorization header
	tokenString, err := extractBearerToken(accessToken)
	if err != nil {
		return 0, err
	}

	// Parse token with custom claims
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return 0, handleParseError(err)
	}

	// Extract and validate claims
	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return 0, &model.ValidationError{
			Code:    constant.ERR_UNATHORIZED_ERROR,
			Message: "Authentication token is invalid",
			Param:   "accessToken",
		}
	}

	return claims.UserId, nil
}

// extractBearerToken extracts the token from "Bearer <token>" format
func extractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", &model.ValidationError{
			Code:    constant.ERR_UNATHORIZED_ERROR,
			Message: "No authentication token is provided",
			Param:   "accessToken",
		}
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", &model.ValidationError{
			Code:    constant.ERR_UNATHORIZED_ERROR,
			Message: "Authentication token format is not match",
			Param:   "accessToken",
		}
	}

	token := strings.TrimPrefix(authHeader, BearerPrefix)
	if token == "" {
		return "", &model.ValidationError{
			Code:    constant.ERR_UNATHORIZED_ERROR,
			Message: "Authentication token is empty",
			Param:   "accessToken",
		}
	}

	return token, nil
}

// handleParseError converts JWT parsing errors to ValidationError
func handleParseError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return &model.ValidationError{
			Code:    constant.ERR_UNATHORIZED_ERROR,
			Message: "Authentication token is malformed",
			Param:   "accessToken",
		}
	case errors.Is(err, jwt.ErrTokenExpired):
		return &model.ValidationError{
			Code:    constant.ERR_UNATHORIZED_ERROR,
			Message: "Authentication token is expired",
			Param:   "accessToken",
		}
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return &model.ValidationError{
			Code:    constant.ERR_UNATHORIZED_ERROR,
			Message: "Authentication token is not valid yet",
			Param:   "accessToken",
		}
	case errors.Is(err, ErrInvalidSigningMethod):
		return &model.ValidationError{
			Code:    constant.ERR_UNATHORIZED_ERROR,
			Message: "Authentication token has invalid signing method",
			Param:   "accessToken",
		}
	default:
		return &model.ValidationError{
			Code:    constant.ERR_UNATHORIZED_ERROR,
			Message: "Authentication token is invalid",
			Param:   "accessToken",
		}
	}
}
