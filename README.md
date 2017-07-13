# Go文件服务说明
基于Go实现的一个简单的文件服务器，只是简单的一个将多服务器中分散存储的文件集中管理，以方便多负载环境下使用。前期暂不支持：集群、自动复制与同步、冗余存储等

## 配置

> Port:端口 默认为 *19860*  
> User:用户名
> Pwd:密码  
> Path:{名称:目录路径} 默认为{"Defaut":"./attachments"  }
> AllowedIP: 英文分号分隔的IP白名单, 默认不限制IP

## 命令

the protocol follow to [redis protocol](https://redis.io/topics/protocol) :
* For Simple Strings the first byte of the reply is "+"
* For Errors the first byte of the reply is "-"
* For Integers the first byte of the reply is ":"
* For Bulk Strings the first byte of the reply is "$" + bytes Count
* For Array the first byte of the reply is "*" + elements count
* **[ext]** For Bytes the first byte of the replay is "!" + byte array count

### `Auth "User" "md5(pwd+timestamp)" "timestamp"`
    身份验证 timestamp: unix time string

### `Select "PathName"`
    切换目录 默认为Default=attachemnts

### `RELOAD`
    重新加载配置

### `Exsits "fileFullName"`
    文件是否在该服务器上存在
    C: $6\r\nEXISTS\r\n
    C: $12\r\nfileFullName\r\n
    S: :0

### `Len "fileFullName"`
    返回文件长度，如果文件不存在返回-1

### `Get "fileFullName"`
    获取文件
    返回文件的字节流

### `CGet "fileFullName" offset length`
    获取文件的一部分
    返回：[文件原始长度][开始位置][实际长度][文件内容]

### `Save "fileFullName" [bytes]`
    保存文件
    成功返回 OK，否则返回错误原因

### `Del "fileFullName"`
    删除文件
    成功返回 OK，否则返回错误原因

### `Ping`
    C: $4\r\nPING\r\n
    S: +PONG

### `PingFile "fileFullName"`
    测试文件是否存在或所在的服务器
    ???

### `DIR "folder"`
    返回对应目录下的所有文件和目录， 目录以/结束

### `Meet IP:Port`
    ??? 连接两个或多个文件服务
