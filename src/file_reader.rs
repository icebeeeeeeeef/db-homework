use std::fs;
use std::path::Path;

use walkdir::WalkDir;

use std::collections::HashMap;
use std::collections::HashSet;

use regex::Regex;
use std::error::Error as StdError;

use crate::config::Config;
use crate::consts;
use crate::{FileStats, Langs, LanguageStats};
use serde::Serialize;

// 用于序列化输出的结构体（JSON 输出）
#[derive(Serialize)]
struct LanguageStatsOutput {
    lang: String,
    files: usize,
    lines: usize,
    blanks: usize,
    comments: usize,
    functions: usize,
    function_lines: Vec<usize>,
    stats: Option<StatsInfo>,
    function_stats: Option<FunctionStatsInfo>,
}

#[derive(Serialize)]
struct StatsInfo {
    min: usize,
    max: usize,
    median: usize,
    mean: f64,
}

#[derive(Serialize)]
struct FunctionStatsInfo {
    min: usize,
    max: usize,
    median: usize,
    mean: f64,
}

pub struct FileReader {
    pub config: Config,
    pub files: Vec<FileStats>,
    pub langs: HashMap<String, LanguageStats>, // 语言名 -> 语言统计
    pub visited: HashSet<String>,          // 已访问的文件路径
}

impl FileReader {
    pub fn new(config: Config) -> Self {
        FileReader {
            config,
            files: vec![],
            langs: HashMap::new(),
            visited: HashSet::new(),
        }
    }

    pub fn is_blank(&self, line: &str) -> bool {
        line.trim().is_empty()
    }

    pub fn is_binary<P: AsRef<Path>>(&self, path: P) -> Result<bool, Box<dyn StdError>> {
        use std::fs::File;
        use std::io::Read;

        let mut file = File::open(path)?;
        let mut buffer = [0u8; 512];

        let bytes = file.read(&mut buffer)?;
        let non_text = buffer
            .iter()
            .filter(|&&b| b < 0x09 || (b > 0x0D && b < 0x20))
            .count();

        Ok(non_text * 10 > bytes)
    }
    pub fn is_comment(&self, line: &str, lang: &Langs, in_block_comment: &mut bool) -> bool {
        let trimmed = line.trim();
        let lang_def = consts::lang_to_lang_def(lang).unwrap();

        if let Some((start, end)) = lang_def.block_comment {
            if trimmed.starts_with(start) && !trimmed.ends_with(end) {
                *in_block_comment = true;
                return true;
            }
            if trimmed.ends_with(end) && *in_block_comment {
                *in_block_comment = false;
                return true;
            }
            if *in_block_comment {
                return true;
            }
        }

        if let Some(start) = lang_def.line_comment {
            if line.starts_with(start) {
                return true;
            }
        }

        if let Some(start) = lang_def.doc_comment {
            if line.starts_with(start) {
                return true;
            }
        }

        false
    }

    pub fn is_function(&self, line: &str, lang: &Langs) -> bool {
/**        let lang_def = consts::lang_to_lang_def(lang).unwrap();
        for pattern in lang_def.function_patterns {
            if let Ok(regex) = regex::Regex::new(pattern) {
                if regex.is_match(line) {
                    return true;
                }
            }
        }
*/
        false
    }

