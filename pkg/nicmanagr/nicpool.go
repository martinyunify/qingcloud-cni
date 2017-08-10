package nicmanagr

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	pool "github.com/jolestar/go-commons-pool"
	"github.com/yunify/qingcloud-cni/pkg/qingactor"
	"github.com/yunify/qingcloud-cni/pkg/common"
	"github.com/yunify/qingcloud-cni/pkg/utils"
	"net"
	"github.com/vishvananda/netlink"
	"syscall"
	"fmt"
)

const NicPoolActorName = "NicPoolActor"


//todo: use allocate policy first

type NicPoolActor struct {
	pool *pool.ObjectPool
}

func NewNicPool(vxnet []string) (*actor.PID, error) {

	// nic pool policy
	config := pool.NewDefaultPoolConfig()
	config.BlockWhenExhausted = false
	config.Lifo = false
	config.MaxIdle = 3
	config.MinIdle = 0
	config.MaxWaitMillis = 0
	config.TimeBetweenEvictionRunsMillis = 30000
	config.MinEvictableIdleTimeMillis = 60000

	nicpool := NicPoolActor{make(map[string]*pool.ObjectPool)}

	qingstub:= actor.NewLocalPID(qingactor.QingCloudActorName)


	for _, subnet := range vxnet {
		factory := pool.NewPooledObjectFactory(
			//Create function
			func() (interface{}, error) {
				result,err :=qingstub.RequestFuture(qingactor.CreateNicMessage{NetworkID:subnet},qingactor.DefaultCreateNicTimeout).Result()
				if err != nil {
					return nil, err
				}
				response := result.(qingactor.CreateNicReplyMessage)
				return pool.NewPooledObject(&response.Nic),nil
			},
			//Destroy func
			func(object *pool.PooledObject) error {
				nic:= object.Object.(*common.Endpoint)
				result,err := qingstub.RequestFuture(qingactor.DeleteNicMessage{Nic:nic.EndpointID},qingactor.DefaultDeleteNicTimeout).Result()
				if err != nil {
					return err
				}
				response := result.(qingactor.DeleteNicReplyMessage)
				return response.Err
			},
			//Validate func
			func(object *pool.PooledObject) bool {
				nic:=object.Object.(*common.Endpoint)
				netnic,err:=utils.LinkByMacAddr(nic.EndpointID)
				if err != nil {
					return false
				}
				if netnic.Attrs().Flags | net.FlagUp !=0{
					return false
				}
				return true
			},
			//
			func(object *pool.PooledObject) error {
				nic:= object.Object.(*common.Endpoint)
				niclink,err := utils.LinkByMacAddr(nic.EndpointID)
				if err != nil {
					return fmt.Errorf("Failed to link iface %v",err)
				}
				addrs, err := netlink.AddrList(niclink, syscall.AF_INET)
				if err == nil && len(addrs) > 0 {
					for _, addr := range addrs {
						err := netlink.AddrDel(niclink, &addr)
						if err != nil {
							return fmt.Errorf("AddrDel err %s addr:%+v, Nic %s", err.Error(), addr, niclink.Attrs().HardwareAddr)
						}
					}
				}
				return netlink.LinkSetDown(niclink)
			},
			func(object *pool.PooledObject) error {
				nic:= object.Object.(*common.Endpoint)
				niclink,err := utils.LinkByMacAddr(nic.EndpointID)
				if err != nil {
					return fmt.Errorf("Failed to link iface %v",err)
				}
				addrs, err := netlink.AddrList(niclink, syscall.AF_INET)
				if err == nil && len(addrs) > 0 {
					for _, addr := range addrs {
						err := netlink.AddrDel(niclink, &addr)
						if err != nil {
							return fmt.Errorf("AddrDel err %s addr:%+v, Nic %s", err.Error(), addr, niclink.Attrs().HardwareAddr)
						}
					}
				}
				return netlink.LinkSetDown(niclink)
			})
		nicpool.poolMap[net] = pool.NewObjectPool(factory, config)
	}

	props := actor.FromInstance(&nicpool)
	return actor.SpawnNamed(props, NicPoolActorName)
}

func (nicPool *NicPoolActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case qingactor.CreateNicMessage:
		reply := qingactor.CreateNicReplyMessage{}
		if pool,exists:= nicPool.poolMap[msg.NetworkID];exists {
			nic,err:= pool.BorrowObject()
			if err!= nil {
				reply.Err =fmt.Errorf("Failed to borrow nic from pool. %v",err)
				context.Respond(reply)
				return
			}
			reply.Nic=*nic.(*common.Endpoint)
			context.Respond(reply)
			return
		}
		reply.Err = fmt.Errorf("Vxnet is not known")
		context.Respond(reply)
	case qingactor.DeleteNicMessage:


	}
}
