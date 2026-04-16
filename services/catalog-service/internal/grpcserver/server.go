// Package grpcserver предоставляет gRPC-сервер для catalog-service.
// Реализует методы CatalogService: управление категориями, производителями и товарами.
package grpcserver

import (
	"context"
	"errors"

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

// Категории
func (s *CatalogServer) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "название категории обязательно")
	}
	category, err := s.catSvc.CreateCategory(ctx, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrCategoryExists) {
			return nil, status.Error(codes.AlreadyExists, "категория с таким названием уже существует")
		}
		return nil, status.Errorf(codes.Internal, "ошибка создания категории: %v", err)
	}
	return &pb.Category{Id: category.ID, Name: category.Name}, nil
}

func (s *CatalogServer) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "название категории обязательно")
	}
	category, err := s.catSvc.UpdateCategory(ctx, req.Id, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "категория не найдена")
		}
		if errors.Is(err, service.ErrCategoryExists) {
			return nil, status.Error(codes.AlreadyExists, "категория с таким названием уже существует")
		}
		return nil, status.Errorf(codes.Internal, "ошибка обновления категории: %v", err)
	}
	return &pb.Category{Id: category.ID, Name: category.Name}, nil
}

func (s *CatalogServer) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if err := s.catSvc.DeleteCategory(ctx, req.Id); err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "категория не найдена")
		}
		return nil, status.Errorf(codes.Internal, "ошибка удаления категории: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *CatalogServer) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "название категории обязательно")
	}
	category, err := s.catSvc.GetCategory(ctx, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "категория не найдена")
		}
		return nil, status.Errorf(codes.Internal, "ошибка поиска категории: %v", err)
	}
	return &pb.Category{Id: category.ID, Name: category.Name}, nil
}

func (s *CatalogServer) GetListCategory(ctx context.Context, req *pb.GetListCategoryRequest) (*pb.ListCategoriesResponse, error) {
	categories, total, err := s.catSvc.GetListCategory(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения списка категорий: %v", err)
	}
	var pbCategories []*pb.Category
	for _, category := range categories {
		pbCategories = append(pbCategories, &pb.Category{Id: category.ID, Name: category.Name})
	}
	return &pb.ListCategoriesResponse{Category: pbCategories, Total: total}, nil
}

// Производители
func (s *CatalogServer) CreateManufacturer(ctx context.Context, req *pb.CreateManufacturerRequest) (*pb.Manufacturer, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "название производителя обязательно")
	}
	manufacturer, err := s.manSvc.CreateManufacturer(ctx, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrManufacturerExists) {
			return nil, status.Error(codes.AlreadyExists, "производитель с таким названием уже существует")
		}
		return nil, status.Errorf(codes.Internal, "ошибка создания производителя: %v", err)
	}
	return &pb.Manufacturer{Id: manufacturer.ID, Name: manufacturer.Name}, nil
}

func (s *CatalogServer) UpdateManufacturer(ctx context.Context, req *pb.UpdateManufacturerRequest) (*pb.Manufacturer, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "название производителя обязательно")
	}
	manufacturer, err := s.manSvc.UpdateManufacturer(ctx, req.Id, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrManufacturerNotFound) {
			return nil, status.Error(codes.NotFound, "производитель не найден")
		}
		if errors.Is(err, service.ErrManufacturerExists) {
			return nil, status.Error(codes.AlreadyExists, "производитель с таким названием уже существует")
		}
		return nil, status.Errorf(codes.Internal, "ошибка обновления производителя: %v", err)
	}
	return &pb.Manufacturer{Id: manufacturer.ID, Name: manufacturer.Name}, nil
}

func (s *CatalogServer) DeleteManufacturer(ctx context.Context, req *pb.DeleteManufacturerRequest) (*emptypb.Empty, error) {
	if err := s.manSvc.DeleteManufacturer(ctx, req.Id); err != nil {
		if errors.Is(err, service.ErrManufacturerNotFound) {
			return nil, status.Error(codes.NotFound, "производитель не найден")
		}
		return nil, status.Errorf(codes.Internal, "ошибка удаления производителя: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *CatalogServer) GetManufacturer(ctx context.Context, req *pb.GetManufacturerRequest) (*pb.Manufacturer, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "название производителя обязательно")
	}
	manufacturer, err := s.manSvc.GetManufacturer(ctx, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrManufacturerNotFound) {
			return nil, status.Error(codes.NotFound, "производитель не найден")
		}
		return nil, status.Errorf(codes.Internal, "ошибка поиска производителя: %v", err)
	}
	return &pb.Manufacturer{Id: manufacturer.ID, Name: manufacturer.Name}, nil
}

