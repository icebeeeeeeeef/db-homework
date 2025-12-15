use std::error::Error;

use std::collections::HashSet;
use std::fmt::Display;
use crate::consts::*;

#[derive(Debug)]
#[derive(Default)]
#[derive(Clone)]
pub struct Config {
    pub paths: Vec<String>,
    pub types: Vec<String>,
    pub ignore_blanks: bool,
    pub ignore_comments: bool,
    pub ignore_files: Vec<String>,
    pub show_stats: bool,
    pub show_function_stats: bool,

    // 输出文件路径或格式（例如: "out.json" 或 "json"）
    pub output: Option<String>,

    pub help: bool,
}

impl Display for Config {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Config {{ paths: {:?}, types: {:?}, ignore_blanks: {}, ignore_comments: {}, ignore_files: {:?} }}", 
            self.paths, self.types, self.ignore_blanks, self.ignore_comments, self.ignore_files)
    }
    
}

impl Config {

    pub fn new() -> Config {
        let exts = LANG_NAMES.to_vec();
        let exts: Vec<String> = exts.iter().map(|s| s.to_string()).collect();
        Config {
            paths: vec![],
            types: exts,
            ignore_blanks: false,
            ignore_comments: false,
            ignore_files: vec![],
            show_stats: false,
            show_function_stats: false,
            help: false,
            output: None,
        }    
    }

    pub fn build(args: &[String]) -> Result<Config, Box<dyn Error>> {
        if args.len() < 2 {
            return Err(From::from("Not enough arguments"));
        }

        let mut paths: Vec<String> = vec![];
        let mut types: Vec<String> = vec![];
        let mut ignore_blanks = false;
        let mut ignore_comments = false;
        let mut ignore_files: Vec<String> = vec![];
        let mut show_stats = false;
        let mut show_function_stats = false;
        let mut help = false;
        let mut output: Option<String> = None;

        let mut i: usize = 1;

        while i < args.len() {
            if args[i].starts_with("-") {
                match args[i].as_str() {
                    "-p" | "--path" => {
                        while i + 1 < args.len() && !args[i+1].starts_with("-") {
                            paths.push(args[i+1].clone());
                            i += 1;
                        }
                    }
                    "-t" | "--type" => {
                        while i + 1 < args.len() && !args[i+1].starts_with("-") {
                            types.push(args[i+1].clone());
                            i += 1;
                        }
                    }
                    "-h" | "--help" => {
                        help = true;
                        i += 1;
                    }
                    "-i" | "--ignore-blanks" => {
                        ignore_blanks = true;
                        i+=1;
                    }
                    "--ignore-comments" => {
                        ignore_comments = true;
                        i+=1;
                    }
                    "--show-stats" => {
                        show_stats = true;
                        i+=1;
                    }
                    "--show-function-stats" => {
                        show_function_stats = true;
                        i+=1;
                    }
                    "--output" => {
                        // 取一个参数作为输出目标（可以是格式名或文件名）
                        if i + 1 < args.len() && !args[i+1].starts_with("-") {
                            output = Some(args[i+1].clone());
                            i += 2;
                        } else {
                            return Err(From::from("--output requires a value"));
                        }
                    }
                    _ => {
                        return Err(From::from("Invalid argument"));
                    }
                }
            }
            else {
                paths.push(args[i].clone());
                i+=1;
            }
        }

        if types.is_empty() {
            types = LANG_NAMES.iter().map(|s| s.to_string()).collect();
        }


        Ok(Config { paths, types, ignore_blanks, ignore_comments, ignore_files, show_stats, show_function_stats, help, output })
    }

    pub fn print_usage() {
        println!("Usage:");
        println!("    {} [options] [file]", env!("CARGO_PKG_NAME"));
        println!("Options:");
        println!("    -h, --help                   Print this help message");
        println!("    -p, --path <path>            Add a path to search for files");
        println!("    -t, --type <type>            Only count lines of the given type");
        println!("    -i, --ignore-blanks          Ignore blank lines");
        println!("    --ignore-comments            Ignore comment lines");
        println!("    --ignore-files <file>        Ignore files matching the given pattern");
        println!("    --show-stats                 Show statistics (mean, median, max, min)");
        println!("    --show-function-stats        Show function statistics (mean, median, max, min)");
        println!("    --output <format|file>       Write output to file. Specify a format (json/csv/xlsx) or a filename. If omitted, print to stdout.");
    }
}