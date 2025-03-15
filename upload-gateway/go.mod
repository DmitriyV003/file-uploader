module upload-gateway

go 1.23.4

replace github.com/dmitriyV003/platform => ../platform/

require (
	github.com/dmitriyV003/platform v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.33.0
	golang.org/x/net v0.34.0
	golang.org/x/sync v0.10.0
	google.golang.org/grpc v1.71.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
