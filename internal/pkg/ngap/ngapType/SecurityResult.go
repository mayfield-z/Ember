package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type SecurityResult struct {
	IntegrityProtectionResult       IntegrityProtectionResult
	ConfidentialityProtectionResult ConfidentialityProtectionResult
	IEExtensions                    *ProtocolExtensionContainerSecurityResultExtIEs `aper:"optional"`
}
