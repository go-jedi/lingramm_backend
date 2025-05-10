package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/go-jedi/lingvogramm_backend/config"
	"github.com/go-jedi/lingvogramm_backend/pkg/uuid"
	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultSecretHashLen = 20
	defaultAccessExpAt   = 5
	defaultRefreshExpAt  = 30
)

var (
	ErrTokenSigningMethod = errors.New("unexpected token signing method")
	ErrTokenInvalid       = errors.New("invalid token")
	ErrTokenClaims        = errors.New("unexpected token claims")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenTelegramID    = errors.New("unexpected token telegram id")
)

// IJWT defines the interface for the jwt.
//
//go:generate mockery --name=IJWT --output=mocks --case=underscore
type IJWT interface {
	Generate(telegramID string) (GenerateResp, error)
	Verify(telegramID string, token string) (VerifyResp, error)
	ParseToken(token string) (VerifyResp, error)
}

type tokenClaims struct {
	TelegramID string `json:"telegram_id"`
	jwt.RegisteredClaims
}

type JWT struct {
	// uuid need for generate crypto hash
	uuid uuid.IUUID
	// secret key need for token signing
	secret []byte
	// secretHashLen need to generate hash
	secretHashLen int
	// accessExpAt expiration time in minutes
	accessExpAt int
	// refreshExpAt expiration time in days
	refreshExpAt int
}

func New(cfg config.JWTConfig, uuid uuid.IUUID) (*JWT, error) {
	j := &JWT{
		uuid:          uuid,
		secretHashLen: cfg.SecretHashLen,
		accessExpAt:   cfg.AccessExpAt,
		refreshExpAt:  cfg.RefreshExpAt,
	}

	if err := j.init(); err != nil {
		return nil, err
	}

	if err := j.generateSecretKey(cfg.SecretPath); err != nil {
		return nil, err
	}

	return j, nil
}

func (j *JWT) init() error {
	if j.secretHashLen == 0 {
		j.secretHashLen = defaultSecretHashLen
	}

	if j.accessExpAt == 0 {
		j.accessExpAt = defaultAccessExpAt
	}

	if j.refreshExpAt == 0 {
		j.refreshExpAt = defaultRefreshExpAt
	}

	return nil
}

type GenerateResp struct {
	AccessToken  string
	RefreshToken string
	AccessExpAt  time.Time
	RefreshExpAt time.Time
}

// Generate token.
func (j *JWT) Generate(telegramID string) (GenerateResp, error) {
	aExpAt := j.getAccessExpAt()
	rExpAt := j.getRefreshExpAt()

	aToken, err := j.createToken(telegramID, aExpAt)
	if err != nil {
		return GenerateResp{}, err
	}

	rToken, err := j.createToken(telegramID, rExpAt)
	if err != nil {
		return GenerateResp{}, err
	}

	return GenerateResp{
		AccessToken:  aToken,
		RefreshToken: rToken,
		AccessExpAt:  aExpAt,
		RefreshExpAt: rExpAt,
	}, nil
}

type VerifyResp struct {
	TelegramID string
	ExpAt      time.Time
}

// Verify token.
func (j *JWT) Verify(telegramID string, token string) (VerifyResp, error) {
	// parse the token
	t, err := jwt.ParseWithClaims(
		token,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrTokenSigningMethod
			}
			return j.secret, nil
		})
	if err != nil {
		return VerifyResp{}, err
	}

	// check token valid
	if !t.Valid {
		return VerifyResp{}, ErrTokenInvalid
	}

	// extract the claims
	c, ok := t.Claims.(*tokenClaims)
	if !ok {
		return VerifyResp{}, ErrTokenClaims
	}

	// check expired token
	if c.ExpiresAt != nil && time.Now().After(c.ExpiresAt.Time) {
		return VerifyResp{}, ErrTokenExpired
	}

	// compare telegram id with telegram id in token
	if telegramID != c.TelegramID {
		return VerifyResp{}, ErrTokenTelegramID
	}

	return VerifyResp{
		TelegramID: c.TelegramID,
		ExpAt:      c.ExpiresAt.Time,
	}, nil
}

// ParseToken parse token.
func (j *JWT) ParseToken(token string) (VerifyResp, error) {
	// parse the token
	t, err := jwt.ParseWithClaims(
		token,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrTokenSigningMethod
			}
			return j.secret, nil
		})
	if err != nil {
		return VerifyResp{}, err
	}

	// check token valid
	if !t.Valid {
		return VerifyResp{}, ErrTokenInvalid
	}

	// extract the claims
	c, ok := t.Claims.(*tokenClaims)
	if !ok {
		return VerifyResp{}, ErrTokenClaims
	}

	// check expired token
	if c.ExpiresAt != nil && time.Now().After(c.ExpiresAt.Time) {
		return VerifyResp{}, ErrTokenExpired
	}

	return VerifyResp{
		TelegramID: c.TelegramID,
		ExpAt:      c.ExpiresAt.Time,
	}, nil
}

// createToken create token.
func (j *JWT) createToken(telegramID string, expAt time.Time) (string, error) {
	// create the claims
	c := tokenClaims{
		TelegramID: telegramID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// create token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateSecretKey generate secret key.
func (j *JWT) generateSecretKey(secretPath string) error {
	ie, err := j.fileExists(secretPath)
	if err != nil {
		return err
	}

	if ie {
		fb, err := os.ReadFile(secretPath)
		if err != nil {
			return err
		}
		j.secret = fb
		return nil
	}

	u, err := j.uuid.Generate()
	if err != nil {
		return err
	}

	j.secret = []byte(u)

	const mode = 0o600
	if err := os.WriteFile(secretPath, j.secret, os.FileMode(mode)); err != nil {
		return err
	}

	return nil
}

// getAccessExpAt get access expires at token time.
func (j *JWT) getAccessExpAt() time.Time {
	return time.Now().Add(time.Duration(j.accessExpAt) * time.Minute)
}

// getRefreshExpAt get refresh expires at token time.
func (j *JWT) getRefreshExpAt() time.Time {
	return time.Now().Add(time.Duration(j.refreshExpAt) * 24 * time.Hour)
}

// fileExists check file exists.
func (j *JWT) fileExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !fi.IsDir(), nil
}
