package discovery

import (
	"context"
	"fmt"
	"sync"

	"github.com/happycurd/rpclib/protojson"
	"github.com/rs/xid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var oncelock sync.Once
var etcdClient *clientv3.Client

func InitEtcdClient(etcdAddress string) {
	var err error
	if etcdClient != nil {
		return
	}
	oncelock.Do(func() {
		etcdClient, err = clientv3.NewFromURL(etcdAddress)
		if err != nil {
			panic(err)
		}
	})
}
func GetEtcdClient() *clientv3.Client {
	return etcdClient
}

type ServiceInstance struct {
	ServiceID       string
	InstanceID      string
	Endpoint        string
	EndpointManager endpoints.Manager
	Etcd            *clientv3.Client
}

func NewServiceInstance(serviceID string, myaddr string, etcd *clientv3.Client) *ServiceInstance {
	i := &ServiceInstance{
		ServiceID:       serviceID,
		InstanceID:      xid.New().String(),
		Endpoint:        myaddr,
		EndpointManager: nil,
		Etcd:            etcd,
	}
	i.EndpointManager, _ = endpoints.NewManager(etcd, i.ServiceID)
	return i
}
func (i *ServiceInstance) DiscoveryID() string {
	return fmt.Sprintf("%s/%s", i.ServiceID, i.InstanceID)
}
func (i *ServiceInstance) Register() error {
	ctx := context.Background()
	lease := clientv3.NewLease(i.Etcd)
	tick, err := lease.Grant(ctx, 30)
	if err != nil {
		return err
	}
	ch, err := lease.KeepAlive(ctx, tick.ID)
	if err != nil {
		return err
	}
	go func() {
		for range ch {
			// fmt.Println("lease ", v)
		}
	}()

	return i.EndpointManager.AddEndpoint(ctx, i.DiscoveryID(), endpoints.Endpoint{Addr: i.Endpoint}, clientv3.WithLease(tick.ID))

}
func (i *ServiceInstance) Deregister() error {
	return i.EndpointManager.DeleteEndpoint(context.Background(), i.DiscoveryID())
}

func NewConn(serviceID string, client *clientv3.Client) (*grpc.ClientConn, error) {
	resolver, err := resolver.NewBuilder(client)
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
