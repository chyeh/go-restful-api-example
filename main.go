package main

func main() {
	appCfg := newApplicationConfig()
	appCfg.bind(commandlineConfig)

	var apiCfg apiServerConfig
	apiCfg.load(appCfg)
	server := newAPIServer(apiCfg)
	server.run()
}
