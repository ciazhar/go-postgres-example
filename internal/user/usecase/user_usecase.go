package usecase

import (
	"context"
	"errors"
	"github.com/ciazhar/golang-example/generated/db"
	"github.com/ciazhar/golang-example/internal/user/model"
	"github.com/ciazhar/golang-example/pkg/crypt"
	"github.com/ciazhar/golang-example/pkg/postgres"
	"github.com/ciazhar/golang-example/pkg/response"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserUseCase interface {
	GetUser(param model.FetchParam) ([]db.FetchUserRow, error)
	GetUserByID(id uuid.UUID) (db.GetUserByIdRow, error)
	UpdateUser(params db.UpdateUserParams) error
}

type userUseCase struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

func (u userUseCase) GetUserByID(id uuid.UUID) (db.GetUserByIdRow, error) {
	res, err := u.queries.GetUserById(context.Background(), id)
	return res, response.Error(err)
}

func (u userUseCase) GetUser(param model.FetchParam) ([]db.FetchUserRow, error) {
	offset, limit := postgres.ToOffsetLimit(param.Page, param.Size)
	res, err := u.queries.FetchUser(context.Background(), db.FetchUserParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, response.Error(err)
	}
	return res, nil
}

func (u userUseCase) UpdateUser(params db.UpdateUserParams) error {
	//validate
	res, err := u.queries.ValidateRegister(context.Background(), db.ValidateRegisterParams{
		ID:          params.ID,
		PhoneNumber: params.PhoneNumber,
		RoleID:      params.RoleID,
	})
	if err != nil {
		return response.Error(err)
	}
	if res != "Validated" {
		return response.Error(errors.New(res))
	}

	if params.Password != "" {
		//decypt password
		if err := crypt.DecryptPassword(&params.Password); err != nil {
			return response.Error(err)
		}
	}

	err = u.queries.UpdateUser(context.Background(), params)
	return response.Error(err)
}

func NewUserUseCase(queries *db.Queries, db *pgxpool.Pool) UserUseCase {
	return userUseCase{
		queries: queries,
		db:      db,
	}
}
