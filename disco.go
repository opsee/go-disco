package disco

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

const (
	RegistratorBasePath = "/opsee.co/registrator"
)

var (
	etcdHost = os.Getenv("ETCD_HOST")
)

type Endpoint struct {
	Host string
	Port string
}

func GetServiceEndpoint(svc string) (*Endpoint, error) {
	// Each host only speaks to its local etcd.
	client := etcd.NewClient([]string{etcdHost})

	svcPath := fmt.Sprintf("%s/%s", RegistratorBasePath, svc)

	resp, err := client.Get(svcPath, false, true)
	defer client.Close()

	if err != nil {
		return nil, err
	}

	nodes := resp.Node.Nodes
	// It's possible that someone will put some bs in the registrator tree.
	// Let's hope that doesn't happen and be lazy for now.
	hp := strings.Split(nodes[rand.Intn(len(nodes))].Value, ":")
	return &Endpoint{
		Host: hp[0],
		Port: hp[1],
	}, nil
}
