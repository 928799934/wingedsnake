# 腾蛇(WingedSnake)
腾蛇是模仿nginx实现的服务端框架
## 支持
热重启  
守护进程  
配置文件热更新
## 效果截图
![效果截图](./res/video.gif)  
## 使用帮助
参考examples内的实现
## 效果对比
![性能变化](./res/pic.png)
左侧为使用腾蛇效果
右侧为使用runtime.GOMAXPROCS效果
## 配置说明
```ini
[Base]
; 使用用户
User=nobody
; 使用组
Group=nobody
; pid 存放路径
PidPath=.
; 多进程数量
Process=2
; CPU亲和 (不设置该参数需要使用runtime.GOMAXPROCS 否则程序将运行于单个CPU)
; 此处个数与process相同 代表各个进程绑定到哪个CPU
; 如果process数量大于affinity 那么剩余的process进程不会指定cpu亲和
; 如果process数量小于affinity 那么剩余的affinity不会被指定到进程
; 建议在affinity 参数与 runtime.GOMAXPROCS 中自由选择一个
Affinity=0001,0010

[Client]
; 子配置文件
Config=logformat.xml
; socket 监听
Sockets=0.0.0.0:9900,0.0.0.0:9920,127.0.0.1:9910
```
## 备注
1.内部自带CPU 亲和配置  自由选择 使用 runtime.GOMAXPROCS 函数  或使用cpu亲和配置
2.不支持windows