    pub fn read_file<P: AsRef<Path>>(&self, path: P) -> Result<FileStats, Box<dyn StdError>> {
        use std::fs::File;
        use std::io::{BufRead, BufReader};

        let p = path.as_ref().display().to_string();
        let ext = path
            .as_ref()
            .extension()
            .and_then(|s| s.to_str())
            .unwrap_or("")
            .to_lowercase();
        let e = Langs::from_string(&ext);

        let ignore_comments = self.config.ignore_comments;
        let ignore_blanks = self.config.ignore_blanks;

        let file = File::open(path.as_ref())?;
        let reader = BufReader::new(file);

        let mut line_count: usize = 0;
        let mut blank_count: usize = 0;
        let mut comment_count: usize = 0;
        let mut in_block_comment: bool = false;
        let mut functions: usize = 0;
        let mut function_lines: Vec<usize> = Vec::new();
        let mut current_function_lines: usize = 0;
        let mut in_function: bool = false;
        let mut brace_count: i32 = 0;

        // 特殊处理 Python 缩进
        let is_python = e == Langs::Python;
        let mut indent_level: i32 = -1;

        for line in reader.lines() {
            // 即使不是 UTF-8，lines() 也会返回 Err，但不会 panic
            if line.is_ok() {
                line_count += 1;

                let line = line.unwrap();
                let trimmed = line.trim();

                if !ignore_blanks {
                    if self.is_blank(&line) {
                        blank_count += 1;
                    }
                }
                if !ignore_comments {
                    if self.is_comment(&line, &e, &mut in_block_comment) {
                        comment_count += 1;
                    }
                }

                // 函数统计
                if self.is_function(&line, &e) {
                    functions += 1;
                    in_function = true;
                    current_function_lines = 1; // 包含函数定义行

                    if is_python {
                        // 计算 Python 缩进级别
                        let indent = line.len() - line.trim_start().len();
                        indent_level = indent as i32;
                    } else {
                        // 对于其他语言，初始化大括号计数
                        brace_count = 0;
                        if trimmed.contains('{') {
                            brace_count += 1;
                        }
                        if trimmed.contains('}') {
                            brace_count -= 1;
                        }
                    }
                } else if in_function {
                    current_function_lines += 1;

                    if is_python {
                        // Python 缩进判断
                        let current_indent = line.len() - line.trim_start().len();
                        // 如果遇到空行或注释，继续统计
                        if self.is_blank(&line) || self.is_comment(&line, &e, &mut false) {
                            // 不改变状态
                        } else if current_indent <= indent_level.try_into().unwrap()
                            && !trimmed.is_empty()
                        {
                            // 缩进级别减小，函数结束
                            in_function = false;
                            function_lines.push(current_function_lines - 1); // 不包含当前行
                            current_function_lines = 0;
                            indent_level = -1;
                        }
                    } else {
                        // 其他语言大括号匹配
                        for c in trimmed.chars() {
                            if c == '{' {
                                brace_count += 1;
                            } else if c == '}' {
                                brace_count -= 1;
                            }
                        }
                        if brace_count == 0 && !trimmed.is_empty() {
                            // 大括号匹配完成，函数结束
                            in_function = false;
                            function_lines.push(current_function_lines);
                            current_function_lines = 0;
                        }
                    }
                }
            }
        }

        // 处理文件末尾未闭合的函数
        if in_function && current_function_lines > 0 {
            function_lines.push(current_function_lines);
        }

        Ok(FileStats {
            path: p,
            lang: e,
            lines: line_count,
            blanks: blank_count,
            comments: comment_count,
            functions,
            function_lines,
        })
    }

    pub fn read_dir<P: AsRef<Path>>(&mut self, path: P) -> Result<usize, Box<dyn StdError>> {
        let mut total = 0;
        let exts = self.config.types.clone();

        for entry in WalkDir::new(path.as_ref())
            .into_iter()
            .filter_map(Result::ok)
            .filter(|e| e.file_type().is_file())
        {
            // 归一化路径，跳过已经访问过的同一物理文件（处理符号链接 / 重复路径导致的重复统计）
            let canonical = fs::canonicalize(entry.path()).ok();
            let real_path_str = canonical
                .as_ref()
                .map(|p| p.to_string_lossy().to_string())
                .unwrap_or_else(|| entry.path().display().to_string());
            if self.visited.contains(&real_path_str) {
                continue;
            }
            self.visited.insert(real_path_str);

            if let Some(ext) = entry.path().extension().and_then(|s| s.to_str()) {
                let ext_lower = ext.to_lowercase();
                if !exts.contains(&ext_lower) {
                    continue;
                }
            } else {
                continue;
            }

            let file_stats = self.read_file(entry.path())?;
            let lang_name = entry
                .path()
                .extension()
                .and_then(|s| s.to_str())
                .unwrap_or("unknown")
                .to_string();
            let lang_name = lang_name.to_lowercase();
            let lang_stats = self
                .langs
                .entry(lang_name.clone())
                .or_insert(LanguageStats {
                    name: lang_name.clone(),
                    files: 0,
                    lines: 0,
                    blanks: 0,
                    comments: 0,
                    functions: 0,
                    function_lines: Vec::new(),
                });
            lang_stats.files += 1;
            lang_stats.lines += file_stats.lines;
            lang_stats.blanks += file_stats.blanks;
            lang_stats.comments += file_stats.comments;
            lang_stats.functions += file_stats.functions;
            lang_stats
                .function_lines
                .extend_from_slice(&file_stats.function_lines);

            total += file_stats.lines;
            self.files.push(file_stats);
        }

        Ok(total)
    }

    pub fn read<P: AsRef<Path>>(&mut self, path: P) -> Result<usize, Box<dyn StdError>> {
        let metadata = fs::metadata(path.as_ref()).map_err(|e| Box::new(e) as Box<dyn StdError>)?;

        if metadata.is_dir() {
            self.read_dir(path)
        } else if metadata.is_file() {
            let fs = self.read_file(path)?;
            Ok(fs.lines)
        } else {
            Err(From::from("Path is neither a file nor a directory"))
        }
    }

