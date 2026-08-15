# 修复前故障复现（Docker）

## 项目与标准命令

Fulfillment Manifest CLI 从标准输入接收一笔订单 JSON，并在标准输出生成履约清单；内存仓储用于覆盖订单规划和会话边界。

在仓库根目录可执行以下标准命令：

```sh
go build ./...
printf '%s\n' '{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":1}]}' | go run ./cmd/manifestctl
go test ./...
```

## 环境构建与编译

在修复前的 base 状态仓库根目录，实际执行以下命令：

```sh
docker build --platform linux/arm64 -f benzhi.Dockerfile -t fulfillment-manifest-cli-bug-005-base:verify .
docker run --rm --platform linux/arm64 fulfillment-manifest-cli-bug-005-base:verify go build ./...
docker build --platform linux/amd64 -f benzhi.Dockerfile -t fulfillment-manifest-cli-bug-005-base:amd64 .
docker run --rm --platform linux/amd64 fulfillment-manifest-cli-bug-005-base:amd64 go build ./...
```

linux/arm64 与 linux/amd64 的镜像构建均成功，两个容器内的 `go build ./...` 均成功。

## 故障触发步骤

在修复前的 base 状态仓库根目录，先按上节构建 `fulfillment-manifest-cli-bug-005-base:verify`，再执行：

```sh
docker run --rm --platform linux/arm64 fulfillment-manifest-cli-bug-005-base:verify go test ./internal/service ./internal/store -run 'TestPlannerReleasesSessionForTheNextOrder|TestClosingCommittedSessionKeepsManifest' -count=1
```

该命令会稳定失败。

## 实际错误输出

```text
--- FAIL: TestPlannerReleasesSessionForTheNextOrder (0.00s)
    planner_test.go:46: Plan("order-5b") error = begin manifest session: another manifest session is active
FAIL
FAIL	github.com/1260124186-cc/fulfillment-manifest-cli/internal/service	0.005s
--- FAIL: TestClosingCommittedSessionKeepsManifest (0.00s)
    memory_test.go:51: Get() did not retain committed manifest after Close()
FAIL
FAIL	github.com/1260124186-cc/fulfillment-manifest-cli/internal/store	0.008s
FAIL
```

## 期望行为

同一个进程连续处理不同订单时，后一笔订单应能正常开始。订单已经确认成功后，即使对应会话结束，查询同一订单仍应返回已保存的履约清单。
