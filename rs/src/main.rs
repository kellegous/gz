use anyhow::Result;
use clap::Parser;
use git2::Repository;

#[derive(Debug, Parser)]
struct Args {
    path: String,
}
fn main() -> Result<()> {
    let args = Args::parse();

    let repo = Repository::open(args.path)?;

    let head = repo.head()?;

    println!("HEAD: {}", head.name()?);
    println!("HEAD: {}", head.shorthand()?);
    println!(
        "HEAD: {}",
        head.target().ok_or(anyhow::anyhow!("No target"))?
    );
    Ok(())
}
