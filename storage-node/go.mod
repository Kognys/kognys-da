module github.com/Kognys/kognys-da/storage-node

go 1.21

require (
	github.com/MOSSV2/dimo-sdk-go v0.0.0-latest
	github.com/ethereum/go-ethereum v1.13.5
	github.com/gin-contrib/cors v1.5.0
	github.com/gin-contrib/zap v0.2.0
	github.com/gin-gonic/gin v1.9.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/urfave/cli/v2 v2.27.1
)

replace github.com/MOSSV2/dimo-sdk-go => ../temp/unibase-sdk-go