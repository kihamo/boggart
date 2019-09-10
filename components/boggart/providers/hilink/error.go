package hilink

const (
	ErrorBusy                                                 int64 = 100004
	ErrorCheckSimCardCanUnuseable                             int64 = 101004
	ErrorCheckSimCardPINLock                                  int64 = 101002
	ErrorCheckSimCardPUNLock                                  int64 = 101003
	ErrorCompressLogFileFailed                                int64 = 103102
	ErrorCradleCodingFailed                                   int64 = 118005
	ErrorCradleGetCrurrentConnectedUserIPFailed               int64 = 118001
	ErrorCradleGetCrurrentConnectedUserMACFailed              int64 = 118002
	ErrorCradleGetWanInformationFailed                        int64 = 118004
	ErrorCradleSetMACFailed                                   int64 = 118003
	ErrorCradleUpdateProfileFailed                            int64 = 118006
	ErrorDefault                                              int64 = -1
	ErrorDeviceATExecuteFailed                                int64 = 103001
	ErrorDeviceCompressLogFileFailed                          int64 = 103015
	ErrorDeviceGetAPIVersionFailed                            int64 = 103006
	ErrorDeviceGetAutorunVersionFailed                        int64 = 103005
	ErrorDeviceGetLogInformatonLevelFailed                    int64 = 103014
	ErrorDeviceGetPCAissstInformationFailed                   int64 = 103012
	ErrorDeviceGetProductInformatonFailed                     int64 = 103007
	ErrorDeviceNotSupportRemoteOperate                        int64 = 103010
	ErrorDevicePINModiffyFailed                               int64 = 103003
	ErrorDevicePINValidateFailed                              int64 = 103002
	ErrorDevicePUKDeadLock                                    int64 = 103011
	ErrorDevicePUKModiffyFailed                               int64 = 103004
	ErrorDeviceRestoreFileDecryptFailed                       int64 = 103016
	ErrorDeviceRestoreFileFailed                              int64 = 103018
	ErrorDeviceRestoreFileVersionMatchFailed                  int64 = 103017
	ErrorDeviceSetLogInformatonLevelFailed                    int64 = 103013
	ErrorDeviceSetTimeFailed                                  int64 = 103101
	ErrorDeviceSimCardBusy                                    int64 = 103008
	ErrorDeviceSimLockInputError                              int64 = 103009
	ErrorDHCPError                                            int64 = 104001
	ErrorDialupAddProfileError                                int64 = 107724
	ErrorDialupDialupManagmentParseError                      int64 = 107722
	ErrorDialupGetAutoAPNMatchError                           int64 = 107728
	ErrorDialupGetConnectFileError                            int64 = 107720
	ErrorDialupGetProfileListError                            int64 = 107727
	ErrorDialupModifyProfileError                             int64 = 107725
	ErrorDialupSetAutoAPNMatchError                           int64 = 107729
	ErrorDialupSetConnectFileError                            int64 = 107721
	ErrorDialupSetDefaultProfileError                         int64 = 107726
	ErrorDisableAutoPINFailed                                 int64 = 101008
	ErrorDisablePINFailed                                     int64 = 101006
	ErrorEnableAutoPINFailed                                  int64 = 101009
	ErrorEnablePinFailed                                      int64 = 101005
	ErrorFirstSend                                            int64 = 1
	ErrorFormatError                                          int64 = 100005
	ErrorGetConfigFileError                                   int64 = 100008
	ErrorGetConnectStatusFailed                               int64 = 102004
	ErrorGetNetTypeFailed                                     int64 = 102001
	ErrorGetRoamStatusFailed                                  int64 = 102003
	ErrorGetServiceStatusFailed                               int64 = 102002
	ErrorLanguageGetFailed                                    int64 = 109001
	ErrorLanguageSetFailed                                    int64 = 109002
	ErrorLoginAlreadyLogined                                  int64 = 108003
	ErrorLoginModifyPasswordFailed                            int64 = 108004
	ErrorLoginNoExistUser                                     int64 = 108001
	ErrorLoginPasswordError                                   int64 = 108002
	ErrorLoginTooManyTimes                                    int64 = 108007
	ErrorLoginTooManyUsersLogined                             int64 = 108005
	ErrorLoginUsernameOrPasswordError                         int64 = 108006
	ErrorNetCurrentNetModeNotSupport                          int64 = 112007
	ErrorNetMemoryAllocFailed                                 int64 = 112009
	ErrorNetNetConnectedOrderNotMatch                         int64 = 112006
	ErrorNetRegisterNetFailed                                 int64 = 112005
	ErrorNetSimCardNotReadyStatus                             int64 = 112008
	ErrorNotSupport                                           int64 = 100002
	ErrorNoDevice                                             int64 = -2
	ErrorNoRight                                              int64 = 100003
	ErrorNoSimCardOrInvalidSimCard                            int64 = 101001
	ErrorOnlineUpdateAlreadyBooted                            int64 = 110002
	ErrorOnlineUpdateCancelDownloding                         int64 = 110007
	ErrorOnlineUpdateConnectError                             int64 = 110009
	ErrorOnlineUpdateGetDeviceInformationFailed               int64 = 110003
	ErrorOnlineUpdateGetLocalGroupCommponentInformationFailed int64 = 110004
	ErrorOnlineUpdateInvalidURLList                           int64 = 110021
	ErrorOnlineUpdateLowBattery                               int64 = 110024
	ErrorOnlineUpdateNeedReconnectServer                      int64 = 110006
	ErrorOnlineUpdateNotBoot                                  int64 = 110023
	ErrorOnlineUpdateNotFindFileOnServer                      int64 = 110005
	ErrorOnlineUpdateNotSupportURLList                        int64 = 110022
	ErrorOnlineUpdateSameFileList                             int64 = 110008
	ErrorOnlineUpdateServerNotAccessed                        int64 = 110001
	ErrorParameterError                                       int64 = 100006
	ErrorPbCallSystemFucntionError                            int64 = 115003
	ErrorPbLocalTelephoneFullError                            int64 = 115199
	ErrorPbNullArgumentOrIllegalArgument                      int64 = 115001
	ErrorPbOvertime                                           int64 = 115002
	ErrorPbReadFileError                                      int64 = 115005
	ErrorPbWriteFileError                                     int64 = 115004
	ErrorSafeError                                            int64 = 106001
	ErrorSaveConfigFileError                                  int64 = 100007
	ErrorSDDirectoryExist                                     int64 = 114002
	ErrorSDFileExist                                          int64 = 114001
	ErrorSDFileIsUploading                                    int64 = 114007
	ErrorSDFileNameTooLong                                    int64 = 114005
	ErrorSDFileOrDirectoryNotExist                            int64 = 114004
	ErrorSDNoRight                                            int64 = 114006
	ErrorSetNetModeAndBandFailed                              int64 = 112003
	ErrorSetNetModeAndBandWhenDailupFailed                    int64 = 112001
	ErrorSetNetSearchModeFailed                               int64 = 112004
	ErrorSetNetSearchModeWhenDailupFailed                     int64 = 112002
	ErrorSMSDeleteSMSFailed                                   int64 = 113036
	ErrorSMSLocalSpaceNotEnough                               int64 = 113053
	ErrorSMSNullArgumentOrIllegalArgument                     int64 = 113017
	ErrorSMSOvertime                                          int64 = 113018
	ErrorSMSQuerySMSIndexListError                            int64 = 113020
	ErrorSMSSaveConfigFileFailed                              int64 = 113047
	ErrorSMSSetSMSCenterNumberFailed                          int64 = 113031
	ErrorSMSTelephoneNumberTooLong                            int64 = 113054
	ErrorSTKCallSystemFucntionError                           int64 = 116003
	ErrorSTKNullArgumentOrIllegalArgument                     int64 = 116001
	ErrorSTKOvertime                                          int64 = 116002
	ErrorSTKReadFileError                                     int64 = 116005
	ErrorSTKWriteFileError                                    int64 = 116004
	ErrorUnknown                                              int64 = 100001
	ErrorUnlockPINFailed                                      int64 = 101007
	ErrorUSSDATSendFailed                                     int64 = 111018
	ErrorUSSDCodingError                                      int64 = 111017
	ErrorUSSDEmptyCommand                                     int64 = 111016
	ErrorUSSDError                                            int64 = 111001
	ErrorUSSDFucntionReturnError                              int64 = 111012
	ErrorUSSDInUSSDSession                                    int64 = 111013
	ErrorUSSDNetNotSupportUSSD                                int64 = 111022
	ErrorUSSDNetNoReturn                                      int64 = 11019
	ErrorUSSDNetOvertime                                      int64 = 111020
	ErrorUSSDTooLongContent                                   int64 = 111014
	ErrorUSSDXMLSpecialCharacterTransferFailed                int64 = 111021
	ErrorWiFiPbcConnectFailed                                 int64 = 117003
	ErrorWiFiStationConnectAPPasswordError                    int64 = 117001
	ErrorWiFiStationConnectAPWISPRPasswordError               int64 = 117004
	ErrorWiFiWebPasswordOrDHCPOvertimeError                   int64 = 117002
)

