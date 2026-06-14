use crate::{BranchRef, store};
use anyhow::Result;
use clap::Parser;
use git2::Repository;
use std::fs;

#[derive(Debug, Parser)]
pub struct Args {
    name: String,
    aliases: Vec<String>,
}

pub fn run(args: Args) -> Result<()> {
    let git_root =
        store::find_git_root()?.ok_or_else(|| anyhow::anyhow!("not in a git repository"))?;

    let (mut store, store_path) = store::load_store_git_root(&git_root)?;

    let repo = Repository::open(&git_root)?;
    let head = repo.head()?;
    let current_branch = head.shorthand()?;

    let name = args
        .name
        .parse::<BranchRef>()?
        .resolve(&store, current_branch)?;

    for alias in args.aliases {
        store.alias_branch(&name, &alias)?;
    }

    store.to_writer(fs::File::create(&store_path)?)?;

    Ok(())
}
