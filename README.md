# Homebrew MCP Server

The project enables you to have a "natural language" conversation with brew. The most important commands are implemented:

[x] Install  
[x] Uninstall  
[x] Get specific package  
[x] List installed packages and their versions  
[x] Update  
[x] Upgrade  
[x] Info  
[x] Search  
[x] Print config  

## install

Build or download the latest binary from the Github releases section.
Move the binary to `/usr/local/bin/`, or to any preferred location in `$PATH`.

## Usage

### Basic

Place the following json into your:
`~/Library/Application\ Support/Claude/claude_desktop_config.json` (or other location depending on your preferred Client).

```json
{
  "mcpServers": {
    "brew": {
      "command": "brew-mcp-server"
    }
  }
}
```

### Advanced

The tool lets you enable only specific tools to save context size.

## Tools:

- install
- uninstall
- get_package_version
- list_installed_package_versions
- update_brew
- upgrade_specific_package
- info
- search
- config

To enable specific tools, use the `-enabled-tools` flag to list tools in a comma separated manner. Example:

```json
{
  "mcpServers": {
    "brew": {
      "command": "brew-mcp-server",
      "args": ["-enabled-tools", "info,search,update_brew"]
    }
  }
}
```
