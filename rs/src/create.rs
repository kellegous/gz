use crate::{
    STORE_FILE,
    model::{Branch, Parent, Sha},
    store::{self, Store},
};
use anyhow::Result;
use clap::{ArgAction, Parser};
use git2::Repository;
use std::fs;

#[derive(Debug, Parser)]
pub struct Args {
    name: String,
    #[arg(short, long)]
    parent: Option<String>,
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

    let store_path = git_root.join(STORE_FILE);
    let mut store = if store_path.exists() {
        Store::from_reader(fs::File::open(&store_path)?)?
    } else {
        Store::default()
    };

    let repo = Repository::open(git_root)?;
    let parent_ref = match &args.parent {
        Some(parent) => repo.find_reference(parent)?,
        None => repo.head()?,
    };

    let parent = store.get_branch(parent_ref.shorthand()?);

    // Create the new branch in git
    let commit = parent_ref.peel_to_commit()?;
    let git_branch = repo.branch(&args.name, &commit, false)?;

    // Set the new branch as the HEAD
    repo.set_head(git_branch.get().name()?)?;
    let mut checkout = git2::build::CheckoutBuilder::new();
    repo.checkout_head(Some(&mut checkout))?;

    let prefix = match &args.prefix {
        Some(prefix) => Some(prefix.clone()),
        None => parent.and_then(|b| b.prefix()).map(|p| p.to_string()),
    };

    store.add_branch(Branch::new(
        args.name.clone(),
        args.description.clone(),
        Parent::new(
            parent_ref.shorthand()?.to_owned(),
            Sha::from(
                parent_ref
                    .target()
                    .ok_or_else(|| anyhow::anyhow!("no target"))?,
            ),
        ),
        prefix,
    ))?;

    for alias in args.aliases {
        store.alias_branch(args.name.as_str(), &alias)?;
    }

    store.to_writer(fs::File::create(store_path)?)?;

    Ok(())
}
