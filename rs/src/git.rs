use anyhow::Result;
use git2::{BranchType, Repository, build};
use std::process::Command;

pub fn checkout(repo: &Repository, name: &str) -> Result<()> {
    let branch = repo.find_branch(name, BranchType::Local)?;
    let commit = branch.get().peel_to_commit()?;
    let mut checkout = build::CheckoutBuilder::new();
    repo.checkout_tree(commit.as_object(), Some(&mut checkout))?;
    repo.set_head(branch.get().name()?)?;
    Ok(())
}

pub fn delete_branch(repo: &Repository, name: &str) -> Result<()> {
    let mut branch = repo.find_branch(name, BranchType::Local)?;
    branch.delete()?;
    Ok(())
}

pub fn pull(repo: &Repository, name: &str) -> Result<()> {
    // uses git directly mainly to piggyback on any credentials helpers
    // that the user might have set up.
    let workdir = repo
        .workdir()
        .ok_or_else(|| anyhow::anyhow!("bare repository"))?;
    checkout(repo, name)?;
    let status = Command::new("git")
        .current_dir(workdir)
        .args(["pull", "origin", name])
        .status()?;
    if !status.success() {
        anyhow::bail!("git pull origin {name} failed");
    }
    Ok(())
}
