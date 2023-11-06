package discovery

import (
	"fmt"
	"os"
	"sync"

	"github.com/happycurd/rpclib/protojson"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var etcdlock sync.Once
var etcdClient *clientv3.Client

func GetEtcdClient() *clientv3.Client {
	if etcdClient != nil {
		return etcdClient
	}
	etcdAddr := os.Getenv("ETCD_ADDRESS")
	if etcdAddr == "" {
		etcdAddr = "http://localhost:2379"
	}
	var err error
	etcdlock.Do(func() {
		etcdClient, err = clientv3.NewFromURL(etcdAddr)
		if err != nil {
			panic(err)
		}
	})
	return etcdClient
}

func NewConn(serviceID string) (*grpc.ClientConn, error) {
	resolver, err := resolver.NewBuilder(GetEtcdClient())
	if err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("etcd:///%s", serviceID),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(resolver),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(protojson.JSON{}.Name())),
	)
}
