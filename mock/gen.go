package mock

//go:generate go run github.com/derision-test/go-mockgen/cmd/go-mockgen --disable-formatting -f github.com/katallaxie/streams -i Source -o source.go
//go:generate go run github.com/derision-test/go-mockgen/cmd/go-mockgen --disable-formatting -f github.com/katallaxie/streams -i Sink -o sink.go
//go:generate go run github.com/derision-test/go-mockgen/cmd/go-mockgen --disable-formatting -f github.com/katallaxie/streams -i Table -o table.go
