package meiwobuxing

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	. "github.com/r3inbowari/zlog"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func InitConfig() {
	if !Exists("res/") {
		_ = os.Mkdir("res/", 0777)
	}
	if !Exists("bili.json") {
		Log.Info("[FILE] init user configuration")
		var config LocalConfig
		var l = "debug"
		config.Finger = "1777945899"
		config.LoggerLevel = &l
		config.APIAddr = ":9090"
		config.AutoUpdate = true
		config.CaCert = ""
		config.CaKey = ""
		config.MaxRetryCount = 3

		config.JwtEnable = true
		config.JwtSecret = "sdasdasdasfasdasd"
		config.JwtMD5 = "48dc65a51b480244c296c44de7be53f5"
		config.JwtTimeout = 72
		_ = config.SetConfig()
	}
}

type RequestResult struct {
	Total   int         `json:"total"`
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"msg"`
}

func ResponseCommon(w http.ResponseWriter, data interface{}, msg string, total int, tag int, code int) {
	var rq RequestResult
	rq.Data = data
	rq.Total = total
	rq.Code = code
	rq.Message = msg
	jsonStr, err := json.Marshal(rq)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err.Error()}).Error("response error")
	}
	w.WriteHeader(tag)
	_, _ = fmt.Fprintf(w, string(jsonStr))
}

func CreateUUID() string {
	u1 := uuid.NewV4()
	return u1.String()
}

var HyperLinkExp = `[a-zA-z]+://[^\s]*`

func GetHyperLinks(s string) ([]string, error) {
	var retArr = make([]string, 2)
	arrRes := strings.Split(s, "\n")

	for _, v := range arrRes {
		if strings.Contains(v, "Android自有证书 下载地址") {
			re := regexp.MustCompile(HyperLinkExp)
			links := re.FindAllStringSubmatch(v, -1)
			for _, value := range links {
				retArr[0] = value[0]
				break
			}
		} else if strings.Contains(v, "IOS自有证书 下载地址") {
			re := regexp.MustCompile(HyperLinkExp)
			links := re.FindAllStringSubmatch(v, -1)
			for _, value := range links {
				retArr[1] = value[0]
				break
			}
		}
	}
	if retArr[0] == "" && retArr[1] == "" {
		return retArr, errors.New("an error occurred during build stage, check the log for more details")
	}
	return retArr, nil
}

func Get(url string) {
	if len(url) < 10 {
		return
	}
	tls := strings.ToLower(url[0:5]) == "https"
	if tls {

	}
	_, err := http.Get(url)
	if err != nil {
		Log.WithFields(logrus.Fields{"url": url}).Info("[CLI] callback failed...")
	}
}

func VerifyUUID(u string) bool {
	_, err := guuid.Parse(u)
	return err == nil
}

var exitChan chan os.Signal
var exitFunc map[string]func()

func InitSignalExit(f func(signal os.Signal)) {
	exitChan = make(chan os.Signal)
	exitFunc = make(map[string]func(), 100)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		s := <-exitChan
		Log.WithFields(logrus.Fields{"signal": s.String()}).Info("[MAIN] exiting...")
		if f != nil {
			f(s)
		}
		for name, fa := range exitFunc {
			Log.WithFields(logrus.Fields{"task": name}).Info("[MAIN] exec task")
			fa()
		}
		os.Exit(77)
	}()
}

// AddExitFunc UNSAFE usage
func AddExitFunc(name string, f func()) {
	Log.WithFields(logrus.Fields{"task": name}).Info("[MAIN] add new exit function...")
	exitFunc[name] = f
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2StringInWindows(byte []byte, charset Charset) string {
	if runtime.GOOS != "windows" {
		return string(byte)
	}
	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

// GetMD5 md5
func GetMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetMD5WithSalt(str string, salt ...interface{}) (CryptStr string) {
	if l := len(salt); l > 0 {
		slice := make([]string, l+1)
		str = fmt.Sprintf(str+strings.Join(slice, "%v"), salt...)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

type JMap map[string]interface{}

// Json2Map json转map函数
func Json2Map(bs []byte) (JMap, error) {
	var tempMap map[string]interface{}
	err := json.Unmarshal(bs, &tempMap)
	if err != nil {
		return nil, err
	}
	return tempMap, nil
}

func (jm *JMap) GenJs() string {
	cnt := len(*jm)
	dim := ""
	export := "export default { "
	for s, i := range *jm {
		dim += fmt.Sprintf("const %s = \"%s\";\r\n", s, i)
		export += fmt.Sprintf("%s: %s", s, s)
		if cnt != 1 {
			export += ", "
			cnt--
		}

	}
	return fmt.Sprintf("%s\r\n%s }\r\n", dim, export)
}

func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}