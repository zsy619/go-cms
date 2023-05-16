
### 一、Go-Stress-Testing

Go-Stress-Testing 是一个用于进行压力测试的 Go 语言库。它可以帮助你模拟多个并发用户对你的应用程序进行负载测试，以评估其性能和稳定性。以下是使用 Go-Stress-Testing 进行压力测试的一般步骤：

1. 安装 Go：首先，你需要在你的系统上安装 Go 语言。你可以从 Go 官方网站（https://golang.org）下载并按照安装说明进行安装。

2. 安装 go-stress-testing：在你的项目中，使用以下命令安装 go-stress-testing：

   ```
   go get -u github.com/yudai/go-stress-testing/...
   ```

3. 编写测试脚本：创建一个 Go 文件，例如 `stress_test.go`，在其中编写你的压力测试脚本。下面是一个简单的示例：

   ```go
   package main

   import (
       "fmt"
       "github.com/yudai/go-stress-testing/stress"
       "net/http"
       "time"
   )

   func main() {
       config := stress.NewConfig()
       config.Timeout = 10 * time.Second // 设置请求超时时间
       config.NumUsers = 100             // 设置并发用户数
       config.Interval = 100 * time.Millisecond // 设置请求间隔时间
       config.Method = "GET"             // 设置请求方法
       config.URL = "http://example.com" // 设置请求的 URL

       // 自定义请求函数，可选
       config.RequestFunc = func(client *http.Client, req *http.Request) (*http.Response, error) {
           // 在这里可以对请求进行自定义处理
           // 例如设置请求头、设置请求体等
           return client.Do(req)
       }

       // 自定义结果处理函数，可选
       config.ResultHandler = func(result *stress.Result) {
           // 在这里可以对测试结果进行自定义处理
           fmt.Printf("请求耗时：%v\n", result.ElapsedTime)
           fmt.Printf("返回状态码：%d\n", result.StatusCode)
       }

       stress.Run(config) // 开始压力测试
   }
   ```

   请根据你的具体需求进行适当的更改。

4. 运行压力测试：在命令行中，使用以下命令运行你的压力测试脚本：

   ```
   go run stress_test.go
   ```

   Go-Stress-Testing 将会开始发送请求并模拟并发用户对目标 URL 进行压力测试。

这只是一个简单的示例，Go-Stress-Testing 还提供了更多的配置选项，例如自定义请求头、请求体、身份验证等。你可以查阅官方文档（https://pkg.go.dev/github.com/yudai/go-stress-testing/stress）以了解更多功能和选项。

请注意，在进行压力测试时，确保你对目标应用程序的性能和可靠性有充分的了解，并在测试过程中注意资源消耗，以避免对生产环境造成负面影响。

### 二、Go-Stress-Testing 命令行参数详解

Go-Stress-Testing 提供了一些命令行参数，可以在运行时对压力测试进行配置。以下是一些常用的命令行参数及其说明：

- `-c, --concurrency <num>`：设置并发用户数。例如，`-c 100` 表示使用 100 个并发用户进行压力测试。

- `-n, --num-requests <num>`：设置要发送的请求数量。例如，`-n 1000` 表示发送 1000 个请求。

- `-m, --method <method>`：设置请求方法。例如，`-m POST` 表示使用 POST 方法发送请求。

- `-H, --header <header>`：设置请求头。可以多次使用该选项来设置多个请求头。例如，`-H "Content-Type: application/json" -H "Authorization: Bearer token"`。

- `-b, --body <body>`：设置请求体。例如，`-b '{"key": "value"}'` 表示将 JSON 格式的请求体作为字符串发送。

- `-d, --data <data>`：设置表单数据。可以多次使用该选项来设置多个表单字段。例如，`-d "username=john" -d "password=secret"`。

- `-t, --timeout <duration>`：设置请求超时时间。可以使用 Go 的时间格式，例如 `1s` 表示 1 秒，`500ms` 表示 500 毫秒。

