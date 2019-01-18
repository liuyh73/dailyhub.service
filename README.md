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

在`db`文件夹下增加`db_user_password.txt`和`db_root_password.txt`，其中内容即为user和root的密码（密码为**dailyhub**）。

在项目根目录下执行：

```bash
$ docker-compose up -d
```

将`db`文件夹下的`data.sql`文件拷贝到`mysql`容器中：

```bash
$ docker cp db/data.sql mysql:/mysql
```

之后，使用`docker exec -it <mysql容器id> /bin/bash`运行容器：

```bash
# 登录mysql
$ mysql -u dailyhub -p 
password: dailyhub
```

初始化数据库：

```bash
mysql> create database dailyhub;
mysql> source /mysql/data.sql
```

