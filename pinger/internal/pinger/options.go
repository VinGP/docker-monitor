package pinger

type Option func(*ContainerPinger)

func AllContainers(val bool) Option {
	return func(p *ContainerPinger) {
		p.allContainers = val
	}
}

func PingerLabel() Option {
	return func(p *ContainerPinger) {
		p.pingerLabel = true
	}
}

func Networks(networks []string) Option {
	return func(p *ContainerPinger) {
		for _, network := range networks {
			p.networks[network] = struct{}{}
		}
	}
}

func ComposeProjects(composeProjects []string) Option {
	return func(p *ContainerPinger) {
		for _, project := range composeProjects {
			p.composeProjects[project] = struct{}{}
		}
	}
}
