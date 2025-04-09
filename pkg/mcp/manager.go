package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"go.uber.org/zap"
)

type Handler = func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)

type Tool struct {
	T       mcp.Tool
	Handler Handler
}

type Manager struct {
	server *server.MCPServer
	tools  []Tool
	logger *zap.SugaredLogger
}

func NewManager(tools []Tool, enabledTools []string, logger *zap.SugaredLogger) *Manager {
	manager := Manager{
		server: server.NewMCPServer(
			"Brew MCP sever",
			"1.0.0",
			server.WithToolCapabilities(true),
		),
		logger: logger,
		tools:  tools,
	}

	for _, tool := range manager.tools {
		if toolEnabled(tool, enabledTools) {
			logger.Debugf("register tool %q", tool.T.Name)
			manager.server.AddTool(tool.T, tool.Handler)
		}
	}

	return &manager
}

func toolEnabled(tool Tool, enabledTools []string) bool {
	// register all tools
	if len(enabledTools) == 0 {
		return true
	}

	for _, enabledTool := range enabledTools {
		if enabledTool == tool.T.Name {
			return true
		}
	}
	return false
}

func (m *Manager) Serve() error {
	if err := server.ServeStdio(m.server); err != nil {
		return fmt.Errorf("mcp server serve: %w", err)
	}
	return nil
}
