package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type SecurityIndication struct {
	IntegrityProtectionIndication       IntegrityProtectionIndication
	ConfidentialityProtectionIndication ConfidentialityProtectionIndication
	MaximumIntegrityProtectedDataRateUL *MaximumIntegrityProtectedDataRate                  `aper:"optional"`
	IEExtensions                        *ProtocolExtensionContainerSecurityIndicationExtIEs `aper:"optional"`
}
