package endpoint

import "errors"

var (
	ErrCannotConnect = errors.New("cannot connect")
	ErrUnknownPeer   = errors.New("unknown peer")
	ErrInvalidPeer   = errors.New("invalid peer")
)
