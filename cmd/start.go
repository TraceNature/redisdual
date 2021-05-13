package cmd

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"redisdual/core"
	"redisdual/global"
	"strconv"
	"syscall"
	"time"
)

func NewStartCommand() *cobra.Command {
	sc := &cobra.Command{
		Use:   "start [-d ]",
		Short: "start server",
		Run:   startCommandFunc,
	}

	sc.Flags().BoolP("daemon", "d", false, "start as daemon")

	return sc
}

func startCommandFunc(cmd *cobra.Command, args []string) {
	cmd.Println("server starting ...")
	daemon, err := cmd.Flags().GetBool("daemon")
	if err != nil {
		panic(err)
	}

	global.RSPViper.Set("daemon", daemon)
	startServer()
}

// 启动服务器
func startServer() {

	// -d 后台启动
	if global.RSPViper.GetBool("daemon") {
		cmd, err := background()
		if err != nil {
			panic(err)
		}

		//根据返回值区分父进程子进程
		if cmd != nil { //父进程
			fmt.Println("PPID: ", os.Getpid(), "; PID:", cmd.Process.Pid, "; Operating parameters: ", os.Args)
			return //父进程退出
		} else { //子进程
			fmt.Println("PID: ", os.Getpid(), "; Operating parameters: ", os.Args)
		}
	}

	global.RSPLog = core.Zap()
	global.RSPLog.Info("server start ... ")

	pidMap := make(map[string]int)
	// 记录pid
	pid := syscall.Getpid()
	pidMap["pid"] = pid

	pidYaml, _ := yaml.Marshal(pidMap)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(dir+"/pid", pidYaml, 0664); err != nil {
		global.RSPLog.Sugar().Error(err)
		panic(err)
	}
	global.RSPLog.Sugar().Infof("Actual pid is %d", pid)

	//redis 读写实例
	redisRW := GetRedisRW()

	//redis 只读实例
	redisRO := GetRedisRO()

	//清理RW
	redisRW.FlushAll()

	global.RSPLog.Sugar().Info("execinterval:", global.RSPViper.GetInt("execinterval"))
	loopstep := global.RSPViper.GetInt("loopstep")
	i := 0
	for {
		if i > loopstep {
			i = 0
		}
		redisRW.Set(global.RSPViper.GetString("localkeyprefix")+"_key"+strconv.Itoa(i), i, 3600*time.Second)
		dual(redisRW, redisRO, global.RSPViper.GetString("localkeyprefix")+"_key"+strconv.Itoa(i))
		dual(redisRW, redisRO, global.RSPViper.GetString("remotekeyprefix")+"_key"+strconv.Itoa(i))
		i++
		time.Sleep(time.Duration(global.RSPViper.GetInt("execinterval")) * time.Millisecond)
	}

}

//双读程序，先读取只读库，若有数据返回，若没有读取读写库
func dual(rw *redis.Client, ro *redis.Client, key string) {
	roResult, err := ro.Get(key).Result()

	if err == nil && roResult != "" {
		global.RSPLog.Sugar().Infof("Get key %s from redisro result is:%s ", key, roResult)
		return
	}

	rwResult, err := rw.Get(key).Result()
	if err != nil || rwResult == "" {
		global.RSPLog.Sugar().Infof("key %s no result return!", key)
		return
	}

	global.RSPLog.Sugar().Infof("Get key %s from redisrw result is: %s ", key, rwResult)

}

func GetRedisRW() *redis.Client {
	redisCfg := global.RSPConfig.RedisRW
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		global.RSPLog.Error("redis connect ping failed, err:", zap.Any("err", err))
	}

	return client
}

func GetRedisRO() *redis.Client {
	redisCfg := global.RSPConfig.RedisRO
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		global.RSPLog.Error("redis connect ping failed, err:", zap.Any("err", err))
	}

	return client
}

func background() (*exec.Cmd, error) {
	envName := "DAEMON"    //环境变量名称
	envValue := "SUB_PROC" //环境变量值

	val := os.Getenv(envName) //读取环境变量的值,若未设置则为空字符串
	if val == envValue {      //监测到特殊标识, 判断为子进程,不再执行后续代码
		return nil, nil
	}

	/*以下是父进程执行的代码*/

	//因为要设置更多的属性, 这里不使用`exec.Command`方法, 直接初始化`exec.Cmd`结构体
	cmd := &exec.Cmd{
		Path: os.Args[0],
		Args: os.Args,      //注意,此处是包含程序名的
		Env:  os.Environ(), //父进程中的所有环境变量
	}

	//为子进程设置特殊的环境变量标识
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envName, envValue))

	//异步启动子进程
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
