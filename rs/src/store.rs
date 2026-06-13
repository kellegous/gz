use crate::model::Branch;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::{collections::HashMap, io};
#[derive(Debug, Serialize, Deserialize, Default)]
pub struct Store {
    branches: HashMap<String, Branch>,
    aliases: HashMap<String, String>,
}

impl Store {
    pub fn from_reader<R: io::Read>(reader: R) -> Result<Self> {
        Ok(serde_json::from_reader(reader)?)
    }

    pub fn to_writer<W: io::Write>(&self, writer: W) -> Result<()> {
        serde_json::to_writer(writer, &self)?;
        Ok(())
    }
}
