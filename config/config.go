package config

type Server struct {
	Zap             Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	RedisRW         RedisRW `mapstructure:"redisrw" json:"redisrw" yaml:"redisrw"`
	RedisRO         RedisRO `mapstructure:"redisro" json:"redisro" yaml:"redisro"`
	ExecInterval    int     `mapstructure:"execinterval" json:"execinterval" yaml:"execinterval"`
	LocalKeyPrefix  string  `mapstructure:"localkeyprefix" json:"localkeyprefix" yaml:"localkeyprefix"`
	RemoteKeyPrefix string  `mapstructure:"remotekeyprefix" json:"remotekeyprefix" yaml:"remotekeyprefix"`
	LoopStep        int     `mapstructure:"loopstep" json:"loopstep" yaml:"loopstep"`
}

