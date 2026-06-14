use crate::{BranchRef, git, store};
use anyhow::Result;
use clap::Parser;
use git2::Repository;
use std::fs;

#[derive(Debug, Parser)]
pub struct Args {
    name: String,
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

    if name == current_branch {
        return Err(anyhow::anyhow!("cannot delete current branch"));
    }

    // TODO(kellegous): We probably shouldn't be able to delete a
    // branch that is listed as a parent.

    if let Some(branch) = store.get_branch_mut(&name) {
        let name = branch.name().to_string();
        git::delete_branch(&repo, &name)?;
        store.delete_branch(&name)?;
        store.to_writer(fs::File::create(&store_path)?)?;
        Ok(())
    } else {
        git::delete_branch(&repo, &name)?;
        Ok(())
    }
}
