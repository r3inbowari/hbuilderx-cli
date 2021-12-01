package test

import (
	"meiwobuxing"
	"testing"
)

func TestHyperLink(t *testing.T) {
	k := "123"
	s := "12:51:57.306 检查云端打包状态...\n12:51:57.726 检查打包资源...\n12:51:58.238 正在编译打包资源...\n12:52:42.982 压缩打包资源...\n12:52:43.920 向云端发送打包请求...\n12:52:54.497 项目 xiguazhu_app [__UNI__B8D0466]的打包状态：时间: 2021-11-07 12:52:54    类型: Android自有证书    准备打包    \n打包成功后会自动返回下载链接。打包过程查询请点菜单发行-查看云打包状态。周五傍晚等高峰期打包排队较长，请耐心等待。如果是为了三方SDK调试，请使用自定义调试基座（菜单运行-手机或模拟器-制作自定义调试基座），不要反复打包。\n12:53:50.818 类型: Android自有证书 下载地址: https://ide.dcloud.net.cn/build/download/913ffdf0-3f86-11ec-af7a-498927f9a93f （注意该地址为临时下载地址，只能下载5次）\n12:52:54.551 项目 xiguazhu_app [__UNI__B8D0466]的打包状态：时间: 2021-11-07 12:52:54    类型: Android自有证书    准备打包    \n12:53:50.818 项目 xiguazhu_app [__UNI__B8D0466]打包成功：\n12:53:50.818 类型: IOS自有证书 下载地址: https://ide.dcloud.net.cn/build/download/913ffdf0-3f86-11ec-af7a-498927f9a93e （注意该地址为临时下载地址，只能下载5次）\n\n"
	_, err := meiwobuxing.GetHyperLinks(k)
	println(err.Error())

	b, err := meiwobuxing.GetHyperLinks(s)
	println(b[0])
}

func TestGet(t *testing.T) {
	// fatal
	meiwobuxing.Get("https://www.qq.com")
}

func TestVerifyUUID(t *testing.T) {
	println(meiwobuxing.VerifyUUID("ea76949e-4a1b-4bc2-91e6-4a9056b780e1"))
}

func TestGenAppConfig(t *testing.T) {

}