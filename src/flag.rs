use anyhow::{Ok, bail};
use clap::Parser;
use std::path::Path;

#[derive(Parser, Debug)]
#[clap(about, long_about = None)]
pub struct Args {
    #[clap(
        default_value = "/opt/homebrew/bin/brew",
        value_name = "brew-override",
        short = 'o',
        value_parser=validate_brew_override,
        help = "override default /opt/homebrew/bin/brew location")]
    pub brew_override: String,
}

fn validate_brew_override(arg: &str) -> anyhow::Result<String> {
    if !Path::new(arg).exists() {
        bail!("brew-override must be valid path, pointing to a binary")
    }

    Ok(arg.to_string())
}
