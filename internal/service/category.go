package service

import (
	"context"
	"io"

	"github.com/lemuelZara/to-grpc/internal/database"
	"github.com/lemuelZara/to-grpc/internal/pb"
	"google.golang.org/grpc"
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

func (c CategoryService) CreateCategoryStream(in grpc.ClientStreamingServer[pb.CreateCategoryRequest, pb.CategoryListResponse]) error {
	categories := pb.CategoryListResponse{}

	for {
		category, err := in.Recv()
		if err == io.EOF {
			return in.SendAndClose(&categories)
		}

		if err != nil {
			return err
		}

		ca, err := c.CategoryDB.Save(category.Name, category.Description)
		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          ca.ID,
			Name:        ca.Name,
			Description: ca.Description,
		})
	}
}

func (c CategoryService) CreateCategoryStreamBidirectional(in grpc.BidiStreamingServer[pb.CreateCategoryRequest, pb.Category]) error {
	for {
		cat, err := in.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		ca, err := c.CategoryDB.Save(cat.Name, cat.Description)
		if err != nil {
			return err
		}

		err = in.Send(&pb.Category{
			Id:          ca.ID,
			Name:        ca.Name,
			Description: ca.Description,
		})
		if err != nil {
			return err
		}
	}
}
