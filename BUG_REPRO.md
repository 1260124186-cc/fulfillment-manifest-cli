# 修复前故障复现（Docker）

## 项目与标准命令

`manifestctl` 从标准输入接收一份 JSON 订单请求，并在标准输出生成履约清单。仓库根目录可执行下列标准命令：

```sh
go build ./...
printf '%s\n' '{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":2}]}' | go run ./cmd/manifestctl
go test ./...
```

其中 `go build ./...` 可以成功完成；本故障会使包含取消和超时场景的测试失败。

## 环境构建与编译

已实际执行下列命令：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t fulfillment-manifest-cli-bug-004-base-amd64 .
docker run --rm --platform linux/amd64 fulfillment-manifest-cli-bug-004-base-amd64 go build ./...
docker build --platform linux/arm64 -f benzhi.Dockerfile -t fulfillment-manifest-cli-bug-004-base-arm64 .
docker run --rm --platform linux/arm64 fulfillment-manifest-cli-bug-004-base-arm64 go build ./...
```

linux/amd64 与 linux/arm64 镜像均构建成功，且两种平台的容器内 `go build ./...` 均成功。

## 故障触发步骤

在仓库根目录执行：

```sh
docker run --rm --platform linux/amd64 fulfillment-manifest-cli-bug-004-base-amd64 go test ./cmd/manifestctl ./internal/service -run "TestRunPropagatesCallerCancellation|TestPlannerStopsWhenStorageContextExpires" -count=1
```

该命令会稳定触发已取消请求仍返回成功、以及存储等待超时后仍未返回超时错误的场景。

## 实际错误输出

```text
--- FAIL: TestRunPropagatesCallerCancellation (0.01s)
    main_test.go:58: run() code = 0, want cancellation code 3; stderr =
FAIL
FAIL	github.com/1260124186-cc/fulfillment-manifest-cli/cmd/manifestctl	0.060s
--- FAIL: TestPlannerStopsWhenStorageContextExpires (0.16s)
    planner_test.go:49: Plan() error = <nil>, want deadline exceeded
FAIL
FAIL	github.com/1260124186-cc/fulfillment-manifest-cli/internal/service	0.201s
FAIL
```

## 期望行为

当调用方已取消订单请求时，命令应结束并返回取消状态，而不是输出成功清单。等待存储操作超过请求时限时，应返回超时错误且不再将该订单作为已成功生成的清单保存。调用方随后重试时，应能根据实际结果决定后续操作，而不会遇到“已经停止却仍产生成功结果”的不一致状态。
