package lessgo

import (
	"encoding/json"
	"mime"
	"os"
	"path/filepath"

	"github.com/lessgo/lessgo/config"
	"github.com/lessgo/lessgo/session"
)

func registerMime() error {
	for k, v := range mimemaps {
		mime.AddExtensionType(k, v)
	}
	return nil
}

func registerConfig() (err error) {
	os.MkdirAll(CONFIG_DIR, 0777)
	fname := CONFIG_DIR + "/" + APP_CONFIG
	appconf, err := config.NewConfig("ini", fname)
	if err != nil {
		file, err := os.Create(fname)
		file.Close()
		appconf, err = config.NewConfig("ini", fname)
		if err != nil {
			panic(err)
		}
		defaultConfig(appconf)
	} else {
		trySet(appconf)
	}
	return appconf.SaveConfigFile(fname)
}

func registerRouter() error {
	// 从数据读取动态配置

	// 与源码配置进行同步

	// 创建真实路由
	ResetRealRoute()

	return nil
}

func registerSession() (err error) {
	if !AppConfig.Session.Enable {
		return
	}
	conf := map[string]interface{}{
		"cookieName":      AppConfig.Session.CookieName,
		"gclifetime":      AppConfig.Session.GcMaxlifetime,
		"providerConfig":  filepath.ToSlash(AppConfig.Session.ProviderConfig),
		"secure":          AppConfig.Listen.EnableHTTPS,
		"enableSetCookie": AppConfig.Session.EnableSetCookie,
		"domain":          AppConfig.Session.Domain,
		"cookieLifeTime":  AppConfig.Session.CookieLifeTime,
	}
	confBytes, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	sessionConfig := string(confBytes)
	GlobalSessions, err = session.NewManager(AppConfig.Session.Provider, sessionConfig)
	if err != nil {
		return
	}
	go GlobalSessions.GC()
	return
}