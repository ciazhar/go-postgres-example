package usecase

import (
	"context"
	"github.com/ciazhar/golang-example/generated/db"
	"github.com/ciazhar/golang-example/internal/auth/model"
	"github.com/ciazhar/golang-example/pkg/auth"
	"github.com/ciazhar/golang-example/pkg/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"time"
)

type TokenUseCase interface {
	GenerateTokens(uuid string, roleName string, name string) (model.TokenResponse, error)
	GetAuthCookies(accessToken, refreshToken string) (*fiber.Cookie, *fiber.Cookie)
}

type tokenUseCase struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

// GenerateTokens generates the access and refresh tokens
func (u tokenUseCase) GenerateTokens(uuid string, roleName string, name string) (model.TokenResponse, error) {
	claim, accessToken := u.GenerateAccessClaims(uuid, roleName, name)

	refreshToken, err := u.GenerateRefreshClaims(claim)
	if err != nil {
		return model.TokenResponse{}, err
	}

	return model.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    time.Unix(claim.ExpiresAt, 0),
		},
		nil
}

// GenerateAccessClaims returns a claim and a acess_token string
func (u tokenUseCase) GenerateAccessClaims(uuid string, roleName string, name string) (*auth.CustomClaims, string) {

	t := time.Now()
	claim := &auth.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    uuid,
			ExpiresAt: t.Add(1 * time.Hour * 24 * 30).Unix(),
			Subject:   "access_token",
			IssuedAt:  t.Unix(),
		},
		User: auth.UserClaims{
			Id:       uuid,
			Name:     name,
			RoleName: roleName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.key")))
	if err != nil {
		panic(err)
	}

	return claim, tokenString
}

// GenerateRefreshClaims returns refresh_token
func (u tokenUseCase) GenerateRefreshClaims(cl *auth.CustomClaims) (string, error) {
	//begin transaction
	tx, err := u.db.Begin(context.Background())
	if err != nil {
		return "", response.Error(err)
	}
	queriesTx := u.queries.WithTx(tx)
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, context.Background())

	//check if issuer exceed token quota
	exceed, err := queriesTx.IsTokenExceedLimitByIssuer(context.Background(), uuid.MustParse(cl.Issuer))
	if err != nil {
		return "", err
	}
	if exceed {
		if err := queriesTx.DeleteTokenByIssuer(context.Background(), uuid.MustParse(cl.Issuer)); err != nil {
			return "", err
		}
	}

	//create token
	t := time.Now()
	params := db.CreateTokenParams{
		Issuer:    uuid.MustParse(cl.Issuer),
		ExpiresAt: t.Add(10 * 24 * time.Hour),
		Subject:   "refresh_token",
		IssuedAt:  t,
	}
	if err := queriesTx.CreateToken(context.Background(), params); err != nil {
		return "", err
	}

	//generate token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    cl.Issuer,
		ExpiresAt: t.Add(10 * 24 * time.Hour).Unix(),
		Subject:   "refresh_token",
		IssuedAt:  t.Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(viper.GetString("jwt.key")))
	if err != nil {
		return "", err
	}

	///end transaction
	return refreshTokenString, response.Error(tx.Commit(context.Background()))
}

// GetAuthCookies sends two cookies of type access_token and refresh_token
func (u tokenUseCase) GetAuthCookies(accessToken, refreshToken string) (*fiber.Cookie, *fiber.Cookie) {
	accessCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	refreshCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(10 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	return accessCookie, refreshCookie
}

func NewTokenUseCase(queries *db.Queries, db *pgxpool.Pool) TokenUseCase {
	return tokenUseCase{
		queries: queries,
		db:      db,
	}
}
