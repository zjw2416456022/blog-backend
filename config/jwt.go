package config

import "time"

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey string
	ExpiresIn time.Duration
}

// GetJWTConfig 获取JWT配置
func GetJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey: "your_secret_key_change_in_production",
		ExpiresIn: 24 * time.Hour, // 24小时过期
	}
}