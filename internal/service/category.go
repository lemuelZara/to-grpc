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

func (c CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryListResponse, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return &pb.CategoryListResponse{}, err
	}

	var clr pb.CategoryListResponse
	for _, category := range categories {
		clr.Categories = append(clr.Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return &clr, nil
}

func (c CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	ca, err := c.CategoryDB.FindByID(in.Id)
	if err != nil {
		return &pb.Category{}, err
	}

	res := &pb.Category{
		Id:          ca.ID,
		Name:        ca.Name,
		Description: ca.Description,
	}

	return res, nil
}
