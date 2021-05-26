package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
)

func (r *locationResolver) Region(ctx context.Context, obj *model.Location) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *locationResolver) Country(ctx context.Context, obj *model.Location) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) PersonalInfo(ctx context.Context, obj *model.Miner) (*model.PersonalInfo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Worker(ctx context.Context, obj *model.Miner) (*model.Worker, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Owner(ctx context.Context, obj *model.Miner) (*model.Owner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Location(ctx context.Context, obj *model.Miner) (*model.Location, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) QualityAdjustedPower(ctx context.Context, obj *model.Miner) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Service(ctx context.Context, obj *model.Miner) (*model.Service, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Pricing(ctx context.Context, obj *model.Miner) (*model.Pricing, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) ReputationScore(ctx context.Context, obj *model.Miner) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) TransparencyScore(ctx context.Context, obj *model.Miner) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ClaimProfile(ctx context.Context, input model.ProfileClaimInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) EditProfile(ctx context.Context, input model.ProfileSettingsInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *personalInfoResolver) Name(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *personalInfoResolver) Bio(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *personalInfoResolver) Email(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *personalInfoResolver) Website(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *personalInfoResolver) Twitter(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *personalInfoResolver) Slack(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pricingResolver) StorageAskPrice(ctx context.Context, obj *model.Pricing) (float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pricingResolver) VerifiedAskPrice(ctx context.Context, obj *model.Pricing) (float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pricingResolver) RetrievalAskPrice(ctx context.Context, obj *model.Pricing) (float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Miner(ctx context.Context, id string) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Miners(ctx context.Context) ([]*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *serviceResolver) ServiceTypes(ctx context.Context, obj *model.Service) (*model.ServiceTypes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *serviceResolver) DataTransferMechanism(ctx context.Context, obj *model.Service) (*model.DataTransferMechanism, error) {
	panic(fmt.Errorf("not implemented"))
}

// Location returns generated.LocationResolver implementation.
func (r *Resolver) Location() generated.LocationResolver { return &locationResolver{r} }

// Miner returns generated.MinerResolver implementation.
func (r *Resolver) Miner() generated.MinerResolver { return &minerResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// PersonalInfo returns generated.PersonalInfoResolver implementation.
func (r *Resolver) PersonalInfo() generated.PersonalInfoResolver { return &personalInfoResolver{r} }

// Pricing returns generated.PricingResolver implementation.
func (r *Resolver) Pricing() generated.PricingResolver { return &pricingResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Service returns generated.ServiceResolver implementation.
func (r *Resolver) Service() generated.ServiceResolver { return &serviceResolver{r} }

type locationResolver struct{ *Resolver }
type minerResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type personalInfoResolver struct{ *Resolver }
type pricingResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type serviceResolver struct{ *Resolver }
