package bot

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
	"log"
	"net/http"
	"reflect"
	"sync"
	"testing"
)

type fakeLogger struct {
	err error
	sync.Mutex
}

func (f *fakeLogger) Printf(format string, args ...interface{}) {
	if len(args) < 1 {
		return
	}
	switch args[0].(type) {
	case error:
		f.Lock()
		f.err = args[0].(error)
		f.Unlock()
	}
}

func (f *fakeLogger) getError() error {
	f.Lock()
	defer f.Unlock()
	return f.err
}

func TestNewController(t *testing.T) {
	type args struct {
		vk     *api.VK
		lp     *longpoll.LongPoll
		log    *log.Logger
		admins []int
	}
	tests := []struct {
		name string
		args args
		want *Controller
	}{
		{"values test", args{
			vk:     &api.VK{Client: http.DefaultClient},
			lp:     &longpoll.LongPoll{Client: http.DefaultClient},
			log:    log.Default(),
			admins: []int{123},
		}, &Controller{
			vk:     &api.VK{Client: http.DefaultClient},
			lp:     &longpoll.LongPoll{Client: http.DefaultClient},
			log:    log.Default(),
			admins: []int{123},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewController(tt.args.vk, tt.args.lp, tt.args.log, tt.args.admins); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inArray(t *testing.T) {
	type args struct {
		id    int
		array []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"in array", args{array: []int{123}, id: 123}, true},
		{"not in array", args{array: []int{123}, id: 321}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inArray(tt.args.id, tt.args.array); got != tt.want {
				t.Errorf("inArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_deleteMessage(t *testing.T) {
	type fields struct {
		vk VK
	}
	type args struct {
		obj events.MessageNewObject
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"normal delete",
			fields{
				vk: &vkMock{},
			},
			args{events.MessageNewObject{
				Message: object.MessagesMessage{FromID: 123, PeerID: 123, ConversationMessageID: 123},
			}}, false,
		},
		{"vk error",
			fields{
				vk: &vkMock{needDeleteError: true},
			},
			args{events.MessageNewObject{
				Message: object.MessagesMessage{FromID: 123, PeerID: 123, ConversationMessageID: 123},
			}}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				vk: tt.fields.vk,
			}
			if err := c.deleteMessage(tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("deleteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestController_Start(t *testing.T) {
	type fields struct {
		vk      VK
		lp      *lpMock
		admins  []int
		message *events.MessageNewObject
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"just start",
			fields{
				vk:     &vkMock{},
				lp:     &lpMock{fastStop: true},
				admins: nil,
			},
			false,
		},
		{"vk error",
			fields{
				vk:      &vkMock{needDeleteError: true},
				lp:      &lpMock{once: true, process: make(chan events.MessageNewObject, 1)},
				admins:  nil,
				message: &events.MessageNewObject{Message: object.MessagesMessage{FromID: 123}},
			},
			true,
		},
		{"skip message",
			fields{
				vk:      &vkMock{needDeleteError: true},
				lp:      &lpMock{once: true, process: make(chan events.MessageNewObject, 1)},
				admins:  []int{123},
				message: &events.MessageNewObject{Message: object.MessagesMessage{FromID: 123}},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &fakeLogger{}
			c := &Controller{
				vk:     tt.fields.vk,
				lp:     tt.fields.lp,
				log:    l,
				admins: tt.fields.admins,
			}
			if tt.fields.message != nil {
				tt.fields.lp.process <- *tt.fields.message
			}
			err := c.Start()
			if err == nil && l.getError() != nil {
				err = l.getError()
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
