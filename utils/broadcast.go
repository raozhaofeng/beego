package utils

// Broadcast 广播
type Broadcast struct {
	chBroadcast chan *Message
	chListeners map[int64]chan *Message
}

// Message 广播消息
type Message struct {
	Type string
	Data interface{}
}

// NewBroadcast 创建广播
func NewBroadcast() *Broadcast {
	return &Broadcast{
		chBroadcast: make(chan *Message),
		chListeners: make(map[int64]chan *Message, 0),
	}
}

// Listener 新增监听者
func (c *Broadcast) Listener(chIndex int64) chan *Message {
	ch := make(chan *Message)
	c.addListener(chIndex, ch)
	return ch
}

// RemoveListener 移除监听者
func (c *Broadcast) RemoveListener(chIndex int64, ch chan *Message) {
	if _, ok := c.chListeners[chIndex]; ok {
		delete(c.chListeners, chIndex)
		close(ch)
	}
}

// addListener 添加监听者
func (c *Broadcast) addListener(chIndex int64, ch chan *Message) {
	if _, ok := c.chListeners[chIndex]; !ok {
		c.chListeners[chIndex] = ch
	}
}

// Send 广播消息
func (c *Broadcast) Send(msg *Message) {
	c.chBroadcast <- msg
}

// Start 启动广播
func (c *Broadcast) Start() chan *Message {
	go func() {
		for {
			select {
			case v, ok := <-c.chBroadcast:
				// 	如果广播通道关闭，则关闭掉所有的消费者通道
				if !ok {
					goto terminate
				}
				// 将值转发到所有的消费者channel
				for _, ch := range c.chListeners {
					ch <- v
				}
			}
		}

	terminate:
		//	关闭所有的消费通道
		for _, ch := range c.chListeners {
			close(ch)
		}
	}()
	return c.chBroadcast
}
