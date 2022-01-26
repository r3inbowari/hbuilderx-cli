package meiwobuxing

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/r3inbowari/common"
	. "github.com/r3inbowari/zlog"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func HandleVersion(w http.ResponseWriter, r *http.Request) {
	ResponseCommon(w, Up.VersionStr+" "+Up.ReleaseTag, "ok", 1, http.StatusOK, 0)
}

type PackRequest struct {
	Package    string     `json:"package"`
	PackConfig PackConfig `json:"packconfig"`
	User
	Callback     string            `json:"callback"`
	State        int               `json:"state"`
	TraceID      string            `json:"traceid"`
	Reason       string            `json:"reason"`
	Safe         HashLink          `json:"safe"`
	Certs        HashLink          `json:"certs"`
	Manifest     string            `json:"manifest"`
	DownloadLink []string          `json:"downloadlink"`
	Log          string            `json:"log"`
	Export       string            `json:"export"`
	Paths        map[string]string `json:"paths"`
}

type ManifestSetting struct {
	Appid                              string
	Name                               string
	Description                        string
	VersionName                        string
	VersionCode                        int
	AndroidHdpi                        string
	AndroidXhdpi                       string
	AndroidXxhdpi                      string
	AndroidXxxhdpi                     string
	Appstore                           string
	IpadApp                            string
	IpadApp2X                          string
	IpadNotification                   string
	IpadNotification2X                 string
	IpadProapp2X                       string
	IpadSettings                       string
	IpadSettings2X                     string
	IpadSpotlight                      string
	IpadSpotlight2X                    string
	IphoneApp2X                        string
	IphoneApp3X                        string
	IphoneNotification2X               string
	IphoneNotification3X               string
	IphoneSettings2X                   string
	IphoneSettings3X                   string
	IphoneSpotlight2X                  string
	IphoneSpotlight3X                  string
	SplashscreenAndroidStyle           string
	SplashscreenAndroidHdpi            string
	SplashscreenAndroidXhdpi           string
	SplashscreenAndroidXxhdpi          string
	AmapPlatform                       []string
	AmapAppkeyIos                      string
	AmapAppkeyAndroid                  string
	MapsAmapAppkeyIos                  string
	MapsAmapAppkeyAndroid              string
	ShareQqAppid                       string
	ShareWeixinAppid                   string
	ShareWeixinUniversalLinks          string
	ComAppleDeveloperAssociatedDomains []string
}

