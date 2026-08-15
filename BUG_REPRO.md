# 修复前故障复现（Docker）

## 项目与标准命令

`manifestctl` 从标准输入接收一个订单 JSON，并输出履约清单。仓库根目录的标准命令如下：

```sh
go build ./...
printf '%s\n' '{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":1}]}' | go run ./cmd/manifestctl
go test ./...
```

## 环境构建与编译

已实际执行以下两种平台的镜像构建与容器内编译命令，镜像构建和容器内 `go build ./...` 均成功：

```sh
./build_benzhi_docker.sh fulfillment-manifest-cli-bug-002-base-arm64 linux/arm64
docker run --rm --platform linux/arm64 fulfillment-manifest-cli-bug-002-base-arm64 go build ./...
./build_benzhi_docker.sh fulfillment-manifest-cli-bug-002-base-amd64 linux/amd64
docker run --rm --platform linux/amd64 fulfillment-manifest-cli-bug-002-base-amd64 go build ./...
```

## 故障触发步骤

在仓库根目录执行以下命令：

```sh
docker run --rm --platform linux/amd64 fulfillment-manifest-cli-bug-002-base-amd64 go test ./internal/service ./internal/store -run 'TestPlannerUsesDefaultWindowForFlexibleDelivery|TestMemoryRepositoryAllowsManifestWithoutDeliveryWindow' -count=1
```

## 实际错误输出

```text
--- FAIL: TestPlannerUsesDefaultWindowForFlexibleDelivery (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x8 pc=0x5e0fd3]

goroutine 17 [running]:
testing.tRunner.func1.2({0x64c620, 0x8932d0})
	/usr/local/go/src/testing/testing.go:1974 +0x232
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1977 +0x349
panic({0x64c620?, 0x8932d0?})
	/usr/local/go/src/runtime/panic.go:860 +0x13a
github.com/1260124186-cc/fulfillment-manifest-cli/internal/store.cloneManifest(...)
	/workspace/internal/store/memory.go:109
github.com/1260124186-cc/fulfillment-manifest-cli/internal/store.(*memorySession).Save(0x10f7bbdf00f0, {0x6a6da8?, 0x8bf0c0?}, {{0x68dfe6, 0x7}, {0x68d6c1, 0x3}, 0x0, {0x10f7bbdf0108, 0x1, ...}, ...})
	/workspace/internal/store/memory.go:60 +0xf3
github.com/1260124186-cc/fulfillment-manifest-cli/internal/service.(*Planner).Plan(0x10f7bbbf0ee8, {0x6a6da8, 0x8bf0c0}, {{0x68dfe6, 0x7}, {0x68d6c1, 0x3}, {0x10f7bbdf00a8, 0x1, 0x1}, ...})
	/workspace/internal/service/planner.go:43 +0x547
github.com/1260124186-cc/fulfillment-manifest-cli/internal/service.TestPlannerUsesDefaultWindowForFlexibleDelivery(0x10f7bbe00248)
	/workspace/internal/service/planner_test.go:39 +0x1e5
testing.tRunner(0x10f7bbe00248, 0x6a10f0)
	/usr/local/go/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:2101 +0x4c5
FAIL	github.com/1260124186-cc/fulfillment-manifest-cli/internal/service	0.057s
--- FAIL: TestMemoryRepositoryAllowsManifestWithoutDeliveryWindow (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x8 pc=0x52ec73]

goroutine 19 [running]:
testing.tRunner.func1.2({0x55ce20, 0x6ced00})
	/usr/local/go/src/testing/testing.go:1974 +0x232
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1977 +0x349
panic({0x55ce20?, 0x6ced00?})
	/usr/local/go/src/runtime/panic.go:860 +0x13a
github.com/1260124186-cc/fulfillment-manifest-cli/internal/store.cloneManifest(...)
	/workspace/internal/store/memory.go:109
github.com/1260124186-cc/fulfillment-manifest-cli/internal/store.(*memorySession).Save(0x2c21d25a60f0, {0x598d08?, 0x6f9c20?}, {{0x589724, 0xc}, {0x0, 0x0}, 0x0, {0x0, 0x0, ...}, ...})
	/workspace/internal/store/memory.go:60 +0xf3
github.com/1260124186-cc/fulfillment-manifest-cli/internal/store.TestMemoryRepositoryAllowsManifestWithoutDeliveryWindow(0x2c21d25f6248)
	/workspace/internal/store/memory_test.go:44 +0x186
testing.tRunner(0x2c21d25f6248, 0x595d48)
	/usr/local/go/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:2101 +0x4c5
FAIL	github.com/1260124186-cc/fulfillment-manifest-cli/internal/store	0.037s
FAIL

EXIT_CODE: 1
```

## 期望行为

客户未填写配送时段时，系统应生成包含默认配送安排的履约清单并正常输出。运营保存和读取未填写配送时段的历史清单时，应保持清单可用且不发生进程崩溃。
