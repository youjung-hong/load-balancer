package main

import (
	"fmt"
	"math/rand"
	"time"
)

// MessageData 구조체는 메시지를 나타냅니다.
type MessageData struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIProvider 인터페이스는 모델과 메시지를 얻는 메서드를 정의합니다.
type AIProvider interface {
	GetModel() string
	GetMessages() []MessageData
}

// ClaudeRequest 및 OpenAIRequest가 AIProvider 인터페이스를 구현합니다.
type ClaudeRequest struct {
	Model    string        `json:"model"`
	Messages []MessageData `json:"messages"`
}

func (c ClaudeRequest) GetModel() string {
	return c.Model
}

func (c ClaudeRequest) GetMessages() []MessageData {
	return c.Messages
}

type OpenAIRequest struct {
	Model    string        `json:"model"`
	Messages []MessageData `json:"messages"`
}

func (o OpenAIRequest) GetModel() string {
	return o.Model
}

func (o OpenAIRequest) GetMessages() []MessageData {
	return o.Messages
}

// Node 구조체는 AIProvider를 포함하여 로드 밸런서가 관리할 각 노드를 나타냅니다.
type Node struct {
	Provider AIProvider
}

// LoadBalancer 구조체는 여러 노드를 관리하고 요청을 분배합니다.
type LoadBalancer struct {
	Nodes []*Node
}

// NewLoadBalancer 함수는 노드를 포함하여 새로운 로드 밸런서를 생성합니다.
func NewLoadBalancer(nodes []*Node) *LoadBalancer {
	return &LoadBalancer{Nodes: nodes}
}

// GetNextNode 함수는 무작위로 다음 노드를 선택합니다.
func (lb *LoadBalancer) GetNextNode() *Node {
	rand.Seed(time.Now().UnixNano())
	nodeIndex := rand.Intn(len(lb.Nodes))
	return lb.Nodes[nodeIndex]
}

// ProcessRequest 함수는 선택된 노드에서 AIProvider 인터페이스를 사용하여 요청을 처리합니다.
func (lb *LoadBalancer) ProcessRequest() {
	node := lb.GetNextNode()
	model := node.Provider.GetModel()
	messages := node.Provider.GetMessages()
	fmt.Printf("Processing request for model: %s\n", model)
	for _, message := range messages {
		fmt.Printf("[%s]: %s\n", message.Role, message.Content)
	}
}

func main() {
	// 각 노드에 ClaudeRequest와 OpenAIRequest를 설정합니다.
	claudeRequest := &ClaudeRequest{
		Model: "claude-v1",
		Messages: []MessageData{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Hello from Claude!"},
		},
	}

	openAIRequest := &OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []MessageData{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Hello from OpenAI!"},
		},
	}

	// Node 리스트에 각 요청을 포함하여 로드 밸런서를 생성합니다.
	nodes := []*Node{
		{Provider: claudeRequest},
		{Provider: openAIRequest},
	}

	// 로드 밸런서 생성 및 요청 분배 테스트
	loadBalancer := NewLoadBalancer(nodes)
	loadBalancer.ProcessRequest()
}

