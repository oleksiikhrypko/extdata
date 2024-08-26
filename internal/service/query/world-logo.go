package query

import (
	"context"
	"fmt"
	"strings"

	"ext-data-domain/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/slyngshot-al/packages/storage"
	"github.com/slyngshot-al/packages/storage/psql"
)

func GetWorldLogoById(ctx context.Context, id string) (model.WorldLogo, error) {
	conn, err := storage.CtxConn(ctx)
	if err != nil {
		return model.WorldLogo{}, err
	}

	q := `select * from world_logo where id=? limit 1`

	var rec model.WorldLogo
	if err = conn.QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), id).StructScan(&rec); err != nil {
		return model.WorldLogo{}, psql.WrapError(err)
	}

	return rec, nil
}

func LockWorldLogoBySrcKey(ctx context.Context, srcKey string) (model.WorldLogo, error) {
	conn, err := storage.CtxConn(ctx)
	if err != nil {
		return model.WorldLogo{}, err
	}

	q := `select * from world_logo where src_key=? for update limit 1`

	var rec model.WorldLogo
	if err = conn.QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), srcKey).StructScan(&rec); err != nil {
		return model.WorldLogo{}, psql.WrapError(err)
	}

	return rec, nil
}

func SaveWorldLogo(ctx context.Context, input model.WorldLogoInput) (err error) {
	conn, err := storage.CtxConn(ctx)
	if err != nil {
		return err
	}

	fields, args := GetInputWorldLogoQueryArgs(input)

	fs, as := genInsertFArgs(fields)
	q := fmt.Sprintf(`insert into world_logo (%s) values (%s) on conflict (src_key) do update set %s`, fs, as, genUpdateOnConflictFArgs(fields))
	if _, err = conn.ExecContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), args...); err != nil {
		return psql.WrapError(err)
	}

	return nil
}

func GetInputWorldLogoQueryArgs(input model.WorldLogoInput) ([]string, []interface{}) {
	fields := []string{"updated_at", "id", "name", "logo_path", "src_key"}
	args := []interface{}{"now()", input.Id, input.Name, input.LogoPath, input.SrcKey}
	return fields, args
}

func DeleteWorldLogo(ctx context.Context, ids ...string) error {
	conn, err := storage.CtxConn(ctx)
	if err != nil {
		return err
	}

	q, qa, err := sqlx.In(`delete from world_logo where id in (?)`, ids)
	if err != nil {
		return psql.WrapError(err)
	}

	if _, err = conn.ExecContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), qa...); err != nil {
		return psql.WrapError(err)
	}

	return nil
}

func GetWorldLogos(ctx context.Context, ops model.WorldLogosQueryOptions, sort []psql.Sort, pg psql.Pagination) ([]model.WorldLogo, error) {
	conn, err := storage.CtxConn(ctx)
	if err != nil {
		return nil, err
	}

	var (
		qs   = []string{"1=1"} // query strings
		args []any             // query args
		q    string
		qa   []any
	)

	// filter by options
	if ops.Ids != nil {
		q, qa, err = sqlx.In(`t.id IN (?)`, ops.Ids)
		if err != nil {
			return nil, psql.WrapError(err)
		}
		qs = append(qs, q)
		args = append(args, qa...)
	}
	if ops.Search != nil {
		qs = append(qs, `t.name ILIKE ?`)
		args = append(args, fmt.Sprintf("%%%s%%", *ops.Search))
	}

	// sort + pagination
	var (
		pgq    string
		pgargs []any
	)
	if len(sort) > 0 {
		s := make([]psql.Sort, 0, len(sort))
		for i := range sort {
			if sort[i].ColumnName == nil {
				continue
			}
			switch *sort[i].ColumnName {
			case "id", "name":
				sk := fmt.Sprintf("t.%s", *sort[i].ColumnName)
				s = append(s, psql.Sort{ColumnName: &sk, Order: sort[i].Order})
			}
		}
		pgq, pgargs = psql.GenPaginationWithSorts("t.id", s, pg)
	} else {
		pgq, pgargs = psql.GenPagination("t.id", pg)
	}

	// build query filter
	q = strings.Join(qs, " AND ")
	if q != "" {
		q = "WHERE " + q
	}

	// build query with filter and pagination
	q = fmt.Sprintf(`select t.* from world_logo t %s %s`, q, pgq)

	var spaces []model.WorldLogo
	err = conn.SelectContext(ctx, &spaces, sqlx.Rebind(sqlx.DOLLAR, q), append(args, pgargs...)...)
	if err != nil {
		return nil, psql.WrapError(err)
	}
	return spaces, nil
}

func GetWorldLogosCount(ctx context.Context, ops model.WorldLogosQueryOptions) (uint64, error) {
	conn, err := storage.CtxConn(ctx)
	if err != nil {
		return 0, err
	}

	var (
		qs   []string // query strings
		args []any    // query args
		q    string
		qa   []any
	)

	// filter by options
	if ops.Ids != nil {
		q, qa, err = sqlx.In(`t.id IN (?)`, ops.Ids)
		if err != nil {
			return 0, psql.WrapError(err)
		}
		qs = append(qs, q)
		args = append(args, qa...)
	}
	if ops.Search != nil {
		qs = append(qs, `t.name ILIKE ?`)
		args = append(args, fmt.Sprintf("%%%s%%", *ops.Search))
	}

	// build query filter
	q = strings.Join(qs, " AND ")
	if q != "" {
		q = "WHERE " + q
	}

	q = fmt.Sprintf(`select count(t.id) from world_logo t %s`, q)

	var count uint64
	err = conn.QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), args...).Scan(&count)
	if err != nil {
		return 0, psql.WrapError(err)
	}

	return count, nil
}
