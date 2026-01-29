# golang-snippets

通用 Go 代码片段合集（偏实用/可直接复制）。目标：
- 快速解决日常工程问题（IO/时间/并发/网络/JSON/日志等）
- 每个片段都 **可运行 / 可测试 / 有最小示例**

## 目录（建议结构）

- `snippets/`
  - `net/`：HTTP 客户端、重试、超时、下载、Webhook
  - `time/`：时间格式、时区、cron 相关
  - `concurrency/`：goroutine 池、限流、超时取消、errgroup
  - `io/`：读写文件、csv/excel、流处理
  - `json/`：struct tag、动态字段、性能注意事项
  - `crypto/`：hash/hmac、AES/RSA、签名验签
  - `testing/`：测试技巧、golden files、httptest、mock

> 目前仓库里已有：`email.go`

## 片段规范（贡献指南）

- 每个片段尽量做到：
  - 文件头部有用途说明、使用示例
  - **最小可运行示例**（`func main()` 或 `*_test.go`）
  - 需要外部依赖的，说明原因并控制在少量

## 快速开始

```bash
go version
# 推荐 Go 1.21+
```

## License

MIT
