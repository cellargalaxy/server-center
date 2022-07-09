# server_center

小小的配置中心，用轮询拉配置。

小小服务，轮询就够了。用法详见server_center_client_test.go

妈妈再也不用担心我一而再，再而三的把配置提交到github上来了。

```
@startuml
sc -> sdk: 导入
sdk -> sdk0: 创建sdk0
sdk0 -> sdk0: 地址为空
sdk0 -> sdk0: 读取&写入本地配置
sdk0 -> sdk0: 本地配置地址为空
sdk0 -> sdk: 
sdk -> sc: 
sc -> sdk: 创建sdk
sdk -> sdk: 地址没有
sdk -> sdk: 读取本地配置
sdk -> sc: 
sc -> sc: 本地服务启动

sdk0 -> sdk0: 地址没有
sdk0 -> sdk0: 读取本地配置
sdk0 -> sdk0: 本地配置地址为空

sdk -> sdk: 地址没有
sdk -> sdk: 读取本地配置
@enduml
```

```
@startuml
sc -> sdk: 导入
sdk -> sdk0: 创建sdk0
sdk0 -> sdk0: 地址有了
sdk0 -> sdk0: 读取远端配置
sdk0 -> sdk0: 更新本地配置
sdk0 -> sdk0: 更新地址
sdk0 -> sdk: 
sdk -> sc: 
sc -> sdk: 创建sdk
sdk -> sdk: 地址有了
sdk -> sdk: 读取远端配置
sdk -> sdk: 更新本地配置
sdk -> sdk: 更新地址
sdk -> sc: 
sc -> sc: 本地服务启动

sdk0 -> sdk0: 地址有了
sdk0 -> sdk0: 读取远端配置
sdk0 -> sdk0: 更新本地配置
sdk0 -> sdk0: 更新地址

sdk -> sdk: 地址有了
sdk -> sdk: 读取远端配置
sdk -> sdk: 更新本地配置
@enduml
```

```
@startuml
other -> sdk: 导入
sdk -> sdk0: 创建sdk0
sdk0 -> sdk0: 地址为空
sdk0 -> sdk0: 读取&写入本地配置
sdk0 -> sdk0: 本地配置地址为空
sdk0 -> sdk: 
sdk -> other: 
other -> sdk: 创建sdk
sdk -> sdk: 地址没有
sdk -> sdk: 读取本地配置
sdk -> sdk: 本地配置非法
sdk -> other: 杀进程

other -> other: 手改本地配置

other -> sdk: 导入
sdk -> sdk0: 创建sdk0
sdk0 -> sdk0: 地址为空
sdk0 -> sdk0: 读取&写入本地配置
sdk0 -> sdk0: 本地配置地址为空
sdk0 -> sdk: 
sdk -> other: 
other -> sdk: 创建sdk
sdk -> sdk: 地址没有
sdk -> sdk: 读取本地配置
sdk -> sdk: 本地配置合法
sdk -> other: 
other -> other: 服务启动

sdk0 -> sdk0: 地址没有
sdk0 -> sdk0: 读取本地服务
sdk0 -> sdk0: 本地配置更新没有

sdk -> sdk: 地址没有
sdk -> sdk: 读取本地服务
@enduml
```