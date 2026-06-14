use anyhow::Result;
use clap::{Parser, Subcommand};
use gdg::{checkout, create, delete};

#[derive(Debug, Parser)]
struct Args {
    #[command(subcommand)]
    command: Command,
}

#[derive(Subcommand, Debug)]
enum Command {
    #[command(alias = "+", alias = "add")]
    Create(create::Args),
    #[command(alias = "co")]
    Checkout(checkout::Args),
    #[command(alias = "rm")]
    Delete(delete::Args),
}

fn main() -> Result<()> {
    let args = Args::parse();

    match args.command {
        Command::Create(args) => create::run(args)?,
        Command::Checkout(args) => checkout::run(args)?,
        Command::Delete(args) => delete::run(args)?,
    }

    Ok(())
}
