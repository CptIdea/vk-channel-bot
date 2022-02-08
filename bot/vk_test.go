package bot

import (
	"context"
	"errors"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"time"
)

type vkMock struct {
	//lp *lpMock

	needDeleteError bool
}

type lpMock struct {
	process chan events.MessageNewObject
	handler func(_ context.Context, obj events.MessageNewObject)

	fastStop bool
	once     bool
}

func (l *lpMock) MessageNew(f func(_ context.Context, obj events.MessageNewObject)) {
	l.handler = f
}

func (l *lpMock) Run() error {
	if l.fastStop {
		return nil
	}
	for r := range l.process {
		l.handler(context.Background(), r)
		if l.once {
			time.Sleep(time.Millisecond * 10)
			return nil
		}
	}
	return nil
}

func (v *vkMock) MessagesDelete(params api.Params) (response api.MessagesDeleteResponse, err error) {
	if v.needDeleteError {
		return nil, errors.New("some error")
	}
	return api.MessagesDeleteResponse{}, nil
}
