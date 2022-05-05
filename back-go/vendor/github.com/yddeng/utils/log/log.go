package log

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	Ldate         = 1 << iota                                   // 日期标记位  2019/01/23
	Ltime                                                       // 时间标记位  01:23:12
	Lmicroseconds                                               // 微秒级标记位 01:23:12.111222
	LUTC                                                        // if Ldate or Ltime is set, use UTC rather than the local time zone
	Llongfile                                                   // 完整文件名称 /home/go/src/zinx/server.go
	Lshortfile                                                  // 最后文件名   server.go
	Llevel                                                      // 当前日志级别： 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
	LstdFlags     = Ldate | Ltime                               // 标准头部日志格式
	LdefFlags     = Ldate | Lmicroseconds | Llevel | Lshortfile // 默认日志头部格式
)

type Level int

const (
	DEBUG Level = iota // DEBUG 用户级调试输出
	INFO               // INFO  用户级重要信息
	WARN               // WARN  用户级警告信息
	ERROR              // ERROR 用户级错误信息
	PANIC              // PANIC
	FATAL              // FATAL
)

var levelString = []string{
	"[DEBUG]",
	"[INFO] ",
	"[WARN] ",
	"[ERROR]",
	"[PANIC]",
	"[FATAL]",
}

type Logger struct {
	flag         int
	prefix       string //日志前缀
	calldepth    int
	debugClosed  bool
	stdOutClosed bool
	buf          []byte
	outFile      *OutFile
	mu           sync.Mutex
}

func NewLogger(basePath, fileName string, maxSize ...int) *Logger {
	if fileName == "" {
		panic("log:New fileName is empty. ")
	}
	return newLogger(newOutFile(basePath, fileName, maxSize...))
}

func newLogger(out *OutFile) *Logger {
	return &Logger{outFile: out, flag: LdefFlags, calldepth: 2}
}

func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}

func (l *Logger) SetOutput(basePath, fileName string, maxSize ...int) {
	outFile := newOutFile(basePath, fileName, maxSize...)

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.outFile != nil {
		l.outFile.close()
	}
	l.outFile = outFile
}

//设置日志的 用户自定义前缀字符串
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// 关闭控制台输出
func (l *Logger) CloseStdOut() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.stdOutClosed = true

}

// 函数调用层数
func (l *Logger) SetCallDepth(calldepth int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.calldepth = calldepth
}

// 关闭debug日志输出
func (l *Logger) CloseDebug() {
	l.debugClosed = true
}

type OutFile struct {
	basePath     string
	fileName     string
	writer       *os.File
	createTime   time.Time //底层文件的创建时间
	writeMaxSize int       //文件写入最大字节数
	writeSize    int       //累计写入文件的字节数量
}

// newOutFile
// maxSize 分割文件字节数
// maxSize = 0 不按照字节大小分割，仅按照日期分割
func newOutFile(basePath, fileName string, maxSize ...int) *OutFile {
	var writeMaxSize = 0
	if len(maxSize) > 0 {
		writeMaxSize = maxSize[0]
	}
	return &OutFile{basePath: basePath, fileName: fileName, writeMaxSize: writeMaxSize}
}

func (out *OutFile) close() {
	if out.writer != nil {
		_ = out.writer.Close()
		out.writer = nil
	}
}

func (out *OutFile) openFile(now *time.Time) {
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	dir := fmt.Sprintf("%s/%04d-%02d-%02d", out.basePath, year, month, day)
	if nil == os.MkdirAll(dir, os.ModePerm) {
		path := fmt.Sprintf("%s/%s.%02d.%02d.%02d.log", dir, out.fileName, hour, min, sec)
		mode := os.O_RDWR | os.O_CREATE | os.O_APPEND
		file, err := os.OpenFile(path, mode, 0666)
		out.close()
		if nil == err {
			out.writer = file
			out.writeSize = 0
			out.createTime = *now
		}
	}
}

func (out *OutFile) checkOutFile(now *time.Time) bool {
	if out.fileName == "" {
		return true
	}

	if out.writer == nil {
		return false
	}

	if out.writeMaxSize != 0 && out.writeSize >= out.writeMaxSize {
		return false
	}

	fyear, fmonth, fday := out.createTime.Date()
	year, month, day := now.Date()
	if fyear != year || fmonth != month || fday != day {
		return false
	}

	return true
}

func (out *OutFile) write(now *time.Time, buff []byte) {
	if false == out.checkOutFile(now) {
		out.openFile(now)
	}

	if out.writer != nil {
		n, err := out.writer.Write(buff)
		if nil == err {
			out.writeSize += n
		}
	}
}

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int, level Level) {
	if l.prefix != "" {
		*buf = append(*buf, l.prefix...)
		*buf = append(*buf, ' ')
	}

	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}

	if l.flag&Llevel != 0 {
		*buf = append(*buf, levelString[level]...)
		*buf = append(*buf, ' ')
	}

	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func (l *Logger) output(lev Level, format string, v ...interface{}) {
	now := time.Now() // get this early.

	text := ""
	if format == "" {
		text = fmt.Sprintln(v...)
	} else {
		text = fmt.Sprintf(format, v...)
	}

	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(l.calldepth)
		if !ok {
			file = "unknown-file"
			line = 0
		}
		l.mu.Lock()
	}

	//清零buf
	l.buf = l.buf[:0]
	//写日志头
	l.formatHeader(&l.buf, now, file, line, lev)
	//写日志内容
	l.buf = append(l.buf, text...)
	//补充回车
	if len(text) > 0 && text[len(text)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	if !l.stdOutClosed {
		_, _ = os.Stderr.Write(l.buf)
	}
	if l.outFile != nil {
		l.outFile.write(&now, l.buf)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l != nil && !l.debugClosed {
		l.output(DEBUG, "", v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l != nil && !l.debugClosed {
		l.output(DEBUG, format, v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l != nil {
		l.output(INFO, "", v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l != nil {
		l.output(INFO, format, v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l != nil {
		l.output(ERROR, "", v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l != nil {
		l.output(ERROR, format, v...)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	if l != nil {
		l.output(FATAL, "", v...)
		os.Exit(1)
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l != nil {
		l.output(FATAL, format, v...)
		os.Exit(1)
	}
}

func (l *Logger) Panic(v ...interface{}) {
	if l != nil {
		l.output(PANIC, "", v...)
		panic(fmt.Sprintln(v...))
	}
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	if l != nil {
		l.output(PANIC, format, v...)
		panic(fmt.Sprintf(format, v...))
	}
}

// Stack
func runStack(v ...interface{}) string {
	s := fmt.Sprintln(v...)
	buf := make([]byte, 64*1024)
	n := runtime.Stack(buf, true) //得到当前堆栈信息
	s += string(buf[:n])
	s += "\n"
	return s
}

func (l *Logger) Stack(v ...interface{}) {
	if l != nil {
		l.output(ERROR, "", runStack(v...))
	}
}
