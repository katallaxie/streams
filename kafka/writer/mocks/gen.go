package mocks

//go:generate go run github.com/golang/mock/mockgen -destination=mocks.go -package mocks github.com/ionos-cloud/streams/kafka/writer Writer
