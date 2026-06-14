use crate::model::Branch;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::path::PathBuf;
use std::{collections::HashMap, io};
#[derive(Debug, Serialize, Deserialize, Default)]
pub struct Store {
    #[serde(skip_serializing_if = "HashMap::is_empty")]
    branches: HashMap<String, Branch>,
    #[serde(skip_serializing_if = "Vec::is_empty")]
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

    pub fn get_branch(&self, name: &str) -> Option<&Branch> {
        if let Some(branch) = self.branches.get(name) {
            Some(branch)
        } else {
            self.aliases
                .iter()
                .find(|(_, alias)| alias == name)
                .map(|(_, alias)| alias)
                .and_then(|name| self.branches.get(name))
        }
    }

    pub fn add_branch(&mut self, branch: Branch) -> Result<()> {
        if self.branches.contains_key(branch.name()) {
            return Err(anyhow::anyhow!("Branch already exists"));
        }
        self.branches.insert(branch.name().to_string(), branch);
        Ok(())
    }

    pub fn alias_branch(&mut self, name: &str, alias: &str) -> Result<()> {
        if !self.branches.contains_key(name) {
            return Err(anyhow::anyhow!("Branch not found"));
        }
        if self.aliases.iter().any(|(_, a)| a == alias) {
            return Err(anyhow::anyhow!("Alias already exists"));
        }
        self.aliases.push((name.to_string(), alias.to_string()));
        Ok(())
    }

    pub fn unalias_branch(&mut self, name: &str) -> Result<()> {
        if !self.aliases.iter().any(|(_, a)| a == name) {
            return Err(anyhow::anyhow!("Alias not found"));
        }
        self.aliases.retain(|(_, a)| a != name);
        Ok(())
    }

    pub fn delete_branch(&mut self, name: &str) -> Result<()> {
        if !self.branches.contains_key(name) {
            return Err(anyhow::anyhow!("Branch not found"));
        }
        self.branches.remove(name);
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
