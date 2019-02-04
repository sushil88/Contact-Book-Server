package conf

import "github.com/koding/multiconfig"

// App is the root-level configuration data for the app
type App struct {
	AppEnv   string `default:"development"`
	Server   string `required:"true"`
	Secret  string `required:"true"`
	Database struct {
		ConnectionString string `required:"true"`
	} `required:"true"`
}

var app *App

// Get either returns an already-loaded configuration or loads config.toml
var Get = func() *App {
	if app == nil {
		appConf := new(App)
		confLoader := multiconfig.DefaultLoader{
			Loader: multiconfig.MultiLoader(
				&multiconfig.TagLoader{},
				&multiconfig.TOMLLoader{Path: "config.toml"},
				&multiconfig.EnvironmentLoader{},
			),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		confLoader.MustLoad(appConf)
		app = appConf
	}

	return app
}
