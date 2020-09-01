package unitdef

import (
	"os"
	"os/user"
	"path"

	"github.com/dnesting/unit"
)

func getSystemSearchPath() []string {
	if exec, err := os.Executable(); err == nil {
		dir := path.Dir(exec)
		if path.Base(dir) == "bin" {
			return []string{path.Join(dir, "..", "data")}
		} else {
			return []string{dir, path.Join(dir, "data")}
		}
	}
	return nil
}

func getUserSearchPath() []string {
	var sp []string
	u, err := user.Current()
	if err == nil {
		sp = []string{u.HomeDir}
	}
	sp = append(sp, ".")
	return sp
}

func expandDirs(fname string, dirs []string) []string {
	r := make([]string, len(dirs))
	for i := range dirs {
		r[i] = path.Join(dirs[i], fname)
	}
	return r
}

func tryFiles(fnames []string) (defs *unit.Registry, err error) {
	for _, fname := range fnames {
		if fname == "" {
			continue
		}
		defs, err = FromFile(fname)
		if err == nil {
			return
		}
	}
	return
}

func FromStandard() (*unit.Registry, error) {
	sys, err := tryFiles(append([]string{os.Getenv("UNITS_FILE")},
		expandDirs("definitions.units", getSystemSearchPath())...))
	pers, err2 := tryFiles(append([]string{os.Getenv("MYUNITSFILE")},
		expandDirs(".units", getUserSearchPath())...))
	if err != nil && err2 != nil {
		return nil, err
	}
	if sys == nil {
		sys = pers
	} else if pers != nil {
		sys.Merge(pers)
	}
	return sys, nil
}
