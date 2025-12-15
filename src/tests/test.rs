#[cfg(test)]
mod tests {
    use std::fs::File;
    use std::io::Write;
    use walkdir::WalkDir;
    use std::collections::HashSet;

    #[test]
    fn test_walkdir_and_write_files() {
        // 指定要遍历的目录
        let test_dir = "F:\\QQ\\Downloads\\testcase5\\testcase5";
        println!("{}", test_dir);
        // 创建输出文件
        let mut output_file = File::create("test_files.txt").expect("Failed to create output file");
        let mut visited = HashSet::<String>::new();
        // 遍历目录并写入文件名
        for entry in WalkDir::new(test_dir)
            .into_iter()
            .filter_map(|e| e.ok())
            .filter(|e| e.file_type().is_file()) 
        {
            if let Some(path) = entry.path().to_str() {
                if visited.contains(path) {
                    println!("Duplicate file: {}", path);
                    continue;
                }
                visited.insert(path.to_string());
                writeln!(output_file, "{}", path).expect("Failed to write to file");
            }
        }
    }
}
