package main

import (
	"check_ebs_check/log"
	"check_ebs_check/tool"
	"check_ebs_check/tpcmon"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

/*

 */

var (
	Project  = "aliyun-diy-alarm"
	Logstore = "check_ebs_write"

	LogFile  = "./check_ebs_write.log"
	ConfFile = "/etc/tpcmon/diy_init/conf.json"

	Hostname string
	HostIp   string

	ProbeCmd = "touch /data/test_dir > /dev/null 2>&1; echo $?"
	Client   sls.ClientInterface

	DoneCh chan struct{}

	Logger *zap.Logger

	Health   = "0"
	UnHealth = "1"
)

func main() {
	// fmt.Printf("hostname:%s hostip:%s\n", Hostname, HostIp)
	for {
		start := time.Now()

		ok, err := ProbeEbsWrite()
		if err != nil {
			msg := fmt.Sprintf("ProbeEbsWrite failed! err:%v", err)
			Logger.Error(msg)
		} else {
			msg := fmt.Sprintln("ProbeEbsWrite success!")
			Logger.Info(msg)
		}

		if ok {
			data := new(Data)
			data.InstanceId = Hostname
			data.InstanceIp = HostIp
			data.Time = tool.CurrentLocalTimeFormatted()
			data.Result = Health

			if err := data.Write2Logstore(); err != nil {
				msg := fmt.Sprintf("data.Write2Logstore() [%s] push sls failed! err:%v", Health, err)
				Logger.Error(msg)
			} else {
				msg := fmt.Sprintf("data.Write2Logstore() [%s] push sls success!", Health)
				Logger.Info(msg)
			}
		} else {
			data := new(Data)
			data.InstanceId = Hostname
			data.InstanceIp = HostIp
			data.Time = tool.CurrentLocalTimeFormatted()
			data.Result = UnHealth

			if err := data.Write2Logstore(); err != nil {
				msg := fmt.Sprintf("data.Write2Logstore() [%s] push sls failed! err:%v", UnHealth, err)
				Logger.Error(msg)
			} else {
				msg := fmt.Sprintf("data.Write2Logstore() [%s] push sls success!", UnHealth)
				Logger.Info(msg)
			}
		}

		end := time.Since(start)
		// msg := fmt.Sprintln(end)

		// Logger.Info(msg)
		sleepDuration := time.Second*10 - end
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
	}
}

// probe ebs

// write log data to sls-inner
// iZk4d01izfdzam0o2vsxygZ
type Data struct {
	Time       string
	InstanceId string
	InstanceIp string
	Result     string
}

func (d *Data) Write2Logstore() error {
	// 日志组: 将多个CollectorLog打包一次性推送到sls
	logGroup := &sls.LogGroup{
		Topic:  proto.String(""),
		Source: proto.String(d.InstanceIp),
	}

	// 将data生成contents
	var contents = []*sls.LogContent{
		{Key: proto.String("time"), Value: proto.String(d.Time)},
		{Key: proto.String("instance_id"), Value: proto.String(d.InstanceId)},
		{Key: proto.String("instance_ip"), Value: proto.String(d.InstanceIp)},
		{Key: proto.String("probe_result"), Value: proto.String(d.Result)},
	}

	// 将contents封装到log对象中
	logIns := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: contents,
	}

	// // 将log对象追加到日志组
	logGroup.Logs = append(logGroup.Logs, logIns)

	// 将日志组推送到sls (三次重试！)
	var err error
	for i := 0; i < 3; i++ {
		err = Client.PutLogs(Project, Logstore, logGroup)
		if err != nil {
			// 重试
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}

	// 三次执行完毕仍然err  -->抛出错误
	if err != nil {
		return err
	}

	return nil
}

// new sls client： default v1 client
func init() {
	// 1. zap logger
	Logger = log.NewZapLogger(LogFile, 100, 3, 30, true)

	// 2. 配置文件 /etc/tpcmon/diy_init/conf.json
	conf, err := LoadConf()
	if err != nil {
		msg := fmt.Sprintf("init load conf failed! err:%v\n", err)
		Logger.Error(msg)
		os.Exit(1)
	}

	// ak_byte, _ := base64.RawStdEncoding.DecodeString(conf.AKSK.SLS.Ak)
	// sk_byte, _ := base64.RawStdEncoding.DecodeString(conf.AKSK.SLS.Sk)

	// 3. （全局变量初始化）
	Client = sls.CreateNormalInterface(conf.AKSK.SLS.Ep, conf.AKSK.SLS.Ak, conf.AKSK.SLS.Sk, "")
	HostIp, err = tool.GetHostIp("eth0")
	if err != nil {
		msg := fmt.Sprintf("init GetHostIp(\"eth0\") failed! err:%v\n", err)
		Logger.Error(msg)
		os.Exit(1)
	}

	Hostname, err = tool.GetHostname()
	if err != nil {
		msg := fmt.Sprintf("init GetHostname() failed! err:%v\n", err)
		Logger.Error(msg)
		os.Exit(1)
	}

	Logger.Info("init success!")

}

// load conf
func LoadConf() (tpcmon.Conf, error) {
	confByte, err := os.ReadFile(ConfFile)
	if err != nil {
		return tpcmon.Conf{}, err
	}

	var conf tpcmon.Conf
	err = json.Unmarshal(confByte, &conf)
	if err != nil {
		return tpcmon.Conf{}, err
	}

	return conf, nil

}

func ProbeEbsWrite() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	stdout, err := tool.ExecCmd(ctx, ProbeCmd)
	if err != nil {
		msg := fmt.Sprintf("exec cmd [%s] failed! err:%v\n", ProbeCmd, err)
		return false, errors.New(msg)
	}

	switch tool.TrimSuffixNewLine(stdout) {
	case "0":
		return true, nil
	default:
		return false, nil
	}
}
