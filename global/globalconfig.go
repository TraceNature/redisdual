package global

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"redisdual/config"
	"sync"
)

var (
	RSPViper  *viper.Viper
	RSPLog    *zap.Logger
	RSPConfig config.Server
	once      sync.Once

)

