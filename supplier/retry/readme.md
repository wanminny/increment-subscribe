
###简单的并发模型



##### 一个goroutine 轮询 是否有消息表；并将消息打入带缓冲的通道；
##### 另外一个goroutine 监听该通道；遍历处理逻辑；如果失败则继续重试两次；





