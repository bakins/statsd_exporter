package main

type lineHandler interface {
	push(line string)
}

type lineProcessor struct {
	eventHandler eventHandler
	queue        chan string
	workers      int
}

func newLineProcessor(size int, eventHandler eventHandler, workers int) *lineProcessor {
	return &lineProcessor{
		eventHandler: eventHandler,
		queue:        make(chan string, size),
		workers:      workers,
	}
}

func (l *lineProcessor) start() {
	for i := 0; i < l.workers; i++ {
		go func() {
			for line := range l.queue {
				events := lineToEvents(line)
				l.eventHandler.queue(events)
			}
		}()
	}
}

func (l *lineProcessor) push(line string) {
	select {
	case l.queue <- line:
	default:
		linesDropped.Inc()
	}
}

func (l *lineProcessor) stop() {
	close(l.queue)
}

type unbufferedLineHandler struct {
	eventHandler eventHandler
}

func (l *unbufferedLineHandler) push(line string) {
	events := lineToEvents(line)
	l.eventHandler.queue(events)
}
