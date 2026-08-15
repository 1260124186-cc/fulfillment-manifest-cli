# 修复前故障复现（Docker）

## 项目与标准命令

`manifestctl` 从标准输入读取一笔 JSON 订单并输出履约清单。以下命令均在仓库根目录执行：

```sh
go build ./...
printf '%s\n' '{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":1}]}' | go run ./cmd/manifestctl
go test ./...
```

## 环境构建与编译

已实际执行以下命令构建 linux/amd64 和 linux/arm64 镜像：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t fulfillment-manifest-cli-bug-001-base-linux-amd64:validation .
docker build --platform linux/arm64 -f benzhi.Dockerfile -t fulfillment-manifest-cli-bug-001-base-linux-arm64:validation .
```

已实际执行以下命令，并确认两个平台的容器内编译均成功：

```sh
docker run --rm --platform linux/amd64 fulfillment-manifest-cli-bug-001-base-linux-amd64:validation go build ./...
docker run --rm --platform linux/arm64 fulfillment-manifest-cli-bug-001-base-linux-arm64:validation go build ./...
```

## 故障触发步骤

在仓库根目录先构建 linux/amd64 镜像，再执行以下命令：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t fulfillment-manifest-cli-bug-001-base-linux-amd64:validation .
docker run --rm --platform linux/amd64 fulfillment-manifest-cli-bug-001-base-linux-amd64:validation go test ./internal/service -run TestPlannerKeepsRequestAndRetrievedManifestIndependent -count=1
```

以下正常订单命令也已在 linux/amd64 容器中实际执行，可输出状态为 `planned` 的履约清单：

```sh
printf '%s\n' '{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":1}]}' | docker run -i --rm --platform linux/amd64 fulfillment-manifest-cli-bug-001-base-linux-amd64:validation
```

## 实际错误输出

```text
--- FAIL: TestPlannerKeepsRequestAndRetrievedManifestIndependent (0.00s)
    planner_test.go:48: Plan() changed the caller request package to "book"
FAIL
FAIL	github.com/1260124186-cc/fulfillment-manifest-cli/internal/service	0.035s
FAIL
```

## 期望行为

同一笔订单完成规划后，调用方保留的原始包裹内容应保持原样。读取已保存的清单后在本地调整预留数量，也不应影响随后读取到的同一订单清单。
