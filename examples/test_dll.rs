use std::ffi::{CStr, CString};
use std::os::raw::c_char;

// 声明DLL函数
#[link(name = "toukei")]
unsafe extern "C" {
    fn toukei_count_path(path: *const c_char) -> *mut c_char;
    fn toukei_count_with_config(config_json: *const c_char) -> *mut c_char;
    fn toukei_free_string(s: *mut c_char);
}

fn main() {
    // 测试简单路径统计
    println!("=== 测试简单路径统计 ===");
    let path = CString::new("./src").unwrap();
    let result_ptr = unsafe { toukei_count_path(path.as_ptr()) };
    
    if !result_ptr.is_null() {
        let result = unsafe { CStr::from_ptr(result_ptr).to_string_lossy() };
        println!("结果: {}", result);
        unsafe { toukei_free_string(result_ptr) };
    } else {
        println!("调用失败");
    }
    
    // 测试JSON配置
    println!("\n=== 测试JSON配置 ===");
    let config_json = r#"{
        "path": "./src",
        "types": ["rs"],
        "show_stats": true,
        "show_function_stats": true
    }"#;
    
    let config = CString::new(config_json).unwrap();
    let result_ptr = unsafe { toukei_count_with_config(config.as_ptr()) };
    
    if !result_ptr.is_null() {
        let result = unsafe { CStr::from_ptr(result_ptr).to_string_lossy() };
        println!("结果: {}", result);
        unsafe { toukei_free_string(result_ptr) };
    } else {
        println!("调用失败");
    }
}