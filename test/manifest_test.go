package test

import (
	"encoding/json"
	"meiwobuxing"
	"testing"
)

func TestManifest(t *testing.T) {
	str := `{
	"name": "花享生活",
	"appid": "__UNI__B8D0466",
	"description": "",
	"versionName": "1.1.9",
	"versionCode": 109,
	"transformPx": false,
	"app-plus": {
		"usingComponents": true,
		"nvueCompiler": "uni-app",
		"compilerVersion": 3,
		"splashscreen": {
			"alwaysShowBeforeRender": true,
			"waiting": true,
			"autoclose": true,
			"delay": 0
		},
		"modules": {
			"Share": {},
			"VideoPlayer": {},
			"Push": {},
			"Geolocation": {}
		},
		"distribute": {
			"android": {
				"autoSdkPermissions": false,
				"permissions": [
					"<uses-feature android:name=\"android.hardware.camera\"/>",
					"<uses-feature android:name=\"android.hardware.camera.autofocus\"/>",
					"<uses-permission android:name=\"android.permission.ACCESS_COARSE_LOCATION\"/>",
					"<uses-permission android:name=\"android.permission.ACCESS_FINE_LOCATION\"/>",
					"<uses-permission android:name=\"android.permission.ACCESS_NETWORK_STATE\"/>",
					"<uses-permission android:name=\"android.permission.ACCESS_WIFI_STATE\"/>",
					"<uses-permission android:name=\"android.permission.CALL_PHONE\"/>",
					"<uses-permission android:name=\"android.permission.CALL_PRIVILEGED\"/>",
					"<uses-permission android:name=\"android.permission.CAMERA\"/>",
					"<uses-permission android:name=\"android.permission.CHANGE_NETWORK_STATE\"/>",
					"<uses-permission android:name=\"android.permission.CHANGE_WIFI_STATE\"/>",
					"<uses-permission android:name=\"android.permission.MODIFY_AUDIO_SETTINGS\"/>",
					"<uses-permission android:name=\"android.permission.MOUNT_UNMOUNT_FILESYSTEMS\"/>",
					"<uses-permission android:name=\"android.permission.READ_LOGS\"/>",
					"<uses-permission android:name=\"android.permission.WAKE_LOCK\"/>",
					"<uses-permission android:name=\"android.permission.WRITE_SETTINGS\"/>"
				],
				"abiFilters": ["armeabi-v7a"]
			},
			"ios": {
				"urltypes": [{
					"urlschemes": ["tbopen32474362"]
				}],
				"urlschemewhitelist": "tbopen,tmall",
				"privacyDescription": {
					"NSPhotoLibraryUsageDescription": "从相册中选择图片作为用户头像",
					"NSPhotoLibraryAddUsageDescription": "写入分享的商品图片",
					"NSCameraUsageDescription": "拍照上传头像",
					"NSLocalNetworkUsageDescription": "应用联网需要",
					"NSLocationWhenInUseUsageDescription": "获取附近站点,提供本地化生活服务",
					"NSLocationAlwaysUsageDescription": "获取附近站点,提供本地化生活服务",
					"NSLocationAlwaysAndWhenInUseUsageDescription": "获取附近站点,提供本地化生活服务",
					"NSUserTrackingUsageDescription": "是否允许花享生活使用您的设备标识(IDFA)信息?以用于向您推荐更感兴趣的内容,并优化我们的技术服务和体验"
				},
				"capabilities": {
					"entitlements": {
						"com.apple.developer.associated-domains": [
							"applinks:static-4dadafa4-717f-465a-bd22-b93ec5ff55cc.bspapp.com",
							"applinks:dccdn.java3.cn"
						]
					}
				}
			},
			"sdkConfigs": {
				"ad": {},
				"share": {
					"qq": {
						"appid": "1111105110"
					},
					"weixin": {
						"appid": "wx7943b5d0592190e3",
						"UniversalLinks": "https://dccdn.java3.cn/uni-universallinks/__UNI__B8D0466"
					}
				},
				"push": {
					"unipush": {}
				},
				"oauth": {},
				"geolocation": {
					"amap": {
						"__platform__": ["ios", "android"],
						"appkey_ios": "d917d16551ead35d22f1e883f0a2c4eb",
						"appkey_android": "10f91611926e85dfaee330ebda3d0f44"
					}
				},
				"maps": {
					"amap": {
						"appkey_ios": "",
						"appkey_android": "0a317eeb54216a9f1d1db65424f43275"
					}
				}
			},
			"icons": {
				"android": {
					"hdpi": "unpackage/res/icons/72x72.png",
					"xhdpi": "unpackage/res/icons/96x96.png",
					"xxhdpi": "unpackage/res/icons/144x144.png",
					"xxxhdpi": "unpackage/res/icons/192x192.png"
				},
				"ios": {
					"appstore": "unpackage/res/icons/1024x1024.png",
					"ipad": {
						"app": "unpackage/res/icons/76x76.png",
						"app@2x": "unpackage/res/icons/152x152.png",
						"notification": "unpackage/res/icons/20x20.png",
						"notification@2x": "unpackage/res/icons/40x40.png",
						"proapp@2x": "unpackage/res/icons/167x167.png",
						"settings": "unpackage/res/icons/29x29.png",
						"settings@2x": "unpackage/res/icons/58x58.png",
						"spotlight": "unpackage/res/icons/40x40.png",
						"spotlight@2x": "unpackage/res/icons/80x80.png"
					},
					"iphone": {
						"app@2x": "unpackage/res/icons/120x120.png",
						"app@3x": "unpackage/res/icons/180x180.png",
						"notification@2x": "unpackage/res/icons/40x40.png",
						"notification@3x": "unpackage/res/icons/60x60.png",
						"settings@2x": "unpackage/res/icons/58x58.png",
						"settings@3x": "unpackage/res/icons/87x87.png",
						"spotlight@2x": "unpackage/res/icons/80x80.png",
						"spotlight@3x": "unpackage/res/icons/120x120.png"
					}
				}
			},
			"splashscreen": {
				"androidStyle": "default",
				"android": {
					"hdpi": "C:/Users/shenlailai/Desktop/桌面备份20210618/O1CN01JdC9582JJi1AjjxR3_!!2053469401.png",
					"xhdpi": "C:/Users/shenlailai/Desktop/桌面备份20210618/O1CN01JdC9582JJi1AjjxR3_!!2053469401.png",
					"xxhdpi": "C:/Users/shenlailai/Desktop/桌面备份20210618/O1CN01JdC9582JJi1AjjxR3_!!2053469401.png"
				}
			}
		},
		"nativePlugins": {
			"Html5app-Baichuan": {
				"__plugin_info__": {
					"name": "Html5app-Baichuan",
					"description": "Android 和 IOS 阿里百川插件",
					"platforms": "Android,iOS",
					"url": "",
					"android_package_name": "",
					"ios_bundle_id": "",
					"isCloud": false,
					"bought": -1,
					"pid": "",
					"parameters": {}
				}
			}
		},
		"privacy": {
			"prompt": "template",
			"template": {
				"title": "服务协议和隐私政策",
				"message": "请你务必审慎阅读、充分理解“服务协议”和“隐私政策”各条款，包括但不限于：为了更好的向你提供服务，我们需要收集你的设备标识、操作日志等信息用于分析、优化应用性能。<br/>　　你可阅读<a  href='http://hx.java3.cn/agreement/str/2'>《服务协议》</a>和<a  href='http://hx.java3.cn/agreement/str/1'>《隐私政策》</a>了解详细信息。如果你同意，请点击下面按钮开始接受我们的服务。",
				"buttonAccept": "同意并继续",
				"buttonRefuse": "不同意",
				"second": {
					"title": "温馨提示",
					"message": "　　进入应用前，你需先同意<a href='http://hx.java3.cn/agreement/str/2'>《服务协议》</a>和<a href='http://hx.java3.cn/agreement/str/1'>《隐私政策》</a>，否则将退出应用。",
					"buttonAccept": "同意并继续",
					"buttonRefuse": "退出应用"
				}
			}
		},
		"uniStatistics": {
			"enable": true
		}
	},
	"quickapp": {},
	"mp-weixin": {
		"appid": "",
		"setting": {
			"urlCheck": false
		},
		"usingComponents": true
	},
	"mp-alipay": {
		"usingComponents": true
	},
	"mp-baidu": {
		"usingComponents": true
	},
	"mp-toutiao": {
		"usingComponents": true
	},
	"uniStatistics": {
		"enable": false
	},
	"h5": {
		"sdkConfigs": {
			"maps": {
				"qqmap": {
					"key": "7ZMBZ-JGFKX-HYR4G-ZUE4Y-DXIJZ-NBBUX"
				}
			}
		}
	},
	"_spaceID": "4dadafa4-717f-465a-bd22-b93ec5ff55cc"
}`
	var st meiwobuxing.Manifest

	json.Unmarshal([]byte(str), &st)

	println(st.Name)
}

func TestConvert(t *testing.T) {
	meiwobuxing.InitFileSystem("12", 100)
	var mf meiwobuxing.Manifest
	mf.AppPlus.Distribute.Icons.Android.Hdpi = "ea76949e-4a1b-4bc2-91e6-4a9056b780e1"
	mf.ConvertPath("appname")
	println(mf.AppPlus.Distribute.Icons.Android.Hdpi)
}