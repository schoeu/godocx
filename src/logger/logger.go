package logger

import(
    "os"
    "bufio"
    "fmt"
    "time"
    "path/filepath"
)

var (
    infoLog = "access.log."
    errorLog = "error.log."
    formatStr = "2006-01-02 15:04:05.000"
    shortFormat = "2006-01-02"
)

type slog struct {
    ty string
    path string
    size int64
    date string
    accessWr *bufio.Writer
    errorWr *bufio.Writer
}

func (l *slog) New(logPath string) *slog{
    dirPath := filepath.Dir(logPath)
    l.path = dirPath
    l.date = time.Now().Format(shortFormat)
    l.newFile("all", l.date)
    
    l.accessWr.Write([]byte("test"))
    l.accessWr.Flush()
    return l
}

func (l *slog) newFile(ty string, date string) *os.File{
    if ty == "all" {
        af, err := os.OpenFile(l.path + infoLog + date , os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
        ef, err := os.OpenFile(l.path + errorLog + date , os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
        aw := bufio.NewWriter(af)
        l.accessWr = aw
        ew := bufio.NewWriter(ef)
        l.errorWr = ew
    }

    f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
    if err != nil {
        fmt.Println(err)
    }
    defer f.Close()
    w := bufio.NewWriter(f)
    l.accessWr = w

    return f
}

func (l *slog) newSimpleFile(p string) *os.File{
    
}


// 普通信息
func (l *slog) Info(info string) {
    now := time.Now().Format(formatStr)
    l.accessWr.WriteString(now + " " + info)
}

// 错误信息
func (l *slog) Error(errorInfo string) {
    now := time.Now().Format(formatStr)
    l.accessWr.WriteString(now + " " + errorInfo)
}

// 文件大小&日期检测
func (l *slog) check() {
    now := time.Now().Format(shortFormat)
    if now != l.date {
        l.newFile()
    }

}
