pub mod alias;
pub mod checkout;
pub mod create;
pub mod delete;
pub mod git;
pub mod model;
pub mod store;
pub mod unalias;

use crate::store::Store;
use anyhow::Result;
use std::str;

pub enum BranchRef {
    Root,
    Parent,
    Last,
    Name(String),
}

impl BranchRef {
    pub fn resolve(&self, store: &Store, current: &str) -> Result<String> {
        match self {
            Self::Root => {
                let mut name = current.to_string();
                loop {
                    match store.get_branch(&name) {
                        Some(branch) => name = branch.parent().name().to_string(),
                        None => return Ok(name),
                    }
                }
            }
            Self::Parent => match store.get_branch(current) {
                Some(branch) => Ok(branch.parent().name().to_string()),
                None => Err(anyhow::anyhow!("{} has no parent", current)),
            },
            Self::Last => {
                match store
                    .branches()
                    .filter(|b| b.name() != current)
                    .max_by_key(|b| b.last_accessed_at())
                {
                    Some(branch) => Ok(branch.name().to_string()),
                    None => Err(anyhow::anyhow!("No branches found")),
                }
            }
            Self::Name(name) => Ok(name.to_string()),
        }
    }
}

impl str::FromStr for BranchRef {
    type Err = anyhow::Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "^root" => Ok(Self::Root),
            "^parent" => Ok(Self::Parent),
            "^last" => Ok(Self::Last),
            _ => {
                if !s.starts_with("^") {
                    Ok(Self::Name(s.to_string()))
                } else {
                    Err(anyhow::anyhow!("Invalid branch reference: {}", s))
                }
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::str::FromStr;

    #[test]
    fn branch_ref_from_str_parses_special_refs() {
        assert!(matches!(
            BranchRef::from_str("^root").unwrap(),
            BranchRef::Root
        ));
        assert!(matches!(
            BranchRef::from_str("^parent").unwrap(),
            BranchRef::Parent
        ));
        assert!(matches!(
            BranchRef::from_str("^last").unwrap(),
            BranchRef::Last
        ));
    }

    #[test]
    fn branch_ref_from_str_parses_name() {
        assert!(matches!(
            BranchRef::from_str("main").unwrap(),
            BranchRef::Name(name) if name == "main"
        ));
    }

    #[test]
    fn branch_ref_from_str_rejects_invalid_special_ref() {
        assert!(BranchRef::from_str("^foo").is_err());
    }
}
