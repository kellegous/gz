use anyhow::Result;
use clap::{Parser, Subcommand};
use gdg::create;

#[derive(Debug, Parser)]
struct Args {
    #[command(subcommand)]
    command: Command,
}

#[derive(Subcommand, Debug)]
enum Command {
    #[command(alias = "+", alias = "add")]
    Create(create::Args),
}

fn main() -> Result<()> {
    let args = Args::parse();

    match args.command {
        Command::Create(args) => {
            create::run(args)?;
        }
    }

    Ok(())
}
