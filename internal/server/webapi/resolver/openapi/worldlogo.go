package openapi

import (
	"net/http"

	"ext-data-domain/internal/model"
	api "ext-data-domain/internal/server/webapi/api/openapi"

	"github.com/Slyngshot-Team/packages/auth"
	"github.com/Slyngshot-Team/packages/log"
	echo "github.com/labstack/echo/v4"
)

func (h Handler) GetWorldLogoById(c echo.Context, id api.IdParam) error {
	ctx := c.Request().Context()
	userId, err := auth.GetUserID(ctx)
	if err != nil {
		return handleError(c, err)
	}
	ctx = log.CtxWithValues(ctx, "user_id", userId)

	data, err := h.worldlogoService.GetWorldLogoById(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, model.ToAPIWorldLogo(data))

}
