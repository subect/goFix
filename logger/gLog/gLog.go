package gLog

type GLog struct {
}

var (
	gLogImp *GLog
)

func Init() {
	gLogImp = NewGLog()
}

func NewGLog() *GLog {
	return &GLog{}
}

func GetGLog() *GLog {
	return gLogImp
}

func SetSysLevel(level int) {
	GetSysLogger().SetLevel(parseLogLevel(level))
}
