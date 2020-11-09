package bus

type CreatedConsumer func(event EventCreated)

func ConsumeCreated(bus EventBus, consumer CreatedConsumer) {
	ch := bus.SubscribeCreated()

	for {
		select {
		case event, open := <-ch:
			if !open {
				return
			}
			go consumer(event)
		}
	}
}
