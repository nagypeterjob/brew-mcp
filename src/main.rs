mod flag;
use brew::cli::CLI;
use clap::Parser;
use log::{error, info};
use rmcp::{ServiceExt, transport::stdio};
mod brew;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    env_logger::init();
    info!("Starting virtualization-cli MCP server");

    let args = flag::Args::parse();
    let cli_path = String::from(args.brew_override);
    let service = CLI::new(cli_path)
        .serve(stdio())
        .await
        .inspect_err(|e| error!("serving error: {:?}", e))?;

    service.waiting().await?;
    Ok(())
}
