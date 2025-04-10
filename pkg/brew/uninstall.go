package brew

import (
	"context"
	"fmt"
	"os/exec"

	mcp_core "github.com/mark3labs/mcp-go/mcp"
)

func Uninstall(
	ctx context.Context,
	request mcp_core.CallToolRequest,
) (*mcp_core.CallToolResult, error) {
	args := []string{"uninstall"}
	formula, ok := request.Params.Arguments["formula"].(string)
	if !ok {
		return mcp_core.NewToolResultError("package name was missing"), nil
	}
	args = append(args, formula)

	_, ok = request.Params.Arguments["force"].(bool)
	if ok {
		args = append(args, "--force")
	}

	binary, err := exec.LookPath(brewBinary)
	if err != nil {
		return nil, fmt.Errorf("find brew binary :%w", err)
	}

	output, err := commandOutputContext(
		ctx,
		binary,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("brew uninstall: %w", err)
	}

	return mcp_core.NewToolResultText(output), nil
}
