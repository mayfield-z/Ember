package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type ProtocolIEID struct {
	Value int64 `aper:"valueLB:0,valueUB:65535"`
}

const (
	ProtocolIEIDAllowedNSSAI                               int64 = 0
	ProtocolIEIDAMFName                                    int64 = 1
	ProtocolIEIDAMFOverloadResponse                        int64 = 2
	ProtocolIEIDAMFSetID                                   int64 = 3
	ProtocolIEIDAMFTNLAssociationFailedToSetupList         int64 = 4
	ProtocolIEIDAMFTNLAssociationSetupList                 int64 = 5
	ProtocolIEIDAMFTNLAssociationToAddList                 int64 = 6
	ProtocolIEIDAMFTNLAssociationToRemoveList              int64 = 7
	ProtocolIEIDAMFTNLAssociationToUpdateList              int64 = 8
	ProtocolIEIDAMFTrafficLoadReductionIndication          int64 = 9
	ProtocolIEIDAMFUENGAPID                                int64 = 10
	ProtocolIEIDAssistanceDataForPaging                    int64 = 11
	ProtocolIEIDBroadcastCancelledAreaList                 int64 = 12
	ProtocolIEIDBroadcastCompletedAreaList                 int64 = 13
	ProtocolIEIDCancelAllWarningMessages                   int64 = 14
	ProtocolIEIDCause                                      int64 = 15
	ProtocolIEIDCellIDListForRestart                       int64 = 16
	ProtocolIEIDConcurrentWarningMessageInd                int64 = 17
	ProtocolIEIDCoreNetworkAssistanceInformation           int64 = 18
	ProtocolIEIDCriticalityDiagnostics                     int64 = 19
	ProtocolIEIDDataCodingScheme                           int64 = 20
	ProtocolIEIDDefaultPagingDRX                           int64 = 21
	ProtocolIEIDDirectForwardingPathAvailability           int64 = 22
	ProtocolIEIDEmergencyAreaIDListForRestart              int64 = 23
	ProtocolIEIDEmergencyFallbackIndicator                 int64 = 24
	ProtocolIEIDEUTRACGI                                   int64 = 25
	ProtocolIEIDFiveGSTMSI                                 int64 = 26
	ProtocolIEIDGlobalRANNodeID                            int64 = 27
	ProtocolIEIDGUAMI                                      int64 = 28
	ProtocolIEIDHandoverType                               int64 = 29
	ProtocolIEIDIMSVoiceSupportIndicator                   int64 = 30
	ProtocolIEIDIndexToRFSP                                int64 = 31
	ProtocolIEIDInfoOnRecommendedCellsAndRANNodesForPaging int64 = 32
	ProtocolIEIDLocationReportingRequestType               int64 = 33
	ProtocolIEIDMaskedIMEISV                               int64 = 34
	ProtocolIEIDMessageIdentifier                          int64 = 35
	ProtocolIEIDMobilityRestrictionList                    int64 = 36
	ProtocolIEIDNASC                                       int64 = 37
	ProtocolIEIDNASPDU                                     int64 = 38
	ProtocolIEIDNASSecurityParametersFromNGRAN             int64 = 39
	ProtocolIEIDNewAMFUENGAPID                             int64 = 40
	ProtocolIEIDNewSecurityContextInd                      int64 = 41
	ProtocolIEIDNGAPMessage                                int64 = 42
	ProtocolIEIDNGRANCGI                                   int64 = 43
	ProtocolIEIDNGRANTraceID                               int64 = 44
	ProtocolIEIDNRCGI                                      int64 = 45
	ProtocolIEIDNRPPaPDU                                   int64 = 46
	ProtocolIEIDNumberOfBroadcastsRequested                int64 = 47
	ProtocolIEIDOldAMF                                     int64 = 48
	ProtocolIEIDOverloadStartNSSAIList                     int64 = 49
	ProtocolIEIDPagingDRX                                  int64 = 50
	ProtocolIEIDPagingOrigin                               int64 = 51
	ProtocolIEIDPagingPriority                             int64 = 52
	ProtocolIEIDPDUSessionResourceAdmittedList             int64 = 53
	ProtocolIEIDPDUSessionResourceFailedToModifyListModRes int64 = 54
	ProtocolIEIDPDUSessionResourceFailedToSetupListCxtRes  int64 = 55
	ProtocolIEIDPDUSessionResourceFailedToSetupListHOAck   int64 = 56
	ProtocolIEIDPDUSessionResourceFailedToSetupListPSReq   int64 = 57
	ProtocolIEIDPDUSessionResourceFailedToSetupListSURes   int64 = 58
	ProtocolIEIDPDUSessionResourceHandoverList             int64 = 59
	ProtocolIEIDPDUSessionResourceListCxtRelCpl            int64 = 60
	ProtocolIEIDPDUSessionResourceListHORqd                int64 = 61
	ProtocolIEIDPDUSessionResourceModifyListModCfm         int64 = 62
	ProtocolIEIDPDUSessionResourceModifyListModInd         int64 = 63
	ProtocolIEIDPDUSessionResourceModifyListModReq         int64 = 64
	ProtocolIEIDPDUSessionResourceModifyListModRes         int64 = 65
	ProtocolIEIDPDUSessionResourceNotifyList               int64 = 66
	ProtocolIEIDPDUSessionResourceReleasedListNot          int64 = 67
	ProtocolIEIDPDUSessionResourceReleasedListPSAck        int64 = 68
	ProtocolIEIDPDUSessionResourceReleasedListPSFail       int64 = 69
	ProtocolIEIDPDUSessionResourceReleasedListRelRes       int64 = 70
	ProtocolIEIDPDUSessionResourceSetupListCxtReq          int64 = 71
	ProtocolIEIDPDUSessionResourceSetupListCxtRes          int64 = 72
	ProtocolIEIDPDUSessionResourceSetupListHOReq           int64 = 73
	ProtocolIEIDPDUSessionResourceSetupListSUReq           int64 = 74
	ProtocolIEIDPDUSessionResourceSetupListSURes           int64 = 75
	ProtocolIEIDPDUSessionResourceToBeSwitchedDLList       int64 = 76
	ProtocolIEIDPDUSessionResourceSwitchedList             int64 = 77
	ProtocolIEIDPDUSessionResourceToReleaseListHOCmd       int64 = 78
	ProtocolIEIDPDUSessionResourceToReleaseListRelCmd      int64 = 79
	ProtocolIEIDPLMNSupportList                            int64 = 80
	ProtocolIEIDPWSFailedCellIDList                        int64 = 81
	ProtocolIEIDRANNodeName                                int64 = 82
	ProtocolIEIDRANPagingPriority                          int64 = 83
	ProtocolIEIDRANStatusTransferTransparentContainer      int64 = 84
	ProtocolIEIDRANUENGAPID                                int64 = 85
	ProtocolIEIDRelativeAMFCapacity                        int64 = 86
	ProtocolIEIDRepetitionPeriod                           int64 = 87
	ProtocolIEIDResetType                                  int64 = 88
	ProtocolIEIDRoutingID                                  int64 = 89
	ProtocolIEIDRRCEstablishmentCause                      int64 = 90
	ProtocolIEIDRRCInactiveTransitionReportRequest         int64 = 91
	ProtocolIEIDRRCState                                   int64 = 92
	ProtocolIEIDSecurityContext                            int64 = 93
	ProtocolIEIDSecurityKey                                int64 = 94
	ProtocolIEIDSerialNumber                               int64 = 95
	ProtocolIEIDServedGUAMIList                            int64 = 96
	ProtocolIEIDSliceSupportList                           int64 = 97
	ProtocolIEIDSONConfigurationTransferDL                 int64 = 98
	ProtocolIEIDSONConfigurationTransferUL                 int64 = 99
	ProtocolIEIDSourceAMFUENGAPID                          int64 = 100
	ProtocolIEIDSourceToTargetTransparentContainer         int64 = 101
	ProtocolIEIDSupportedTAList                            int64 = 102
	ProtocolIEIDTAIListForPaging                           int64 = 103
	ProtocolIEIDTAIListForRestart                          int64 = 104
	ProtocolIEIDTargetID                                   int64 = 105
	ProtocolIEIDTargetToSourceTransparentContainer         int64 = 106
	ProtocolIEIDTimeToWait                                 int64 = 107
	ProtocolIEIDTraceActivation                            int64 = 108
	ProtocolIEIDTraceCollectionEntityIPAddress             int64 = 109
	ProtocolIEIDUEAggregateMaximumBitRate                  int64 = 110
	ProtocolIEIDUEAssociatedLogicalNGConnectionList        int64 = 111
	ProtocolIEIDUEContextRequest                           int64 = 112
	ProtocolIEIDUENGAPIDs                                  int64 = 114
	ProtocolIEIDUEPagingIdentity                           int64 = 115
	ProtocolIEIDUEPresenceInAreaOfInterestList             int64 = 116
	ProtocolIEIDUERadioCapability                          int64 = 117
	ProtocolIEIDUERadioCapabilityForPaging                 int64 = 118
	ProtocolIEIDUESecurityCapabilities                     int64 = 119
	ProtocolIEIDUnavailableGUAMIList                       int64 = 120
	ProtocolIEIDUserLocationInformation                    int64 = 121
	ProtocolIEIDWarningAreaList                            int64 = 122
	ProtocolIEIDWarningMessageContents                     int64 = 123
	ProtocolIEIDWarningSecurityInfo                        int64 = 124
	ProtocolIEIDWarningType                                int64 = 125
	ProtocolIEIDAdditionalULNGUUPTNLInformation            int64 = 126
	ProtocolIEIDDataForwardingNotPossible                  int64 = 127
	ProtocolIEIDDLNGUUPTNLInformation                      int64 = 128
	ProtocolIEIDNetworkInstance                            int64 = 129
	ProtocolIEIDPDUSessionAggregateMaximumBitRate          int64 = 130
	ProtocolIEIDPDUSessionResourceFailedToModifyListModCfm int64 = 131
	ProtocolIEIDPDUSessionResourceFailedToSetupListCxtFail int64 = 132
	ProtocolIEIDPDUSessionResourceListCxtRelReq            int64 = 133
	ProtocolIEIDPDUSessionType                             int64 = 134
	ProtocolIEIDQosFlowAddOrModifyRequestList              int64 = 135
	ProtocolIEIDQosFlowSetupRequestList                    int64 = 136
	ProtocolIEIDQosFlowToReleaseList                       int64 = 137
	ProtocolIEIDSecurityIndication                         int64 = 138
	ProtocolIEIDULNGUUPTNLInformation                      int64 = 139
	ProtocolIEIDULNGUUPTNLModifyList                       int64 = 140
	ProtocolIEIDWarningAreaCoordinates                     int64 = 141
	ProtocolIEIDPDUSessionResourceSecondaryRATUsageList    int64 = 142
	ProtocolIEIDHandoverFlag                               int64 = 143
	ProtocolIEIDSecondaryRATUsageInformation               int64 = 144
	ProtocolIEIDPDUSessionResourceReleaseResponseTransfer  int64 = 145
	ProtocolIEIDRedirectionVoiceFallback                   int64 = 146
	ProtocolIEIDUERetentionInformation                     int64 = 147
	ProtocolIEIDSNSSAI                                     int64 = 148
	ProtocolIEIDPSCellInformation                          int64 = 149
	ProtocolIEIDLastEUTRANPLMNIdentity                     int64 = 150
	ProtocolIEIDMaximumIntegrityProtectedDataRateDL        int64 = 151
	ProtocolIEIDAdditionalDLForwardingUPTNLInformation     int64 = 152
	ProtocolIEIDAdditionalDLUPTNLInformationForHOList      int64 = 153
	ProtocolIEIDAdditionalNGUUPTNLInformation              int64 = 154
	ProtocolIEIDAdditionalDLQosFlowPerTNLInformation       int64 = 155
	ProtocolIEIDSecurityResult                             int64 = 156
	ProtocolIEIDENDCSONConfigurationTransferDL             int64 = 157
	ProtocolIEIDENDCSONConfigurationTransferUL             int64 = 158
)
