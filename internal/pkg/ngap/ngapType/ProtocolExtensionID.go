package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type ProtocolExtensionID struct {
	Value int64 `aper:"valueLB:0,valueUB:65535"`
}
