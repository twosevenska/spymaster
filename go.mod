module twosevenska.com/spymaster

go 1.17

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/kelseyhightower/envconfig v1.4.0
)

require (
	github.com/smartystreets/goconvey v1.7.2
	internal/mongo v1.0.0
	internal/server v1.0.0
	types v1.0.0
)

replace (
	internal/controllers => ./internal/controllers
	internal/mongo => ./internal/mongo
	internal/server => ./internal/server
	internal/spymaster => ./internal/spymaster
	types => ./types
)

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	internal/controllers v1.0.0 // indirect
	internal/spymaster v1.0.0 // indirect
)
