package obcore

import (
	"fmt"
	"path/filepath"
	"time"
)

const (
	//	Framework/platform title. Who knows, it might change..
	OB_TITLE = "OpenBase"
)

var (
	//	Runtime options
	Opt struct {
		//	Set via Init(), never nil (even if logging is disabled)
		Log Logger

		//	Set this to true before calling Init() if the runtime is a sandboxed environment (such
		//	as Google App Engine) with security restrictions (no syscall, no unsafe, no file-writes)
		Sandboxed bool

		//	Set to true before Init() in cmd/ob-server/main.go.
		//	Should remain false in practically all other scenarios.
		//	(If true, much additional logic is executed and server-related resources allocated that
		//	are unneeded when importing this package in a "server-side but server-less client" scenario.)
		Server bool

		initTime time.Time
	}

	RootHub *ObHub
)

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

//	Clean-up. Call this when you're done working with this package and all allocated resources should be released.
func Dispose() {
	if Hive.Watch != nil {
		Hive.Watch.Close()
		Hive.Watch = nil
	}
}

//	Initialization. Call this before working with this package.
//	Before calling Init(), you may need to set Opt.Sandboxed, see Opt for details.
//	If logger is nil, Log is set to a no-op dummy and logging is disabled.
//	In any event, Init() doesn't log the err being returned, so be sure to check it.
//	If err is not nil, this package is not in a usable state and must not be used.
func Init(hiveDirPath string, logger Logger) (err error) {
	Opt.initTime = time.Now()
	if Opt.Log = logger; Opt.Log == nil {
		Opt.Log = NewLogger(nil)
	}
	if Opt.Server {
		Opt.Log.Infof("INIT @ hive = '%s', server = %v, sandboxed = %v", hiveDirPath, Opt.Server, Opt.Sandboxed)
	}
	if !Opt.Sandboxed {
		if hiveDirPath, err = filepath.Abs(hiveDirPath); (err == nil) && !Hive.IsHive(hiveDirPath) {
			err = fmt.Errorf("Not a valid %s Hive directory installation: '%s'.", OB_TITLE, hiveDirPath)
		}
	}
	if err == nil {
		Hive.DirPath = hiveDirPath
		err = Hive.init()
	}
	return
}
