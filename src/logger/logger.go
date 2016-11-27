package logger

import(
    // "io/ioutil"
    "os"
    "bufio"
    "fmt"
)

type slog struct {
    ty string
    path string
    size int64
    date string
    w *bufio.Writer
}

func (l *slog) New() *slog{
    f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
    if err != nil {
        fmt.Println(err)
    }
    defer f.Close()
    w := bufio.NewWriter(f)
    l.w = w
    w.Write([]byte("test"))
    w.Flush()

    return l
}

// 普通信息
func (l *slog) Info(info string) {
    l.w.WriteString(info)
}

// 错误信息
func (l *slog) Error(info string) {
    l.w.WriteString(info)
}
