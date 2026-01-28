package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	db "url_shortner_backend/db/output"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, params db.CreateUserParams) (db.CreateUserRow, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (db.GetUserByIDRow, error)
	CreateRefreshToken(ctx context.Context, params db.CreateRefreshTokenParams) (db.CreateRefreshTokenRow, error)
	GetRefreshToken(ctx context.Context, tokenHash string) (db.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userID pgtype.UUID) error
}

type AuthRepoImp struct {
	Queries *db.Queries
}

type AuthRepoParams struct {
	Queries *db.Queries
}

func NewAuthRepoImp(p AuthRepoParams) *AuthRepoImp {
	return &AuthRepoImp{
		Queries: p.Queries,
	}
}

func (r *AuthRepoImp) CreateUser(ctx context.Context, params db.CreateUserParams) (db.CreateUserRow, error) {
	return r.Queries.CreateUser(ctx, params)
}

func (r *AuthRepoImp) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.Queries.GetUserByEmail(ctx, email)
}

func (r *AuthRepoImp) GetUserByID(ctx context.Context, id pgtype.UUID) (db.GetUserByIDRow, error) {
	return r.Queries.GetUserByID(ctx, id)
}

func (r *AuthRepoImp) CreateRefreshToken(ctx context.Context, params db.CreateRefreshTokenParams) (db.CreateRefreshTokenRow, error) {
	return r.Queries.CreateRefreshToken(ctx, params)
}

func (r *AuthRepoImp) GetRefreshToken(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
	return r.Queries.GetRefreshToken(ctx, tokenHash)
}

func (r *AuthRepoImp) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	return r.Queries.DeleteRefreshToken(ctx, tokenHash)
}

func (r *AuthRepoImp) DeleteAllUserRefreshTokens(ctx context.Context, userID pgtype.UUID) error {
	return r.Queries.DeleteAllUserRefreshTokens(ctx, userID)
}
