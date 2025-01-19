module github.com/malinatrash/kartinki-gateway

go 1.22.4

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.25.1
	github.com/malinatrash/kartinki-proto v0.0.0-00010101000000-000000000000
)

require github.com/joho/godotenv v1.5.1 // indirect

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/grpc v1.69.4
	google.golang.org/protobuf v1.36.3 // indirect
)

replace github.com/malinatrash/kartinki-proto => ../proto
