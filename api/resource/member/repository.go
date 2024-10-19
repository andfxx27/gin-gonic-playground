package member

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func NewRepositorier(dbConnPool *pgxpool.Pool) Repositorier {
	return &repositorier{
		dbConnPool,
	}
}

type Repositorier interface {
	CreateMember(member *Member, ctx *gin.Context) (*Member, error)
	GetMemberByID(memberID string, ctx *gin.Context) (*Member, error)
	GetMemberByUsernameOrEmail(identifier string, ctx *gin.Context) (*Member, error)
}

type repositorier struct {
	dbConnPool *pgxpool.Pool
}

func (r *repositorier) CreateMember(member *Member, ctx *gin.Context) (*Member, error) {
	errMsg := "member.Repo query error, failed to create member."

	tag, err := r.dbConnPool.Exec(ctx, `insert into members (id, username, email, password) values ($1, $2, $3, $4)`,
		member.ID,
		member.Username,
		member.Email,
		member.Password,
	)
	if err != nil || tag.RowsAffected() == 0 {
		log.Info().Err(err).Msg(errMsg)
		return nil, err
	}

	createdMember, err := r.GetMemberByID(member.ID, ctx)
	if err != nil {
		log.Info().Err(err).Msg(errMsg)
		return nil, err
	}

	return createdMember, nil
}

func (r *repositorier) GetMemberByID(memberID string, ctx *gin.Context) (*Member, error) {
	errMsg := "member.Repo query error, failed to get member by id."

	rows, err := r.dbConnPool.Query(ctx, "select * from members where ID = $1", memberID)
	if err != nil {
		log.Error().Err(err).Msg(errMsg)
		return nil, err
	}
	member, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Member])
	if err != nil {
		log.Error().Err(err).Msg(errMsg)
		return nil, err
	}

	return &member, nil
}

func (r *repositorier) GetMemberByUsernameOrEmail(identifier string, ctx *gin.Context) (*Member, error) {
	errMsg := "member.Repo query error, failed to get member by username or email."

	rows, err := r.dbConnPool.Query(ctx, "select * from members where username = $1 or email = $2", identifier, identifier)
	if err != nil {
		log.Error().Err(err).Msg(errMsg)
		return nil, err
	}
	member, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Member])
	if err != nil {
		log.Error().Err(err).Msg(errMsg)
		return nil, err
	}

	return &member, nil
}
