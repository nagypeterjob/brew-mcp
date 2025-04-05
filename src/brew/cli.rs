use rmcp::{Error as McpError, ServerHandler, model::*, tool};
use std::process::Command;

#[derive(Clone)]
pub struct CLI {
    path: String,
}

#[tool(tool_box)]
impl CLI {
    pub fn new(path: String) -> Self {
        Self { path }
    }

    #[tool(description = "List all installed formulae and casks")]
    async fn list_installed_packages_and_versions(
        &self,
        #[tool(param)]
        #[schemars(description = "list versions for casks only")]
        only_cask: bool,
    ) -> Result<CallToolResult, McpError> {
        let cmd_path = self.path.clone();
        let cmd = Command::new(cmd_path)
            .arg("list")
            .arg("--version")
            .arg(if only_cask { "--cask" } else { "" })
            .output()
            .expect("failed to execute list process");

        let output_string = String::from_utf8_lossy(&cmd.stdout).to_string();

        Ok(CallToolResult::success(vec![Content::text(output_string)]))
    }

    #[tool(description = "Get version for specific installed formulae or casks")]
    async fn get_package_version(
        &self,
        #[tool(param)]
        #[schemars(description = "name of the specifc formula or cask")]
        formula: String,
    ) -> Result<CallToolResult, McpError> {
        if formula.is_empty() {
            return Err(McpError::invalid_params("package name was missing", None));
        }
        let cmd = Command::new(&self.path)
            .arg("list")
            .arg("--version")
            .arg(formula)
            .output()
            .expect("failed to execute list process");

        let cmd_output = String::from_utf8_lossy(&cmd.stdout).to_string();

        Ok(CallToolResult::success(vec![Content::text(cmd_output)]))
    }

    #[tool(
        description = "Fetch the newest version of Homebrew and all formulae from GitHub using git(1)
    and perform any necessary migrations."
    )]
    async fn update(&self) -> Result<CallToolResult, McpError> {
        let cmd = Command::new(&self.path)
            .arg("update")
            .output()
            .expect("failed to execute update");

        match cmd.status.success() {
            true => {
                let cmd_stdout = String::from_utf8_lossy(&cmd.stdout).to_string();
                return Ok(CallToolResult::success(vec![Content::text(cmd_stdout)]));
            }
            _ => {
                let cmd_stderr = String::from_utf8_lossy(&cmd.stderr).to_string();

                Err(McpError::invalid_params(
                    format!("error while updating brew: {}", cmd_stderr),
                    None,
                ))
            }
        }
    }

    #[tool(description = "Upgrade outdated casks and outdated, unpinned formulae")]
    async fn upgrade(
        &self,
        #[tool(param)]
        #[schemars(description = "name of the specifc formula or cask")]
        formula: String,
    ) -> Result<CallToolResult, McpError> {
        if formula.is_empty() {
            return Err(McpError::invalid_params("package name was missing", None));
        }
        let cmd = Command::new(&self.path)
            .arg("upgrade")
            .arg(&formula)
            .status()
            .expect("failed to execute upgrade");

        match cmd.success() {
            true => Ok(CallToolResult::success(vec![Content::text(format!(
                "{} was upgraded successfully",
                formula
            ))])),
            _ => Err(McpError::invalid_params("error while updating brew", None)),
        }
    }
}

#[tool(tool_box)]
impl ServerHandler for CLI {
    fn get_info(&self) -> ServerInfo {
        ServerInfo {
            protocol_version: ProtocolVersion::V_2024_11_05,
            capabilities: ServerCapabilities::builder()
                .enable_prompts()
                .enable_resources()
                .enable_tools()
                .build(),
            server_info: Implementation::from_build_env(),
            instructions: Some(
                "This server provides a tool to interact with Brew command line tool".to_string(),
            ),
        }
    }
}
