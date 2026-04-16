package errs

import "errors"

var (
	ErrInventoryNotFound    = errors.New("inventory not found")
	ErrInventoryExists      = errors.New("inventory already exists")
	ErrInvalidQuantity      = errors.New("invalid quantity")
	ErrInvalidProductID     = errors.New("invalid product id")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrInvalidLocation      = errors.New("invalid warehouse location")
	ErrInvalidReservationID = errors.New("invalid reservation id")
)

var (
	ErrWarehouseNotFound = errors.New("warehouse not found")
	ErrWarehouseExists   = errors.New("warehouse already exists")
	ErrWarehouseRequired = errors.New("warehouse name, city and street are required")
)
