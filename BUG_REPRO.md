# 修复前故障复现（Docker）

## 项目与标准命令

`manifestctl` 从标准输入接收一笔订单 JSON，并输出履约清单 JSON。在仓库根目录可执行以下标准命令：

```sh
go build ./...
printf '%s\n' '{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":1}]}' | go run ./cmd/manifestctl
go test ./...
```

## 环境构建与编译

已实际执行以下命令。`linux/arm64` 与 `linux/amd64` 的镜像均构建成功，且各自在容器内执行 `go build ./...` 均成功：

```sh
docker build --platform linux/arm64 -f benzhi.Dockerfile -t fulfillment-manifest-cli-bug003-base:arm64 .
docker run --rm --platform linux/arm64 fulfillment-manifest-cli-bug003-base:arm64 go build ./...
docker build --platform linux/amd64 -f benzhi.Dockerfile -t fulfillment-manifest-cli-bug003-base:amd64 .
docker run --rm --platform linux/amd64 fulfillment-manifest-cli-bug003-base:amd64 go build ./...
```

## 故障触发步骤

在仓库根目录执行以下命令，使用 `linux/amd64` 镜像稳定触发重复订单冲突未被正确识别的问题：

```sh
docker run --rm --platform linux/amd64 fulfillment-manifest-cli-bug003-base:amd64 go test ./cmd/manifestctl ./internal/service -run 'TestExitCodeForDuplicateOrder|TestPlannerPreservesDuplicateOrderError' -count=1
```

## 实际错误输出

```text
--- FAIL: TestExitCodeForDuplicateOrder (0.01s)
    main_test.go:40: exitCodeFor() = 1, want 2
FAIL
FAIL	github.com/1260124186-cc/fulfillment-manifest-cli/cmd/manifestctl	0.066s
--- FAIL: TestPlannerPreservesDuplicateOrderError (0.01s)
    planner_test.go:50: second Plan() error = persist manifest: manifest already exists for order, want duplicate-order classification
FAIL
FAIL	github.com/1260124186-cc/fulfillment-manifest-cli/internal/service	0.053s
FAIL
```

退出结果为 `1`。

## 期望行为

同一订单已生成履约清单后再次提交时，调用方应能观察到可识别的业务冲突：命令行返回冲突退出状态 `2`，并输出 `order already has a manifest`。库存不足仍应作为普通失败，取消或超时仍应返回取消状态，正常订单继续输出履约清单 JSON。
