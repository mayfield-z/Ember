package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PriorityLevelARP struct {
	Value int64 `aper:"valueLB:1,valueUB:15"`
}