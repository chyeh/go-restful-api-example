package main

func main() {
	appCfg := newApplicationConfig()
	appCfg.bind(newEnvironmentVariableConfig())
	appCfg.bind(newCommandLineConfig())

	var apiCfg apiServerConfig
	apiCfg.load(appCfg)
	server := newAPIServer(apiCfg)
	server.run()
}
