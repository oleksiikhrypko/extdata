package query

import (
	"context"
	"fmt"
	"os"
	"testing"

	"ext-data-domain/internal/model"

	"github.com/Slyngshot-Team/packages/storage"
	"github.com/Slyngshot-Team/packages/storage/psql"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

const DriverName = "postgres"

var db *sqlx.DB

func TestMain(m *testing.M) {
	if os.Getenv("INTEGRATION_TEST") != "YES" {
		os.Exit(0)
	}
	var err error

	db, err = sqlx.Connect(DriverName, os.Getenv("DATABASE_DSN"))
	if err != nil {
		fmt.Println("DB connection error:", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func ptr[T any](v T) *T {
	return &v
}

func Test_WorldLogo(t *testing.T) {
	ctx := context.Background()
	ctx = storage.InitCtxConn(ctx, db)

	id := ulid.Make().String()
	err := SaveWorldLogo(ctx, model.WorldLogoInput{
		Id:       id,
		Name:     "name 1",
		LogoPath: "logo url 1",
		SrcKey:   "key1",
	})
	require.NoError(t, err)

	rec, err := GetWorldLogoById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, rec.Id)
	require.Equal(t, "name 1", rec.Name)
	require.Equal(t, "logo url 1", rec.LogoPath)
	require.Equal(t, "key1", rec.SrcKey)

	rec, err = LockWorldLogoBySrcKey(ctx, "key1")
	require.NoError(t, err)
	require.Equal(t, id, rec.Id)
	require.Equal(t, "name 1", rec.Name)
	require.Equal(t, "logo url 1", rec.LogoPath)
	require.Equal(t, "key1", rec.SrcKey)

	err = DeleteWorldLogo(ctx, id)
	require.NoError(t, err)
	_, err = GetWorldLogoById(ctx, id)
	require.Error(t, err)
	require.ErrorIs(t, err, storage.ErrNotFound)

	for i := 0; i < 10; i++ {
		err = SaveWorldLogo(ctx, model.WorldLogoInput{
			Id:       ulid.Make().String(),
			Name:     fmt.Sprintf("name %d", i),
			LogoPath: fmt.Sprintf("logo url %d", i),
			SrcKey:   fmt.Sprintf("key%d", i),
		})
		require.NoError(t, err)
	}

	ops := model.WorldLogosQueryOptions{
		Search: ptr("name 2"),
		Ids:    nil,
	}
	recs, err := GetWorldLogos(ctx, ops, nil, psql.Pagination{
		Limit: 10,
	})
	require.NoError(t, err)
	require.Len(t, recs, 1)

	c, err := GetWorldLogosCount(ctx, ops)
	require.NoError(t, err)
	require.Equal(t, uint64(1), c)

	recs, err = GetWorldLogos(ctx, model.WorldLogosQueryOptions{}, []psql.Sort{
		{ColumnName: ptr("name"), Order: nil},
	}, psql.Pagination{Limit: 10})
	require.NoError(t, err)
	require.Len(t, recs, 10)

	c, err = GetWorldLogosCount(ctx, model.WorldLogosQueryOptions{})
	require.NoError(t, err)
	require.True(t, c >= 10)
}
