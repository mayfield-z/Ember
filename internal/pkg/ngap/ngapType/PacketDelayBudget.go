package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PacketDelayBudget struct {
	Value int64 `aper:"valueExt,valueLB:0,valueUB:1023"`
}
