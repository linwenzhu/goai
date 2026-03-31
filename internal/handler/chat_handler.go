package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"goAi/pkg/ai"
	"goAi/pkg/response"
	"io"
)

type ChatHandler struct{}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误"+err.Error())
		return
	}
	resp, err := ai.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "moonshot-v1-8k",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: req.Message,
				},
			},
		},
	)
	if err != nil {
		response.Fail(c, 500, "AI 调用失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"reply": resp.Choices[0].Message.Content,
	})
}

func (h *ChatHandler) ChatStream(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误"+err.Error())
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	stream, err := ai.Client.CreateChatCompletionStream(
		c.Request.Context(),
		openai.ChatCompletionRequest{
			Model:  "moonshot-v1-8k",
			Stream: true,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: req.Message,
				},
			},
		},
	)
	if err != nil {
		response.Fail(c, 500, "AI 调用失败: "+err.Error())
		return
	}
	defer stream.Close()
	c.Stream(func(w io.Writer) bool {
		recv, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// 流结束，发送结束标记
			c.SSEvent("message", "[DONE]")
			return false
		}
		if err != nil {
			c.SSEvent("message", "[ERROR]")
			return false
		}

		content := recv.Choices[0].Delta.Content
		if content != "" {
			c.SSEvent("message", content)
		}
		return true
	})
}
