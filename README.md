# ShuYou Go

ShuYou Go AI agents platform — 基于 Gin 的分层 Web 服务。

## 技术栈

| 组件 | 包 |
|------|-----|
| Go | 1.26 |
| Web 框架 | [gin-gonic/gin](https://github.com/gin-gonic/gin) |
| 配置 | [spf13/viper](https://github.com/spf13/viper) |
| 日志 | [uber-go/zap](https://github.com/uber-go/zap) + [lumberjack](https://github.com/natefinch/lumberjack) |
| 数据库 | [MongoDB Go Driver v2](https://www.mongodb.com/docs/drivers/go/current/) |
| 校验 | [go-playground/validator](https://github.com/go-playground/validator) |
| JWT | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) |

## 目录结构

```
.
├── cmd/server/              # 应用入口
├── configs/                 # 配置文件
│   ├── config.example.yaml  # 配置模板（入库）
│   └── config.yaml          # 本地配置（gitignore）
├── internal/
│   ├── app/                 # 依赖组装与生命周期
│   ├── config/              # 配置加载与校验
│   ├── dto/                 # 请求/响应 DTO
│   ├── errors/              # 业务错误码
│   ├── handler/             # HTTP 处理器
│   ├── infra/               # 基础设施（db/jwt/logger）
│   ├── middleware/          # Gin 中间件
│   ├── model/               # 数据模型
│   ├── repository/          # 数据访问层
│   ├── router/              # 路由注册
│   └── service/             # 业务逻辑层
└── pkg/
    ├── response/            # 统一 HTTP 响应
    └── validator/           # 参数校验
```

## 快速开始

```bash
# 1. 复制配置
make setup

# 2. 启动 MongoDB（Docker 示例）
docker run -d -p 27017:27017 --name mongo mongo:7

# 3. 安装依赖并启动
go mod tidy
make run
```

服务默认监听 `http://0.0.0.0:8080`。

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/health/live` | 存活探针 |
| GET | `/api/v1/health/ready` | 就绪探针（含 DB 检查） |
| POST | `/api/v1/users/register` | 用户注册（返回 token） |
| POST | `/api/v1/users/login` | 用户登录（返回 token） |
| GET | `/api/v1/users/me` | 当前用户（需 Bearer Token） |

## 配置

默认读取 `configs/config.yaml`，支持 `SHUYOU_` 环境变量覆盖：

```bash
export SHUYOU_SERVER_PORT=9090
export SHUYOU_JWT_SECRET=your-secret
export SHUYOU_DATABASE_URI=mongodb://localhost:27017
export SHUYOU_DATABASE_DATABASE=shuyou
```

## 开发

```bash
make fmt      # 格式化
make test     # 运行测试
make build    # 编译 bin/server
```
