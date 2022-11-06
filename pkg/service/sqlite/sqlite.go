package sqlite

import (
	"context"
	"database/sql"

	logging "github.com/ipfs/go-log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var sqLog = logging.Logger("sqlite")

type YetiSQLiteService struct {
	ctx context.Context
	db  *bun.DB
}

func (r *YetiSQLiteService) GetCtx() context.Context {
	return r.ctx
}

func (r *YetiSQLiteService) GetDB() *bun.DB {
	return r.db
}

func NewYetiSQLiteService(file string) (*YetiSQLiteService, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, file)
	if err != nil {
		return nil, err
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())

	r := YetiSQLiteService{ctx: context.Background(), db: db}
	return &r, nil
}
