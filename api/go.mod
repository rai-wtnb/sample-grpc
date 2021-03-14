module github.com/rai-wtnb/sample-grpc

go 1.16

require (
	github.com/golang/protobuf v1.4.3 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210312152112-fc591d9ea70f // indirect
	google.golang.org/grpc v1.36.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

replace (
	github.com/rai-wtnb/sample-grpc/gen/api => ./gen/api
	github.com/rai-wtnb/sample-grpc/handler => ./handler
)