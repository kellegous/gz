use anyhow::Result;
use chrono::{DateTime, Utc};
use git2::Oid;
use serde::de::Error as DeError;
use serde::{Deserialize, Deserializer, Serialize, Serializer};
use std::fmt;
use std::str::FromStr;

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Branch {
    name: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    description: Option<String>,
    parent: Parent,
    #[serde(skip_serializing_if = "Option::is_none")]
    prefix: Option<String>,
    created_at: DateTime<Utc>,
    updated_at: DateTime<Utc>,
    last_accessed_at: DateTime<Utc>,
}

impl Branch {
    pub fn new(
        name: String,
        description: Option<String>,
        parent: Parent,
        prefix: Option<String>,
    ) -> Self {
        Self {
            name,
            description,
            parent,
            prefix,
            created_at: Utc::now(),
            updated_at: Utc::now(),
            last_accessed_at: Utc::now(),
        }
    }

    pub fn prefix(&self) -> Option<&str> {
        self.prefix.as_deref()
    }

    pub fn name(&self) -> &str {
        &self.name
    }

    pub fn parent(&self) -> &Parent {
        &self.parent
    }

    pub fn description(&self) -> Option<&str> {
        self.description.as_deref()
    }

    pub fn created_at(&self) -> DateTime<Utc> {
        self.created_at
    }

    pub fn updated_at(&self) -> DateTime<Utc> {
        self.updated_at
    }

    pub fn last_accessed_at(&self) -> DateTime<Utc> {
        self.last_accessed_at
    }

    pub fn update_last_accessed_at(&mut self, now: DateTime<Utc>) {
        self.last_accessed_at = now;
    }

    pub fn update(
        &mut self,
        description: Option<String>,
        parent: Parent,
        prefix: Option<String>,
        now: DateTime<Utc>,
    ) {
        self.description = description;
        self.parent = parent;
        self.prefix = prefix;
        self.updated_at = now;
        self.last_accessed_at = now;
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Parent {
    name: String,
    sha: Sha,
}

impl Parent {
    pub fn new(name: String, sha: Sha) -> Self {
        Self { name, sha }
    }

    pub fn name(&self) -> &str {
        &self.name
    }

    pub fn sha(&self) -> &Sha {
        &self.sha
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub struct Sha(Oid);

impl Serialize for Sha {
    fn serialize<S: Serializer>(&self, serializer: S) -> Result<S::Ok, S::Error> {
        serializer.collect_str(self)
    }
}

impl<'de> Deserialize<'de> for Sha {
    fn deserialize<D: Deserializer<'de>>(deserializer: D) -> Result<Self, D::Error> {
        let s = String::deserialize(deserializer)?;
        Oid::from_str(&s).map(Sha).map_err(DeError::custom)
    }
}

impl From<Oid> for Sha {
    fn from(oid: Oid) -> Self {
        Sha(oid)
    }
}

impl FromStr for Sha {
    type Err = anyhow::Error;

    fn from_str(s: &str) -> Result<Self> {
        Ok(Sha(Oid::from_str(s)?))
    }
}

impl fmt::Display for Sha {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl Sha {
    pub fn oid(&self) -> Oid {
        self.0
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn sha_serializes_as_lowercase_hex_string() {
        let oid = Oid::from_str("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef").unwrap();
        let sha = Sha::from(oid);

        let json = serde_json::to_string(&sha).unwrap();

        assert_eq!(json, "\"deadbeefdeadbeefdeadbeefdeadbeefdeadbeef\"");
    }

    #[test]
    fn sha_deserializes_from_hex_string() {
        let sha: Sha =
            serde_json::from_str("\"0123456789abcdef0123456789abcdef01234567\"").unwrap();

        assert_eq!(
            sha.oid(),
            Oid::from_str("0123456789abcdef0123456789abcdef01234567").unwrap()
        );
    }

    #[test]
    fn sha_roundtrips_through_json() {
        let oid = Oid::from_str("abcdef0123456789abcdef0123456789abcdef01").unwrap();
        let sha = Sha::from(oid);

        let json = serde_json::to_string(&sha).unwrap();
        let restored: Sha = serde_json::from_str(&json).unwrap();

        assert_eq!(restored.oid(), oid);
    }
}
