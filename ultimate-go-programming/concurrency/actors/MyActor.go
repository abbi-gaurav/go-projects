package actors

type action func()

type MyActor struct {
	counter       int
	actionChannel chan action
}

func (myActor *MyActor) receive(ch <-chan action) {
	for {
		select {
		case actionItem := <-ch:
			actionItem()
		}
	}
}

func CreateActor() *MyActor {
	channel := make(chan action)
	actor := &MyActor{
		actionChannel: channel,
		counter:       0,
	}
	go actor.receive(channel)

	return actor
}

func (myActor *MyActor) Increment(by int) {
	myActor.actionChannel <- func() {
		myActor.counter = myActor.counter + by
		println(myActor.counter)
	}
}

func (myActor *MyActor) tryToIncrement(by int) {
	select {
	case myActor.actionChannel <- func() {
		myActor.counter = myActor.counter + by
	}:
	default:
		println("channel is full")
	}
}
