// Package config 提供配置验证功能
package config

import (
	"fmt"
	"strings"
)

func (c *Config) Validate() error {
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required - generate with: make generate-jwt-secret")
	}

	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf(
			"JWT_SECRET must be at least 32 characters (current: %d)\nGenerate secure secret: make generate-jwt-secret",
			len(c.JWT.Secret),
		)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}

	if c.Server.ReadTimeout < 0 {
		return fmt.Errorf("server.readtimeout must be non-negative")
	}

	if c.Server.WriteTimeout < 0 {
		return fmt.Errorf("server.writetimeout must be non-negative")
	}

	if c.Server.IdleTimeout < 0 {
		return fmt.Errorf("server.idletimeout must be non-negative")
	}

	if c.Server.ShutdownTimeout < 0 {
		return fmt.Errorf("server.shutdowntimeout must be non-negative")
	}

	if c.Server.MaxHeaderBytes < 0 {
		return fmt.Errorf("server.maxheaderbytes must be non-negative")
	}

	if c.App.Environment == "production" {
		if c.Database.Password == "" {
			return fmt.Errorf("database.password is required in production")
		}

		if c.Database.SSLMode == "disable" {
			return fmt.Errorf("database SSL mode cannot be 'disable' in production")
		}
	}

	// Redis 配置验证（如果启用）
	if c.Redis.Enabled {
		if c.Redis.Host == "" {
			return fmt.Errorf("redis.host is required when Redis is enabled")
		}
		if c.Redis.Port == 0 {
			return fmt.Errorf("redis.port is required when Redis is enabled")
		}
		if c.Redis.DB < 0 || c.Redis.DB > 15 {
			return fmt.Errorf("redis.db must be between 0-15")
		}
	}

	// MongoDB 配置验证（如果启用）
	if c.MongoDB.Enabled {
		if c.MongoDB.URI == "" {
			return fmt.Errorf("mongodb.uri is required when MongoDB is enabled")
		}
		if c.MongoDB.Database == "" {
			return fmt.Errorf("mongodb.database is required when MongoDB is enabled")
		}
	}

	// RabbitMQ 配置验证（如果启用）
	if c.RabbitMQ.Enabled {
		if c.RabbitMQ.URL == "" {
			return fmt.Errorf("rabbitmq.url is required when RabbitMQ is enabled")
		}
	}

	// 安全配置验证
	if c.Security.BcryptCost < 10 || c.Security.BcryptCost > 14 {
		fmt.Printf("⚠️  Warning: bcrypt cost factor (%d) should be between 10-14 for optimal security\n", c.Security.BcryptCost)
	}

	if c.Security.PasswordMinLength < 8 {
		fmt.Printf("⚠️  Warning: password minimum length (%d) should be at least 8 characters\n", c.Security.PasswordMinLength)
	}

	if c.Security.MaxLoginAttempts < 3 || c.Security.MaxLoginAttempts > 10 {
		fmt.Printf("⚠️  Warning: max login attempts (%d) should be between 3-10\n", c.Security.MaxLoginAttempts)
	}

	return nil
}

// ValidateOrPanic 验证配置，如果失败则 panic
// 用于应用启动时的配置验证
func (c *Config) ValidateOrPanic() {
	if err := c.Validate(); err != nil {
		fmt.Printf("\n❌ 配置验证失败:\n  %s\n\n", strings.ReplaceAll(err.Error(), "\n", "\n  "))
		panic("配置验证失败，应用无法启动")
	}
}
