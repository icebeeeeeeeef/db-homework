#include <iostream>
#include <fstream>
#include <filesystem>
#include <string>
#include <map>
#include <vector>
#include <algorithm>
#include <iomanip>
#include <sstream>
#include <numeric>
#include <set>
#include <cctype>
#include <regex>

namespace fs = std::filesystem;

static std::string trim(const std::string& str) {
    size_t start = str.find_first_not_of(" \t\r\n");
    if (start == std::string::npos) return "";
    size_t end = str.find_last_not_of(" \t\r\n");
    return str.substr(start, end - start + 1);
}

static std::string toLower(std::string s) {
    std::transform(s.begin(), s.end(), s.begin(), [](unsigned char c) {
        return static_cast<char>(std::tolower(c));
    });
    return s;
}

static int indentLevel(const std::string& line) {
    int count = 0;
    for (char ch : line) {
        if (ch == ' ') count++;
        else if (ch == '\t') count += 4;
        else break;
    }
    return count;
}

// 编程语言定义
struct Language {
    std::string name;
    std::vector<std::string> extensions;
    std::vector<std::string> singleLineComments;
    std::vector<std::string> multiLineCommentStart;
    std::vector<std::string> multiLineCommentEnd;
};

class CodeStats {
private:
    std::map<std::string, Language> languages;
    std::map<std::string, int> fileCounts;
    std::map<std::string, int> lineCounts;
    std::map<std::string, int> codeLines;
    std::map<std::string, int> commentLines;
    std::map<std::string, int> blankLines;
    std::map<std::string, std::vector<int>> functionLengths;
    bool collectFunctionStats = false;
    std::set<std::string> languageFilter;
    std::map<std::string, std::vector<std::string>> functionPatternSources;
    std::map<std::string, std::vector<std::regex>> compiledFunctionPatterns;
    std::set<std::string> indentBasedFunctionLanguages;
    
    void initializeLanguages() {
        // C/C++
        languages["cpp"] = {
            "C/C++", 
            {".cpp", ".c", ".cc", ".cxx", ".c++", ".h", ".hpp", ".hxx"},
            {"//"},
            {"/*"},
            {"*/"}
        };
        
        // Java
        languages["java"] = {
            "Java",
            {".java"},
            {"//"},
            {"/*"},
            {"*/"}
        };
        
        // Python
        languages["python"] = {
            "Python",
            {".py", ".pyw"},
            {"#"},
            {"\"\"\"", "'''"},
            {"\"\"\"", "'''"}
        };
        
        // JavaScript
        languages["javascript"] = {
            "JavaScript",
            {".js", ".jsx", ".mjs"},
            {"//"},
            {"/*"},
            {"*/"}
        };
        
        // TypeScript
        languages["typescript"] = {
            "TypeScript",
            {".ts", ".tsx"},
            {"//"},
            {"/*"},
            {"*/"}
        };
        
        // Go
        languages["go"] = {
            "Go",
            {".go"},
            {"//"},
            {"/*"},
            {"*/"}
        };
        
        // Rust
        languages["rust"] = {
            "Rust",
            {".rs"},
            {"//"},
            {"/*"},
            {"*/"}
        };
        
        // C#
        languages["csharp"] = {
            "C#",
            {".cs"},
            {"//"},
            {"/*"},
            {"*/"}
        };
        
        // PHP
        languages["php"] = {
            "PHP",
            {".php", ".phtml"},
            {"//", "#"},
            {"/*"},
            {"*/"}
        };
        
        // Ruby
        languages["ruby"] = {
            "Ruby",
            {".rb"},
            {"#"},
            {"=begin"},
            {"=end"}
        };
        
        // Swift
        languages["swift"] = {
            "Swift",
            {".swift"},
            {"//"},
            {"/*"},
            {"*/"}
        };
        
        // Kotlin
        languages["kotlin"] = {
            "Kotlin",
            {".kt", ".kts"},
            {"//"},
            {"/*"},
            {"*/"}
        };
        
        // HTML
        languages["html"] = {
            "HTML",
            {".html", ".htm"},
            {},
            {"<!--"},
            {"-->"}
        };
        
        // CSS
        languages["css"] = {
            "CSS",
            {".css", ".scss", ".sass", ".less"},
            {},
            {"/*"},
            {"*/"}
        };
        
        // Shell/Bash
        languages["shell"] = {
            "Shell/Bash",
            {".sh", ".bash", ".zsh", ".fish"},
            {"#"},
            {},
            {}
        };
        
        // SQL
        languages["sql"] = {
            "SQL",
            {".sql"},
            {"--"},
            {"/*"},
            {"*/"}
        };
        
        // YAML
        languages["yaml"] = {
            "YAML",
            {".yml", ".yaml"},
            {"#"},
            {},
            {}
        };
        
        // JSON
        languages["json"] = {
            "JSON",
            {".json"},
            {},
            {},
            {}
        };
        
        // XML
        languages["xml"] = {
            "XML",
            {".xml"},
            {},
            {"<!--"},
            {"-->"}
        };
    }
    
