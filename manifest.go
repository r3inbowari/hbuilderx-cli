package meiwobuxing

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/r3inbowari/common"
	. "github.com/r3inbowari/zlog"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Manifest struct {
	UUID        string
	Name        string `json:"name"`
	Appid       string `json:"appid"`
	Description string `json:"description"`
	VersionName string `json:"versionName"`
	VersionCode int    `json:"versionCode"`
	TransformPx bool   `json:"transformPx"`
	AppPlus     struct {
		UsingComponents bool   `json:"usingComponents"`
		NvueCompiler    string `json:"nvueCompiler"`
		CompilerVersion int    `json:"compilerVersion"`
		Splashscreen    struct {
			AlwaysShowBeforeRender bool `json:"alwaysShowBeforeRender"`
			Waiting                bool `json:"waiting"`
			Autoclose              bool `json:"autoclose"`
			Delay                  int  `json:"delay"`
		} `json:"splashscreen"`
		Modules struct {
			Share struct {
			} `json:"Share"`
			VideoPlayer struct {
			} `json:"VideoPlayer"`
			Push struct {
			} `json:"Push"`
			Geolocation struct {
			} `json:"Geolocation"`
		} `json:"modules"`
		Distribute struct {
			Android struct {
				AutoSdkPermissions bool     `json:"autoSdkPermissions"`
				Permissions        []string `json:"permissions"`
				AbiFilters         []string `json:"abiFilters"`
			} `json:"android"`
			Ios struct {
				Urltypes []struct {
					Urlschemes []string `json:"urlschemes"`
				} `json:"urltypes"`
				Urlschemewhitelist string `json:"urlschemewhitelist"`
				PrivacyDescription struct {
					NSPhotoLibraryUsageDescription               string `json:"NSPhotoLibraryUsageDescription"`
					NSPhotoLibraryAddUsageDescription            string `json:"NSPhotoLibraryAddUsageDescription"`
					NSCameraUsageDescription                     string `json:"NSCameraUsageDescription"`
					NSLocalNetworkUsageDescription               string `json:"NSLocalNetworkUsageDescription"`
					NSLocationWhenInUseUsageDescription          string `json:"NSLocationWhenInUseUsageDescription"`
					NSLocationAlwaysUsageDescription             string `json:"NSLocationAlwaysUsageDescription"`
					NSLocationAlwaysAndWhenInUseUsageDescription string `json:"NSLocationAlwaysAndWhenInUseUsageDescription"`
					NSUserTrackingUsageDescription               string `json:"NSUserTrackingUsageDescription"`
				} `json:"privacyDescription"`
				Capabilities struct {
					Entitlements struct {
						ComAppleDeveloperAssociatedDomains []string `json:"com.apple.developer.associated-domains"`
					} `json:"entitlements"`
				} `json:"capabilities"`
			} `json:"ios"`
			SdkConfigs struct {
				Ad struct {
				} `json:"ad"`
				Share struct {
					Qq struct {
						Appid string `json:"appid"`
					} `json:"qq"`
					Weixin struct {
						Appid          string `json:"appid"`
						UniversalLinks string `json:"UniversalLinks"`
					} `json:"weixin"`
				} `json:"share"`
				Push struct {
					Unipush struct {
					} `json:"unipush"`
				} `json:"push"`
				Oauth struct {
				} `json:"oauth"`
				Geolocation struct {
					Amap struct {
						Platform      []string `json:"__platform__"`
						AppkeyIos     string   `json:"appkey_ios"`
						AppkeyAndroid string   `json:"appkey_android"`
					} `json:"amap"`
				} `json:"geolocation"`
				Maps struct {
					Amap struct {
						AppkeyIos     string `json:"appkey_ios"`
						AppkeyAndroid string `json:"appkey_android"`
					} `json:"amap"`
				} `json:"maps"`
			} `json:"sdkConfigs"`
			Icons struct {
				Android struct {
					Hdpi    string `json:"hdpi"`
					Xhdpi   string `json:"xhdpi"`
					Xxhdpi  string `json:"xxhdpi"`
					Xxxhdpi string `json:"xxxhdpi"`
				} `json:"android"`
				Ios struct {
					Appstore string `json:"appstore"`
					Ipad     struct {
						App            string `json:"app"`
						App2X          string `json:"app@2x"`
						Notification   string `json:"notification"`
						Notification2X string `json:"notification@2x"`
						Proapp2X       string `json:"proapp@2x"`
						Settings       string `json:"settings"`
						Settings2X     string `json:"settings@2x"`
						Spotlight      string `json:"spotlight"`
						Spotlight2X    string `json:"spotlight@2x"`
					} `json:"ipad"`
					Iphone struct {
						App2X          string `json:"app@2x"`
						App3X          string `json:"app@3x"`
						Notification2X string `json:"notification@2x"`
						Notification3X string `json:"notification@3x"`
						Settings2X     string `json:"settings@2x"`
						Settings3X     string `json:"settings@3x"`
						Spotlight2X    string `json:"spotlight@2x"`
						Spotlight3X    string `json:"spotlight@3x"`
					} `json:"iphone"`
				} `json:"ios"`
			} `json:"icons"`
			Splashscreen struct {
				AndroidStyle string `json:"androidStyle"`
				Android      struct {
					Hdpi   string `json:"hdpi"`
					Xhdpi  string `json:"xhdpi"`
					Xxhdpi string `json:"xxhdpi"`
				} `json:"android"`
			} `json:"splashscreen"`
		} `json:"distribute"`
		NativePlugins struct {
			HTML5AppBaichuan struct {
				PluginInfo struct {
					Name               string `json:"name"`
					Description        string `json:"description"`
					Platforms          string `json:"platforms"`
					URL                string `json:"url"`
					AndroidPackageName string `json:"android_package_name"`
					IosBundleID        string `json:"ios_bundle_id"`
					IsCloud            bool   `json:"isCloud"`
					Bought             int    `json:"bought"`
					Pid                string `json:"pid"`
					Parameters         struct {
					} `json:"parameters"`
				} `json:"__plugin_info__"`
			} `json:"Html5app-Baichuan"`
		} `json:"nativePlugins"`
		Privacy struct {
			Prompt   string `json:"prompt"`
			Template struct {
				Title        string `json:"title"`
				Message      string `json:"message"`
				ButtonAccept string `json:"buttonAccept"`
				ButtonRefuse string `json:"buttonRefuse"`
				Second       struct {
					Title        string `json:"title"`
					Message      string `json:"message"`
					ButtonAccept string `json:"buttonAccept"`
					ButtonRefuse string `json:"buttonRefuse"`
				} `json:"second"`
			} `json:"template"`
		} `json:"privacy"`
		UniStatistics struct {
			Enable bool `json:"enable"`
		} `json:"uniStatistics"`
	} `json:"app-plus"`
	Quickapp struct {
	} `json:"quickapp"`
	MpWeixin struct {
		Appid   string `json:"appid"`
		Setting struct {
			URLCheck bool `json:"urlCheck"`
		} `json:"setting"`
		UsingComponents bool `json:"usingComponents"`
	} `json:"mp-weixin"`
	MpAlipay struct {
		UsingComponents bool `json:"usingComponents"`
	} `json:"mp-alipay"`
	MpBaidu struct {
		UsingComponents bool `json:"usingComponents"`
	} `json:"mp-baidu"`
	MpToutiao struct {
		UsingComponents bool `json:"usingComponents"`
	} `json:"mp-toutiao"`
	UniStatistics struct {
		Enable bool `json:"enable"`
	} `json:"uniStatistics"`
	H5 struct {
		SdkConfigs struct {
			Maps struct {
				Qqmap struct {
					Key string `json:"key"`
				} `json:"qqmap"`
			} `json:"maps"`
		} `json:"sdkConfigs"`
	} `json:"h5"`
	SpaceID string `json:"_spaceID"`
}

