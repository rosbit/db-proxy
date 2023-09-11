package db

import (
	djs "github.com/rosbit/dukgo"
	"db-proxy/conf"
	"fmt"
)

const (
	name2dsn_driver = "dbproxy"
)

func init() {
	djs.InitCache()
}

type fnGetDsnByName func(name string) (dsnWithDriverName string)
var getDsnByName fnGetDsnByName

func GetDSNByName(originDSN string) (dsnWithDriverName string, err error) {
	driverName, realDSN, e := breakDSN(originDSN)
	if e != nil {
		err = e
		return
	}
	if driverName != name2dsn_driver {
		dsnWithDriverName = originDSN
		return
	}

	if err = reloadJS(); err != nil {
		return
	}
	name := realDSN
	dsnWithDriverName = getDsnByName(name)
	if len(dsnWithDriverName) == 0 {
		err = fmt.Errorf("no dsn found for %s", name)
		return
	}
	return
}

func reloadJS() error {
	jsConf := &conf.ServiceConf.DsnParams
	vm, existing, err := djs.LoadFileFromCache(jsConf.JsFile, nil)
	if err != nil {
		return err
	}
	if !existing {
		err = vm.BindFunc(jsConf.JsFunc, &getDsnByName)
	}
	return err
}
