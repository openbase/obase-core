package obcore

const (
	//	Framework/platform title. Who knows, it might change..
	OB_TITLE = "OpenBase"
)

//	Global access. ONLY valid when initialized via `NewCtx`.
type Ctx struct {
	//	Represents access to the `Hive`-directory.
	Hive HiveRoot

	//	Set via `NewCtx`, never `nil` (even if logging is disabled).
	Log Logger

	bundles BundleRegistry
}

//	Initializes and returns a new `*Ctx` providing access to the specified `hiveDir`.
//
//	- `hiveDir`: the `Hive`-directory path accessed by `me`.
//
//	- If `logger` is `nil`, `me.Log` is set to a no-op dummy and logging is disabled.
//	In any event, `NewCtx` doesn't log the `err` being returned (if any), so be sure to handle it.
//
//	Whenever `err` is `nil`, `me` is non-`nil` and vice versa.
func NewCtx(hiveDir string, logger Logger) (me *Ctx, err error) {
	me = &Ctx{Log: logger}
	if me.Log == nil {
		me.Log = NewLogger(nil)
	}
	if !IsHive(hiveDir) {
		err = errf("Not a valid %s Hive directory installation: '%s'.", OB_TITLE, hiveDir)
	}
	if err == nil {
		err = me.Hive.init(hiveDir)
	}
	if err != nil {
		me.Dispose()
		me = nil
	}
	return
}

//	Returns the bundle package registry for `me`.
func (me *Ctx) Bundles() *BundleRegistry {
	me.ensureBundles()
	return &me.bundles
}

func (me *Ctx) ensureBundles() {
	defer me.bundles.mutex.UnlockIf(me.bundles.mutex.Lock())
	if me.bundles.Ctx != me { // fancier `if me.bundles.Ctx == nil`
		me.bundles.init(me)
	}
}

//	Clean-up when you're shutting down.
func (me *Ctx) Dispose() (err error) {
	err = me.Hive.dispose()
	return
}
