# 聊天服压力测试指标

## 压测程序使用手册

#### 1. 启动压测程序

```
./start_test.sh
```

#### 2.参数详解

```
-t=3 -addr=http://127.0.0.1:8200 -playerNum=30000 -chatCount=50 -local=1 -loginAddr=http://127.0.0.1:8000

-local: 非0表示本地压测
-t  :0 默认压测 发送消息, 阻塞接口消息
    :1 发送消息
    :2 阻塞接口消息  
    :3 非阻塞接收消息 旨在测试ws连接数量
-addr=http://127.0.0.1:8200  :默认服务器地址
-playerNum: 玩家数量
-chatCount: 聊天数量(每个玩家发送的消息数量)
-loginAddr=http://127.0.0.1:8000    :登录服务器地址
-accountAddr=http://127.0.0.1:8001  :账号服务器地址
-percentChatPlayers=0.1             :聊天玩家占所有玩家数量的百分比
-c  --------------------------------:线程数
```

## go-stress-testing压测工具

````
Go-stress-testing 测试工具   这个牛逼
https://blog.csdn.net/GoNewWay/article/details/130887182
./go-stress-testing-win.exe -n 100 -c 10 -H 'content-type:application/json' -H 'Authorization: Bearer eyJ1c2VyaWQiOjEwMjQwNDE1fQ==.e25d5b8e9e84da6029ea550070ee8f28' -u http://8.219.160.79:81/xy3-cross/new-chat/SendMessage -data '{"RoomType":1,"Msg":"我是xxx"}'

````
