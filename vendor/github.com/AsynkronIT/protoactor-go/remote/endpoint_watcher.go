package remote

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
)

func newEndpointWatcher(address string) actor.Producer {
	return func() actor.Actor {
		return &endpointWatcher{
			address: address,
		}
	}
}

type endpointWatcher struct {
	address string
	watched map[string]*actor.PID //key is the watching PID string, value is the watched PID
}

func (state *endpointWatcher) initialize() {
	plog.Info("Started EndpointWatcher", log.String("address", state.address))
	state.watched = make(map[string]*actor.PID)
}

func (state *endpointWatcher) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.initialize()

	case *remoteTerminate:
		delete(state.watched, msg.Watcher.Id)

		terminated := &actor.Terminated{
			Who:               msg.Watchee,
			AddressTerminated: false,
		}
		ref, ok := actor.ProcessRegistry.GetLocal(msg.Watcher.Id)
		if ok {
			ref.SendSystemMessage(msg.Watcher, terminated)
		}

	case *EndpointTerminatedEvent:
		plog.Info("EndpointWatcher handling terminated", log.String("address", state.address))

		for id, pid := range state.watched {

			//try to find the watcher ID in the local actor registry
			ref, ok := actor.ProcessRegistry.GetLocal(id)
			if ok {

				//create a terminated event for the Watched actor
				terminated := &actor.Terminated{
					Who:               pid,
					AddressTerminated: true,
				}

				watcher := actor.NewLocalPID(id)
				//send the address Terminated event to the Watcher
				ref.SendSystemMessage(watcher, terminated)
			}
		}

		ctx.SetBehavior(state.Terminated)

	case *remoteWatch:

		state.watched[msg.Watcher.Id] = msg.Watchee

		//recreate the Watch command
		w := &actor.Watch{
			Watcher: msg.Watcher,
		}

		//pass it off to the remote PID
		SendMessage(msg.Watchee, w, nil, -1)

	case *remoteUnwatch:

		//delete the watch entries
		delete(state.watched, msg.Watcher.Id)

		//recreate the Unwatch command
		uw := &actor.Unwatch{
			Watcher: msg.Watcher,
		}

		//pass it off to the remote PID
		SendMessage(msg.Watchee, uw, nil, -1)

	default:
		plog.Error("EndpointWatcher received unknown message", log.String("address", state.address), log.Message(msg))
	}
}

func (state *endpointWatcher) Terminated(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *remoteWatch:

		//try to find the watcher ID in the local actor registry
		ref, ok := actor.ProcessRegistry.GetLocal(msg.Watcher.Id)
		if ok {

			//create a terminated event for the Watched actor
			terminated := &actor.Terminated{
				Who:               msg.Watchee,
				AddressTerminated: true,
			}
			//send the address Terminated event to the Watcher
			ref.SendSystemMessage(msg.Watcher, terminated)
		}

	case *remoteTerminate, *EndpointTerminatedEvent, *remoteUnwatch:
		// pass

	default:
		plog.Error("EndpointWatcher received unknown message", log.String("address", state.address), log.Message(msg))
	}
}
