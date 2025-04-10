package brew

import (
	"context"
	"fmt"
	"os/exec"

	mcp_core "github.com/mark3labs/mcp-go/mcp"
)

func Info(ctx context.Context, request mcp_core.CallToolRequest) (*mcp_core.CallToolResult, error) {
	formula, ok := request.Params.Arguments["formula"].(string)
	if !ok {
		return mcp_core.NewToolResultError("package name was missing"), nil
	}

	binary, err := exec.LookPath(brewBinary)
	if err != nil {
		return nil, fmt.Errorf("find brew binary :%w", err)
	}

	output, err := commandOutputContext(ctx, binary, "info", formula)
	if err != nil {
		return nil, fmt.Errorf("brew info: %w", err)
	}

	return mcp_core.NewToolResultText(output), nil
}
