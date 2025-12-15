use std::fmt;
use std::error::Error;

#[derive(Debug)]
pub enum LangError {
    UnsupportedExtension(String),
    UnsupportedLanguage(String),
}

impl fmt::Display for LangError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            LangError::UnsupportedExtension(ext) => write!(f, "unsupported extension: {}", ext),
            LangError::UnsupportedLanguage(lang) => write!(f, "unsupported language: {}", lang),
        }
    }
}

impl Error for LangError {}