    std::string getLanguageByExtension(const std::string& extension) {
        for (const auto& [key, lang] : languages) {
            for (const auto& ext : lang.extensions) {
                if (ext == extension) {
                    return key;
                }
            }
        }
        return "";
    }
    
    bool isCommentLine(const std::string& line, const Language& lang) {
        if (lang.singleLineComments.empty() && lang.multiLineCommentStart.empty()) {
            return false;
        }
        
        std::string trimmed = line;
        trimmed.erase(0, trimmed.find_first_not_of(" \t"));
        
        // 检查单行注释
        for (const auto& comment : lang.singleLineComments) {
            if (trimmed.find(comment) == 0) {
                return true;
            }
        }
        
        return false;
    }
    
    bool isBlankLine(const std::string& line) {
        return line.find_first_not_of(" \t\n\r") == std::string::npos;
    }
    
    bool isLanguageEnabled(const std::string& langKey) const {
        if (languageFilter.empty()) return true;
        std::string lowerKey = toLower(langKey);
        auto it = languages.find(langKey);
        std::string lowerName = it != languages.end() ? toLower(it->second.name) : "";
        for (const auto& token : languageFilter) {
            if (token.empty()) continue;
            if (lowerKey.find(token) != std::string::npos) return true;
            if (!lowerName.empty() && lowerName.find(token) != std::string::npos) return true;
        }
        return false;
    }
    
    void initializeFunctionPatterns() {
        functionPatternSources["cpp"] = {
            R"(^\s*[A-Za-z0-9_\*&<>\[\]]+\s+[A-Za-z0-9_]+\s*\([^)]*\)\s*(const\s*)?\s*\{)",
            R"(^\s*void\s+[A-Za-z0-9_]+\s*\([^)]*\)\s*(const\s*)?\s*\{)"
        };
        functionPatternSources["python"] = {
            R"(^\s*def\s+[A-Za-z0-9_]+\s*\([^)]*\)\s*:)",
            R"(^\s*async\s+def\s+[A-Za-z0-9_]+\s*\([^)]*\)\s*:)"
        };
        indentBasedFunctionLanguages.insert("python");
    }

    bool supportsFunctionStats(const std::string& langKey) const {
        return collectFunctionStats && functionPatternSources.find(langKey) != functionPatternSources.end();
    }

    const std::vector<std::regex>& getFunctionPatterns(const std::string& langKey) {
        auto it = compiledFunctionPatterns.find(langKey);
        if (it != compiledFunctionPatterns.end()) {
            return it->second;
        }
        std::vector<std::regex> compiled;
        auto src = functionPatternSources.find(langKey);
        if (src != functionPatternSources.end()) {
            for (const auto& pattern : src->second) {
                try {
                    compiled.emplace_back(pattern, std::regex::ECMAScript);
                } catch (const std::regex_error&) {
                    // ignore bad pattern
                }
            }
        }
        auto inserted = compiledFunctionPatterns.emplace(langKey, std::move(compiled));
        return inserted.first->second;
    }

    bool isFunctionSignatureLine(const std::string& line, const std::string& langKey) {
        const auto& patterns = getFunctionPatterns(langKey);
        for (const auto& pattern : patterns) {
            if (std::regex_search(line, pattern)) {
                return true;
            }
        }
        return false;
    }