    pub fn run(&mut self) -> Result<(), Box<dyn StdError>> {
        if self.config.help {
            Config::print_usage();
            return Ok(());
        }

        let paths = self.config.paths.clone();
        let mut total = 0;

        for path in paths {
            total += self.read_dir(&path)?;
        }

        // 如果指定了输出，则写入文件；否则维持原有行为打印到 stdout
        if let Some(ref out) = self.config.output {
            write_output(&self.config, &self.langs, total, out)?;
        } else {
            print(&self.config, &self.langs, total);
        }

        Ok(())
    }
}

#[inline]
// 计算统计信息
fn calculate_stats(data: &Vec<usize>) -> Option<(usize, usize, usize, f64)> {
    if data.is_empty() {
        return None;
    }

    let mut sorted = data.clone();
    sorted.sort();

    let min = sorted[0];
    let max = sorted[sorted.len() - 1];

    let median = if sorted.len() % 2 == 0 {
        (sorted[sorted.len() / 2 - 1] + sorted[sorted.len() / 2]) / 2
    } else {
        sorted[sorted.len() / 2]
    };

    let mean: f64 = sorted.iter().sum::<usize>() as f64 / sorted.len() as f64;

    Some((min, max, median, mean))
}

pub fn print_divider() {
    println!("{}", "=".repeat(80));
}

pub fn write_output(
    config: &Config,
    langs: &HashMap<String, LanguageStats>,
    total: usize,
    output_param: &str,
) -> Result<(), Box<dyn StdError>> {
    use std::fs::File;
    use std::path::Path;

    // 决定输出文件路径：如果传入值看起来像文件名则直接使用，否则当作格式名，生成默认文件名
    let path_str = if output_param.contains('.')
        || output_param.contains(std::path::MAIN_SEPARATOR)
        || output_param.contains('/')
        || output_param.contains('\\')
    {
        output_param.to_string()
    } else {
        format!("toukei_output.{}", output_param)
    };

    let path = Path::new(&path_str);
    let ext = path
        .extension()
        .and_then(|s| s.to_str())
        .unwrap_or("")
        .to_lowercase();

    // 将数据转换为可序列化的结构
    let mut lang_list: Vec<(&String, &LanguageStats)> = langs.iter().collect();
    lang_list.sort_by(|a, b| b.1.lines.cmp(&a.1.lines));

    if ext == "json" {
        let mut out_vec: Vec<LanguageStatsOutput> = Vec::new();
        for (name, stats) in lang_list {
            let mut l = LanguageStatsOutput {
                lang: name.clone(),
                files: stats.files,
                lines: stats.lines,
                blanks: stats.blanks,
                comments: stats.comments,
                functions: stats.functions,
                function_lines: stats.function_lines.clone(),
                stats: None,
                function_stats: None,
            };

            if config.show_stats && stats.files > 0 {
                let min = stats.lines / stats.files;
                let max = stats.lines;
                let median = stats.lines / stats.files;
                let mean = stats.lines as f64 / stats.files as f64;
                l.stats = Some(StatsInfo {
                    min,
                    max,
                    median,
                    mean,
                });
            }

            if config.show_function_stats && !stats.function_lines.is_empty() {
                if let Some((min, max, median, mean)) = calculate_stats(&stats.function_lines) {
                    l.function_stats = Some(FunctionStatsInfo {
                        min,
                        max,
                        median,
                        mean,
                    });
                }
            }

            out_vec.push(l);
        }

        let file = File::create(path)?;
        serde_json::to_writer_pretty(file, &out_vec)?;
        println!("Wrote JSON output to {}", path.display());
        return Ok(());
    }

    if ext == "csv" {
        let mut wtr = csv::Writer::from_path(path)?;
        // header
        wtr.write_record(&[
            "Language",
            "Files",
            "Lines",
            "Blanks",
            "Comments",
            "Functions",
            "FuncMin",
            "FuncMax",
            "FuncMedian",
            "FuncMean",
        ])?;

        for (name, stats) in lang_list {
            let (fmin, fmax, fmedian, fmean) = if !stats.function_lines.is_empty() {
                if let Some((min, max, median, mean)) = calculate_stats(&stats.function_lines) {
                    (
                        min.to_string(),
                        max.to_string(),
                        median.to_string(),
                        format!("{:.2}", mean),
                    )
                } else {
                    (
                        "".to_string(),
                        "".to_string(),
                        "".to_string(),
                        "".to_string(),
                    )
                }
            } else {
                (
                    "".to_string(),
                    "".to_string(),
                    "".to_string(),
                    "".to_string(),
                )
            };

            wtr.write_record(&[
                name,
                &stats.files.to_string(),
                &stats.lines.to_string(),
                &stats.blanks.to_string(),
                &stats.comments.to_string(),
                &stats.functions.to_string(),
                &fmin,
                &fmax,
                &fmedian,
                &fmean,
            ])?;
        }
        wtr.flush()?;
        println!("Wrote CSV output to {}", path.display());
        return Ok(());
    }

    if ext == "xlsx" || ext == "xls" {
        // 使用 umya_spreadsheet 创建简单的表格
        let mut book = umya_spreadsheet::new_file();
        let sheet_name = "Sheet1";
        // 获取工作表（new_file 已默认创建 Sheet1，这里直接获取即可）
        let sheet = book.get_sheet_by_name_mut(sheet_name).unwrap();

        // 写 header
        let headers = vec![
            "Language",
            "Files",
            "Lines",
            "Blanks",
            "Comments",
            "Functions",
            "FuncMin",
            "FuncMax",
            "FuncMedian",
            "FuncMean",
        ];
        for (col_idx, h) in headers.iter().enumerate() {
            // 明确转为 u32 类型，列和行均从 1 开始
            let col = (col_idx + 1) as u32;
            let row = 1u32;
            sheet.get_cell_mut((col, row)).set_value(*h);
        }

        for (row_idx, (name, stats)) in lang_list.iter().enumerate() {
            let row = (row_idx + 2) as u32; // 数据从第 2 行开始（表头在第 1 行），转为 u32
            sheet.get_cell_mut((1u32, row)).set_value(name.as_str());
            sheet
                .get_cell_mut((2u32, row))
                .set_value(&stats.files.to_string());
            sheet
                .get_cell_mut((3u32, row))
                .set_value(&stats.lines.to_string());
            sheet
                .get_cell_mut((4u32, row))
                .set_value(&stats.blanks.to_string());
            sheet
                .get_cell_mut((5u32, row))
                .set_value(&stats.comments.to_string());
            sheet
                .get_cell_mut((6u32, row))
                .set_value(&stats.functions.to_string());

            if let Some((min, max, median, mean)) = calculate_stats(&stats.function_lines) {
                sheet.get_cell_mut((7u32, row)).set_value(&min.to_string());
                sheet.get_cell_mut((8u32, row)).set_value(&max.to_string());
                sheet
                    .get_cell_mut((9u32, row))
                    .set_value(&median.to_string());
                sheet
                    .get_cell_mut((10u32, row))
                    .set_value(&format!("{:.2}", mean));
            }
        }

        // 写入文件
        umya_spreadsheet::writer::xlsx::write(&book, path)
            .map_err(|e| -> Box<dyn std::error::Error> { Box::new(e) })?;
        println!("Wrote XLSX output to {}", path.display());
        return Ok(());
    }

    Err(From::from(format!("Unsupported output format: {}", ext)))
}

