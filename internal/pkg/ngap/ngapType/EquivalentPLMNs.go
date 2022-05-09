package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct EquivalentPLMNs */
/* PLMNIdentity */
type EquivalentPLMNs struct {
	List []PLMNIdentity `aper:"sizeLB:1,sizeUB:15"`
}