var errorText = map[int64]string{
	ErrorBusy:                                                 "",
	ErrorCheckSimCardCanUnuseable:                             "",
	ErrorCheckSimCardPINLock:                                  "",
	ErrorCheckSimCardPUNLock:                                  "",
	ErrorCompressLogFileFailed:                                "",
	ErrorCradleCodingFailed:                                   "",
	ErrorCradleGetCrurrentConnectedUserIPFailed:               "",
	ErrorCradleGetCrurrentConnectedUserMACFailed:              "",
	ErrorCradleGetWanInformationFailed:                        "",
	ErrorCradleSetMACFailed:                                   "",
	ErrorCradleUpdateProfileFailed:                            "",
	ErrorDefault:                                              "",
	ErrorDeviceATExecuteFailed:                                "",
	ErrorDeviceCompressLogFileFailed:                          "",
	ErrorDeviceGetAPIVersionFailed:                            "",
	ErrorDeviceGetAutorunVersionFailed:                        "",
	ErrorDeviceGetLogInformatonLevelFailed:                    "",
	ErrorDeviceGetPCAissstInformationFailed:                   "",
	ErrorDeviceGetProductInformatonFailed:                     "",
	ErrorDeviceNotSupportRemoteOperate:                        "",
	ErrorDevicePINModiffyFailed:                               "",
	ErrorDevicePINValidateFailed:                              "",
	ErrorDevicePUKDeadLock:                                    "",
	ErrorDevicePUKModiffyFailed:                               "",
	ErrorDeviceRestoreFileDecryptFailed:                       "",
	ErrorDeviceRestoreFileFailed:                              "",
	ErrorDeviceRestoreFileVersionMatchFailed:                  "",
	ErrorDeviceSetLogInformatonLevelFailed:                    "",
	ErrorDeviceSetTimeFailed:                                  "",
	ErrorDeviceSimCardBusy:                                    "",
	ErrorDeviceSimLockInputError:                              "",
	ErrorDHCPError:                                            "",
	ErrorDialupAddProfileError:                                "",
	ErrorDialupDialupManagmentParseError:                      "",
	ErrorDialupGetAutoAPNMatchError:                           "",
	ErrorDialupGetConnectFileError:                            "",
	ErrorDialupGetProfileListError:                            "",
	ErrorDialupModifyProfileError:                             "",
	ErrorDialupSetAutoAPNMatchError:                           "",
	ErrorDialupSetConnectFileError:                            "",
	ErrorDialupSetDefaultProfileError:                         "",
	ErrorDisableAutoPINFailed:                                 "",
	ErrorDisablePINFailed:                                     "",
	ErrorEnableAutoPINFailed:                                  "",
	ErrorEnablePinFailed:                                      "",
	ErrorFirstSend:                                            "",
	ErrorFormatError:                                          "",
	ErrorGetConfigFileError:                                   "",
	ErrorGetConnectStatusFailed:                               "",
	ErrorGetNetTypeFailed:                                     "",
	ErrorGetRoamStatusFailed:                                  "",
	ErrorGetServiceStatusFailed:                               "",
	ErrorLanguageGetFailed:                                    "",
	ErrorLanguageSetFailed:                                    "",
	ErrorLoginAlreadyLogined:                                  "",
	ErrorLoginModifyPasswordFailed:                            "",
	ErrorLoginNoExistUser:                                     "",
	ErrorLoginPasswordError:                                   "",
	ErrorLoginTooManyTimes:                                    "",
	ErrorLoginTooManyUsersLogined:                             "",
	ErrorLoginUsernameOrPasswordError:                         "",
	ErrorNetCurrentNetModeNotSupport:                          "",
	ErrorNetMemoryAllocFailed:                                 "",
	ErrorNetNetConnectedOrderNotMatch:                         "",
	ErrorNetRegisterNetFailed:                                 "",
	ErrorNetSimCardNotReadyStatus:                             "",
	ErrorNotSupport:                                           "",
	ErrorNoDevice:                                             "",
	ErrorNoRight:                                              "",
	ErrorNoSimCardOrInvalidSimCard:                            "",
	ErrorOnlineUpdateAlreadyBooted:                            "",
	ErrorOnlineUpdateCancelDownloding:                         "",
	ErrorOnlineUpdateConnectError:                             "",
	ErrorOnlineUpdateGetDeviceInformationFailed:               "",
	ErrorOnlineUpdateGetLocalGroupCommponentInformationFailed: "",
	ErrorOnlineUpdateInvalidURLList:                           "",
	ErrorOnlineUpdateLowBattery:                               "",
	ErrorOnlineUpdateNeedReconnectServer:                      "",
	ErrorOnlineUpdateNotBoot:                                  "",
	ErrorOnlineUpdateNotFindFileOnServer:                      "",
	ErrorOnlineUpdateNotSupportURLList:                        "",
	ErrorOnlineUpdateSameFileList:                             "",
	ErrorOnlineUpdateServerNotAccessed:                        "",
	ErrorParameterError:                                       "",
	ErrorPbCallSystemFucntionError:                            "",
	ErrorPbLocalTelephoneFullError:                            "",
	ErrorPbNullArgumentOrIllegalArgument:                      "",
	ErrorPbOvertime:                                           "",
	ErrorPbReadFileError:                                      "",
	ErrorPbWriteFileError:                                     "",
	ErrorSafeError:                                            "",
	ErrorSaveConfigFileError:                                  "",
	ErrorSDDirectoryExist:                                     "",
	ErrorSDFileExist:                                          "",
	ErrorSDFileIsUploading:                                    "",
	ErrorSDFileNameTooLong:                                    "",
	ErrorSDFileOrDirectoryNotExist:                            "",
	ErrorSDNoRight:                                            "",
	ErrorSetNetModeAndBandFailed:                              "",
	ErrorSetNetModeAndBandWhenDailupFailed:                    "",
	ErrorSetNetSearchModeFailed:                               "",
	ErrorSetNetSearchModeWhenDailupFailed:                     "",
	ErrorSMSDeleteSMSFailed:                                   "",
	ErrorSMSLocalSpaceNotEnough:                               "",
	ErrorSMSNullArgumentOrIllegalArgument:                     "",
	ErrorSMSOvertime:                                          "",
	ErrorSMSQuerySMSIndexListError:                            "",
	ErrorSMSSaveConfigFileFailed:                              "",
	ErrorSMSSetSMSCenterNumberFailed:                          "",
	ErrorSMSTelephoneNumberTooLong:                            "",
	ErrorSTKCallSystemFucntionError:                           "",
	ErrorSTKNullArgumentOrIllegalArgument:                     "",
	ErrorSTKOvertime:                                          "",
	ErrorSTKReadFileError:                                     "",
	ErrorSTKWriteFileError:                                    "",
	ErrorUnknown:                                              "",
	ErrorUnlockPINFailed:                                      "",
	ErrorUSSDATSendFailed:                                     "",
	ErrorUSSDCodingError:                                      "",
	ErrorUSSDEmptyCommand:                                     "",
	ErrorUSSDError:                                            "",
	ErrorUSSDFucntionReturnError:                              "",
	ErrorUSSDInUSSDSession:                                    "",
	ErrorUSSDNetNotSupportUSSD:                                "",
	ErrorUSSDNetNoReturn:                                      "",
	ErrorUSSDNetOvertime:                                      "",
	ErrorUSSDTooLongContent:                                   "",
	ErrorUSSDXMLSpecialCharacterTransferFailed:                "",
	ErrorWiFiPbcConnectFailed:                                 "",
	ErrorWiFiStationConnectAPPasswordError:                    "",
	ErrorWiFiStationConnectAPWISPRPasswordError:               "",
	ErrorWiFiWebPasswordOrDHCPOvertimeError:                   "",
}

func ErrorMessage(code int64) string {
	if text, ok := errorText[code]; ok {
		return text
	}

	return ""
}