pub fn print(config: &Config, langs: &HashMap<String, LanguageStats>, total: usize) {
    let mut lang_list: Vec<&LanguageStats> = langs.values().collect();
    lang_list.sort_by(|a, b| b.lines.cmp(&a.lines));

    print_divider();
    println!(
        "{:<10} {:<10} {:<10} {:<10} {:<10} {:<10}",
        "Language", "Files", "Lines", "Blanks", "Comments", "Functions"
    );
    print_divider();

    for lang in &lang_list {
        println!(
            "{:<10} {:<10} {:<10} {:<10} {:<10} {:<10}",
            lang.name, lang.files, lang.lines, lang.blanks, lang.comments, lang.functions
        );
    }

    print_divider();
    println!("Total files: {}", total);

    // 显示行数统计信息
    if config.show_stats {
        print_divider();
        println!("行数统计信息:");
        print_divider();
        println!(
            "{:<10} {:<10} {:<10} {:<10} {:<10}",
            "Language", "最小值", "最大值", "中位数", "平均值"
        );
        print_divider();

        for lang in &lang_list {
            if lang.files > 0 {
                let avg_lines = lang.lines as f64 / lang.files as f64;
                println!(
                    "{:<10} {:<10} {:<10} {:<10.2} {:<10.2}",
                    lang.name,
                    lang.lines / lang.files,
                    lang.lines,
                    lang.lines / lang.files,
                    avg_lines
                );
            }
        }
    }

    // 显示函数统计信息
    if config.show_function_stats {
        print_divider();
        println!("函数行数统计信息:");
        print_divider();
        println!(
            "{:<10} {:<10} {:<10} {:<10} {:<10}",
            "Language", "最小值", "最大值", "中位数", "平均值"
        );
        print_divider();

        for lang in &lang_list {
            if !lang.function_lines.is_empty() {
                if let Some((min, max, median, mean)) = calculate_stats(&lang.function_lines) {
                    println!(
                        "{:<10} {:<10} {:<10} {:<10} {:<10.2}",
                        lang.name, min, max, median, mean
                    );
                }
            }
        }
    }
}
