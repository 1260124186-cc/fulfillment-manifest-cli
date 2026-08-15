# Docker 交付与验收说明

本项目是一个履约清单 CLI：它读取一笔订单的 JSON 请求，校验包裹和配送时间，并输出对应的库存预留清单。项目不依赖外部数据库或网络服务，适合用 Docker 在不同 CPU 架构下从源码进行一致性验证。

## 本机验证

以下命令都在仓库根目录执行：

```sh
go build ./...
printf '%s\n' '{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":1}]}' | go run ./cmd/manifestctl
go test ./...
```

第一条命令成功退出表示所有 Go 包可以编译；第二条命令会输出一行 JSON 清单；第三条命令成功退出表示公开行为测试全部通过。

## 固定环境

实际 Dockerfile 为 `benzhi.Dockerfile`。它使用 `golang:1.26.2` 完整工具链，先执行 `go mod download`，再将仓库源码复制进容器并执行 `go build ./...`。`go.mod` 同样固定 `go 1.26.2` 和 `toolchain go1.26.2`，Docker 中设置 `GOTOOLCHAIN=local`，因此容器不会自动切换到其他 Go 工具链。

镜像用于从源码构建和运行项目，不复制宿主机二进制，也不跳过容器内编译。`.dockerignore` 会排除 Git 历史、缓存、补丁、二进制和本地验证证据。

## 标准脚本验收

脚本会依次构建指定架构镜像、启动容器执行 `go build ./...`，再启动 CLI 并向标准输入传入一笔订单。分别执行以下两条命令：

```sh
./build_benzhi_docker.sh fulfillment-manifest-cli:amd64 linux/amd64
./build_benzhi_docker.sh fulfillment-manifest-cli:arm64 linux/arm64
```

每次脚本以退出码 `0` 结束，且最后输出一行含有 `status":"planned"` 的 JSON，即表示该架构下镜像构建、容器内编译和项目启动均通过。

## 手工 Docker 验证

如需分步执行，可先构建并验证 `linux/amd64`：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t fulfillment-manifest-cli:amd64 .
docker run --rm --platform linux/amd64 fulfillment-manifest-cli:amd64 go build ./...
printf '%s\n' '{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":1}]}' | docker run --rm -i --platform linux/amd64 fulfillment-manifest-cli:amd64
docker run --rm --platform linux/amd64 fulfillment-manifest-cli:amd64 go test ./...
```

将命令中的 `linux/amd64` 与镜像标签中的 `amd64` 分别替换为 `linux/arm64` 和 `arm64`，即可进行 arm64 验证。四条手工命令都应以退出码 `0` 结束；运行入口时应产生 JSON 清单，运行测试时应显示全部包通过。
