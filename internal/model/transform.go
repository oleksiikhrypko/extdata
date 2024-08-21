package model

import (
	"net/url"
	"strings"

	graph "ext-data-domain/internal/server/webapi/api/graph/model"
	api "ext-data-domain/internal/server/webapi/api/openapi"

	"github.com/Slyngshot-Team/packages/storage/psql"
)

func Ptr[T any](v T) *T {
	return &v
}

func Val[T any](v *T) T {
	var res T
	if v == nil {
		return res
	}
	return *v
}

func UrlJoinPath(base string, path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	res, _ := url.JoinPath(base, path)
	return res
}

func FromGraphPagination(input *graph.Pagination) psql.Pagination {
	if input == nil {
		return psql.Pagination{
			Limit: 25,
		}
	}

	return psql.Pagination{
		OffSetKey: input.OffsetKey,
		Offset:    input.Offset,
		Limit:     input.Limit,
	}
}

func FromGraphOrderWorldLogoOps(input []graph.OrderWorldLogoOps) []psql.Sort {
	if len(input) == 0 {
		return nil
	}
	res := make([]psql.Sort, len(input))
	for i, order := range input {
		var d *string
		if order.Direction.IsValid() {
			d = Ptr(order.Direction.String())
		}
		var c *string
		if order.Field.IsValid() {
			c = Ptr(order.Field.String())
		}
		res[i] = psql.Sort{
			ColumnName: c,
			Order:      d,
		}
	}
	return res
}

func ToGraphWorldLogos(input []WorldLogo) []graph.WorldLogo {
	if len(input) == 0 {
		return nil
	}
	res := make([]graph.WorldLogo, len(input))
	for i, item := range input {
		res[i] = ToGraphWorldLogo(item)
	}
	return res
}

func ToGraphWorldLogo(input WorldLogo) graph.WorldLogo {
	return graph.WorldLogo{
		ID:        input.Id,
		Name:      input.Name,
		LogoPath:  input.LogoPath,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
}

func ToAPIWorldLogo(input WorldLogo) api.WorldLogo {
	return api.WorldLogo{
		Id:        input.Id,
		Name:      input.Name,
		LogoPath:  input.LogoPath,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
}
