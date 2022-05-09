package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct QoSFlowsUsageReportList */
/* QoSFlowsUsageReportItem */
type QoSFlowsUsageReportList struct {
	List []QoSFlowsUsageReportItem `aper:"valueExt,sizeLB:1,sizeUB:64"`
}
