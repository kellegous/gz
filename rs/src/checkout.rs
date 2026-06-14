use crate::store;
use anyhow::Result;
use chrono::Utc;
use clap::Parser;
use git2::{BranchType, Repository};
use std::fs;

#[derive(Debug, Parser)]
pub struct Args {
    name: String,
}

pub fn run(args: Args) -> Result<()> {
    let git_root =
        store::find_git_root()?.ok_or_else(|| anyhow::anyhow!("not in a git repository"))?;

    let (mut store, store_path) = store::load_store_git_root(&git_root)?;

    let branch = store
        .get_branch_mut(args.name.as_str())
        .ok_or(anyhow::anyhow!("branch not found"))?;

    let repo = Repository::open(&git_root)?;
    let git_branch = repo.find_branch(branch.name(), BranchType::Local)?;
    repo.set_head(git_branch.get().name()?)?;
    let mut checkout = git2::build::CheckoutBuilder::new();
    repo.checkout_head(Some(&mut checkout))?;

    branch.update_last_accessed_at(Utc::now());
    store.to_writer(fs::File::create(&store_path)?)?;
    Ok(())
}
