pub struct LangDef {
    pub name: &'static str,
    pub extensions: &'static [&'static str],
    pub line_comment: Option<&'static str>,
    pub block_comment: Option<(&'static str, &'static str)>,
    pub doc_comment: Option<&'static str>,
    pub function_patterns: &'static [&'static str],
    pub class_pattern: Option<&'static str>,
}