type Mp struct {
	UsingComponents bool `json:"usingComponents"`
}

func OpenManifest(uuid string) (*Manifest, error) {
	var mf Manifest
	err := common.OpenJsonFromRes(uuid, &mf)
	mf.UUID = uuid
	return &mf, err
}

func (ts *Manifest) Save() (string, error) {
	id := CreateUUID()
	return id, SaveJsonFromHtml(common.GetPath(id), ts)
}

func (ts *Manifest) Update() error {
	if ts.UUID != "" {
		return SaveJsonFromHtml(common.GetPath(ts.UUID), ts)
	} else {
		return errors.New("can not update a new manifest")
	}
}

func SaveJsonFromHtml(path string, v interface{}) error {
	fp, err := os.Create(path)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("loading file failed")
	}
	defer fp.Close()

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(v)

	//data, err := json.Marshal(v)

	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("marshal file failed")
	}
	n, err := fp.Write(bf.Bytes())
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("write file failed")
	}
	Log.WithFields(logrus.Fields{"size": n, "path": path}).Info("[FILE] save file")
	return nil
}

//manifest.AppPlus.Distribute.Icons.Android.Hdpi = c.ManifestSetting.AndroidHdpi
//manifest.AppPlus.Distribute.Icons.Android.Xhdpi = c.ManifestSetting.AndroidXhdpi
//manifest.AppPlus.Distribute.Icons.Android.Xxhdpi = c.ManifestSetting.AndroidXxhdpi
//manifest.AppPlus.Distribute.Icons.Android.Xxxhdpi = c.ManifestSetting.AndroidXxxhdpi
//
//manifest.AppPlus.Distribute.Icons.Ios.Appstore = c.ManifestSetting.Appstore
//manifest.AppPlus.Distribute.Icons.Ios.Ipad.App = c.ManifestSetting.IpadApp
//manifest.AppPlus.Distribute.Icons.Ios.Ipad.App2X = c.ManifestSetting.IpadApp2X
//manifest.AppPlus.Distribute.Icons.Ios.Ipad.Notification = c.ManifestSetting.IpadNotification
//manifest.AppPlus.Distribute.Icons.Ios.Ipad.Notification2X = c.ManifestSetting.IpadNotification2X
//manifest.AppPlus.Distribute.Icons.Ios.Ipad.Proapp2X = c.ManifestSetting.IpadProapp2X
//manifest.AppPlus.Distribute.Icons.Ios.Ipad.Settings = c.ManifestSetting.IpadSettings
//manifest.AppPlus.Distribute.Icons.Ios.Ipad.Settings2X = c.ManifestSetting.IpadSettings2X
//manifest.AppPlus.Distribute.Icons.Ios.Ipad.Spotlight = c.ManifestSetting.IpadSpotlight
//manifest.AppPlus.Distribute.Icons.Ios.Ipad.Spotlight2X = c.ManifestSetting.IpadSpotlight2X
//
//manifest.AppPlus.Distribute.Icons.Ios.Iphone.App2X = c.ManifestSetting.IphoneApp2X
//manifest.AppPlus.Distribute.Icons.Ios.Iphone.App3X = c.ManifestSetting.IphoneApp3X
//manifest.AppPlus.Distribute.Icons.Ios.Iphone.Notification2X = c.ManifestSetting.IphoneNotification2X
//manifest.AppPlus.Distribute.Icons.Ios.Iphone.Notification3X = c.ManifestSetting.IphoneNotification3X
//manifest.AppPlus.Distribute.Icons.Ios.Iphone.Settings2X = c.ManifestSetting.IphoneSettings2X
//manifest.AppPlus.Distribute.Icons.Ios.Iphone.Settings3X = c.ManifestSetting.IphoneSettings3X
//manifest.AppPlus.Distribute.Icons.Ios.Iphone.Spotlight2X = c.ManifestSetting.IphoneSpotlight2X
//manifest.AppPlus.Distribute.Icons.Ios.Iphone.Spotlight3X = c.ManifestSetting.IphoneSpotlight3X
//
//manifest.AppPlus.Distribute.Splashscreen.AndroidStyle = c.ManifestSetting.SplashscreenAndroidStyle
//manifest.AppPlus.Distribute.Splashscreen.Android.Hdpi = c.ManifestSetting.SplashscreenAndroidHdpi
//manifest.AppPlus.Distribute.Splashscreen.Android.Xhdpi = c.ManifestSetting.SplashscreenAndroidXhdpi
//manifest.AppPlus.Distribute.Splashscreen.Android.Xxhdpi = c.ManifestSetting.SplashscreenAndroidXxhdpi

