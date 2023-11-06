package discovery

import (
	"context"
	"fmt"

	"github.com/rs/xid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type ServiceInstance struct {
	ServiceID       string
	InstanceID      string
	Endpoint        string
	EndpointManager endpoints.Manager
	Etcd            *clientv3.Client
}

func GetServiceInstance(serviceID string, myaddr string) *ServiceInstance {
	i := &ServiceInstance{
		ServiceID:       serviceID,
		InstanceID:      xid.New().String(),
		Endpoint:        myaddr,
		EndpointManager: nil,
		Etcd:            GetEtcdClient(),
	}
	i.EndpointManager, _ = endpoints.NewManager(i.Etcd, i.ServiceID)
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
