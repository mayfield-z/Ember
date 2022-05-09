package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type AveragingWindow struct {
	Value int64 `aper:"valueExt,valueLB:0,valueUB:4095"`
}