func (s *CatalogServer) GetListManufacturers(ctx context.Context, req *pb.GetListManufacturerRequest) (*pb.ListManufacturerResponse, error) {
	manufacturer, total, err := s.manSvc.GetListManufacturers(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения списка производителей: %v", err)
	}
	var pbManufacturer []*pb.Manufacturer
	for _, m := range manufacturer {
		pbManufacturer = append(pbManufacturer, &pb.Manufacturer{Id: m.ID, Name: m.Name})
	}
	return &pb.ListManufacturerResponse{Manufacturer: pbManufacturer, Total: total}, nil
}

// Товары
func (s *CatalogServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "название товара обязательно")
	}
	if req.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "цена должна быть больше нуля")
	}
	product, err := s.prodSvc.CreateProduct(ctx, req.Name, req.ManufacturerId, req.CategoryId, req.Price)
	if err != nil {
		if errors.Is(err, service.ErrProductExists) {
			return nil, status.Error(codes.AlreadyExists, "товар с таким названием уже существует")
		}
		if errors.Is(err, service.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "категория не найдена")
		}
		if errors.Is(err, service.ErrManufacturerNotFound) {
			return nil, status.Error(codes.NotFound, "производитель не найден")
		}
		return nil, status.Errorf(codes.Internal, "ошибка создания товара: %v", err)
	}
	return &pb.Product{Id: product.ID, Name: product.Name, ManufacturerId: product.ManufacturersID, CategoryId: product.CategoryID, Price: product.Price}, nil
}

func (s *CatalogServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "название товара обязательно")
	}
	if req.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "цена должна быть больше нуля")
	}
	product, err := s.prodSvc.UpdateProduct(ctx, req.Id, req.Name, req.ManufacturerId, req.CategoryId, req.Price)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "товар не найден")
		}
		if errors.Is(err, service.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "категория не найдена")
		}
		if errors.Is(err, service.ErrManufacturerNotFound) {
			return nil, status.Error(codes.NotFound, "производитель не найден")
		}
		return nil, status.Errorf(codes.Internal, "ошибка обновления товара: %v", err)
	}
	return &pb.Product{Id: product.ID, Name: product.Name, ManufacturerId: product.ManufacturersID, CategoryId: product.CategoryID, Price: product.Price}, nil
}

func (s *CatalogServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	if err := s.prodSvc.DeleteProduct(ctx, req.Id); err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "товар не найден")
		}
		return nil, status.Errorf(codes.Internal, "ошибка удаления товара: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *CatalogServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "название товара обязательно")
	}
	product, err := s.prodSvc.GetProduct(ctx, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "товар не найден")
		}
		return nil, status.Errorf(codes.Internal, "ошибка поиска товара: %v", err)
	}
	return &pb.Product{Id: product.ID, Name: product.Name, ManufacturerId: product.ManufacturersID, CategoryId: product.CategoryID, Price: product.Price}, nil
}

func (s *CatalogServer) GetProductByID(ctx context.Context, req *pb.GetProductByIDRequest) (*pb.Product, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id товара обязателен")
	}

	product, err := s.prodSvc.GetProductByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "товар не найден")
		}
		return nil, status.Errorf(codes.Internal, "ошибка поиска товара: %v", err)
	}

	return &pb.Product{
		Id:             product.ID,
		Name:           product.Name,
		ManufacturerId: product.ManufacturersID,
		CategoryId:     product.CategoryID,
		Price:          product.Price,
	}, nil
}

func (s *CatalogServer) GetListProduct(ctx context.Context, req *pb.GetListProductRequest) (*pb.ListProductsResponse, error) {
	filter := models.ProductFilter{}
	if req.CategoryId != 0 {
		filter.CategoryID = &req.CategoryId
	}
	if req.ManufacturerId != 0 {
		filter.ManufacturersID = &req.ManufacturerId
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
	products, total, err := s.prodSvc.GetListProduct(ctx, filter, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения списка товаров: %v", err)
	}
	var pbProducts []*pb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:             product.ID,
			Name:           product.Name,
			ManufacturerId: product.ManufacturersID,
			CategoryId:     product.CategoryID,
			Price:          product.Price,
		})
	}
	return &pb.ListProductsResponse{Product: pbProducts, Total: total}, nil
}
