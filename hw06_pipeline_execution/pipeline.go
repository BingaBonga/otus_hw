package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(gracefulStop(in, done))
	}

	return gracefulStop(in, done)
}

func gracefulStop(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
			//nolint:revive
			for range in {
			}
		}()

		for {
			select {
			case <-done:
				return
			case v := <-in:
				if v == nil {
					return
				}

				out <- v
			}
		}
	}()

	return out
}
