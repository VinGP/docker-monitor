package pinger

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"log/slog"
	"pinger/internal/model"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	probing "github.com/prometheus-community/pro-bing"
)

type ContainerPinger struct {
	backend         StatusSaver
	ctx             context.Context
	networks        map[string]struct{}
	composeProjects map[string]struct{}
	cancelFn        context.CancelFunc
	dockerCli       *client.Client
	interval        time.Duration
	allContainers   bool
	pingerLabel     bool
}

type StatusSaver interface {
	SaveContainerStatus(containerStatus model.ContainerStatus) error
}

const pingerLabel = "pinger"

func NewPinger(backend StatusSaver, interval time.Duration, opts ...Option) (*ContainerPinger, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		cancel()
		return nil, err
	}

	p := &ContainerPinger{
		backend:   backend,
		ctx:       ctx,
		cancelFn:  cancel,
		interval:  interval,
		dockerCli: cli,
	}
	if len(opts) == 0 {
		p.allContainers = true

		return p, nil
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}

func (p *ContainerPinger) Run() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.pingAllContainers()
		}
	}
}

func (p *ContainerPinger) Stop() {
	p.cancelFn()
}

func (p *ContainerPinger) pingAllContainers() {
	containers, err := p.dockerCli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		slog.Error("error getting containers list", slog.Any("error", err))
		return
	}

	var wg sync.WaitGroup

	for i := range containers {
		c := &containers[i]
		if p.neededPing(c) {
			wg.Add(1)
			go p.pingAndSaveContainer(c, &wg)
		}
	}
	wg.Wait()
}

func (p *ContainerPinger) pingAndSaveContainer(dockerContainer *types.Container, wg *sync.WaitGroup) {
	defer wg.Done()

	status, err := p.pingContainer(dockerContainer)

	if err != nil {
		slog.Error("error pinging container",
			slog.Any("error", err),
			slog.Any("container", dockerContainer))
		return
	}

	for _, s := range status {
		if err := p.backend.SaveContainerStatus(s); err != nil {
			slog.Error("error saving container status",
				slog.Any("error", err),
				slog.Any("status", status))
		}
	}
}

func (p *ContainerPinger) neededPing(c *types.Container) bool {
	if p.allContainers {
		return true
	}

	for network := range c.NetworkSettings.Networks {
		if _, ok := p.networks[network]; ok {
			return true
		}
	}

	if p.pingerLabel {
		if c.Labels[pingerLabel] == "true" {
			return true
		}
	}

	if _, ok := p.composeProjects[c.Labels["com.docker.compose.project"]]; ok {
		return true
	}

	return false
}

func (p *ContainerPinger) pingContainer(c *types.Container) ([]model.ContainerStatus, error) {
	var statuses []model.ContainerStatus

	inspected, err := p.dockerCli.ContainerInspect(context.Background(), c.ID)

	if err != nil {
		return statuses, err
	}

	sandboxKey := inspected.NetworkSettings.SandboxKey

	for _, network := range inspected.NetworkSettings.Networks {
		pingTime, errNetnsPing := p.netnsPing(network.IPAddress, sandboxKey)

		if errNetnsPing == nil {
			timeNow := time.Now()
			statuses = append(statuses, model.ContainerStatus{
				IPAddress:   network.IPAddress,
				PingTime:    &pingTime,
				LastSuccess: &timeNow,
			})
			continue
		}

		pingTime, errPing := p.pingIP(network.IPAddress)

		if errPing == nil {
			timeNow := time.Now()
			statuses = append(statuses, model.ContainerStatus{
				IPAddress:   network.IPAddress,
				PingTime:    &pingTime,
				LastSuccess: &timeNow,
			})
			continue
		}
		slog.Error("cant ping container",
			slog.Any("ping", errPing),
			slog.Any("netns ping", errNetnsPing),
			slog.String("containerID", c.ID))
	}

	return statuses, err
}

var ErrPingFailed = errors.New("can't ping")

func (p *ContainerPinger) pingIP(ip string) (float64, error) {
	const pingerTimeout = 3 * time.Second
	const pingerCount = 1

	pinger, err := probing.NewPinger(ip)

	if err != nil {
		return 0, err
	}

	pinger.Count = pingerCount
	pinger.Timeout = pingerTimeout

	err = pinger.Run()

	if err != nil {
		return 0, err
	}
	stats := pinger.Statistics()

	if int(stats.AvgRtt.Nanoseconds()) == 0 {
		return 0, ErrPingFailed
	}

	pingTime := float64(stats.AvgRtt.Nanoseconds()) / float64(time.Millisecond)
	return pingTime, nil
}

var ErrGetDockerPidFailed = errors.New("can't get docker pid")
var ErrParsePingOutputFailed = errors.New("can't parse ping output")

func (p *ContainerPinger) netnsPing(ip, netns string) (float64, error) {
	pid, err := exec.Command("cat", "/host/var/run/docker.pid").Output()

	if err != nil {
		return 0, ErrGetDockerPidFailed
	}

	cmd := exec.Command("nsenter",
		"--mount=/host/proc/"+string(pid)+"/ns/mnt",
		"nsenter",
		"--net="+netns,
		"ping", "-c", "3", "-W", "1", ip,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(`rtt min/avg/max/mdev = ([\d.]+)/([\d.]+)/([\d.]+)/([\d.]+) ms`)
	matches := re.FindStringSubmatch(string(output))

	const pingOutputMatchesCount = 3

	if len(matches) < pingOutputMatchesCount {
		return 0, ErrParsePingOutputFailed
	}

	var avg float64
	_, err = fmt.Sscanf(matches[2], "%f", &avg)

	if err != nil {
		return 0, err
	}

	return avg, nil
}
