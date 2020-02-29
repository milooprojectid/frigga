package bot

// 	service "frigga/modules/service"

// Bot ...
type Bot struct {
	Provider Provider
	Commands *commands
}

// Provider ...
type Provider struct {
	Name    string
	Handler func()
}

// New ...
func New(provider Provider) Bot {
	return Bot{
		Provider: provider,
		Commands: &Commands,
	}
}