    void analyzeFile(const fs::path& filePath) {
        std::ifstream file(filePath);
        if (!file.is_open()) {
            return;
        }
        
        std::string extension = filePath.extension().string();
        std::string langKey = getLanguageByExtension(extension);
        
        if (langKey.empty()) {
            return;
        }

        if (!isLanguageEnabled(langKey)) {
            return;
        }
        
        const Language& lang = languages[langKey];
        int totalLines = 0;
        int codeLines = 0;
        int commentLines = 0;
        int blankLines = 0;
        
        std::string line;
        bool inMultiLineComment = false;
        bool trackFunctions = supportsFunctionStats(langKey);
        bool indentBased = indentBasedFunctionLanguages.count(langKey) > 0;
        bool inFunction = false;
        int currentFunctionLines = 0;
        int functionIndent = -1;
        int braceDepth = 0;
        
        while (std::getline(file, line)) {
            totalLines++;
            std::string trimmed = trim(line);
            bool blankLine = isBlankLine(line);
            bool blockStartDetected = false;
            size_t blockStartPos = std::string::npos;
            std::string blockStartToken;
            if (!inMultiLineComment) {
                for (const auto& startComment : lang.multiLineCommentStart) {
                    size_t pos = line.find(startComment);
                    if (pos != std::string::npos) {
                        blockStartDetected = true;
                        blockStartPos = pos;
                        blockStartToken = startComment;
                        break;
                    }
                }
            }
            bool singleLineComment = false;
            if (!inMultiLineComment) {
                singleLineComment = isCommentLine(line, lang);
            }
            bool signatureLine = trackFunctions && isFunctionSignatureLine(line, langKey);
            bool lineIsCommentForFunction = singleLineComment || inMultiLineComment || blockStartDetected;

            if (trackFunctions) {
                if (signatureLine) {
                    inFunction = true;
                    currentFunctionLines = 1;
                    if (indentBased) {
                        functionIndent = indentLevel(line);
                    } else {
                        braceDepth = 0;
                        braceDepth += static_cast<int>(std::count(line.begin(), line.end(), '{'));
                        braceDepth -= static_cast<int>(std::count(line.begin(), line.end(), '}'));
                    }
                } else if (inFunction) {
                    currentFunctionLines++;
                    if (indentBased) {
                        if (!blankLine && !lineIsCommentForFunction) {
                            int indent = indentLevel(line);
                            if (indent <= functionIndent) {
                                inFunction = false;
                                if (currentFunctionLines > 1) {
                                    functionLengths[langKey].push_back(currentFunctionLines - 1);
                                }
                                currentFunctionLines = 0;
                                functionIndent = -1;
                            }
                        }
                    } else {
                        braceDepth += static_cast<int>(std::count(line.begin(), line.end(), '{'));
                        braceDepth -= static_cast<int>(std::count(line.begin(), line.end(), '}'));
                        if (braceDepth == 0 && !trimmed.empty()) {
                            if (currentFunctionLines > 0) {
                                functionLengths[langKey].push_back(currentFunctionLines);
                            }
                            inFunction = false;
                            currentFunctionLines = 0;
                        }
                    }
                }
            }

            if (blankLine) {
                blankLines++;
                continue;
            }
            
            if (inMultiLineComment) {
                commentLines++;
                // 检查多行注释结束
                for (const auto& endComment : lang.multiLineCommentEnd) {
                    if (line.find(endComment) != std::string::npos) {
                        inMultiLineComment = false;
                        break;
                    }
                }
                continue;
            }
            
            if (blockStartDetected) {
                commentLines++;
                bool endsSameLine = false;
                for (const auto& endComment : lang.multiLineCommentEnd) {
                    size_t endPos = line.find(endComment, blockStartPos + blockStartToken.size());
                    if (endPos != std::string::npos) {
                        endsSameLine = true;
                        break;
                    }
                }
                inMultiLineComment = !endsSameLine && !lang.multiLineCommentEnd.empty();
                continue;
            }
            
            if (singleLineComment) {
                commentLines++;
            } else {
                codeLines++;
            }
        }
        
        if (trackFunctions && inFunction && currentFunctionLines > 0) {
            functionLengths[langKey].push_back(currentFunctionLines);
        }
        
        // 更新统计
        fileCounts[langKey]++;
        lineCounts[langKey] += totalLines;
        this->codeLines[langKey] += codeLines;
        this->commentLines[langKey] += commentLines;
        this->blankLines[langKey] += blankLines;

    }

    double computeAverage(const std::vector<int>& values) const {
        if (values.empty()) return 0.0;
        int sum = std::accumulate(values.begin(), values.end(), 0);
        return static_cast<double>(sum) / values.size();
    }

    double computeMedian(std::vector<int> values) const {
        if (values.empty()) return 0.0;
        std::sort(values.begin(), values.end());
        size_t mid = values.size() / 2;
        if (values.size() % 2 == 0) {
            return (values[mid - 1] + values[mid]) / 2.0;
        }
        return values[mid];
    }
    
public:
    CodeStats() {
        initializeLanguages();
        initializeFunctionPatterns();
    }

