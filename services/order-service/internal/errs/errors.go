package errs

import "errors"

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidOrderID    = errors.New("invalid order id")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrEmptyCart         = errors.New("cart is empty")
	ErrCannotCancelOrder = errors.New("order cannot be cancelled in current status")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInventoryNotFound = errors.New("inventory not found")
)
