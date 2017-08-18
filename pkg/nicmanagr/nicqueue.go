package nicmanagr

import (
	"github.com/yunify/qingcloud-cni/pkg/common"
)

type Nicqueue struct {
	queue []*string
	size int
	dict map[string]*common.Endpoint
}

func (queue *Nicqueue) Enqueue(nics ...*common.Endpoint) []*string{
	appendableSize := queue.size-len(queue.queue)
	var niclist []*string

	for _,nic := range nics {
		niclist = append(niclist, &nic.EndpointID)
	}

	if appendableSize >0 {
		if appendableSize > len(nics) {
			queue.queue = append(queue.queue,niclist...)
			queue.AddNewEntries(nics...)
			return []*string{}
		} else {
			queue.queue = append(queue.queue, niclist[:appendableSize]...)
			queue.AddNewEntries(nics[:appendableSize]...)
			return niclist[appendableSize:]
		}
	}
	return niclist
}

func (queue *Nicqueue) Dequeue() ( *common.Endpoint) {
	nicid := queue.queue[0]
	queue.queue = queue.queue[1:]
	return queue.dict[*nicid]
}

func (queue *Nicqueue) IsEmpty() bool {
	return len(queue.queue) == 0
}

func (queue *Nicqueue) GetPoolShortFall() int {
	return queue.size-len(queue.queue)
}

func (queue *Nicqueue) Size() int {
	return len(queue.queue)
}

func (queue *Nicqueue) AddNewEntries(endpoint ... *common.Endpoint)  {
	for _, nic := range endpoint {
		if _,ok:=queue.dict[nic.EndpointID]; !ok {
			queue.dict[nic.EndpointID] = nic
		}
	}
}

func (queue *Nicqueue) RemoveEntries(endpoint ... *string) {
	for _, nicid := range endpoint{
		if _, ok := queue.dict[*nicid];ok{
			delete(queue.dict,*nicid)
		}
	}
}

func (queue *Nicqueue) GetEntry (nicid *string) *common.Endpoint{
	nic, ok:= queue.dict[*nicid]
	if ok {
		return nic
	}
	return nil
}