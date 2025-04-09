package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"syscall"

	mcp_core "github.com/mark3labs/mcp-go/mcp"

	zlogger "github.com/bitrise-io/brew-mcp/internal/logger"
	"github.com/bitrise-io/brew-mcp/pkg/mcp"
	"github.com/bitrise-io/brew-mcp/pkg/mcp/brew"
	"go.uber.org/zap"
)

const development = "development"

// BuildVersion is overwritten with go build flags.
var BuildVersion = development // nolint:gochecknoglobals

func main() {
	ctx := context.Background()
	done := make(chan struct{})

	// Cancel context on SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		select {
		case <-ctx.Done():
			close(done)
			return
		}
	}()

	logLevel, tools := parseFlags()
	logger, err := zlogger.NewStructuredLogger(logLevel)
	if err != nil {
		log.Fatalf("setup logger: %s", err.Error())
	}

	logger.Debugf("Server version: %s", BuildVersion)

	if err := run(tools, logger); err != nil {
		log.Fatalf("run server: %s", err.Error())
	}
	<-done
}

func run(enabledTools []string, logger *zap.SugaredLogger) error {
	tools := []mcp.Tool{
		{
			T: mcp_core.NewTool("update_brew",
				mcp_core.WithDescription("Fetch the newest version of Homebrew and all formulae from GitHub using git(1) and perform any necessary migrations."),
			),
			Handler: brew.Update,
		},
		{
			T: mcp_core.NewTool("upgrade_specific_package",
				mcp_core.WithDescription("Upgrade outdated casks and outdated, unpinned formulae."),
				mcp_core.WithString("formula",
					mcp_core.Required(),
					mcp_core.Description("Name of the package to upgrade"),
				),
			),
			Handler: brew.Upgrade,
		},
		{
			T: mcp_core.NewTool("list_installed_package_versions",
				mcp_core.WithDescription("List all installed formulae and casks"),
				mcp_core.WithBoolean("only_cask",
					mcp_core.Description("list versions for casks only")),
			),
			Handler: brew.ListInstalledPackagesAndVersions,
		},
		{
			T: mcp_core.NewTool("get_package_version",
				mcp_core.WithDescription("Get version for specific installed formulae or casks"),
				mcp_core.WithString("formula",
					mcp_core.Required(),
					mcp_core.Description("name of the specifc formula or cask")),
			),
			Handler: brew.GetPackageVersion,
		},
	}

	ntools := len(tools)
	if len(enabledTools) > 0 {
		ntools = len(enabledTools)
	}
	logger.Debugf("Create manager with %d tools", ntools)
	mcpManager := mcp.NewManager(tools, enabledTools, logger)

	logger.Debug("Start to serve")
	if err := mcpManager.Serve(); err != nil {
		return fmt.Errorf("run server :%w", err)
	}

	return nil
}

func parseFlags() (string, []string) {
	var tools string
	flag.StringVar(&tools, "enabled-tools", "", "comma separated list of enabled tool(s)")

	var level string
	flag.StringVar(&level, "log-level", "info", "specify log granularity")

	flag.Parse()

	if len(tools) == 0 {
		return level, []string{}
	}

	toolsArr := make([]string, 0)
	for _, tool := range strings.Split(tools, ",") {
		toolsArr = append(toolsArr, tool)
	}

	return level, toolsArr
}
