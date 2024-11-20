package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(chan interface{})

	go func() {
		bufferChannel := make(chan interface{}, len(stages))
		go readChannelWithDone(in, bufferChannel, done)

		for _, stage := range stages {
			outChannel := stage(bufferChannel)

			bufferChannel = make(chan interface{}, len(stages))
			go readChannel(outChannel, bufferChannel)
		}

		go readChannelWithDone(bufferChannel, out, done)
	}()

	return out
}

func readChannelWithDone(from <-chan interface{}, to chan<- interface{}, done <-chan interface{}) {
	defer close(to)

	for {
		select {
		case <-done:
			return
		case v := <-from:
			if v == nil {
				return
			}

			to <- v
		}
	}
}

func readChannel(from <-chan interface{}, to chan<- interface{}) {
	defer close(to)
	for v := range from {
		to <- v
	}
}
