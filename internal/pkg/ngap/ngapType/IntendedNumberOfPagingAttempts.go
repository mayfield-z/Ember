package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type IntendedNumberOfPagingAttempts struct {
	Value int64 `aper:"valueExt,valueLB:1,valueUB:16"`
}