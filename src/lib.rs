pub mod file_reader;
pub mod config;
pub mod utils;
pub mod ffi;
pub mod core;
pub mod consts;
pub mod tests;
 
pub use std::collections::HashMap;
pub use crate::config::Config;
pub use crate::file_reader::FileReader;

#[derive(Debug, Copy, Clone, PartialEq)]
pub enum Langs {
    Rust,
    Python,
    JavaScript,
    Go,
    C,
    Cpp,
    H,
    Hpp,
    CSharp,
    Java,
    Html,
    Css,
    Scss,
    Json,
    Xml,
    Toml,
    Yaml,
    Unknown,
}

impl Langs {
    pub fn from_string(s: &str) -> Langs {
        match s.to_lowercase().as_str() {
            "rs" => Langs::Rust,
            "py" => Langs::Python,
            "js" | "jsx" | "ts" | "tsx" => Langs::JavaScript,
            "go" => Langs::Go,
            "c" => Langs::C,
            "cpp" => Langs::Cpp,
            "h" => Langs::H,
            "hpp" => Langs::Hpp,
            "cs" => Langs::CSharp,
            "java" => Langs::Java,
            "html" => Langs::Html,
            "css" => Langs::Css,
            "scss" => Langs::Scss,
            "json" => Langs::Json,
            "xml" => Langs::Xml,
            "toml" => Langs::Toml,
            "yml" | "yaml" => Langs::Yaml,
            _ => Langs::Unknown,
        }
    }
    pub fn to_string(&self) -> &str {
        match self {
            Langs::Rust => "Rust",
            Langs::Python => "Python",
            Langs::JavaScript => "JavaScript",
            Langs::Go => "Go",
            Langs::C => "C",
            Langs::Cpp => "C++",
            Langs::H => "C",
            Langs::Hpp => "C++",
            Langs::CSharp => "C#",
            Langs::Java => "Java",
            Langs::Html => "HTML",
            Langs::Css => "CSS",
            Langs::Scss => "SCSS",
            Langs::Json => "JSON",
            Langs::Xml => "XML",
            Langs::Toml => "TOML",
            Langs::Yaml => "YAML",
            Langs::Unknown => "Unknown",
        }
    }
}

pub struct FileStats {
    pub path: String,
    pub lang: Langs,
    pub lines: usize,
    pub blanks: usize,
    pub comments: usize,
    pub functions: usize,
    pub function_lines: Vec<usize>,
}

// 语言统计
#[derive(Clone, Debug)]
pub struct LanguageStats {
    name: String,
    files: usize,
    lines: usize,
    blanks: usize,
    comments: usize,
    functions: usize,
    function_lines: Vec<usize>,
}

impl Default for LanguageStats {
    fn default() -> Self {
        LanguageStats {
            name: String::new(),
            files: 0,
            lines: 0,
            blanks: 0,
            comments: 0,
            functions: 0,
            function_lines: Vec::new(),
        }
    }
}
