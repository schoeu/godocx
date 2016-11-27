package logger

import(
    // "io/ioutil"
    "os"
    "bufio"
    "fmt"
    "time"
)

var (
    formatStr = "2006-01-02 15:04:05.000"
    shortFormat = "2006-01-02"
)

type slog struct {
    ty string
    path string
    size int64
    date string
    w *bufio.Writer
}

func (l *slog) New() *slog{
    f := l.newFile()
    l.date = time.Now().Format(shortFormat)

    w := bufio.NewWriter(f)
    l.w = w
    w.Write([]byte("test"))
    w.Flush()

    return l
}

func (l *slog) newFile() *os.File{
    f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
    if err != nil {
        fmt.Println(err)
    }
    defer f.Close()
    return f
}

// 普通信息
func (l *slog) Info(info string) {
    now := time.Now().Format(formatStr)
    l.w.WriteString(now + " " + info)
}

// 错误信息
func (l *slog) Error(errorInfo string) {
    now := time.Now().Format(formatStr)
    l.w.WriteString(now + " " + errorInfo)
}

// 文件大小&日期检测
func (l *slog) check() {
    now := time.Now().Format(shortFormat)
    if now != l.date {

    }

}
