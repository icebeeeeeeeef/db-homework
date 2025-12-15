use std::collections::HashMap;
use crate::core::lang_def::LangDef;
use crate::Langs;
use crate::core::lang_error::LangError;

pub static LANG_NAMES: &[&str] = &[
    "c","cpp","h","hpp","rs","cs","go","py","java","js","ts","jsx","tsx","html","css","scss","json","xml","yml","yaml","toml"
];

pub static LANG_DEFS: &[LangDef] = &[
    LangDef {
        name: "Rust",
        extensions: &["rs"],
        line_comment: Some("//"),
        block_comment: Some(("/*", "*/")),
        doc_comment: Some("///"),
        function_patterns: &[
            r"^\s*fn\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*(->\s*[^\{]*)?\s*\{",
            r"^\s*pub\s+fn\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*(->\s*[^\{]*)?\s*\{",
        ],
        class_pattern: Some(r"^\s*(pub\s+)?struct\s+[a-zA-Z0-9_]+"),
    },
    LangDef {
        name: "C++",
        extensions: &["cpp", "cc", "cxx", "hpp", "h"],
        line_comment: Some("//"),
        block_comment: Some(("/*", "*/")),
        doc_comment: Some("///"),
        function_patterns: &[
            r"^\s*[a-zA-Z0-9_\*&<>\[\]]+\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*(const\s*)?\s*\{",
            r"^\s*void\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*(const\s*)?\s*\{",
        ],
        class_pattern: Some(r"^\s*class\s+[a-zA-Z0-9_]+"),
    },
    LangDef {
        name: "Python",
        extensions: &["py"],
        line_comment: Some("#"),
        block_comment: Some(("\"\"\"", "\"\"\"")),
        doc_comment: Some("'''"),
        function_patterns: &[
            r"^\s*def\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*:",
            r"^\s*async\s+def\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*:",
        ],
        class_pattern: Some(r"^\s*class\s+[a-zA-Z0-9_]+")
    },
    LangDef {
        name: "JavaScript",
        extensions: &["js"],
        line_comment: Some("//"),
        block_comment: Some(("/*", "*/")),
        doc_comment: Some("/**"),
        function_patterns: &[
            r"^\s*function\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*\{",
            r"^\s*const\s+[a-zA-Z0-9_]+\s*=\s*\([^)]*\)\s*=>",
            r"^\s*let\s+[a-zA-Z0-9_]+\s*=\s*\([^)]*\)\s*=>",
            r"^\s*var\s+[a-zA-Z0-9_]+\s*=\s*\([^)]*\)\s*=>",
            r"^\s*[a-zA-Z0-9_]+\s*=\s*\([^)]*\)\s*=>",
        ],
        class_pattern: Some(r"^\s*(export\s+)?class\s+[a-zA-Z0-9_]+"),
    },
    LangDef {
        name: "HTML",
        extensions: &["html"],
        line_comment: Some("<!--"),
        block_comment: Some(("<!--", "-->")),
        function_patterns: &[],
        class_pattern: None,
        doc_comment: None,
    },
    LangDef {
        name: "Go",
        extensions: &["go"],
        line_comment: Some("//"),
        block_comment: Some(("/*", "*/")),
        doc_comment: Some("///"),
        function_patterns: &[
            r"^\s*func\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*\{",
            r"^\s*func\s+\([^)]*\)\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*\{",
        ],
        class_pattern: Some(r"^\s*type\s+[a-zA-Z0-9_]+\s+(struct|interface)\s*\{"),
    },
    LangDef {
        name: "C",
        extensions: &["c"],
        line_comment: Some("//"),
        block_comment: Some(("/*", "*/")),
        doc_comment: Some("///"),
        function_patterns: &[
            r"^\s*[a-zA-Z0-9_\*&\[\]]+\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*\{",
            r"^\s*void\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*\{",
        ],
        class_pattern: None,
    },
    LangDef {
        name: "C#",
        extensions: &["cs"],
        line_comment: Some("//"),
        block_comment: Some(("/*", "*/")),
        doc_comment: Some("///"),
        function_patterns: &[
            r"^\s*(public\s+|private\s+|protected\s+)?(static\s+)?[a-zA-Z0-9_<>,\[\]]+\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*\{",
            r"^\s*(public\s+|private\s+|protected\s+)?(static\s+)?void\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*\{",
        ],
        class_pattern: Some(r"^\s*(public\s+|private\s+|protected\s+)?class\s+[a-zA-Z0-9_]+"),
    },
    LangDef {
        name: "Java",
        extensions: &["java"],
        line_comment: Some("//"),
        block_comment: Some(("/*", "*/")),
        doc_comment: Some("/**"),
        function_patterns: &[
            r"^\s*(public\s+|private\s+|protected\s+)?(static\s+)?[a-zA-Z0-9_<>,\[\]]+\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*(\{|\{[^;]*\{)",
            r"^\s*(public\s+|private\s+|protected\s+)?(static\s+)?void\s+[a-zA-Z0-9_]+\s*\([^)]*\)\s*(\{|\{[^;]*\{)",
        ],
        class_pattern: Some(r"^\s*(public\s+|private\s+|protected\s+)?class\s+[a-zA-Z0-9_]+"),
    },
    LangDef {
        name: "CSS",
        extensions: &["css"],
        line_comment: Some("//"),
        block_comment: Some(("/*", "*/")),
        doc_comment: Some("///"),
        function_patterns: &[],
        class_pattern: None,
    },
    LangDef {
        name: "XML",
        extensions: &["xml"],
        line_comment: Some("<!--"),
        block_comment: Some(("<!--", "-->")),
        doc_comment: None,
        class_pattern: None,
        function_patterns: &[],
    },
    LangDef {
        name: "SCSS",
        extensions: &["scss"],
        line_comment: Some("//"),
        block_comment: Some(("/*", "*/")),
        doc_comment: Some("///"),
        class_pattern: None,
        function_patterns: &[],
    },
    LangDef {
        name: "JSON",
        extensions: &["json"],
        line_comment: None,
        block_comment: None,
        doc_comment: None,
        class_pattern: None,
        function_patterns: &[],
    },
    LangDef {
        name: "Markdown",
        extensions: &["md"],
        line_comment: Some("<!--"),
        block_comment: Some(("<!--", "-->")),
        doc_comment: None,
        class_pattern: None,
        function_patterns: &[],
    },
    LangDef {
        name: "TOML",
        extensions: &["toml"],
        line_comment: Some("#"),
        block_comment: None,
        doc_comment: None,
        class_pattern: None,
        function_patterns: &[],
    },
    LangDef {
        name: "YAML",
        extensions: &["yaml"],
        line_comment: Some("#"),
        block_comment: None,
        doc_comment: None,
        class_pattern: None,
        function_patterns: &[],
    },

];

/// 根据扩展名查找语言定义，找不到时返回错误
pub fn get_lang_def(ext: &str) -> Result<&'static LangDef, LangError> {
    LANG_DEFS
        .iter()
        .find(|def| def.extensions.contains(&ext))
        .ok_or_else(|| LangError::UnsupportedExtension(ext.to_string()))
}

/// 根据 Langs 查找语言定义，找不到时返回错误
pub fn lang_to_lang_def(lang: &Langs) -> Result<&'static LangDef, LangError> {
    let name = lang.to_string().to_string();
    LANG_DEFS
        .iter()
        .find(|def| def.name == name)
        .ok_or_else(|| LangError::UnsupportedLanguage(name))
}