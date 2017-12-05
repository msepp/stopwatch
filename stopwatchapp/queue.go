package stopwatchapp

import (
	"log"
)

// Send sends an message into queue for delivery to renderer.
func (a *App) Send(m *Message) {
	a.msgQueue <- *m
}

// messageQueueFlusher sends out message from the queue to given window.
func (a *App) messageQueueFlusher() {
	for {
		select {
		case m, ok := <-a.msgQueue:
			if !ok {
				// Closed. Time to exit.
				return
			}

			if m.Key == RequestWindowClose {
				a.Window.Close()
				a.Renderer.Stop()

			} else if m.Key == RequestWindowMinimize {
				if err := a.Window.Minimize(); err != nil {
					log.Printf("Error minimizing window: %s", err)
				}

			} else {
				if err := a.Window.SendMessage(m); err != nil {
					log.Printf("While sending to client: %s", err)
				}
			}
		}
	}
}
