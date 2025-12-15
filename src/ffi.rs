use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use std::ptr;
use std::collections::HashMap;

use serde::{Serialize, Deserialize};
use serde_json;

use crate::file_reader::FileReader;
use crate::config::Config;

// 用于JSON序列化的LanguageStats结构体
#[derive(Serialize)]
struct LanguageStatsOutput {
    lang: String,
    files: usize,
    lines: usize,
    blanks: usize,
    comments: usize,
    functions: usize,
    function_lines: Vec<usize>,
    // 统计信息
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

// 输入配置结构体
#[derive(Deserialize)]
struct ToukeiConfig {
    path: String,
    types: Option<Vec<String>>,
    ignore_blanks: Option<bool>,
    ignore_comments: Option<bool>,
    ignore_files: Option<Vec<String>>,
    show_stats: Option<bool>,
    show_function_stats: Option<bool>,
    output: Option<String>,
}

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

// 主要接口 - 接受JSON配置
#[unsafe(no_mangle)]
pub extern "C" fn toukei_count_with_config(config_json: *const c_char) -> *mut c_char {
    if config_json.is_null() { return ptr::null_mut(); }
    
    let cstr = unsafe { CStr::from_ptr(config_json) };
    let config_str = match cstr.to_str() {
        Ok(s) => s,
        Err(_) => return ptr::null_mut(),
    };

    // 解析JSON配置
    let toukei_config: Result<ToukeiConfig, _> = serde_json::from_str(config_str);
    let toukei_config = match toukei_config {
        Ok(config) => config,
        Err(_) => return ptr::null_mut(),
    };

    // 创建Config实例
    let mut config = Config::new();
    config.paths.push(toukei_config.path.clone());
    
    if let Some(types) = toukei_config.types {
        config.types = types;
    }
    
    if let Some(ignore_blanks) = toukei_config.ignore_blanks {
        config.ignore_blanks = ignore_blanks;
    }
    
    if let Some(ignore_comments) = toukei_config.ignore_comments {
        config.ignore_comments = ignore_comments;
    }
    
    if let Some(ignore_files) = toukei_config.ignore_files {
        config.ignore_files = ignore_files;
    }
    
    if let Some(show_stats) = toukei_config.show_stats {
        config.show_stats = show_stats;
    }
    
    if let Some(show_function_stats) = toukei_config.show_function_stats {
        config.show_function_stats = show_function_stats;
    }

    if let Some(output) = toukei_config.output {
        config.output = Some(output);
    }

    // 执行统计
    let mut reader = FileReader::new(config.clone());
    
    let mut total_files = 0;
    for path in &config.paths {
        if let Err(_) = reader.read_dir(path) {
            return ptr::null_mut();
        }
        total_files += 1;
    }

    // 准备输出数据
    let mut output: Vec<LanguageStatsOutput> = Vec::new();
    
    for (lang_name, stats) in &reader.langs {
        let mut lang_output = LanguageStatsOutput {
            lang: lang_name.clone(),
            files: stats.files,
            lines: stats.lines,
            blanks: stats.blanks,
            comments: stats.comments,
            functions: stats.functions,
            function_lines: stats.function_lines.clone(),
            stats: None,
            function_stats: None,
        };

        // 添加行数统计信息
        if config.show_stats && stats.files > 0 {
            let min = stats.lines / stats.files;
            let max = stats.lines;
            let median = stats.lines / stats.files;
            let mean = stats.lines as f64 / stats.files as f64;
            
            lang_output.stats = Some(StatsInfo {
                min,
                max,
                median,
                mean,
            });
        }

        // 添加函数统计信息
        if config.show_function_stats && !stats.function_lines.is_empty() {
            if let Some((min, max, median, mean)) = calculate_stats(&stats.function_lines) {
                lang_output.function_stats = Some(FunctionStatsInfo {
                    min,
                    max,
                    median,
                    mean,
                });
            }
        }

        output.push(lang_output);
    }

    // 如果配置中指定了输出文件，则尝试写出并返回写入结果（包含输出路径）
    if let Some(ref out) = config.output {
        // 调用 file_writer 写文件
        match crate::file_reader::write_output(&config, &reader.langs, total_files, out) {
            Ok(()) => {
                let result_json = serde_json::json!({"status": "ok", "output": out});
                if let Ok(s) = CString::new(result_json.to_string()) {
                    return s.into_raw();
                } else {
                    return ptr::null_mut();
                }
            }
            Err(_) => {
                return ptr::null_mut();
            }
        }
    }

    // 否则返回 JSON 序列化的结果数组（原有行为）
    let json_result = serde_json::to_string(&output);
    let json = match json_result {
        Ok(json_str) => json_str,
        Err(_) => return ptr::null_mut(),
    };

    match CString::new(json) {
        Ok(s) => s.into_raw(),
        Err(_) => ptr::null_mut(),
    }
}

// 简化的接口 - 仅传入路径
#[unsafe(no_mangle)]
pub extern "C" fn toukei_count_path(path: *const c_char) -> *mut c_char {
    if path.is_null() { return ptr::null_mut(); }
    
    let cstr = unsafe { CStr::from_ptr(path) };
    let path_str = match cstr.to_str() {
        Ok(s) => s,
        Err(_) => return ptr::null_mut(),
    };

    // 创建简单的JSON配置
    let config_json = format!(r#"{{"path":"{}"}}"#, path_str);
    
    // 调用主函数
    match CString::new(config_json) {
        Ok(config_cstr) => {
            toukei_count_with_config(config_cstr.as_ptr())
        },
        Err(_) => ptr::null_mut(),
    }
}

#[unsafe(no_mangle)]
pub extern "C" fn toukei_free_string(s: *mut c_char) {
    if s.is_null() { return; }
    unsafe { let _ = CString::from_raw(s); } // drop -> 释放内存
}