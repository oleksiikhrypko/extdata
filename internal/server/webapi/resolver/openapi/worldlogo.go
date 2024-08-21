package openapi

import (
	"net/http"

	"ext-data-domain/internal/model"
	api "ext-data-domain/internal/server/webapi/api/openapi"

	"github.com/Slyngshot-Team/packages/auth"
	"github.com/Slyngshot-Team/packages/log"
	"github.com/Slyngshot-Team/packages/storage/psql"
	echo "github.com/labstack/echo/v4"
)

func (h Handler) GetWorldLogoById(c echo.Context, id api.IdParam) error {
	ctx := c.Request().Context()
	if userId, err := auth.GetUserID(ctx); err != nil {
		ctx = log.CtxWithValues(ctx, "user_id", userId)
	}

	data, err := h.worldlogoService.GetWorldLogoById(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, model.ToAPIWorldLogo(data))
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
