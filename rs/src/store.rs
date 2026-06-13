use crate::model::Branch;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::PathBuf;

#[derive(Debug, Serialize, Deserialize)]
pub struct Store {
    path: PathBuf,
    branches: HashMap<String, Branch>,
    aliases: HashMap<String, String>,
}
