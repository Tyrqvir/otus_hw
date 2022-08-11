package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// ExecutePipeline if stage is nil -> skip that stage.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := decorateWithDoneChannel(in, done)

	for _, stage := range stages {
		if nil == stage {
			continue
		}
		out = stage(decorateWithDoneChannel(out, done))
	}

	return out
}

func decorateWithDoneChannel(in In, done In) Out {
	biDirectionChan := make(Bi)

	go func(biDirectionChan Bi) {
		defer close(biDirectionChan)
		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				// If channel closed
				if !ok {
					return
				}
				biDirectionChan <- value
			}
		}
	}(biDirectionChan)

	return biDirectionChan
}
