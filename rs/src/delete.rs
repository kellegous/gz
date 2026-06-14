use crate::{BranchRef, git, store};
use anyhow::Result;
use clap::Parser;
use git2::Repository;
use std::fs;

#[derive(Debug, Parser)]
pub struct Args {
    names: Vec<String>,
}

pub fn run(args: Args) -> Result<()> {
    let git_root =
        store::find_git_root()?.ok_or_else(|| anyhow::anyhow!("not in a git repository"))?;

    let (mut store, store_path) = store::load_store_git_root(&git_root)?;

    let repo = Repository::open(&git_root)?;
    let head = repo.head()?;
    let current_branch = head.shorthand()?;

    for name in args.names {
        let name = &name.parse::<BranchRef>()?.resolve(&store, current_branch)?;

        if name == current_branch {
            return Err(anyhow::anyhow!("cannot delete current branch"));
        }

        if let Some(branch) = store.get_branch_mut(name) {
            let name = branch.name().to_string();
            store.delete_branch(&name)?;
            git::delete_branch(&repo, &name)?;
            store.to_writer(fs::File::create(&store_path)?)?;
        } else {
            if store.children_of(name).count() > 0 {
                return Err(anyhow::anyhow!(r#"Branch "{}" has children"#, name,));
            }
            git::delete_branch(&repo, name)?;
        }
    }

    Ok(())
}