type HashLink struct {
	Android string `json:"android"`
	IOS     string `json:"ios"`
	IOSEx   string `json:"iosex,omitempty"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PackConfig struct {
	Project     string         `json:"project"`
	Platform    string         `json:"platform"`
	IsCustom    bool           `json:"iscustom"`
	SafeMode    bool           `json:"safemode"`
	Android     SettingAndroid `json:"android"`
	IOS         SettingIOS     `json:"ios"`
	IsConfusion bool           `json:"isconfusion"`
	SplashAds   bool           `json:"splashads"`
	RpAds       bool           `json:"rpads"`
	PushAds     bool           `json:"pushads"`
	Exchange    bool           `json:"exchange"`
}

type SettingAndroid struct {
	PackageName     string `json:"packagename"`
	AndroidPackType string `json:"androidpacktype"`
	CertAlias       string `json:"certalias"`
	CertFile        string `json:"certfile"`
	CertPassword    string `json:"certpassword"`
	Channels        string `json:"channels"`
}

type SettingIOS struct {
	Bundle          string `json:"bundle"`
	SupportedDevice string `json:"supporteddevice"`
	Profile         string `json:"profile"`
	CertFile        string `json:"certfile"`
	CertPassword    string `json:"certpassword"`
}

func HandlePackRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}

	var conf PackRequest
	err = json.Unmarshal(body, &conf)
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6402)
		return
	}

	traceid := conf.PackEnqueue()
	ResponseCommon(w, traceid, "request succeed", 1, http.StatusOK, 6400)
}

func HandlePackState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceID := vars["traceID"]
	if len(traceID) != 36 {
		ResponseCommon(w, "error traceID format", "request failed", 1, http.StatusOK, 6403)
		return
	}
	res := GetBuildState(traceID)
	if res == nil {
		ResponseCommon(w, "", "nonexistent trace", 1, http.StatusOK, 6404)
		return
	}
	ResponseCommon(w, res, "query succeed", 1, http.StatusOK, 6400)
}

func FileUpload(w http.ResponseWriter, r *http.Request) {
	uuid := CreateUUID()
	uploadFile, handle, err := r.FormFile("file")
	saveFile, err := os.OpenFile(GetPath(uuid), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil || handle == nil {
		ResponseCommon(w, "upload failed", "request failed", 1, http.StatusOK, 6405)
		return
	}
	_, _ = io.Copy(saveFile, uploadFile)

	defer func() {
		saveFile.Close()
		uploadFile.Close()
	}()
	ResponseCommon(w, uuid, "upload succeed", 1, http.StatusOK, 6400)
}

func FileDownload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["id"]
	path := GetPath(fileID)
	file, _ := os.Open(path)
	defer file.Close()
	fileHeader := make([]byte, 512)
	_, err := file.Read(fileHeader)
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	fileStat, _ := file.Stat()
	w.Header().Set("Content-Disposition", "attachment; filename="+fileID)
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	_, err = file.Seek(0, 0)
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	return
}

func HandleManifest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}

	var manifest Manifest
	err = json.Unmarshal(body, &manifest)
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6402)
		return
	}

	id, err := manifest.Save()
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	ResponseCommon(w, id, "request succeed", 1, http.StatusOK, 6400)
}

func ExportElements(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	if !json.Valid(body) {
		ResponseCommon(w, "not a json", "request failed", 1, http.StatusOK, 6433)
		return
	}
	//c := CreateUUID()
	//err = SaveBytesToRes(c, body)
	//if err != nil {
	//	ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6433)
	//	return
	//}
	//ResponseCommon(w, c, "request succeed", 1, http.StatusOK, 6400)

	j2m, err := Json2Map(body)
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	s := j2m.GenJs()
	c := CreateUUID()
	err = SaveBytesToRes(c, []byte(s))
	ResponseCommon(w, c, "request succeed", 1, http.StatusOK, 6400)
}

func CallbackTest(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	id := r.Form.Get("trace")
	Log.WithFields(logrus.Fields{"id": id}).Warn("[TEST] 回调到了内部测试接口，你的打包请求已经完成，请及时查看结果。")
	ResponseCommon(w, "done", "request succeed", 1, http.StatusOK, 6400)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	var u User
	err = json.Unmarshal(body, &u)
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	token, err := u.Login()
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	ResponseCommon(w, token, "request succeed", 1, http.StatusOK, 6400)
}

func HandleAddPathMapping(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	var pm PathMapping
	err = json.Unmarshal(body, &pm)
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}

	// validate data at first
	// path, name or type

	// merge data and save to leveldb
	key := common.CreateUUID()
	err = SetJson("path", key, &pm)
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	ResponseCommon(w, key, "request succeed", 1, http.StatusOK, 6400)
}

func HandleAllPathMapping(w http.ResponseWriter, r *http.Request) {
	items, err := Iter("path")
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	ResponseCommon(w, items, "request succeed", 1, http.StatusOK, 6400)
}

func HandleDelPathMapping(w http.ResponseWriter, r *http.Request) {
	// only try!
	_ = r.ParseForm()
	id := r.Form.Get("id")
	if !common.VerifyUUID(id) {
		ResponseCommon(w, "invalid id", "request failed", 1, http.StatusOK, 6491)
		return
	}

	err := Delete("path", id)
	if err != nil {
		ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
		return
	}
	ResponseCommon(w, "ok", "request succeed", 1, http.StatusOK, 6400)
}