type ManifestConvertor struct {
	PackageName string
	err         error
}

// 处理并试图转换所有的图像指针数据
func (m *ManifestConvertor) convert(vs ...*string) {
	if m.err != nil {
		// 不再做后续转换
		return
	}

	for _, v := range vs {
		if VerifyUUID(*v) {
			// 覆盖资源-uuid转换
			*v = common.GetPath(*v)
			Log.WithFields(logrus.Fields{"item-trace": *v}).Info("[CLI] convert a uuid manifest")
		} else if !filepath.IsAbs(*v) {
			// 原始资源-相对路径转换
			abs, err := filepath.Abs(m.PackageName + "/" + *v)
			if err != nil {
				return
			}
			Log.WithFields(logrus.Fields{"item-trace": *v}).Info("[CLI] convert a absolute origin tag")
			*v = abs
		} else {
			// 原始资源-绝对路径不处理
			Log.WithFields(logrus.Fields{"item-trace": *v}).Info("[CLI] convert a relative origin tag")
		}
		// 统一检查文件是否存在
		if !Exists(*v) {
			Log.WithFields(logrus.Fields{"item-trace": *v}).Error("[CLI] not found target file")
			m.err = errors.New("error path param injected during convert manifest: " + *v)
			// 不再做后续转换
			return
		}
	}
}

