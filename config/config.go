package config

type Server struct {
	Zap      Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	RedisRW   RedisRW   `mapstructure:"redisrw" json:"redisrw" yaml:"redisrw"`
	RedisRO   RedisRO   `mapstructure:"redisro" json:"redisro" yaml:"redisro"`
}
