package auth

import "github.com/pkg/errors"

// ErrNotAuthorized ...
var ErrNotAuthorized = errors.New("not access to this resource")

// ErrNotOwner ...
var ErrNotOwner = errors.New("this resources is not in your ownership")

// ErrTokenExpired ...
var ErrTokenExpired = errors.New("token is expired")
