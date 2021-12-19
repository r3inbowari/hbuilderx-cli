# hbuilderx-cli
a simple wrapper around hbuilderx cli in Go   
这是一个使用go封装hbuilderx cli的批量云打包工具    

## 配置文件
```
{
    "log_level": "debug",
    "check_link": "",
    "auto_update": true,
    "max_retry_count": 3,
    "api_addr": ":9090",
    "ca_key": "",
    "ca_crt": "",
    "log_path": "null",                            
    "jwt_enable": false,
    "jwt_secret": "r3inbowari",
    "jwt_timeout": 72,                              // 有效期
    "jwt_md5": "48dc65a51b480244c296c44de7be53f5"   // 使用 md5 命令生成
}
```

