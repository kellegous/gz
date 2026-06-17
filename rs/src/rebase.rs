use crate::store;
use anyhow::Result;
use clap::Parser;
use std::str::FromStr;

#[derive(Debug, Parser)]
pub struct Args {
    #[arg(long = "at-root", default_value = "nothing", value_parser = AtRoot::from_str)]
    at_root: AtRoot,
}

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

    let (mut _store, _store_path) = store::load_store_git_root(&git_root)?;

    // Sketch this out.
    // 1. We will build a path to the root.
    // 2. We will start at the root, and apply the relevant strategy (fetch or leave as is)
    // 3. We will then move up the stack and we'll do the following at each level:
    //    a. Count the number of commits between HEAD and the parent.
    //    b. Check if the parent HEAD matches our parent ref
    //    c. if not, we will do a rebase --onto parent HEAD~N
    //    d. if there is a merge conflict, we will stop and do what?

    // rebase --continue and --abort should probably be available.

    Ok(())
}
