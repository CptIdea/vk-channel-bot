package bot

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"strings"
)

type Controller struct {
	vk  VK
	lp  LongPoll
	log Logger

	admins []int
}

type VK interface {
	MessagesDelete(params api.Params) (response api.MessagesDeleteResponse, err error)
}

type LongPoll interface {
	MessageNew(f func(_ context.Context, obj events.MessageNewObject))
	Run() error
}

type Logger interface {
	Printf(format string, args ...interface{})
}

func NewController(vk VK, lp LongPoll, log Logger, admins []int) *Controller {
	return &Controller{vk: vk, lp: lp, log: log, admins: admins}
}

func (c *Controller) Start() error {
	c.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		c.log.Printf("%d: %s (%d:%d)", obj.Message.FromID, obj.Message.Text, obj.Message.PeerID, obj.Message.ConversationMessageID)

		if !inArray(obj.Message.FromID, c.admins) {
			err := c.deleteMessage(obj)
			if err != nil {
				c.log.Printf("error: %s", err)
			}
		}
	})

	c.log.Printf("Бот запущен со списком админов: %v", c.admins)
	return c.lp.Run()
}

func (c *Controller) deleteMessage(obj events.MessageNewObject) error {
	_, err := c.vk.MessagesDelete(params.NewMessagesDeleteBuilder().DeleteForAll(true).ConversationMessageIDs([]int{obj.Message.ConversationMessageID}).PeerID(obj.Message.PeerID).Params)
	if err != nil && !strings.Contains(err.Error(), "Access denied: message can not be deleted (admin message)") {
		return err
	}
	return nil
}

func inArray(id int, array []int) bool {
	for _, i := range array {
		if i == id {
			return true
		}
	}
	return false
}
