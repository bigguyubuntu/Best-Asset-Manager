package api

type readRoutes struct {
	Ping     string
	Ping2    string
	PingCors string
	OpenAPI  string
	OpenAPI2 string
}

var SystemReadRoutes = readRoutes{
	Ping:     "/ping",
	Ping2:    Prefix + "/ping",
	PingCors: Prefix + "/ping_cors",
	// we need the same route twice once with slash and once without
	// this is because sometimes the browser can add a trailing slash
	// which causes an invalid route error to appear.
	OpenAPI:  Prefix + "/openapi",
	OpenAPI2: Prefix + "/openapi/",
}
