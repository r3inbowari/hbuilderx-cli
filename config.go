package meiwobuxing

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"time"
)

// LocalConfig 配置
type LocalConfig struct {
	Finger        string    `json:"finger"`          // canvas指纹
	APIAddr       string    `json:"api_addr"`        // API服务ADDR
	CacheTime     time.Time `json:"-"`               // 缓存时间
	LoggerLevel   *string   `json:"log_level"`       // 日志等级
	CheckLink     string    `json:"check_link"`      // 检查更新地址
	AutoUpdate    bool      `json:"auto_update"`     // 自动更新
	CaKey         string    `json:"ca_key"`          // CA 密钥
	CaCert        string    `json:"ca_crt"`          // CA 证书
	MaxRetryCount int       `json:"max_retry_count"` // 重试次数
	LogPath       string    `json:"log_path"`        // hbx log
	JwtEnable     bool      `json:"jwt_enable"`      // 鉴权
	JwtSecret     string    `json:"jwt_secret"`
	JwtTimeout    int       `json:"jwt_timeout"`
	JwtMD5        string    `json:"jwt_md5"`
}

var config = new(LocalConfig)
var configPath = "bili.json"

// GetConfig 返回配置文件
// imm 立即返回
func GetConfig(imm bool) *LocalConfig {
	if config.CacheTime.Before(time.Now()) || imm {
		if err := LoadConfig(configPath, config); err != nil {
			Log.Error("loading file failed")
			time.Sleep(time.Second * 5)
			os.Exit(76)
			return nil
		}
		config.CacheTime = time.Now().Add(time.Second * 60)
	}
	return config
}

// SetConfig 更新
func (lc *LocalConfig) SetConfig() error {
	fp, err := os.Create(configPath)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("loading file failed")
	}
	defer fp.Close()
	data, err := json.Marshal(lc)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("marshal file failed")
	}
	n, err := fp.Write(data)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("write file failed")
	}
	Log.WithFields(logrus.Fields{"size": n}).Info("[FILE] update user configuration")
	return nil
}

const configFileSizeLimit = 10 << 20

// LoadConfig path 文件路径 dist 存放目标
func LoadConfig(path string, dist interface{}) error {
	configFile, err := os.Open(path)
	if err != nil {
		Log.WithFields(logrus.Fields{"path": path, "err": err}).Error("failed to open config file.")
		return err
	}

	fi, _ := configFile.Stat()
	if size := fi.Size(); size > (configFileSizeLimit) {
		Log.WithFields(logrus.Fields{"path": path, "size": size}).Error("config file size exceeds reasonable limited")
		return errors.New("limited")
	}

	if fi.Size() == 0 {
		Log.WithFields(logrus.Fields{"path": path, "size": 0}).Error("config file is empty, skipping")
		return errors.New("empty")
	}

	buffer := make([]byte, fi.Size())
	_, err = configFile.Read(buffer)
	buffer, err = StripComments(buffer)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("failed to strip comments from json")
		return err
	}

	buffer = []byte(os.ExpandEnv(string(buffer)))

	err = json.Unmarshal(buffer, &dist)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("failed unmarshalling json")
		return err
	}
	return nil
}

// StripComments 注释清除
func StripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0)
	lines := bytes.Split(data, []byte("\n"))
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}
	return bytes.Join(filtered, []byte("\n")), nil
}
