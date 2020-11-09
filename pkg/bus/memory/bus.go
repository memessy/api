package memory

import (
	"memessy-api/pkg"
	"memessy-api/pkg/bus"
	"sync"
)

type EventBus struct {
	createdChannels []chan bus.EventCreated
	createdMutex    sync.Mutex
}

func (b *EventBus) SubscribeCreated() chan bus.EventCreated {
	println(len(b.createdChannels))
	ch := make(chan bus.EventCreated)
	b.createdMutex.Lock()
	defer b.createdMutex.Unlock()
	b.createdChannels = append(b.createdChannels, ch)
	return ch
}

func (b *EventBus) Created(meme pkg.Meme) {
	for _, channel := range b.createdChannels {
		go func(c chan bus.EventCreated) { c <- bus.EventCreated{Meme: meme} }(channel)
	}
}

func (b *EventBus) Close() {
	for _, channel := range b.createdChannels {
		close(channel)
	}
}
