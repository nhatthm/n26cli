package cli

type globalConfig struct {
	Verbose    bool
	Debug      bool
	ConfigFile string
}

type apiConfig struct {
	Username string
	Password string
	Format   string
}
