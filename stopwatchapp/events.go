package stopwatchapp

import (
	astilectron "github.com/asticode/go-astilectron"
)

// MessageHandlerFn defines an callback for incoming GUI messages
type MessageHandlerFn func(*Message) (interface{}, error)

// OnWindowMessage creates a handler for incoming EventNameWindowEventMessage
// events.
func (a *App) onWindowMessage() astilectron.ListenerMessage {
	return func(msg *astilectron.EventMessage) interface{} {
		var err error
		var res interface{}
		var req Message
		var resp *Message

		if err = msg.Unmarshal(&req); err != nil {
			resp = NewError(req.ID, err.Error())
			a.msgQueue <- *resp
			return nil
		}

		switch req.Key {
		// We pass these on for handling elsewhere. For one reason or another,
		// eg. w.Minimize() will never return here, maybe something to do with
		// the context we're in.
		case RequestWindowClose, RequestWindowMinimize:
			res = "ok"

		default:
			res, err = a.msgHandler(&req)
		}

		// Generate response Message, use Error type if erros were found
		if err != nil {
			resp = NewError(req.ID, err.Error())
		} else {
			resp = NewResponse(req.ID, req.Key, res)
		}

		// Queue response
		a.msgQueue <- *resp
		return nil
	}
}
