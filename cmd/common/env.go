package cmn

type environment string
type environments struct {
	Development environment
	Test        environment
	Production  environment
}

func stringToEnvironment(s string) environment {
	return environment(s)
}

var Envs = environments{Development: "dev", Test: "test", Production: "prod"}
var currentEnv environment

func setEnvironment(e environment) {
	if currentEnv != "" {
		return
	}
	currentEnv = e
}

func GetEnvironment() environment {
	return currentEnv
}

func initEnvironment(env environment) {
	setEnvironment(env)
}
