package scan

type App struct {
	config *Config
}

type Config struct {
	BugReportDirector   string
	GenerateFakeService bool
}

type DeploymentInfo struct {
	Name     string
	Replicas int32
}

type NamespaceInfo struct {
	Name        string
	Deployments map[string]int
}
