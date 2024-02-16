package usecase

import (
	"context"
	"errors"
	"github.com/ciazhar/golang-example/generated/db"
	"github.com/ciazhar/golang-example/internal/auth/model"
	"github.com/ciazhar/golang-example/internal/token/usecase"
	model2 "github.com/ciazhar/golang-example/internal/user/model"
	"github.com/ciazhar/golang-example/pkg/crypt"
	"github.com/ciazhar/golang-example/pkg/response"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(params db.RegisterParams) (model.TokenResponse, error)
	Login(c model.LoginForm) (model.TokenResponse, error)
	FetchRole() (row []db.FetchRoleRow, err error)
	ForgotPassword(params db.ForgotPasswordParams) error
	CheckPhoneNumber(phoneNumber string) (bool, error)
}

type authUseCase struct {
	queries      *db.Queries
	db           *pgxpool.Pool
	tokenUseCase usecase.TokenUseCase
}

// CheckPhoneNumber implements AuthUseCase
func (u authUseCase) CheckPhoneNumber(phoneNumber string) (bool, error) {
	_, err := u.queries.FindByPhone(context.Background(), phoneNumber)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, response.ErrorS("Nomer Telp Tidak sesuai")
		} else {
			return false, response.Error(err)
		}
	}
	return true, nil
}

func (u authUseCase) ForgotPassword(params db.ForgotPasswordParams) error {

	if params.Password != params.RePassword {
		return response.ErrorS("Password dan Re-Password tidak sesuai")
	}

	_, err := u.queries.FindByPhone(context.Background(), params.PhoneNumber)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.ErrorS("Nomer Telp Tidak sesuai")
		} else {
			return response.Error(err)
		}
	}

	//decypt password
	if err := crypt.DecryptPassword(&params.Password); err != nil {
		return response.Error(err)
	}
	params.RePassword = params.Password

	if err := u.queries.ForgotPassword(context.Background(), params); err != nil {
		return response.Error(err)
	}
	return response.Error(err)
}

func (u authUseCase) FetchRole() (row []db.FetchRoleRow, err error) {
	res, err := u.queries.FetchRole(context.Background())
	return res, response.Error(err)
}

func (u authUseCase) Register(params db.RegisterParams) (model.TokenResponse, error) {
	//begin transaction
	tx, err := u.db.Begin(context.Background())
	if err != nil {
		return model.TokenResponse{}, response.Error(err)
	}
	queriesTx := u.queries.WithTx(tx)
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, context.Background())

	//get param
	if params.Name == "" {
		return model.TokenResponse{}, response.Error(errors.New("name is required"))
	}
	if params.PhoneNumber == "" {
		return model.TokenResponse{}, response.Error(errors.New("phone number is required"))
	}

	//validate
	res, err := queriesTx.ValidateRegister(context.Background(), db.ValidateRegisterParams{
		PhoneNumber: params.PhoneNumber,
	})
	if err != nil {
		return model.TokenResponse{}, response.Error(err)
	}
	if res != "Validated" {
		return model.TokenResponse{}, response.Error(errors.New(res))
	}

	//decypt password
	if err := crypt.DecryptPassword(&params.Password); err != nil {
		return model.TokenResponse{}, response.Error(err)
	}

	//register auth
	res2, err := queriesTx.Register(context.Background(), params)
	if err != nil {
		return model.TokenResponse{}, response.Error(errors.New(res))
	}

	//generate token
	res3, err := u.tokenUseCase.GenerateTokens(res2.ID.String(), model2.User, res2.Name)
	if err != nil {
		return model.TokenResponse{}, response.Error(err)
	}

	//end transaction
	return res3, response.Error(tx.Commit(context.Background()))
}

func (u authUseCase) Login(c model.LoginForm) (model.TokenResponse, error) {

	//get auth by phone
	res, err := u.queries.FetchUserByPhone(context.Background(), c.PhoneNumber)
	if err != nil {
		return model.TokenResponse{}, response.Error(err)
	}

	//check if password same
	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(c.Password)); err != nil {
		return model.TokenResponse{}, response.Error(err, response.CodeUnauthorized)
	}

	//update fcm tokem auth
	if err := u.queries.UpdateUser(context.Background(), db.UpdateUserParams{
		FcmToken: c.FCMToken,
		ID:       res.ID,
	}); err != nil {
		return model.TokenResponse{}, response.Error(err)
	}

	//generate token
	res2, err := u.tokenUseCase.GenerateTokens(res.ID.String(), res.RoleName, res.Name)
	if err != nil {
		return model.TokenResponse{}, response.Error(err)
	}

	//response
	return res2, nil
}

func NewAuthUseCase(queries *db.Queries, db *pgxpool.Pool, tokenUseCase usecase.TokenUseCase) AuthUseCase {
	return authUseCase{
		queries:      queries,
		db:           db,
		tokenUseCase: tokenUseCase,
	}
}
