package brew

import (
	"context"
	"fmt"
	"os/exec"

	mcp_core "github.com/mark3labs/mcp-go/mcp"
)

func ListInstalledPackagesAndVersions(
	ctx context.Context,
	request mcp_core.CallToolRequest,
) (*mcp_core.CallToolResult, error) {
	onlyCask, ok := request.Params.Arguments["only_cask"].(bool)
	if !ok {
		return mcp_core.NewToolResultError("provided variable 'only_cask' was missing"), nil
	}

	binary, err := exec.LookPath(brewBinary)
	if err != nil {
		return nil, fmt.Errorf("find brew binary :%w", err)
	}

	args := []string{
		"list",
		"--version",
	}

	if onlyCask {
		args = append(args, "--cask")
	}

	output, err := commandOutputContext(
		ctx,
		binary,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("brew list installed: %w", err)
	}

	return mcp_core.NewToolResultText(output), nil
}
