package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct QosFlowListWithDataForwarding */
/* QosFlowItemWithDataForwarding */
type QosFlowListWithDataForwarding struct {
	List []QosFlowItemWithDataForwarding `aper:"valueExt,sizeLB:1,sizeUB:64"`
}
