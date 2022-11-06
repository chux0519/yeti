package sqlite

import (
	"time"

	"github.com/uptrace/bun"
)

type MapleGG struct {
	bun.BaseModel `bun:"table:maplegg,alias:m"`

	IGN        string                 `bun:",pk"`
	Data       map[string]interface{} `bun:"type:json"` // json format raw data
	ProfileImg []byte                 `bun:"type:blob"`

	CreatedAt string `bun:",nullzero,notnull"` //ISO8601 utc time
	UpdatedAt string `bun:",nullzero,notnull"`
}

func (r *YetiSQLiteService) UpsertMapleGG(
	ign string,
	data map[string]interface{},
	img []byte,
) (*MapleGG, error) {
	exists, err := r.CheckMapleGGByIGN(ign)
	if err != nil {
		return nil, err
	}
	if exists {
		rank, err := r.GetMapleGGByIGNWithoutExpire(ign)
		if err != nil {
			return nil, err
		}
		rank.Data = data
		rank.ProfileImg = img
		return r.UpdateMapleGG(rank)
	} else {
		return r.NewMapleGG(ign, data, img)
	}
}

func (r *YetiSQLiteService) NewMapleGG(
	ign string,
	data map[string]interface{},
	img []byte,
) (*MapleGG, error) {
	now := time.Now().Format(time.RFC3339)
	record := &MapleGG{
		IGN:        ign,
		Data:       data,
		ProfileImg: img,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if _, err := r.db.NewInsert().Model(record).Exec(r.ctx); err != nil {
		sqLog.Error(err)
		return nil, err
	}
	return record, nil
}

func (r *YetiSQLiteService) UpdateMapleGG(record *MapleGG) (*MapleGG, error) {
	record.UpdatedAt = time.Now().Format(time.RFC3339)
	q := r.db.NewUpdate().Model(record).WherePK().Returning("*")
	if _, err := q.Exec(r.ctx); err != nil {
		sqLog.Debug(err)
		return nil, err
	}
	return record, nil
}

func (r *YetiSQLiteService) GetMapleGGByIGNWithoutExpire(ign string) (*MapleGG, error) {
	var rank MapleGG

	sql := r.db.NewSelect().
		Model(&rank).
		Where("m.ign = ?", ign).
		Limit(1)

	sqLog.Debug(sql.String())

	if err := sql.Scan(r.ctx); err != nil {
		sqLog.Error(err)
		return nil, err
	}

	return &rank, nil
}

func (r *YetiSQLiteService) GetMapleGGByIGN(ign string) (*MapleGG, error) {
	now := time.Now()
	expireTime := now.Add(-8 * time.Hour)

	var rank MapleGG

	sql := r.db.NewSelect().
		Model(&rank).
		Where("ign = ?", ign).
		Where("updated_at > ?", expireTime.Format(time.RFC3339)).
		Limit(1)

	sqLog.Debug(sql.String())

	if err := sql.Scan(r.ctx); err != nil {
		sqLog.Error(err)
		return nil, err
	}

	return &rank, nil
}

func (r *YetiSQLiteService) CheckMapleGGByIGN(ign string) (bool, error) {
	now := time.Now()
	expireTime := now.Add(-8 * time.Hour)

	var rank MapleGG

	return r.db.NewSelect().
		Model(&rank).
		Where("m.ign = ?", ign).
		Where("m.updated_at > ?", expireTime.Format(time.RFC3339)).
		Exists(r.ctx)
}
