package common

import "os"

var HostingEnvName = os.Getenv("GO_ENVIRONMENT")

func IsProduction() bool {
	return HostingEnvName == "Production" || len(HostingEnvName) == 0
}
func IsDevelopment() bool {
	return HostingEnvName == "Development"
}
