package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = doneWrapStage(stage(out), done)
	}
	return out
}

func doneWrapStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}
