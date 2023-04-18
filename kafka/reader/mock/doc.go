package mock

//go:generate go run github.com/golang/mock/mockgen -destination=mocks.go -package mock github.com/katallaxie/streams/kafka/reader Reader
