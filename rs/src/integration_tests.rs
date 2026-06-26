use anyhow::Result;
use chrono::{DateTime, NaiveDate, Utc};
use std::{env, fs, path::PathBuf, process::Command, time::Duration};
use tempfile::TempDir;

#[test]
fn create_simple() -> Result<()> {
    let repo = GitRepo::init("simple")?;
    repo.run_in(|| {
        let dir = env::current_dir()?;
        println!("dir: {:?}", dir);
        Ok(())
    })?;
    // 2. run create with a name
    // 3. verify that we're chained
    Ok(())
}

struct GitRepo {
    dir: TempDir,
    time: DateTime<Utc>,
}

impl GitRepo {
    fn init(title: &str) -> Result<Self> {
        let dir = tempfile::tempdir()?;
        let time = NaiveDate::from_ymd_opt(2026, 6, 24)
            .ok_or_else(|| anyhow::anyhow!("invalid date"))?
            .and_hms_opt(16, 20, 0)
            .ok_or_else(|| anyhow::anyhow!("invalid time"))?
            .and_utc();
        let repo = Self { dir, time };
        repo.run("init", &[])?;
        fs::write(repo.path().join("README.md"), format!("# {title}\n"))?;
        repo.run("add", &["README.md"])?;
        repo.run("commit", &["-m", "initial commit"])?;
        Ok(repo)
    }

    #[allow(dead_code)]
    fn inc_time(&mut self) {
        self.time += Duration::from_mins(1);
    }

    fn path(&self) -> PathBuf {
        self.dir.path().to_path_buf()
    }

    fn run_in<V, F>(&self, f: F) -> Result<V>
    where
        F: FnOnce() -> Result<V>,
    {
        let previous = std::env::current_dir()?;
        std::env::set_current_dir(&self.dir)?;
        let result = f();
        std::env::set_current_dir(previous)?;
        result
    }

    fn run(&self, cmd: &str, args: &[&str]) -> Result<()> {
        let args: Vec<&str> = std::iter::once(cmd).chain(args.iter().copied()).collect();
        let status = Command::new("git")
            .args(&args)
            .current_dir(&self.dir)
            .env("GIT_AUTHOR_DATE", self.time.to_rfc3339())
            .env("GIT_COMMITTER_DATE", self.time.to_rfc3339())
            .status()?;
        if !status.success() {
            anyhow::bail!("command {cmd} {args:?} failed");
        }
        Ok(())
    }
}
