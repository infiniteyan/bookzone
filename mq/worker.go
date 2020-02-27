package mq

import (
	"bookzone/models"
	"encoding/json"
	"fmt"
	"github.com/lunny/log"
	"sync"
)

const MSG_QUEUE_SIZE int = 100000

var GlobalSynWorker *SyncWorker
var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		GlobalSynWorker = NewSyncWorker()
		GlobalSynWorker.Run()
	})
}

type SyncWorker struct {
	messageQueue 			IMessageQueue
	messageChan 			chan string
}

func NewSyncWorker() *SyncWorker {
	worker := &SyncWorker{}
	worker.messageQueue =  NewRabbitMQ(MQ_MODEL_PUBSUB, "", "bookzone")
	worker.messageChan = make(chan string, MSG_QUEUE_SIZE)
	return worker
}

func (this *SyncWorker) Push(msg string) {
	if msg != "" {
		this.messageChan <- msg
	}
}

func (this *SyncWorker) handleMessage(msg string) {
	var msgEntity MsgEntity
	err := json.Unmarshal([]byte(msg), &msgEntity)
	if err != nil {
		log.Infof(err.Error())
		return
	}

	switch msgEntity.Type {
	case MSG_TYPE_COMMENT:
		var commentEntity CommentEntity
		err = json.Unmarshal([]byte(msgEntity.Data), &commentEntity)
		if err != nil {
			log.Infof(err.Error())
			return
		}
		log.Infof("comment data:%+v", commentEntity)

		if err := models.NewComments().AddComments(commentEntity.MemberId, commentEntity.BookId, commentEntity.Content); err != nil {
			log.Error("评论失败:", err.Error())
			return
		}
		log.Infof("评论成功, %+v", commentEntity)
	default:
	}
}

func (this *SyncWorker) Run() {
	go func() {
		for {
			select {
			case msg := <- this.messageChan:
				this.messageQueue.Publish("", msg)
			}
		}
	}()

	msgDeliever, err := this.messageQueue.Consume("")
	if err != nil {
		log.Info(err.Error())
		panic(err)
	}

	go func() {
		for {
			var msg string
			for d := range msgDeliever {
				msg = string(d.Body)
				fmt.Printf("received a msg: %s", msg)
				this.handleMessage(msg)
			}
		}
	}()

	log.Info("SyncWorker run...")
}