use crate::store;
use anyhow::Result;
use clap::Parser;
use git2::Repository;

#[derive(Debug, Parser)]
pub struct Args {}

pub fn run(args: Args) -> Result<()> {
    let git_root =
        store::find_git_root()?.ok_or_else(|| anyhow::anyhow!("not in a git repository"))?;

    let (mut store, store_path) = store::load_store_git_root(&git_root)?;

    let repo = Repository::open(&git_root)?;

    Ok(())
}
