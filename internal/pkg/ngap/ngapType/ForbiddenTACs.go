package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct ForbiddenTACs */
/* TAC */
type ForbiddenTACs struct {
	List []TAC `aper:"sizeLB:1,sizeUB:4096"`
}
