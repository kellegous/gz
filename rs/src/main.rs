use anyhow::Result;
use clap::{Parser, Subcommand};
use gdg::{checkout, create};

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
}

fn main() -> Result<()> {
    let args = Args::parse();

    match args.command {
        Command::Create(args) => create::run(args)?,
        Command::Checkout(args) => checkout::run(args)?,
    }

    Ok(())
}
