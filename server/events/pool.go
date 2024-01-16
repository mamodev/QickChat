package events

import (
	"context"
	"sync"
	"time"
)

type Message struct {
	ID string `json:"id"`
	ChatID string `json:"chat_id"`
	Message string `json:"message"`
	SenderID string `json:"sender_id"`
	SenderUsername string `json:"sender_username"`
	Timestamp time.Time `json:"timestamp"`
}

type UserMessagePool struct {
	UserID string
	Messages []Message
	TTL int64
}

func (pool *UserMessagePool) AddMessage(message Message) {
	pool.Messages = append(pool.Messages, message)
}

func (pool *UserMessagePool) GetMessages() []Message {
	msg := pool.Messages
	pool.Messages = []Message{}
	return msg
}

type MessagePool struct {
	Pools map[string]*UserMessagePool
	OldestTTL int64
	lock sync.Mutex
	cancel context.CancelFunc
}

func (pool *MessagePool) GarbageCollect(timer *time.Timer, ctx context.Context) {
	select {
	case <-timer.C:
		pool.lock.Lock()
		defer pool.lock.Unlock()

		now := time.Now().Unix()

		for _, userPool := range pool.Pools {
			if userPool.TTL < now {
				delete(pool.Pools, userPool.UserID)
			}
		}

		if len(pool.Pools) > 0 {
			pool.OldestTTL = -1
			for _, userPool := range pool.Pools {
				if pool.OldestTTL == -1 || userPool.TTL < pool.OldestTTL {
					pool.OldestTTL = userPool.TTL
				}
			}
		}

		return
	case <-ctx.Done():
		return
	}
}

func (pool *MessagePool) startGarbageCollection() {
	if (pool.OldestTTL == -1) {
		return
	}

	if pool.cancel != nil {
		pool.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	pool.cancel = cancel

	now := time.Now().Unix()
	timer := time.NewTimer(time.Duration(pool.OldestTTL - now) * time.Second)

	go pool.GarbageCollect(timer, ctx)
}

func (pool *MessagePool) addUser (userID string) {
	ttl := time.Now().Unix() + 60
	pool.Pools[userID] = &UserMessagePool{
		UserID: userID,
		Messages: []Message{},
		TTL: ttl,
	}

	if pool.OldestTTL == -1 || pool.OldestTTL < ttl {
		pool.startGarbageCollection()
		pool.OldestTTL = ttl
	}
}

func (pool *MessagePool) AddMessage (userID []string, message Message) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	// log.Printf("Sending message to users: %v, chat: %s, sender: %s\n", userID, message.ChatID, message.SenderID)

	for _, user := range userID {
		userPool := pool.Pools[user]

		// if userPool == nil {
		// 	log.Println("User not found in pool:", user)
		// }

		if userPool != nil {
			userPool.AddMessage(message)
		}
	}
}

func (pool *MessagePool) GetUsersMessages (userID string) []Message {
	pool.lock.Lock()
	defer pool.lock.Unlock()


	userPool := pool.Pools[userID]
		if userPool == nil {
			pool.addUser(userID)
			return []Message{}
		}
		
	messages := userPool.GetMessages()
	// if len(messages) != 0 {
	// 	fmt.Printf("Sending messages to user: %s\n", userID)
	// 	for _, msg := range messages {
	// 		fmt.Printf("Message id: %s, chat: %s, sender: %s\n", msg.ID, msg.ChatID, msg.SenderID)
	// 	}

	// }
	return messages
}

var Pool *MessagePool = &MessagePool{
	Pools: map[string]*UserMessagePool{},
	OldestTTL: -1,
	lock: sync.Mutex{},
}







	





