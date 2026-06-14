use crate::{
    git,
    model::{Branch, Parent, Sha},
    store,
};
use anyhow::Result;
use clap::{ArgAction, Parser};
use git2::Repository;
use std::fs;

#[derive(Debug, Parser)]
pub struct Args {
    name: String,
    #[arg(short, long)]
    from: Option<String>,
    #[arg(short, long)]
    description: Option<String>,
    #[arg(short, long)]
    prefix: Option<String>,
    #[arg(short = 'a', long = "alias", action = ArgAction::Append)]
    aliases: Vec<String>,
}

pub fn run(args: Args) -> Result<()> {
    let git_root =
        store::find_git_root()?.ok_or_else(|| anyhow::anyhow!("not in a git repository"))?;

    let (mut store, store_path) = store::load_store_git_root(&git_root)?;

    let repo = Repository::open(git_root)?;
    let from_ref = match &args.from {
        Some(from) => repo.find_reference(from)?,
        None => repo.head()?,
    };

    let parent = store.get_branch(from_ref.shorthand()?);

    let prefix = match &args.prefix {
        Some(prefix) => Some(prefix.clone()),
        None => parent.and_then(|b| b.prefix()).map(|p| p.to_string()),
    };

    let new_name = match &prefix {
        Some(prefix) => format!("{}/{}", prefix, args.name),
        None => args.name.clone(),
    };

    // Create the new branch in git
    let commit = from_ref.peel_to_commit()?;
    let git_branch = repo.branch(&new_name, &commit, false)?;
    git::checkout(&repo, git_branch.get().name()?)?;

    store.add_branch(Branch::new(
        new_name.clone(),
        args.description.clone(),
        Parent::new(
            from_ref.shorthand()?.to_owned(),
            Sha::from(
                from_ref
                    .target()
                    .ok_or_else(|| anyhow::anyhow!("no target"))?,
            ),
        ),
        prefix,
    ))?;

    for alias in args.aliases {
        store.alias_branch(new_name.as_str(), &alias)?;
    }

    store.to_writer(fs::File::create(store_path)?)?;

    Ok(())
}
