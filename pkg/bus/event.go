package bus

import "memessy-api/pkg"

type EventCreated struct {
	Meme pkg.Meme
}