func (ts *Manifest) ConvertPath(packageName string) error {
	var mc ManifestConvertor
	mc.PackageName = packageName

	mc.convert(
		&ts.AppPlus.Distribute.Icons.Android.Hdpi,
		&ts.AppPlus.Distribute.Icons.Android.Xhdpi,
		&ts.AppPlus.Distribute.Icons.Android.Xxhdpi,
		&ts.AppPlus.Distribute.Icons.Android.Xxxhdpi,

		&ts.AppPlus.Distribute.Icons.Ios.Appstore,
		&ts.AppPlus.Distribute.Icons.Ios.Ipad.App,
		&ts.AppPlus.Distribute.Icons.Ios.Ipad.App2X,
		&ts.AppPlus.Distribute.Icons.Ios.Ipad.Notification,
		&ts.AppPlus.Distribute.Icons.Ios.Ipad.Notification2X,
		&ts.AppPlus.Distribute.Icons.Ios.Ipad.Proapp2X,
		&ts.AppPlus.Distribute.Icons.Ios.Ipad.Settings,
		&ts.AppPlus.Distribute.Icons.Ios.Ipad.Settings2X,
		&ts.AppPlus.Distribute.Icons.Ios.Ipad.Spotlight,
		&ts.AppPlus.Distribute.Icons.Ios.Ipad.Spotlight2X,

		&ts.AppPlus.Distribute.Icons.Ios.Iphone.App2X,
		&ts.AppPlus.Distribute.Icons.Ios.Iphone.App3X,
		&ts.AppPlus.Distribute.Icons.Ios.Iphone.Notification2X,
		&ts.AppPlus.Distribute.Icons.Ios.Iphone.Notification3X,
		&ts.AppPlus.Distribute.Icons.Ios.Iphone.Settings2X,
		&ts.AppPlus.Distribute.Icons.Ios.Iphone.Settings3X,
		&ts.AppPlus.Distribute.Icons.Ios.Iphone.Spotlight2X,
		&ts.AppPlus.Distribute.Icons.Ios.Iphone.Spotlight3X,

		&ts.AppPlus.Distribute.Splashscreen.Android.Hdpi,
		&ts.AppPlus.Distribute.Splashscreen.Android.Xhdpi,
		&ts.AppPlus.Distribute.Splashscreen.Android.Xxhdpi,
	)

	//mc.convert(&ts.AppPlus.Distribute.Icons.Android.Hdpi)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Android.Xhdpi)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Android.Xxhdpi)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Android.Xxxhdpi)
	//
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Appstore)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Ipad.App)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Ipad.App2X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Ipad.Notification)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Ipad.Notification2X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Ipad.Proapp2X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Ipad.Settings)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Ipad.Settings2X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Ipad.Spotlight)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Ipad.Spotlight2X)
	//
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Iphone.App2X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Iphone.App3X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Iphone.Notification2X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Iphone.Notification3X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Iphone.Settings2X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Iphone.Settings3X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Iphone.Spotlight2X)
	//mc.convert(&ts.AppPlus.Distribute.Icons.Ios.Iphone.Spotlight3X)
	//
	//// mc.convert(&ts.AppPlus.Distribute.Splashscreen.AndroidStyle)
	//mc.convert(&ts.AppPlus.Distribute.Splashscreen.Android.Hdpi)
	//mc.convert(&ts.AppPlus.Distribute.Splashscreen.Android.Xhdpi)
	//mc.convert(&ts.AppPlus.Distribute.Splashscreen.Android.Xxhdpi)

	if mc.err != nil {
		return mc.err
	}
	return ts.Update()
}
