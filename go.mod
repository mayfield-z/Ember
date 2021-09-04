module github.com/mayfield-z/ember

go 1.17

replace (
	github.com/mayfield-z/ember/internal/pkg/gnb => ./internal/pkg/gnb
	github.com/mayfield-z/ember/internal/pkg/logger => ./internal/pkg/logger
	github.com/mayfield-z/ember/internal/pkg/mq => ./internal/pkg/rule_queue
	github.com/mayfield-z/ember/internal/pkg/packet_driver => ./internal/pkg/packet_driver
	github.com/mayfield-z/ember/internal/pkg/timer => ./internal/pkg/timer
	github.com/mayfield-z/ember/internal/pkg/ue => ./internal/pkg/ue
	github.com/mayfield-z/ember/internal/pkg/utils => ./internal/pkg/utils
)

require (
	git.cs.nctu.edu.tw/calee/sctp v1.1.0
	github.com/antonfisher/nested-logrus-formatter v1.3.0
	github.com/free5gc/aper v1.0.1
	github.com/free5gc/nas v1.0.0
	github.com/free5gc/ngap v1.0.2
	github.com/google/gopacket v1.1.19
	github.com/looplab/fsm v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/free5gc/logger_conf v1.0.0 // indirect
	github.com/free5gc/logger_util v1.0.0 // indirect
	github.com/free5gc/openapi v1.0.0 // indirect
	github.com/free5gc/path_util v1.0.0 // indirect
	github.com/free5gc/util_3gpp v1.0.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.6.3 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/ugorji/go/codec v1.2.1 // indirect
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9 // indirect
	golang.org/x/sys v0.0.0-20201211090839-8ad439b19e0f // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
