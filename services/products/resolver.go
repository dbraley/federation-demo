package products

import (
	"context"
	"fmt"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Entity() EntityResolver {
	return &entityResolver{r}
}

type entityResolver struct{ *Resolver }

func (r *entityResolver) FindProductByUpc(ctx context.Context, upc string) (*Product, error) {
	for _, p := range products {
		if p.Upc == upc {
			return p, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) TopProducts(ctx context.Context, first *int) ([]*Product, error) {
	if first == nil || *first > len(products) {
		first = num(len(products))
	}
	return products[0:*first], nil
}

var products = []*Product{
	&Product{
		Upc:    "1",
		Name:   str("Table"),
		Price:  num(899),
		Weight: num(100),
	},
	&Product{
		Upc:    "2",
		Name:   str("Couch"),
		Price:  num(1299),
		Weight: num(1000),
	},
	&Product{
		Upc:    "3",
		Name:   str("Chair"),
		Price:  num(54),
		Weight: num(50),
	},
}

func str(s string) *string { return &s }

func num(i int) *int { return &i }
