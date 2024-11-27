package service

import (
	"context"

	"github.com/lemuelZara/to-grpc/internal/database"
	"github.com/lemuelZara/to-grpc/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) CategoryService {
	return CategoryService{CategoryDB: categoryDB}
}

func (c CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	ca, err := c.CategoryDB.Save(in.Name, in.Description)
	if err != nil {
		return &pb.CreateCategoryResponse{}, err
	}

	res := &pb.CreateCategoryResponse{
		Category: &pb.Category{
			Id:          ca.ID,
			Name:        ca.Name,
			Description: ca.Description,
		},
	}

	return res, nil
}
