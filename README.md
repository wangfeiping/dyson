# lever
Command wrapper capable of being invoked remotely or automatically executed

本项目目的是为了提供自动维护和管理(周期自动)账户资产(比如查询并委托账户资产)，同时提供外部调用(单次)功能(比如提交交易)，但无需每次输入密码解锁秘钥。

目前考虑的主要功能：

* 密码缓存
* 自动输入密码
* 程序命令代理执行

为了实现项目目标，需考虑如下几个问题：

* 安全
  * 每次启动程序时都需要输入密码，密码仅缓存在内存中；(密码指私钥库的解锁密码，通过密码以获得账户私钥)
  * 私钥库由原程序管理，本身不会对私钥解锁或通过其他手段获取，仅在程序需要签名确认交易时自动输入缓存的密码；
  * 因为是通过系统调用实现程序的代理调用，处于安全考虑对外部调用的API需要进行白名单配置，没有配置的默认不可执行；
* 内部命令执行
  * 开放全功能，可通过配置文件配置保存；
* 外部调用API
  * health检测接口
  * 仅可执行单次调用命令；
* 日志
  * 记录每次执行的命令；
  * 记录执行警告及异常；
* 监控
  * prometheus exporter统计执行警告及异常数量；
  * 各功能可用性，健康度；
  * 内部命令及外部调用执行统计数；
  * 通过配置命令生成监控指标？
