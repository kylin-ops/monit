# 1 编译
```shell script
go build main.go
```

# 2 启动服务
```shell script
nohup main &
```

# 3 配置说明
配置文件与主程序在统一目录下
```json
{
  # http服务配置
  "http_server": {
    "address": "0.0.0.0:8090",                   # http服务监听地址
    "ngx_log_path": "/tmp/log",                  # 分析的nginx服务日志的文件路径
    "exclude_file": ".+\\.(js|css|html|ico)$"    # 过滤掉指定的日志
  },
  # 阿里云短信配置,配置自己阿里云的相关内容
  "aliyun_sms": {
    "region_id": "",
    "access_key_id": "",
    "access_key_secret": "",
    "sign_name": "",
    "template_code": "",                         # 短信模板必须接收content参数
    "phones": "13770519341,13874592241"
  },
     # http检测配置
     "http_detectors": [
       {
         "url": "http://master.dev.kangyitong.cn:18081/uaa/oauth/token",  # http探测地址
         "username": "10000000000",                                       # 登录用户名
         "password": "e10adc3949ba59abbe56e057f20f883e",                  # 登录密码
         "interval": 300,                                                 # http探测的间隔时间(秒)
         "failed_retry": 3,                                               # 失败重试多少次
         "alert_pre_hour": 3,                                             # 每个小时告警多少次
         "start_alert_time": "08:00:00",                                  # 开发告警时间
         "end_alert_time": "20:00:00"                                     # 结束告警时间
       }
     ]
   },
  # 服务日志配置
  "logger": {
    "level": "info",
    "path": "/tmp/detector.log",
    "roll_time": 24,
    "roll_count": 5,
    "is_console": true

}
```