package brew

import (
	"context"
	"fmt"
	"os/exec"

	mcp_core "github.com/mark3labs/mcp-go/mcp"
)

func Update(ctx context.Context, request mcp_core.CallToolRequest) (*mcp_core.CallToolResult, error) {
	binary, err := exec.LookPath(brewBinary)
	if err != nil {
		return nil, fmt.Errorf("find brew binary :%w", err)
	}

	output, err := commandOutputContext(ctx, binary, "update")
	if err != nil {
		return nil, fmt.Errorf("brew update: %w", err)
	}

	return mcp_core.NewToolResultText(output), nil
}
