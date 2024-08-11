package model

import "errors"

var ErrEmailAlreadyExists = errors.New("email already exists")

var ErrExpiredRefreshToken = errors.New("Expired refresh token")

var ErrInvalidRefreshToken = errors.New("Invalid refresh token")

var ErrAddBlacklistTokenFailed = errors.New("Add blacklist token failed")

var ErrFailedToHashPassword = errors.New("Failed to hash password")
