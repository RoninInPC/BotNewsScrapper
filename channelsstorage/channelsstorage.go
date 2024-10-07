package channelsstorage

type ChannelsStorage interface {
	Add(int64)
	GetChatsId() []int64
	Size() int
}
