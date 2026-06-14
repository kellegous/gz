use crate::store;
use anyhow::Result;
use clap::Parser;
use std::fs;

#[derive(Debug, Parser)]
pub struct Args {
    aliases: Vec<String>,
}

pub fn run(args: Args) -> Result<()> {
    let git_root =
        store::find_git_root()?.ok_or_else(|| anyhow::anyhow!("not in a git repository"))?;

    let (mut store, store_path) = store::load_store_git_root(&git_root)?;

    for alias in args.aliases {
        store.unalias_branch(&alias)?;
    }

    store.to_writer(fs::File::create(&store_path)?)?;

    Ok(())
}
