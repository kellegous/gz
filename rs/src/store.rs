use crate::model::Branch;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::{
    collections::HashMap,
    fs, io,
    path::{Path, PathBuf},
};

pub const DEFAULT_STORE_PATH: &str = ".git/gdg.json";

#[derive(Debug, Serialize, Deserialize, Default)]
pub struct Store {
    #[serde(default, skip_serializing_if = "HashMap::is_empty")]
    branches: HashMap<String, Branch>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    aliases: Vec<(String, String)>,
}

impl Store {
    pub fn from_reader<R: io::Read>(reader: R) -> Result<Self> {
        Ok(serde_json::from_reader(reader)?)
    }

    pub fn to_writer<W: io::Write>(&self, writer: W) -> Result<()> {
        serde_json::to_writer(writer, &self)?;
        Ok(())
    }

    fn resolve_branch_name(&self, name: &str) -> Option<String> {
        if self.branches.contains_key(name) {
            Some(name.to_string())
        } else {
            self.aliases
                .iter()
                .find(|(_, alias)| alias == name)
                .map(|(branch_name, _)| branch_name.clone())
        }
    }

    pub fn get_branch(&self, name: &str) -> Option<&Branch> {
        let branch_name = self.resolve_branch_name(name)?;
        self.branches.get(&branch_name)
    }

    pub fn get_branch_mut(&mut self, name: &str) -> Option<&mut Branch> {
        let branch_name = self.resolve_branch_name(name)?;
        self.branches.get_mut(&branch_name)
    }

    pub fn add_branch(&mut self, branch: Branch) -> Result<()> {
        if self.branches.contains_key(branch.name()) {
            return Err(anyhow::anyhow!(
                r#"Branch "{}" already exists"#,
                branch.name()
            ));
        }
        self.branches.insert(branch.name().to_string(), branch);
        Ok(())
    }

    pub fn alias_branch(&mut self, name: &str, alias: &str) -> Result<()> {
        if !self.branches.contains_key(name) {
            return Err(anyhow::anyhow!(r#"Branch "{}" not found"#, name));
        }
        if self.aliases.iter().any(|(_, a)| a == alias) {
            return Err(anyhow::anyhow!(r#"Alias "{}" already exists"#, alias));
        }
        self.aliases.push((name.to_string(), alias.to_string()));
        Ok(())
    }

    pub fn unalias_branch(&mut self, name: &str) -> Result<()> {
        let n = self.aliases.len();
        self.aliases.retain(|(_, a)| a != name);
        if self.aliases.len() == n {
            return Err(anyhow::anyhow!(r#"Alias "{}" not found"#, name));
        }
        Ok(())
    }

    pub fn children_of(&self, name: &str) -> impl Iterator<Item = &Branch> {
        self.branches
            .values()
            .filter(move |b| b.parent().name() == name)
    }

    pub fn delete_branch(&mut self, name: &str) -> Result<()> {
        let name = self
            .resolve_branch_name(name)
            .ok_or_else(|| anyhow::anyhow!(r#"Branch "{}" not found"#, name))?;

        if let Some(child) = self.children_of(&name).next() {
            return Err(anyhow::anyhow!(
                r#"Branch "{}" is a parent of "{}""#,
                name,
                child.name(),
            ));
        }

        self.branches.remove(&name);
        self.aliases.retain(|(n, _)| n != &name);
        Ok(())
    }

    pub fn branches(&self) -> impl Iterator<Item = &Branch> {
        self.branches.values()
    }

    pub fn aliases(&self) -> impl Iterator<Item = &(String, String)> {
        self.aliases.iter()
    }
}

pub fn find_git_root() -> Result<Option<PathBuf>> {
    let mut path = std::env::current_dir()?;
    loop {
        if path.join(".git").is_dir() {
            return Ok(Some(path));
        }
        match path.parent() {
            Some(parent) => path = parent.to_path_buf(),
            None => return Ok(None),
        }
    }
}

pub fn load_store_git_root<P: AsRef<Path>>(root: P) -> Result<(Store, PathBuf)> {
    let store_path = root.as_ref().join(DEFAULT_STORE_PATH);
    if store_path.exists() {
        Ok((
            Store::from_reader(fs::File::open(&store_path)?)?,
            store_path,
        ))
    } else {
        Ok((Store::default(), store_path))
    }
}
