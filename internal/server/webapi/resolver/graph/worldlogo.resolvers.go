package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.46

import (
	"context"

	"ext-data-domain/internal/model"
	"ext-data-domain/internal/server/webapi/api/graph"
	apimodel "ext-data-domain/internal/server/webapi/api/graph/model"
	"ext-data-domain/internal/service"

	"github.com/Slyngshot-Team/packages/auth"
	"github.com/Slyngshot-Team/packages/log"
)

// WorldLogos is the resolver for the world_logos field.
func (r *queryResolver) WorldLogos(ctx context.Context, filterOptions *apimodel.SpaceFilterOptions, search *string, orderOps []apimodel.OrderWorldLogoOps, pagination *apimodel.Pagination) (*apimodel.PaginatedWorldLogos, error) {
	userId, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, service.ErrNotFound.Consume(err)
	}
	ctx = log.CtxWithValues(ctx, "user_id", userId)

	ops := model.WorldLogosQueryOptions{
		Search: search,
	}
	if filterOptions != nil {
		ops.Ids = filterOptions.Ids
	}

	items, err := r.worldlogoService.GetWorldLogos(ctx, ops, model.FromGraphOrderWorldLogoOps(orderOps), model.FromGraphPagination(pagination))
	if err != nil {
		return nil, err
	}

	count, err := r.worldlogoService.GetWorldLogosCount(ctx, ops)
	if err != nil {
		return nil, err
	}

	return &apimodel.PaginatedWorldLogos{
		Items: model.ToGraphWorldLogos(items),
		Total: count,
	}, nil
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
