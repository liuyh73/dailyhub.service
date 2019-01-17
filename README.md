## dailyhub.service
此项目为dailyhub应用的后台服务实现。

### 源代码部署应用

**新建数据库**

创建数据库`dailyhub`，并使用`db`文件夹下的`data.sql`创建数据库中的表。

**运行服务器**

将`db`文件夹下的`conf.example.yml`更名为`conf.yml`，并进行相应的配置（数据库用户名和密码）。然后在项目根目录下运行：

```bash
$ go run main.go
[negroni] listening on :9090
......
```

### docker部署服务

在项目根目录下执行：

```bash
$ docker-compose up -d
```

