# 更新记录

## 2023-01-28

- 更新 更换logrotate工具包为lumberjack，重点支持log文件gzip压缩，以减小日志文件过大问题
- 更新 调整部分logger配置项
  ```yaml
  loggerCfg:
  caller: true
  console: true
  level: "debug"
  logs:
    # 文件路径
    path: "logs/service-manager.log"
    # 最大保留天数，单位：天，默认：7天
    maxAge: 30
    # 文件分割大小，单位：MB，默认：100MB
    size: 100
    # 最大保留文件数量，默认：0为全部保留
    count: 0
    # 是否使用UTC时间
    isUTC: false 
    # 是否启用gzip压缩
    isGzip: true
  file:
  # 文件模式默认启用压缩且不使用UTC时间，如需支持请完善相应代码
    - path: "logs/debug.log"
      # 保存日志级别
      level: "debug"
      # 最大保留天数，单位：天，默认：7天
      maxAge: 7
      # 文件分割大小，单位：MB，默认：100MB
      size: 500
      # 最大保留文件数量，默认：0为全部保留
      count: 0
    - path: "logs/info.log"
      level: "info"
      maxAge: 30
  ```