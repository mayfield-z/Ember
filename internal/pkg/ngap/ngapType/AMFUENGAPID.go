package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type AMFUENGAPID struct {
	Value int64 `aper:"valueLB:0,valueUB:1099511627775"`
}
