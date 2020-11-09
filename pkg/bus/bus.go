package bus

import "memessy-api/pkg"

type EventBus interface {
	SubscribeCreated() chan EventCreated

	Created(meme pkg.Meme)

	Close()
}
