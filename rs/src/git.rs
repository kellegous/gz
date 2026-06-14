use anyhow::Result;
use git2::{BranchType, Repository, build};

pub fn checkout(repo: &Repository, name: &str) -> Result<()> {
    let branch = repo.find_branch(name, BranchType::Local)?;
    let commit = branch.get().peel_to_commit()?;
    let mut checkout = build::CheckoutBuilder::new();
    repo.checkout_tree(commit.as_object(), Some(&mut checkout))?;
    repo.set_head(branch.get().name()?)?;
    Ok(())
}
