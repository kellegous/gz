use crate::{git, model::Branch, store};
use anyhow::Result;
use clap::Parser;
use git2::Repository;
use std::str::FromStr;

#[derive(Debug, Parser)]
pub struct Args {
    #[arg(long = "at-root", default_value = "nothing", value_parser = AtRoot::from_str)]
    at_root: AtRoot,
}

// TODO(kellegous): There are three strategies possible here:
// 1. do nothing to the root and do not update the first branch to the
//    HEAD of root. This is like a pinned chain.
// 2. fetch the root, fast-forward it, rebase the chain. This is what
//    gs rebase does.
// 3. do nothing to the root, but update the first branch to the HEAD of
//    root.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum AtRoot {
    Pull,
    Nothing,
}

impl FromStr for AtRoot {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "pull" => Ok(Self::Pull),
            "nothing" => Ok(Self::Nothing),
            other => Err(format!(
                "invalid value '{other}' for --at-root (expected: pull, nothing)"
            )),
        }
    }
}

pub fn run(args: Args) -> Result<()> {
    println!("rebase: {:?}", args);
    let git_root =
        store::find_git_root()?.ok_or_else(|| anyhow::anyhow!("not in a git repository"))?;

    let (store, _store_path) = store::load_store_git_root(&git_root)?;

    let repo = Repository::open(&git_root)?;
    let head = repo.head()?;

    if let Some(current) = store.get_branch(head.shorthand()?) {
        // TODO(kellegous): write the chain into a recovery file
        let chain = load_chain(&store, current)?;
        println!("{:?}", chain);
    } else {
        update_root(&repo, head.shorthand()?, args.at_root)?;
    }

    // Sketch this out.
    // 1. We will build a path to the root.
    // 2. We will start at the root, and apply the relevant strategy (fetch or leave as is)
    // 3. We will then move up the stack and we'll do the following at each level:
    //    a. Count the number of commits between HEAD and the parent.
    //    b. Check if the parent HEAD matches our parent ref
    //    c. if not, we will do a rebase --onto parent HEAD~N
    //    d. if there is a merge conflict, we will stop and do what?

    // we will likely need to write a state file for the node that is being rebased
    // otherwise we will not know how to continue when a merge conflict happens.

    Ok(())
}

fn update_root(repo: &Repository, name: &str, strategy: AtRoot) -> Result<()> {
    if strategy == AtRoot::Pull {
        git::pull(repo, name)?;
    }
    Ok(())
}

fn load_chain(store: &store::Store, current: &Branch) -> Result<Vec<Branch>> {
    let mut chain = vec![current.clone()];
    let mut current = current;
    while let Some(parent) = store.get_branch(current.parent().name()) {
        chain.push(parent.clone());
        current = parent;
    }
    Ok(chain)
}
