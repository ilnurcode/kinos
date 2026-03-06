// Package grpcserver предоставляет gRPC-сервер для catalog-service.
// Реализует методы CatalogService: управление категориями, производителями и товарами.
package grpcserver

import (
	"context"
	"kinos/catalog-service/internal/models"
	"kinos/catalog-service/internal/service"
	pb "kinos/proto/catalog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CatalogServer struct {
	pb.UnimplementedCatalogServiceServer
	prodSvc *service.ProductService
	catSvc  *service.CategoryService
	manSvc  *service.ManufacturersService
}

func NewCatalogServer(prodSvc *service.ProductService, catSvc *service.CategoryService, manSvc *service.ManufacturersService) *CatalogServer {
	return &CatalogServer{
		prodSvc: prodSvc,
		catSvc:  catSvc,
		manSvc:  manSvc,
	}
}

//Категории

func (s *CatalogServer) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := s.catSvc.CreateCategory(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create category failed: %v", err.Error())
	}
	return &pb.Category{Id: category.Id, Name: category.Name}, nil
}

func (s *CatalogServer) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error) {
	category, err := s.catSvc.UpdateCategory(ctx, req.Id, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update category failed: %v", err.Error())
	}
	return &pb.Category{Id: category.Id, Name: category.Name}, nil
}

func (s *CatalogServer) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if err := s.catSvc.DeleteCategory(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete category failed: %v", err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *CatalogServer) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
	category, err := s.catSvc.GetCategory(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "category not found")
	}
	return &pb.Category{Id: category.Id, Name: category.Name}, nil
}

func (s *CatalogServer) GetListCategory(ctx context.Context, req *pb.GetListCategoryRequest) (*pb.ListCategoriesResponse, error) {
	categories, total, err := s.catSvc.GetListCategory(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list category failed: %v", err.Error())
	}
	var pbCategories []*pb.Category
	for _, category := range categories {
		pbCategories = append(pbCategories, &pb.Category{Id: category.Id, Name: category.Name})
	}
	return &pb.ListCategoriesResponse{Category: pbCategories, Total: total}, nil
}

//Производители

func (s *CatalogServer) CreateManufacturer(ctx context.Context, req *pb.CreateManufacturerRequest) (*pb.Manufacturer, error) {
	manufacturer, err := s.manSvc.CreateManufacturer(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create manufacturer failed: %v", err.Error())
	}
	return &pb.Manufacturer{Id: manufacturer.Id, Name: manufacturer.Name}, nil
}

func (s *CatalogServer) UpdateManufacturer(ctx context.Context, req *pb.UpdateManufacturerRequest) (*pb.Manufacturer, error) {
	manufacturer, err := s.manSvc.UpdateManufacturer(ctx, req.Id, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update manufacturer failed: %v", err.Error())
	}
	return &pb.Manufacturer{Id: manufacturer.Id, Name: manufacturer.Name}, nil
}

func (s *CatalogServer) DeleteManufacturer(ctx context.Context, req *pb.DeleteManufacturerRequest) (*emptypb.Empty, error) {
	if err := s.manSvc.DeleteManufacturer(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete manufacturer failed: %v", err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *CatalogServer) GetManufacturer(ctx context.Context, req *pb.GetManufacturerRequest) (*pb.Manufacturer, error) {
	manufacturer, err := s.manSvc.GetManufacturer(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "manufacturer not found")
	}
	return &pb.Manufacturer{Id: manufacturer.Id, Name: manufacturer.Name}, nil
}

func (s *CatalogServer) GetListManufacturers(ctx context.Context, req *pb.GetListManufacturerRequest) (*pb.ListManufacturerResponse, error) {
	manufacturer, total, err := s.manSvc.GetListManufacturers(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list manufacturer failed: %v", err.Error())
	}
	var pbManufacturer []*pb.Manufacturer
	for _, category := range manufacturer {
		pbManufacturer = append(pbManufacturer, &pb.Manufacturer{Id: category.Id, Name: category.Name})
	}
	return &pb.ListManufacturerResponse{Manufacturer: pbManufacturer, Total: total}, nil
}

//Товары

func (s *CatalogServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	product, err := s.prodSvc.CreateProduct(ctx, req.Name, req.ManufacturerId, req.CategoryId, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create product failed: %v", err.Error())
	}
	return &pb.Product{Id: product.Id, Name: product.Name, ManufacturerId: product.ManufacturersId, CategoryId: product.CategoryId, Price: product.Price}, nil
}

func (s *CatalogServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	product, err := s.prodSvc.UpdateProduct(ctx, req.Id, req.Name, req.ManufacturerId, req.CategoryId, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update product failed: %v", err.Error())
	}
	return &pb.Product{Id: product.Id, Name: product.Name, ManufacturerId: product.ManufacturersId, CategoryId: product.CategoryId, Price: product.Price}, nil
}

func (s *CatalogServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	if err := s.prodSvc.DeleteProduct(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete product failed: %v", err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *CatalogServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	product, err := s.prodSvc.GetProduct(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}
	return &pb.Product{Id: product.Id, Name: product.Name, ManufacturerId: product.ManufacturersId, CategoryId: product.CategoryId, Price: product.Price}, nil
}

func (s *CatalogServer) GetListProduct(ctx context.Context, req *pb.GetListProductRequest) (*pb.ListProductsResponse, error) {
	filter := models.ProductFilter{}
	if req.CategoryId != 0 {
		filter.CategoryId = &req.CategoryId
	}
	if req.ManufacturerId != 0 {
		filter.ManufacturersId = &req.ManufacturerId
	}
	if req.PriceMin != 0 {
		filter.PriceMin = &req.PriceMin
	}
	if req.PriceMax != 0 {
		filter.PriceMax = &req.PriceMax
	}
	if req.NameContains != "" {
		filter.NameContains = &req.NameContains
	}
	product, total, err := s.prodSvc.GetListProduct(ctx, filter, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list product failed: %v", err.Error())
	}
	var pbProducts []*pb.Product
	for _, product := range product {
		pbProducts = append(pbProducts, &pb.Product{
			Id:             product.Id,
			Name:           product.Name,
			ManufacturerId: product.ManufacturersId,
			CategoryId:     product.CategoryId,
			Price:          product.Price,
		})
	}
	return &pb.ListProductsResponse{Product: pbProducts, Total: total}, nil
}
