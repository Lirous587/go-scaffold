package httpserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		// 处理错误
		if len(ctx.Errors) > 0 {
			for _, e := range ctx.Errors {
				// 记录详细错误日志
				//log.Printf("Error: %+v\n", e.Err)

				// 使用自定义格式化错误栈
				printBusinessStack(e.Err)
			}
		}
	}
}

func printBusinessStack(err error) {
	// 获取完整错误栈
	stackTrace := fmt.Sprintf("%+v", err)
	lines := strings.Split(stackTrace, "\n")

	// 错误消息
	if len(lines) > 0 {
		log.Printf("\n\n")
		log.Printf("业务逻辑错误: %s\n", lines[0])
	}

	// 记录已打印的栈帧数量
	framePrinted := 0
	maxBusinessFrames := 3 // 最多打印栈帧条数

	// 逐行检查并不做任何修改，保持原始格式
	for i := 0; i < len(lines)-1 && framePrinted < maxBusinessFrames; i++ {
		currentLine := lines[i]
		nextLine := lines[i+1]

		// 只检查是否为业务相关行，但完全保持原始格式
		if strings.Contains(currentLine, "internal") &&
			!strings.Contains(currentLine, "github.com/gin-gonic") &&
			!strings.Contains(currentLine, "net/http") &&
			strings.Contains(nextLine, ".go:") {
			log.Println(currentLine)
			log.Println(nextLine)
			framePrinted++
		}
	}

	// 如果还有更多栈帧但已达到限制
	totalBusinessFrames := countBusinessFrames(lines)
	if framePrinted == maxBusinessFrames && framePrinted < totalBusinessFrames {
		log.Printf("一共%d条栈帧,实际打印%d条 (更多栈帧已省略)\n", totalBusinessFrames, maxBusinessFrames)
	}
}

// 计算业务栈帧总数
func countBusinessFrames(lines []string) int {
	count := 0
	for i := 0; i < len(lines)-1; i++ {
		currentLine := lines[i]
		nextLine := lines[i+1]

		if strings.Contains(currentLine, "internal") &&
			!strings.Contains(currentLine, "github.com/gin-gonic") &&
			!strings.Contains(currentLine, "net/http") &&
			strings.Contains(nextLine, ".go:") {
			count++
		}
	}
	return count
}
