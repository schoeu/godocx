package logger

import(
    "os"
    "bufio"
    "fmt"
    "time"
    "path/filepath"
    "strconv"  
)

var (
    crtAccessIndex = 0
    crtErrorIndex = 0
)

const (
    infoLog = "access.log."
    errorLog = "error.log."
    formatStr = "2006-01-02 15:04:05.000"
    shortFormat = "2006-01-02"
    maxfileSize int64 = 1024 * 1024 * 10
)

type Logger struct {
    ty string
    path string
    size int64
    date string
    accessWr *bufio.Writer
    errorWr *bufio.Writer
    crtAccessLogPath string
    crtErrorLogPath string
}

var Slogger = &Logger{
    ty: "file",
}

// 日志模块初始化
func (l *Logger) New(logPath string) *Logger{
    dirPath := filepath.Dir(logPath)
    l.path = dirPath
    l.date = time.Now().Format(shortFormat)
    l.newFile("all", l.date)
    
    l.accessWr.Write([]byte("test"))
    l.accessWr.Flush()
    return l
}

// 日志文件管理
func (l *Logger) newFile(ty string, date string){
    accessLogPath := l.path + infoLog + date
    errorLogPath := l.path + infoLog + date
    l.crtAccessLogPath = accessLogPath
    l.crtErrorLogPath = errorLogPath

    if ty == "all" {
        l.newSimpleFile(l.crtAccessLogPath, "access")
        l.newSimpleFile(l.crtErrorLogPath, "error")
    } else if ty == "access" {
        l.newSimpleFile(l.crtAccessLogPath, "access")
    } else if ty == "error" {
        l.newSimpleFile(l.crtErrorLogPath, "error")
    }
}

// 处理单个日志文件
func (l *Logger) newSimpleFile(p , ty string){
    f, err := os.OpenFile(p , os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
    if err != nil {
        fmt.Println(err)
    }
    defer f.Close()
    w := bufio.NewWriter(f)

    if ty == "access" {
        l.accessWr = w
    } else if ty == "error" {
        l.errorWr = w
    }
}

// 普通信息
func (l *Logger) Info(info string) {
    now := time.Now().Format(formatStr)
    l.accessWr.WriteString(now + " " + info)
}

// 错误信息
func (l *Logger) Error(errorInfo string) {
    now := time.Now().Format(formatStr)
    l.accessWr.WriteString(now + " " + errorInfo)
}

// 文件大小&日期检测
func (l *Logger) check() {
    now := time.Now().Format(shortFormat)
    if now != l.date {
        l.newFile("all", l.date)
        crtAccessIndex = 0
        crtErrorIndex = 0
    }

    l.overSize("access")
    l.overSize("error")
}

func  (l *Logger) overSize(ty string) {
    path := ""

    if ty == "access" {
        crtAccessIndex ++
        path = l.crtAccessLogPath + strconv.Itoa(crtAccessIndex)
    } else if ty == "error" {
        crtErrorIndex ++
        path = l.crtErrorLogPath + strconv.Itoa(crtErrorIndex)
    }

    fileSize, err := fileSize(path)
    if err != nil {
        fmt.Println(err)
    }

    if fileSize > maxfileSize {
        l.newSimpleFile(path, ty)
    }
}

// 检测文件大小
func fileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}