- `-i, --interval <duration>`：设置请求间隔时间。可以使用 Go 的时间格式。例如，`100ms` 表示每个请求之间间隔 100 毫秒。

- `-u, --url <url>`：设置目标 URL。例如，`-u http://example.com` 表示对 `http://example.com` 进行压力测试。

- `--insecure`：禁用 SSL 证书验证。使用该选项可以在测试时忽略 SSL 证书错误。

- `-h, --help`：显示命令行帮助信息。

这些命令行参数可以与 Go-Stress-Testing 的配置结合使用。例如，你可以通过命令行参数设置并发用户数和请求方法，而在代码中通过配置文件或代码中的配置变量设置其他选项。

请注意，命令行参数的使用可能因不同的版本或 fork 的变体而有所不同。在使用特定版本的 Go-Stress-Testing 时，请查阅其文档或运行 `go-stress-testing --help` 命令以获取准确的命令行参数说明。

以下是一个使用 Go-Stress-Testing 进行压力测试的示例命令行：

```
go-stress-testing -c 100 -n 1000 -m POST -H "Content-Type: application/json" -H "Authorization: Bearer token" -b '{"key": "value"}' -t 5s -i 100ms -u http://example.com
```

这个示例命令行将执行以下操作：

- `-c 100`：使用 100 个并发用户进行压力测试。
- `-n 1000`：发送 1000 个请求。
- `-m POST`：使用 POST 方法发送请求。
- `-H "Content-Type: application/json" -H "Authorization: Bearer token"`：设置两个请求头，一个是 "Content-Type"，值为 "application/json"，另一个是 "Authorization"，值为 "Bearer token"。
- `-b '{"key": "value"}'`：将 JSON 格式的请求体作为字符串发送。
- `-t 5s`：设置请求超时时间为 5 秒。
- `-i 100ms`：设置每个请求之间的间隔时间为 100 毫秒。
- `-u http://example.com`：对 `http://example.com` 进行压力测试。

请根据你的具体需求进行适当的更改。你可以根据实际情况设置并发用户数、请求数量、请求方法、请求头、请求体、超时时间、间隔时间和目标 URL。

### 三、使用Go-Stress-Testing测试本网站

./go-stress-testing-mac -c 1 -n 1000 -u http://localhost:8125

 开始启动  并发数:1 请求数:1000 请求参数:
request:
 form:http
 url:http://localhost:8125
 method:GET
 headers:map[Content-Type:application/x-www-form-urlencoded; charset=utf-8]
 data:
 verify:statusCode
 timeout:30s
 debug:false
 http2.0：false
 keepalive：false
 maxCon:1


─────┬───────┬───────┬───────┬────────┬────────┬────────┬────────┬────────┬────────┬────────
 耗时│ 并发数│ 成功数│ 失败数│   qps  │最长耗时│最短耗时│平均耗时│下载字节│字节每秒│ 状态码
─────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────┼────────┼────────
   1s│      1│    153│      0│  164.66│   57.66│    4.07│    6.07│2,953,818│2,951,981│200:153
   2s│      1│    309│      0│  165.35│   57.66│    4.07│    6.05│5,965,554│2,982,677│200:309
   3s│      1│    463│      0│  164.77│   57.66│    3.66│    6.07│8,938,678│2,979,520│200:463
   4s│      1│    655│      0│  175.06│   57.66│    3.42│    5.71│12,645,430│3,161,288│200:655
   5s│      1│    798│      0│  170.42│   63.16│    3.42│    5.87│15,406,188│3,081,182│200:798
   6s│      1│    952│      0│  169.51│   63.16│    3.42│    5.90│18,379,312│3,063,123│200:952
   6s│      1│   1000│      0│  167.85│   63.16│    3.42│    5.96│19,306,000│3,033,172│200:1000


*************************  结果 stat  ****************************
处理协程数量: 1
请求总数（并发数*请求数 -c * -n）: 1000 总请求时间: 6.365 秒 successNum: 1000 failureNum: 0
tp90: 7.000
tp95: 8.000
tp99: 26.000
*************************  结果 end   ****************************