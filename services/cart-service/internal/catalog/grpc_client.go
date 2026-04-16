// Package catalog предоставляет gRPC-клиент для взаимодействия с catalog-service.
package catalog

import (
	"context"

	pb "kinos/proto/catalog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CatalogClient интерфейс для взаимодействия с Catalog Service
type CatalogClient interface {
	GetProductByID(ctx context.Context, id uint64) (*pb.Product, error)
	Close() error
}

// CatalogClientImpl реализация CatalogClient
type CatalogClientImpl struct {
	client pb.CatalogServiceClient
	conn   *grpc.ClientConn
}

func NewCatalogClient(address string) (*CatalogClientImpl, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &CatalogClientImpl{
		client: pb.NewCatalogServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close закрывает соединение
func (c *CatalogClientImpl) Close() error {
	return c.conn.Close()
}

// GetProductByID получает товар по ID
func (c *CatalogClientImpl) GetProductByID(ctx context.Context, id uint64) (*pb.Product, error) {
	req := &pb.GetProductByIDRequest{
		Id: id,
	}
	return c.client.GetProductByID(ctx, req)
}

// Остальные методы (для совместимости)
func (c *CatalogClientImpl) CreateCategory(ctx context.Context, name string) (*pb.Category, error) {
	req := &pb.CreateCategoryRequest{Name: name}
	return c.client.CreateCategory(ctx, req)
}

func (c *CatalogClientImpl) UpdateCategory(ctx context.Context, id uint64, name string) (*pb.Category, error) {
	req := &pb.UpdateCategoryRequest{Id: id, Name: name}
	return c.client.UpdateCategory(ctx, req)
}

func (c *CatalogClientImpl) DeleteCategory(ctx context.Context, id uint64) (*emptypb.Empty, error) {
	req := &pb.DeleteCategoryRequest{Id: id}
	return c.client.DeleteCategory(ctx, req)
}

func (c *CatalogClientImpl) GetCategory(ctx context.Context, name string) (*pb.Category, error) {
	req := &pb.GetCategoryRequest{Name: name}
	return c.client.GetCategory(ctx, req)
}

func (c *CatalogClientImpl) GetListCategory(ctx context.Context, limit, offset int32) (*pb.ListCategoriesResponse, error) {
	req := &pb.GetListCategoryRequest{Limit: limit, Offset: offset}
	return c.client.GetListCategory(ctx, req)
}

func (c *CatalogClientImpl) CreateManufacturer(ctx context.Context, name string) (*pb.Manufacturer, error) {
	req := &pb.CreateManufacturerRequest{Name: name}
	return c.client.CreateManufacturer(ctx, req)
}

func (c *CatalogClientImpl) UpdateManufacturer(ctx context.Context, id uint64, name string) (*pb.Manufacturer, error) {
	req := &pb.UpdateManufacturerRequest{Id: id, Name: name}
	return c.client.UpdateManufacturer(ctx, req)
}

func (c *CatalogClientImpl) DeleteManufacturer(ctx context.Context, id uint64) (*emptypb.Empty, error) {
	req := &pb.DeleteManufacturerRequest{Id: id}
	return c.client.DeleteManufacturer(ctx, req)
}

func (c *CatalogClientImpl) GetManufacturer(ctx context.Context, name string) (*pb.Manufacturer, error) {
	req := &pb.GetManufacturerRequest{Name: name}
	return c.client.GetManufacturer(ctx, req)
}

func (c *CatalogClientImpl) GetListManufacturer(ctx context.Context, limit, offset int32) (*pb.ListManufacturerResponse, error) {
	req := &pb.GetListManufacturerRequest{Limit: limit, Offset: offset}
	return c.client.GetListManufacturers(ctx, req)
}

func (c *CatalogClientImpl) CreateProduct(ctx context.Context, name string, manufacturersID, categoryID uint64, price float64) (*pb.Product, error) {
	req := &pb.CreateProductRequest{
		Name: name, ManufacturerId: manufacturersID, CategoryId: categoryID, Price: price,
	}
	return c.client.CreateProduct(ctx, req)
}

func (c *CatalogClientImpl) UpdateProduct(ctx context.Context, id uint64, name string, manufacturersID, categoryID uint64, price float64) (*pb.Product, error) {
	req := &pb.UpdateProductRequest{
		Id: id, Name: name, ManufacturerId: manufacturersID, CategoryId: categoryID, Price: price,
	}
	return c.client.UpdateProduct(ctx, req)
}

func (c *CatalogClientImpl) DeleteProduct(ctx context.Context, id uint64) (*emptypb.Empty, error) {
	req := &pb.DeleteProductRequest{Id: id}
	return c.client.DeleteProduct(ctx, req)
}

func (c *CatalogClientImpl) GetProduct(ctx context.Context, name string) (*pb.Product, error) {
	req := &pb.GetProductRequest{Name: name}
	return c.client.GetProduct(ctx, req)
}

func (c *CatalogClientImpl) GetListProduct(ctx context.Context, limit, offset int32, categoryID, manufacturersID uint64, priceMax, priceMin float64, nameContains string) (*pb.ListProductsResponse, error) {
	req := &pb.GetListProductRequest{
		Limit: limit, Offset: offset, CategoryId: categoryID, ManufacturerId: manufacturersID,
		PriceMax: priceMax, PriceMin: priceMin, NameContains: nameContains,
	}
	return c.client.GetListProduct(ctx, req)
}
