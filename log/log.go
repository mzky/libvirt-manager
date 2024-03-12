package log

import (
	"fmt"

	"github.com/bingoohuang/golog"
)

func InitLog(level string, logPath string) {
	layout := `%t{yyyy-MM-dd_HH:mm:ss.SSS} [%-5l{length=5}] %msg %caller{skip=5} %fields%n`
	spec := fmt.Sprintf("level=%s,file=%s,maxSize=10M,maxAge=1095d,gzipAge=3d,stdout=true", level, logPath)
	golog.Setup(golog.Layout(layout), golog.Spec(spec))
}
