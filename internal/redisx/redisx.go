package redisx

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/ZooLearn/file/internal/log"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var ErrNotConfigured = errors.New("Redis is not configured")

type RedisConfig struct {
	Address           string   `yaml:"address,omitempty"`
	Username          string   `yaml:"username,omitempty"`
	Password          string   `yaml:"password,omitempty"`
	DB                int      `yaml:"db,omitempty"`
	UseTLS            bool     `yaml:"use_tls,omitempty"`
	MasterName        string   `yaml:"sentinel_master_name,omitempty"`
	SentinelUsername  string   `yaml:"sentinel_username,omitempty"`
	SentinelPassword  string   `yaml:"sentinel_password,omitempty"`
	SentinelAddresses []string `yaml:"sentinel_addresses,omitempty"`
	ClusterAddresses  []string `yaml:"cluster_addresses,omitempty"`
	DialTimeout       int      `yaml:"dial_timeout,omitempty"`
	ReadTimeout       int      `yaml:"read_timeout,omitempty"`
	WriteTimeout      int      `yaml:"write_timeout,omitempty"`
	MaxRedirects      *int     `yaml:"max_redirects,omitempty"`
}

func (r *RedisConfig) IsConfigured() bool {
	if r.Address != "" {
		return true
	}
	if len(r.SentinelAddresses) > 0 {
		return true
	}
	if len(r.ClusterAddresses) > 0 {
		return true
	}
	return false
}

func (r *RedisConfig) GetMaxRedirects() int {
	if r.MaxRedirects != nil {
		return *r.MaxRedirects
	}
	return 2
}

func GetRedisClient(conf *RedisConfig) (redis.UniversalClient, error) {
	if conf == nil {
		return nil, nil
	}

	if !conf.IsConfigured() {
		return nil, ErrNotConfigured
	}

	var rcOptions *redis.UniversalOptions
	var rc redis.UniversalClient
	var tlsConfig *tls.Config

	if conf.UseTLS {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	if len(conf.SentinelAddresses) > 0 {
		log.Infof("connecting to redis sentinel addr: %s masterName: %s", conf.SentinelAddresses, conf.MasterName)

		// By default DialTimeout set to 2s
		if conf.DialTimeout == 0 {
			conf.DialTimeout = 2000
		}
		// By default ReadTimeout set to 0.2s
		if conf.ReadTimeout == 0 {
			conf.ReadTimeout = 200
		}
		// By default WriteTimeout set to 0.2s
		if conf.WriteTimeout == 0 {
			conf.WriteTimeout = 200
		}

		rcOptions = &redis.UniversalOptions{
			Addrs:            conf.SentinelAddresses,
			SentinelUsername: conf.SentinelUsername,
			SentinelPassword: conf.SentinelPassword,
			MasterName:       conf.MasterName,
			Username:         conf.Username,
			Password:         conf.Password,
			DB:               conf.DB,
			TLSConfig:        tlsConfig,
			DialTimeout:      time.Duration(conf.DialTimeout) * time.Millisecond,
			ReadTimeout:      time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout:     time.Duration(conf.WriteTimeout) * time.Millisecond,
		}
	} else if len(conf.ClusterAddresses) > 0 {
		log.Infof("connecting to redis cluster addr: %s", conf.ClusterAddresses)
		rcOptions = &redis.UniversalOptions{
			Addrs:        conf.ClusterAddresses,
			Username:     conf.Username,
			Password:     conf.Password,
			DB:           conf.DB,
			TLSConfig:    tlsConfig,
			MaxRedirects: conf.GetMaxRedirects(),
		}
	} else {
		log.Infof("connecting to redis simple addr: %s", conf.Address)
		rcOptions = &redis.UniversalOptions{
			Addrs:     []string{conf.Address},
			Username:  conf.Username,
			Password:  conf.Password,
			DB:        conf.DB,
			TLSConfig: tlsConfig,
		}
	}
	rc = redis.NewUniversalClient(rcOptions)

	if err := rc.Ping(context.Background()).Err(); err != nil {
		err = errors.Wrap(err, "unable to connect to redis")
		return nil, err
	}

	return rc, nil
}

// TODO: function convert data input to bytes, don't save it in string.