    void enableFunctionStats(bool enabled) {
        collectFunctionStats = enabled;
    }

    void setLanguageFilter(const std::set<std::string>& filter) {
        languageFilter = filter;
    }
    
    void analyzeDirectory(const std::string& directoryPath) {
        try {
            for (const auto& entry : fs::recursive_directory_iterator(directoryPath)) {
                if (entry.is_regular_file()) {
                    analyzeFile(entry.path());
                }
            }
        } catch (const fs::filesystem_error& e) {
            std::cerr << "Error accessing directory: " << e.what() << std::endl;
        }
    }
    
    void printStatistics() {
        std::cout << "\n=== 代码统计报告 ===\n\n";
        
        // 按语言排序
        std::vector<std::pair<std::string, int>> sortedLanguages;
        for (const auto& [langKey, count] : fileCounts) {
            if (count > 0) {
                sortedLanguages.push_back({langKey, count});
            }
        }
        
        std::sort(sortedLanguages.begin(), sortedLanguages.end(),
            [this](const auto& a, const auto& b) {
                return lineCounts[a.first] > lineCounts[b.first];
            });
        
        // 打印表头
        std::cout << std::left << std::setw(15) << "语言"
                  << std::setw(8) << "文件数"
                  << std::setw(10) << "总行数"
                  << std::setw(10) << "代码行"
                  << std::setw(10) << "注释行"
                  << std::setw(10) << "空行"
                  << std::endl;
        std::cout << std::string(75, '-') << std::endl;
        
        // 打印每种语言的统计
        for (const auto& [langKey, _] : sortedLanguages) {
            const Language& lang = languages[langKey];
            std::cout << std::left << std::setw(15) << lang.name
                      << std::setw(8) << fileCounts[langKey]
                      << std::setw(10) << lineCounts[langKey]
                      << std::setw(10) << codeLines[langKey]
                      << std::setw(10) << commentLines[langKey]
                      << std::setw(10) << blankLines[langKey]
                      << std::endl;
        }
        
        // 打印总计
        int totalFiles = 0, totalLines = 0, totalCodeLines = 0, totalCommentLines = 0, totalBlankLines = 0;
        for (const auto& [langKey, _] : fileCounts) {
            totalFiles += fileCounts[langKey];
            totalLines += lineCounts[langKey];
            totalCodeLines += codeLines[langKey];
            totalCommentLines += commentLines[langKey];
            totalBlankLines += blankLines[langKey];
        }
        
        std::cout << std::string(75, '-') << std::endl;
        std::cout << std::left << std::setw(15) << "总计"
                  << std::setw(8) << totalFiles
                  << std::setw(10) << totalLines
                  << std::setw(10) << totalCodeLines
                  << std::setw(10) << totalCommentLines
                  << std::setw(10) << totalBlankLines
                  << std::endl;
        
        std::cout << "\n=== 详细统计 ===\n";
        for (const auto& [langKey, _] : sortedLanguages) {
            const Language& lang = languages[langKey];
            if (fileCounts[langKey] > 0) {
                std::cout << "\n" << lang.name << ":\n";
                std::cout << "  文件数: " << fileCounts[langKey] << "\n";
                std::cout << "  总行数: " << lineCounts[langKey] << "\n";
                std::cout << "  代码行: " << codeLines[langKey] << " ("
                          << std::fixed << std::setprecision(1)
                          << (lineCounts[langKey] > 0 ? (double)codeLines[langKey] / lineCounts[langKey] * 100 : 0)
                          << "%)\n";
                std::cout << "  注释行: " << commentLines[langKey] << " ("
                          << std::fixed << std::setprecision(1)
                          << (lineCounts[langKey] > 0 ? (double)commentLines[langKey] / lineCounts[langKey] * 100 : 0)
                          << "%)\n";
                std::cout << "  空行: " << blankLines[langKey] << " ("
                          << std::fixed << std::setprecision(1)
                          << (lineCounts[langKey] > 0 ? (double)blankLines[langKey] / lineCounts[langKey] * 100 : 0)
                          << "%)\n";
            }
        }

        if (collectFunctionStats) {
            std::cout << "\n=== 函数长度统计 (行) ===\n";
            for (const auto& entry : functionLengths) {
                if (entry.second.empty()) continue;
                const Language& lang = languages[entry.first];
                std::vector<int> values = entry.second;
                std::cout << "\n" << lang.name << ":\n";
                std::cout << "  函数数量: " << values.size() << "\n";
                std::cout << "  平均: " << std::fixed << std::setprecision(2) << computeAverage(values) << "\n";
                std::cout << "  最大: " << *std::max_element(values.begin(), values.end()) << "\n";
                std::cout << "  最小: " << *std::min_element(values.begin(), values.end()) << "\n";
                std::cout << "  中位数: " << computeMedian(values) << "\n";
            }
        }
    }

