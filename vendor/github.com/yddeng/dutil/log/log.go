package log

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type Level int

const (
	TRACE Level = iota // TRACE 用户级基本输出
	DEBUG              // DEBUG 用户级调试输出
	INFO               // INFO  用户级重要信息
	WARN               // WARN  用户级警告信息
	ERROR              // ERROR 用户级错误信息
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

const calldepth = 2

var (
	LevelString = [...]string{
		"[TRACE]",
		"[DEBUG]",
		"[INFO] ",
		"[WARN] ",
		"[ERROR]",
	}

	defOutLevel = map[Level]struct{}{
		TRACE: {},
		DEBUG: {},
		INFO:  {},
		WARN:  {},
		ERROR: {},
	}

	stdOut     = true //控制台输出,默认开启
	stdConsole = os.Stdout
)

//关闭控制台输出
func CloseStdOut() {
	stdOut = false
	//stdConsole.Close()
}

type Logger struct {
	mu     sync.Mutex
	flag   int
	buf    []byte
	outLev map[Level]struct{}

	async  bool //异步输出，默认同步
	logOut *OutFile
}

func NewLogger(basePath, fileName string, maxSize ...int) *Logger {
	var writeMaxSize = 0
	if len(maxSize) > 0 {
		writeMaxSize = maxSize[0]
	}
	return newLogger(basePath, fileName, writeMaxSize)
}

// newLogger
// maxSize 分割文件字节数
// maxSize = 0 不按照字节大小分割，仅按照日期分割
func newLogger(basePath, fileName string, maxSize int) *Logger {
	out := newOutFile(basePath, fileName, maxSize)
	return &Logger{
		flag:   Ldate | Lmicroseconds | Lshortfile,
		outLev: defOutLevel, // 默认所有日志类型输出
		logOut: out,
	}
}

//开启异步输出
func (l *Logger) AsyncOut() {
	l.async = true

	l.logOut.chCache = make(chan *message, 512)
	go l.logOut.run()
}

//设置输出等级
func (l *Logger) SetOutLevel(levels ...Level) {
	if len(levels) != 0 {
		l.outLev = map[Level]struct{}{}
		for _, lev := range levels {
			l.outLev[lev] = struct{}{}
		}
	}
}

type OutFile struct {
	basePath     string
	fileName     string
	writer       *os.File
	createTime   time.Time     //底层文件的创建时间
	writeMaxSize int           //文件写入最大字节数
	writeSize    int           //累计写入文件的字节数量
	chCache      chan *message //异步缓存
}

type message struct {
	now  *time.Time
	data []byte
}

func newOutFile(basePath, fileName string, maxSize int) *OutFile {
	return &OutFile{basePath: basePath, fileName: fileName, writeMaxSize: maxSize}
}

func (out *OutFile) openFile(now *time.Time) {
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	dir := fmt.Sprintf("%s/%04d-%02d-%02d", out.basePath, year, month, day)
	if nil == os.MkdirAll(dir, os.ModePerm) {
		path := fmt.Sprintf("%s/%s.%02d.%02d.%02d.log", dir, out.fileName, hour, min, sec)
		mode := os.O_RDWR | os.O_CREATE | os.O_APPEND
		file, err := os.OpenFile(path, mode, 0666)
		if out.writer != nil {
			_ = out.writer.Close()
		}
		if nil == err {
			out.writer = file
			out.writeSize = 0
			out.createTime = *now
		}
	}
}

func (out *OutFile) checkOutFile(now *time.Time) bool {
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

func (out *OutFile) flush(now *time.Time, buff []byte) {
	data := make([]byte, len(buff))
	copy(data, buff[:])
	out.chCache <- &message{
		now:  now,
		data: data,
	}
}

func (out *OutFile) run() {
	for {
		msg := <-out.chCache
		out.write(msg.now, msg.data)
	}
}

func (out *OutFile) write(now *time.Time, buff []byte) {
	if false == out.checkOutFile(now) {
		out.openFile(now)
	}

	if stdOut {
		_, _ = stdConsole.Write(buff)
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

func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	*buf = append(*buf, ' ')
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
	if _, ok := l.outLev[lev]; !ok {
		return
	}

	now := time.Now() // get this early.

	prefix := LevelString[lev]
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
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.buf = append(l.buf, prefix...)
	l.formatHeader(&l.buf, now, file, line)
	l.buf = append(l.buf, text...)
	if len(text) == 0 || text[len(text)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	if l.async {
		l.logOut.flush(&now, l.buf)
	} else {
		l.logOut.write(&now, l.buf)
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	l.output(DEBUG, "", v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.output(DEBUG, format, v...)
}

func (l *Logger) Infoln(v ...interface{}) {
	l.output(INFO, "", v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(INFO, format, v...)
}

func (l *Logger) Warnln(v ...interface{}) {
	l.output(WARN, "", v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.output(WARN, format, v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	l.output(ERROR, "", v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output(ERROR, format, v...)
}
