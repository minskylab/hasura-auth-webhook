package main

// func runServer(conf *config.Config, service server.Service, errCollector chan error) {
// 	serverURI := fmt.Sprintf("%s:%d", conf.API.Internal.Hostname, conf.API.Internal.Port)

// 	router := mux.NewRouter()
// 	server.SetupRoutes(router, service)
// 	handler := cors.AllowAll().Handler(router)

// 	helpers.PrintLogo(serverURI)
// 	errCollector <- http.ListenAndServe(serverURI, handler)
// }
