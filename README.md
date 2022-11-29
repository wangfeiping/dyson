# Dyson
Command wrapper capable of being invoked remotely or automatically executed

The purpose of this project is to provide decentralized automatic maintenance and management of account assets (e.g., Queries; Allows agencies to provide reverse management of their clients' private assets). Includes features like: command executes automatically; configurable prometheus-based data metrics monitoring; APIs for remote calls, etc.

## Key features under consideration:

* (Implemented) Program command executes automatically
* (Implemented) Configurable prometheus-based data metrics monitoring
* (Not yet implemented) APIs for remote calls
* (Not yet implemented) Password caching and automatically entering passwords if necessary

In order to achieve the project objectives, the following issues need to be considered:

* Security
  * 每次启动程序时都需要输入密码，密码仅缓存在内存中；(密码指私钥库的解锁密码，通过密码以获得账户私钥)
  * 私钥库由原程序管理，本身不会对私钥解锁或通过其他手段获取，仅在程序需要签名确认交易时自动输入缓存的密码；
  * 因为是通过系统调用实现程序的代理调用，处于安全考虑对外部调用的API需要进行白名单配置，没有配置的默认不可执行；
* Internal command execution
  * 开放全功能，可通过配置文件配置保存；
  * 可借助原程序的私钥库功能，管理执行多账户操作；
* Monitor
  * prometheus exporter统计执行警告及异常数量；
  * 各功能可用性，健康度；
  * 内部命令及外部调用执行统计数；
  * 通过配置命令生成监控指标？如：新提案
* Logs
  * 记录每次执行的命令；
  * 记录执行警告及异常；
* APIs for remote calls
  * health检测接口
  * 仅可执行单次调用命令；

## Usage and configuration instructions


