use crate::{
    model::{Branch, Parent, Sha},
    store::{self, Store},
};
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

    let store_path = git_root.join(".git/gtg.json");
    let mut store = if store_path.exists() {
        Store::from_reader(fs::File::open(&store_path)?)?
    } else {
        Store::default()
    };

    let repo = Repository::open(git_root)?;
    let head = repo.head()?;
    let parent = store.get_branch(head.shorthand()?);

    let commit = head.peel_to_commit()?;
    repo.branch(&args.name, &commit, false)?;

    let prefix = parent.and_then(|b| b.prefix()).map(|p| p.to_string());

    let new_branch = Branch::new(
        args.name,
        None,
        Parent::new(
            head.shorthand()?.to_owned(),
            Sha::from(head.target().ok_or_else(|| anyhow::anyhow!("no target"))?),
        ),
        prefix,
    );

    store.add_branch(new_branch)?;

    store.to_writer(fs::File::create(store_path)?)?;

    Ok(())
}
