# Go + MySQL 云托管服务

## 快速开始

### 配置 MySQL 连接参数

1. 初始化 MySQL 参考 [开通及使用文档](https://docs.cloudbase.net/database/configuration/db/tdsql/initialization).
2. 前往云开发控制台, 选择 `数据库 - MySQL`, 在右上角点击 `连接数据库`, 在 `连接字符串` 标签页可以获得数据库实例可以获得连接参数 `host`, `port`, `database` 参数.
3. [可选] 如果不希望使用管理员账号, 在弹窗中点击 `忘记密码` 按钮, 可以创建一套新的用户名密码.

### 部署服务

前往 [腾讯云云开发控制台](https://tcb.cloud.tencent.com/dev), 选择 `云托管`, 在 `创建服务` 区块选择 `通过模板部署`, 选择模版, 进入模板代码预览页面.
在代码预览区, 找到 `.env` 文件, 点击右上角 `编辑` 按钮, 修改 MySQL 链接相关参数, 并点击 `保存` 按钮.

```
MYSQL_ADDRESS=<mysql-instance-host>:<mysql-instance-port>
MYSQL_USERNAME=<mysql-username>
MYSQL_PASSWORD=<mysql-password>
MYSQL_DATABASE=<mysql-database>
```

之后点击 `部署` 按钮, 等待部署完成. 部署完成后, 在云托管的 `服务列表` 页面, 点击 `服务名称` 进入服务详情页, 通过服务域名即可访问服务.

## 本地调试

### 环境要求

- Go 1.23+
- MySQL 5.7+ 或 MySQL 8.0+

### 1. 克隆项目并安装依赖

```bash
# 克隆项目
git clone https://github.com/TencentCloudBase/awesome-cloudbase-examples.git
cd cloudbaserun/golang-mysql

# 使用 Go modules 安装依赖
go mod download

# [可选] 使用国内镜像源以提高下载速度
go env -w GOPROXY=https://mirrors.cloud.tencent.com/go/,direct
```

### 2. 配置 MySQL 数据库

创建 `.env` 文件，配置你的本地 MySQL 连接参数：

```bash
MYSQL_USERNAME=your_username
MYSQL_PASSWORD=your_password
MYSQL_ADDRESS=localhost:3306
MYSQL_DATABASE=your_database
# 可以用 PORT 环境变量指定其他端口
PORT=8080
```

### 3. 启动本地开发服务器

确保本地 MySQL 服务已启动

```bash
# 直接运行
go run ./cmd/server

# 或者先构建再运行
go build -o server ./cmd/server
./server
```

应用将在 `http://localhost:8080` 启动。

### 4. 测试应用功能

- 访问 `http://localhost:8080` 查看前端页面
- 使用页面上的按钮测试计数器功能：
  - 点击 "计数+1" 按钮增加计数
  - 点击 "清零" 按钮重置计数
- 或者直接调用 API 接口：

  ```bash
  # 获取当前计数
  curl http://localhost:8080/api/count

  # 增加计数
  curl -X POST http://localhost:8080/api/count \
    -H "Content-Type: application/json" \
    -d '{"action": "inc"}'

  # 清空计数
  curl -X POST http://localhost:8080/api/count \
    -H "Content-Type: application/json" \
    -d '{"action": "clear"}'
  ```

### 常见问题排查

- **数据库连接失败**：检查 `.env` 文件中的数据库配置和 MySQL 服务状态
- **端口被占用**：修改 `PORT` 环境变量或直接修改代码中的端口号
- **权限问题**：确保 MySQL 用户有对应数据库的读写权限
- **Go 版本问题**：确保使用 Go 1.23 或更高版本

### 使用 Docker 进行本地调试（可选）

```bash
# 构建镜像
docker build -t golang-mysql-demo .

# 运行容器（需要先启动 MySQL 容器）
docker run --rm -p 8080:8080 golang-mysql-demo
```

## Dockerfile 最佳实践

请参考 [优化容器镜像](https://docs.cloudbase.net/run/develop/image-optimization)
