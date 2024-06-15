package redis

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
)

// Config represents redis config
type Config struct {
	Addr     string
	Password string
	DB       int

	SSHTunnel bool
	SSHHost   string // included port
	SSHUser   string
	SSHPwd    string
}

// Redis represents redis service
type Redis struct {
	cfg *Config
	ctx context.Context
	rd  *redis.Client
}

// New initializes redis service with default config
func New(cfg *Config) *Redis {
	ctx := context.Background()

	rdClient := &redis.Client{}
	// Create a Redis connection URL for the remote Redis server
	redisURL := &url.URL{
		Scheme: "redis",
		Host:   cfg.Addr,
	}

	if cfg.SSHTunnel {
		// start ssh tunnel
		sshConfig := &ssh.ClientConfig{
			User:            cfg.SSHUser,
			Auth:            []ssh.AuthMethod{ssh.Password(cfg.SSHPwd)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         15 * time.Second,
		}

		sshClient, err := ssh.Dial("tcp", cfg.SSHHost, sshConfig)
		if err != nil {
			panic(err)
		}
		defer sshClient.Close()

		// Connect to the Redis server on the remote host
		remoteConn, err := sshClient.Dial("tcp", redisURL.Host)
		if err != nil {
			panic(err)
		}
		defer remoteConn.Close()

		rdClient = redis.NewClient(&redis.Options{
			Addr:     redisURL.Host,
			Password: cfg.Password,
			DB:       cfg.DB,
			Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return sshClient.Dial(network, addr)
			},
			// Disable timeouts, because SSH does not support deadlines.
			ReadTimeout:  -2,
			WriteTimeout: -2,
		})
	} else {
		// start client
		rdClient = redis.NewClient(&redis.Options{
			Addr:     redisURL.Host,
			Password: cfg.Password,
			DB:       cfg.DB,
		})
	}

	// ping for 1st start
	if _, err := rdClient.Ping(context.Background()).Result(); err != nil {
		// panic()
		panic(fmt.Errorf("redis service failed to start %s", err.Error()))
	}

	return &Redis{
		cfg: cfg,
		ctx: ctx,
		rd:  rdClient,
	}
}
