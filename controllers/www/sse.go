package www

import (
	"time"
)

type SSEController struct{ BaseController }

func (ctrl *SSEController) Message() {
	ctrl.Ctx.ResponseWriter.Header().Set("Content-Type", "text/event-stream")
	ctrl.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache")
	ctrl.Ctx.ResponseWriter.Header().Set("Connection", "keep-alive")
	ctrl.Ctx.ResponseWriter.Header().Set("Transfer-Encoding", "chunked")

	for {
		// 模拟ChatGPT生成的文本
		text := generateText()

		ctrl.Ctx.ResponseWriter.Write([]byte("data: " + text + "\n\n"))
		ctrl.Ctx.ResponseWriter.Flush()

		// 休眠一段时间，控制打字速度
		time.Sleep(time.Millisecond * 50)
	}
}

func generateText() string {
	// 在这里实现ChatGPT生成文本的逻辑
	// 返回生成的文本
	return "Hello World!"
}
