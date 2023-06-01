package simplerouting

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// Config is a structure containing all the options that can be used when
// constructing a SimpleRouting.
type Config struct {
	QueryConcurrency      int
	QueryTimeout          time.Duration
	MaxConcurrentRequests int
	ProtocolID            protocol.ID
}

// Apply applies the given options to this Option
func (cfg *Config) Apply(opts ...Option) error {
	for i, opt := range opts {
		if err := opt(cfg); err != nil {
			return fmt.Errorf("SimpleRouting option %d failed: %s", i, err)
		}
	}
	return nil
}

// Validate validates the configuration options (fool proofing)
func (cfg *Config) Validate() error {
	return nil
}

// Option type for SimpleRouting
type Option func(*Config) error

// DefaultConfig is the default options for SimpleRouting. This option is always
// prepended to the list of options passed to the SimpleRouting constructor.
var DefaultConfig = func(cfg *Config) error {
	cfg.QueryConcurrency = 3
	cfg.QueryTimeout = 10 * time.Second
	cfg.MaxConcurrentRequests = 5
	cfg.ProtocolID = consts.ProtocolDHT

	return nil
}

// QueryConcurrency sets the maximum number of concurrent inflight requests
// that can be made by a single query.
func QueryConcurrency(n int) Option {
	return func(cfg *Config) error {
		cfg.QueryConcurrency = n
		return nil
	}
}

// QueryTimeout sets the timeout for a query.
func QueryTimeout(t time.Duration) Option {
	return func(cfg *Config) error {
		cfg.QueryTimeout = t
		return nil
	}
}

// MaxConcurrentRequests sets the maximum number of concurrent requests that
// can be made to the network.
func MaxConcurrentRequests(n int) Option {
	return func(cfg *Config) error {
		cfg.MaxConcurrentRequests = n
		return nil
	}
}

// ProtocolID sets the protocol ID to use for the network.
func ProtocolID(p protocol.ID) Option {
	return func(cfg *Config) error {
		cfg.ProtocolID = p
		return nil
	}
}
