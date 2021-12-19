package meiwobuxing

import (
	"encoding/json"
	. "github.com/r3inbowari/zlog"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

// FileSystemPath 资源执行目录
var fileSystemPath = "./"

// CopyFormRes 从 res 中拷贝资源到指定路径
// uuid res中的索引
// dst 目标文件夹
func CopyFormRes(uuid string, dst string) error {
	if !VerifyUUID(uuid) {
		return nil
	}
	var srcFile, dstFile *os.File
	srcFile, err := os.Open(GetPath(uuid))
	if err != nil {
		return err
	}
	// 覆盖写入
	dstFile, err = os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	defer func() {
		if srcFile != nil {
			srcFile.Close()
		}
		if dstFile != nil {
			dstFile.Close()
		}
	}()
	return nil
}

func SaveBytesToRes(uuid string, v []byte) error {
	if !VerifyUUID(uuid) {
		uuid = CreateUUID()
	}
	fp, err := os.Create(GetPath(uuid))
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("loading file failed")
		return err
	}
	defer fp.Close()
	n, err := fp.Write(v)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("write file failed")
		return err
	}
	Log.WithFields(logrus.Fields{"size": n, "id": uuid}).Info("[FILE] save file")
	return err
}

// SaveJsonToRes 保存一个json对象到res
// 处理命名和未命名对象
func SaveJsonToRes(uuid string, v interface{}) error {
	if !VerifyUUID(uuid) {
		uuid = CreateUUID()
	}
	fp, err := os.Create(GetPath(uuid))
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("loading file failed")
	}
	defer fp.Close()
	data, err := json.Marshal(v)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("marshal file failed")
	}
	n, err := fp.Write(data)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("write file failed")
	}
	Log.WithFields(logrus.Fields{"size": n, "id": uuid}).Info("[FILE] save file")
	return nil
}

func GetResPath() string {
	return fileSystemPath
}

func OpenBytesFromRes(uuid string) ([]byte, error) {
	filePtr, err := os.Open(GetPath(uuid))
	if err != nil {
		return nil, err
	}
	defer filePtr.Close()

	return ioutil.ReadAll(filePtr)
}

func OpenJsonFromRes(uuid string, v interface{}) error {
	filePtr, err := os.Open(GetPath(uuid))
	if err != nil {
		return err
	}
	defer filePtr.Close()

	decoder := json.NewDecoder(filePtr)
	return decoder.Decode(v)
}

var sized int

func InitFileSystem(root string, size int) {
	fileSystemPath = root + "/res/"

	DirMap = make([]string, size)
	sized = size
	for i := 0; i < size; i++ {
		m := GetMD5("r3inbowari" + strconv.Itoa(i))
		DirMap[i] = m + "/"
		if !Exists(fileSystemPath + m) {
			_ = os.Mkdir(fileSystemPath+m, os.ModePerm)
		}
	}
}

var DirMap []string

func GetPath(uuid string) string {
	return fileSystemPath + DirMap[Calc(uuid)] + uuid
}

func Calc(uuid string) int {
	var ret int
	for _, v := range GetMD5(uuid) {
		ret += int(v)
	}
	return ret % sized
}
