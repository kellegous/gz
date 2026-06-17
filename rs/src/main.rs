use anyhow::Result;
use clap::{Parser, Subcommand};
use gdg::{alias, checkout, create, delete, rebase, rebase_abort, rebase_continue, unalias};

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
    Alias(alias::Args),
    Unalias(unalias::Args),
    Rebase(rebase::Args),
    #[command(name = "rebase:abort")]
    RebaseAbort,
    #[command(name = "rebase:continue")]
    RebaseContinue,
}

fn main() -> Result<()> {
    let args = Args::parse();

    match args.command {
        Command::Create(args) => create::run(args)?,
        Command::Checkout(args) => checkout::run(args)?,
        Command::Delete(args) => delete::run(args)?,
        Command::Alias(args) => alias::run(args)?,
        Command::Unalias(args) => unalias::run(args)?,
        Command::Rebase(args) => rebase::run(args)?,
        Command::RebaseAbort => rebase_abort::run()?,
        Command::RebaseContinue => rebase_continue::run()?,
    }

    Ok(())
}
