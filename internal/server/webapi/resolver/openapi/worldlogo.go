package openapi

import (
	"net/http"

	"ext-data-domain/internal/model"
	api "ext-data-domain/internal/server/webapi/api/openapi"

	echo "github.com/labstack/echo/v4"
	"github.com/slyngshot-al/packages/auth"
	"github.com/slyngshot-al/packages/log"
	"github.com/slyngshot-al/packages/storage/psql"
)

func (h Handler) GetWorldLogoById(c echo.Context, id api.IdParam) error {
	ctx := c.Request().Context()
	if userId, err := auth.GetUserID(ctx); err != nil {
		ctx = log.CtxWithValues(ctx, "user_id", userId)
	}

	rec, err := h.worldlogoService.GetWorldLogoById(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, model.ToAPIWorldLogo(rec))
}

func (h Handler) GetWorldLogos(c echo.Context, params api.GetWorldLogosParams) error {
	ctx := c.Request().Context()
	if userId, err := auth.GetUserID(ctx); err != nil {
		ctx = log.CtxWithValues(ctx, "user_id", userId)
	}

	// query options
	var ids []string
	if params.Ids != nil {
		ids = *params.Ids
	}
	query := model.WorldLogosQueryOptions{
		Search: params.Search,
		Ids:    ids,
	}
	// order options
	var sort []psql.Sort
	if params.SortBy != nil {
		switch {
		case params.SortOrder == nil:
			sort = append(sort, psql.Sort{ColumnName: model.Ptr(string(*params.SortBy))})
		default:
			sort = append(sort, psql.Sort{ColumnName: model.Ptr(string(*params.SortBy)), Order: model.Ptr(string(*params.SortOrder))})
		}
	}
	// pagination
	p := psql.Pagination{OffSetKey: params.OffsetKey, Offset: params.Offset, Limit: params.Limit}

	recs, err := h.worldlogoService.GetWorldLogos(ctx, query, sort, p)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, model.ToAPIWorldLogos(recs))
}

func (h Handler) CreateWorldLogo(c echo.Context, params api.CreateWorldLogoParams) error {
	ctx := c.Request().Context()
	if userId, err := auth.GetUserID(ctx); err != nil {
		ctx = log.CtxWithValues(ctx, "user_id", userId)
	}

	var input api.WorldLogoInput
	err := c.Bind(&input)
	if err != nil {
		return bindError(c, err)
	}

	id, err := h.worldlogoService.SaveWorldLogo(ctx, params.XAPIKEY, model.SaveWorldLogoInput{
		Name:          input.Name,
		LogoBase64Str: input.LogoBase64Str,
		SrcKey:        input.SrcKey,
		ContentType:   input.ContentType,
		FileExtension: input.FileExtension,
	})
	if err != nil {
		return handleError(c, err)
	}

	rec, err := h.worldlogoService.GetWorldLogoById(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, model.ToAPIWorldLogo(rec))
}

func (h Handler) DeleteWorldLogoById(c echo.Context, id api.IdParam, params api.DeleteWorldLogoByIdParams) error {
	ctx := c.Request().Context()
	if userId, err := auth.GetUserID(ctx); err != nil {
		ctx = log.CtxWithValues(ctx, "user_id", userId)
	}

	err := h.worldlogoService.DeleteWorldLogo(ctx, params.XAPIKEY, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusOK)
}