    // 以TSV格式输出: 语言\t文件数\t总行数\t代码行\t注释行\t空行
    void printStatisticsTSV() {
        // 收集并排序
        std::vector<std::pair<std::string, int>> sortedLanguages;
        for (const auto& [langKey, count] : fileCounts) {
            if (count > 0) sortedLanguages.push_back({langKey, count});
        }
        std::sort(sortedLanguages.begin(), sortedLanguages.end(),
            [this](const auto& a, const auto& b) { return lineCounts[a.first] > lineCounts[b.first]; });

        int totalFiles = 0, totalLines_ = 0, totalCode = 0, totalComment = 0, totalBlank = 0;
        for (const auto& [langKey, _] : sortedLanguages) {
            const Language& lang = languages[langKey];
            std::cout << lang.name << '\t'
                      << fileCounts[langKey] << '\t'
                      << lineCounts[langKey] << '\t'
                      << codeLines[langKey] << '\t'
                      << commentLines[langKey] << '\t'
                      << blankLines[langKey] << '\n';
            totalFiles += fileCounts[langKey];
            totalLines_ += lineCounts[langKey];
            totalCode += codeLines[langKey];
            totalComment += commentLines[langKey];
            totalBlank += blankLines[langKey];
        }
        // 总计行
        std::cout << "TOTAL\t" << totalFiles << '\t' << totalLines_ << '\t'
                  << totalCode << '\t' << totalComment << '\t' << totalBlank << '\n';

        if (collectFunctionStats) {
            for (const auto& entry : functionLengths) {
                if (entry.second.empty()) continue;
                const Language& lang = languages[entry.first];
                const std::vector<int>& values = entry.second;
                std::cout << "FUNC\t" << lang.name << '\t'
                          << values.size() << '\t'
                          << std::fixed << std::setprecision(2) << computeAverage(values) << '\t'
                          << *std::min_element(values.begin(), values.end()) << '\t'
                          << *std::max_element(values.begin(), values.end()) << '\t'
                          << computeMedian(values) << '\n';
            }
        }
    }
};

int main(int argc, char* argv[]) {
    bool tsv = false;
    bool functions = false;
    std::string directoryPath = ".";
    std::set<std::string> filters;

    for (int i = 1; i < argc; ++i) {
        std::string arg = argv[i];
        if (arg == "--tsv") {
            tsv = true;
        } else if (arg == "--functions") {
            functions = true;
        } else if (arg.rfind("--dir=", 0) == 0) {
            directoryPath = arg.substr(6);
        } else if (arg.rfind("--languages=", 0) == 0) {
            std::stringstream ss(arg.substr(12));
            std::string token;
            while (std::getline(ss, token, ',')) {
                token = trim(token);
                if (!token.empty()) {
                    filters.insert(toLower(token));
                }
            }
        } else if (!arg.empty() && arg[0] != '-') {
            directoryPath = arg;
        } else if (arg == "--help" || arg == "-h") {
            std::cout << "用法: " << argv[0] << " [--tsv] [--functions] [--languages=a,b] [--dir=path] [目录]\n";
            return 0;
        }
    }
    
    if (!fs::exists(directoryPath)) {
        std::cerr << "错误: 目录不存在: " << directoryPath << std::endl;
        return 1;
    }
    
    if (!fs::is_directory(directoryPath)) {
        std::cerr << "错误: 不是目录: " << directoryPath << std::endl;
        return 1;
    }
    
    CodeStats stats;
    stats.enableFunctionStats(functions);
    if (!filters.empty()) {
        stats.setLanguageFilter(filters);
    }
    stats.analyzeDirectory(directoryPath);
    if (tsv) {
        stats.printStatisticsTSV();
    } else {
        std::cout << "正在分析目录: " << directoryPath << std::endl;
        std::cout << "请稍候..." << std::endl;
        stats.printStatistics();
    }
    
    return 0;
}
