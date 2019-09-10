package hilink

const (
	ErrorBusy                                                 int64 = 100004
	ErrorCheckSimCardCanUnuseable                             int64 = 101004
	ErrorCheckSimCardPINLock                                  int64 = 101002
	ErrorCheckSimCardPUNLock                                  int64 = 101003
	ErrorCompressLogFileFailed                                int64 = 103102
	ErrorCradleCodingFailed                                   int64 = 118005
	ErrorCradleGetCurrentConnectedUserIPFailed                int64 = 118001
	ErrorCradleGetCurrentConnectedUserMACFailed               int64 = 118002
	ErrorCradleGetWANInformationFailed                        int64 = 118004
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
	ErrorDevicePINModifyFailed                                int64 = 103003
	ErrorDevicePINValidateFailed                              int64 = 103002
	ErrorDevicePUKDeadLock                                    int64 = 103011
	ErrorDevicePUKModifyFailed                                int64 = 103004
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
	ErrorEnablePINFailed                                      int64 = 101005
	ErrorFirstSend                                            int64 = 1
	ErrorFormatError                                          int64 = 100005
	ErrorGetConfigFileError                                   int64 = 100008
	ErrorGetConnectStatusFailed                               int64 = 102004
	ErrorGetNetTypeFailed                                     int64 = 102001
	ErrorGetRoamStatusFailed                                  int64 = 102003
	ErrorGetServiceStatusFailed                               int64 = 102002
	ErrorGetLanguageFailed                                    int64 = 109001
	ErrorSetLanguageFailed                                    int64 = 109002
	ErrorLoginAlreadyLogined                                  int64 = 108003
	ErrorLoginModifyPasswordFailed                            int64 = 108004
	ErrorLoginUserNoExist                                     int64 = 108001
	ErrorLoginPasswordError                                   int64 = 108002
	ErrorLoginTooManyTimes                                    int64 = 108007
	ErrorLoginTooManyUsersLogined                             int64 = 108005
	ErrorLoginUsernameOrPasswordError                         int64 = 108006
	ErrorNetCurrentNetModeNotSupport                          int64 = 112007
	ErrorNetMemoryAllocFailed                                 int64 = 112009
	ErrorNetConnectedOrderNotMatch                            int64 = 112006
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
	ErrorBusy:                                                 "busy",
	ErrorCheckSimCardCanUnuseable:                             "sim card can unuseable",
	ErrorCheckSimCardPINLock:                                  "sim card PIN lock",
	ErrorCheckSimCardPUNLock:                                  "sim card PUN lock",
	ErrorCompressLogFileFailed:                                "compress log file failed",
	ErrorCradleCodingFailed:                                   "cradle coding failed",
	ErrorCradleGetCurrentConnectedUserIPFailed:                "get current connected user IP failed",
	ErrorCradleGetCurrentConnectedUserMACFailed:               "get current connected user MAC failed",
	ErrorCradleGetWANInformationFailed:                        "get WAN information failed",
	ErrorCradleSetMACFailed:                                   "get MAC failed",
	ErrorCradleUpdateProfileFailed:                            "update profile failed",
	ErrorDefault:                                              "default error",
	ErrorDeviceATExecuteFailed:                                "AT execute failed",
	ErrorDeviceCompressLogFileFailed:                          "compress log file failed",
	ErrorDeviceGetAPIVersionFailed:                            "get API version failed",
	ErrorDeviceGetAutorunVersionFailed:                        "get autorun version failed",
	ErrorDeviceGetLogInformatonLevelFailed:                    "get log information level failed",
	ErrorDeviceGetPCAissstInformationFailed:                   "get PC aisst information failed",
	ErrorDeviceGetProductInformatonFailed:                     "get product information failed",
	ErrorDeviceNotSupportRemoteOperate:                        "device not support remote operate",
	ErrorDevicePINModifyFailed:                                "PIN modify failed",
	ErrorDevicePINValidateFailed:                              "PIN validate failed",
	ErrorDevicePUKDeadLock:                                    "PUK dead lock",
	ErrorDevicePUKModifyFailed:                                "PUK modify failed",
	ErrorDeviceRestoreFileDecryptFailed:                       "restore file decrypt failed",
	ErrorDeviceRestoreFileFailed:                              "restore file failed",
	ErrorDeviceRestoreFileVersionMatchFailed:                  "restore file version match failed",
	ErrorDeviceSetLogInformatonLevelFailed:                    "set log information level failed",
	ErrorDeviceSetTimeFailed:                                  "set time failed",
	ErrorDeviceSimCardBusy:                                    "sim card busy",
	ErrorDeviceSimLockInputError:                              "sim lock input error",
	ErrorDHCPError:                                            "DHCP error",
	ErrorDialupAddProfileError:                                "add profile error",
	ErrorDialupDialupManagmentParseError:                      "managment parse error",
	ErrorDialupGetAutoAPNMatchError:                           "get auto APN match error",
	ErrorDialupGetConnectFileError:                            "get connect file error",
	ErrorDialupGetProfileListError:                            "get profile list error",
	ErrorDialupModifyProfileError:                             "modify profile error",
	ErrorDialupSetAutoAPNMatchError:                           "set auto APN match error",
	ErrorDialupSetConnectFileError:                            "set connect file error",
	ErrorDialupSetDefaultProfileError:                         "set default profile error",
	ErrorDisableAutoPINFailed:                                 "disable auto PIN failed",
	ErrorDisablePINFailed:                                     "disable PIN failed",
	ErrorEnableAutoPINFailed:                                  "enable auto PIN failed",
	ErrorEnablePINFailed:                                      "enable PIN failed",
	ErrorFirstSend:                                            "first send error",
	ErrorFormatError:                                          "format error",
	ErrorGetConfigFileError:                                   "get config file error",
	ErrorGetConnectStatusFailed:                               "get connect status failed",
	ErrorGetNetTypeFailed:                                     "get net type failed",
	ErrorGetRoamStatusFailed:                                  "get roam status failed",
	ErrorGetServiceStatusFailed:                               "get service status failed",
	ErrorGetLanguageFailed:                                    "get language failed",
	ErrorSetLanguageFailed:                                    "set language failed",
	ErrorLoginAlreadyLogined:                                  "already logined",
	ErrorLoginModifyPasswordFailed:                            "modify password failed",
	ErrorLoginUserNoExist:                                     "user no exist",
	ErrorLoginPasswordError:                                   "password error",
	ErrorLoginTooManyTimes:                                    "too many times",
	ErrorLoginTooManyUsersLogined:                             "too many users logined",
	ErrorLoginUsernameOrPasswordError:                         "username or password error",
	ErrorNetCurrentNetModeNotSupport:                          "current net mode not support",
	ErrorNetMemoryAllocFailed:                                 "memory alloc failed",
	ErrorNetConnectedOrderNotMatch:                            "connected order not match",
	ErrorNetRegisterNetFailed:                                 "register net failed",
	ErrorNetSimCardNotReadyStatus:                             "sim card not ready status",
	ErrorNotSupport:                                           "not support",
	ErrorNoDevice:                                             "no device",
	ErrorNoRight:                                              "no right",
	ErrorNoSimCardOrInvalidSimCard:                            "no sim card or invalid sim card",
	ErrorOnlineUpdateAlreadyBooted:                            "already booted",
	ErrorOnlineUpdateCancelDownloding:                         "cancel downloading",
	ErrorOnlineUpdateConnectError:                             "connect error",
	ErrorOnlineUpdateGetDeviceInformationFailed:               "get device information failed",
	ErrorOnlineUpdateGetLocalGroupCommponentInformationFailed: "get local group component information failed",
	ErrorOnlineUpdateInvalidURLList:                           "invalid URL list",
	ErrorOnlineUpdateLowBattery:                               "low battery",
	ErrorOnlineUpdateNeedReconnectServer:                      "need reconnect server",
	ErrorOnlineUpdateNotBoot:                                  "not boot",
	ErrorOnlineUpdateNotFindFileOnServer:                      "not find file on server",
	ErrorOnlineUpdateNotSupportURLList:                        "not support URL list",
	ErrorOnlineUpdateSameFileList:                             "same file list",
	ErrorOnlineUpdateServerNotAccessed:                        "server not accessed",
	ErrorParameterError:                                       "parameter error",
	ErrorPbCallSystemFucntionError:                            "call system function error",
	ErrorPbLocalTelephoneFullError:                            "local telephone full error",
	ErrorPbNullArgumentOrIllegalArgument:                      "null argument or illegal argument",
	ErrorPbOvertime:                                           "overtime",
	ErrorPbReadFileError:                                      "read file error",
	ErrorPbWriteFileError:                                     "write file error",
	ErrorSafeError:                                            "safe error",
	ErrorSaveConfigFileError:                                  "save config file error",
	ErrorSDDirectoryExist:                                     "directory exist",
	ErrorSDFileExist:                                          "file exist",
	ErrorSDFileIsUploading:                                    "file is uploading",
	ErrorSDFileNameTooLong:                                    "file name too long",
	ErrorSDFileOrDirectoryNotExist:                            "file or directory not exist",
	ErrorSDNoRight:                                            "no right",
	ErrorSetNetModeAndBandFailed:                              "set net mode and band failed",
	ErrorSetNetModeAndBandWhenDailupFailed:                    "set net mode and band when dailup failed",
	ErrorSetNetSearchModeFailed:                               "set net search mode failed",
	ErrorSetNetSearchModeWhenDailupFailed:                     "set net search mode when dailup failed",
	ErrorSMSDeleteSMSFailed:                                   "delete SMS failed",
	ErrorSMSLocalSpaceNotEnough:                               "local space not enough",
	ErrorSMSNullArgumentOrIllegalArgument:                     "null argument or illegal argument",
	ErrorSMSOvertime:                                          "overtime",
	ErrorSMSQuerySMSIndexListError:                            "query SMS index list error",
	ErrorSMSSaveConfigFileFailed:                              "save config file failed",
	ErrorSMSSetSMSCenterNumberFailed:                          "set SMS center number failed",
	ErrorSMSTelephoneNumberTooLong:                            "telephone number too long",
	ErrorSTKCallSystemFucntionError:                           "call system function error",
	ErrorSTKNullArgumentOrIllegalArgument:                     "null argument or illegal argument",
	ErrorSTKOvertime:                                          "overtime",
	ErrorSTKReadFileError:                                     "read file error",
	ErrorSTKWriteFileError:                                    "write file error",
	ErrorUnknown:                                              "unknown",
	ErrorUnlockPINFailed:                                      "unlock PIN failed",
	ErrorUSSDATSendFailed:                                     "AT send failed",
	ErrorUSSDCodingError:                                      "coding error",
	ErrorUSSDEmptyCommand:                                     "empty USSD command",
	ErrorUSSDError:                                            "USSD error",
	ErrorUSSDFucntionReturnError:                              "function return error",
	ErrorUSSDInUSSDSession:                                    "in USSD session",
	ErrorUSSDNetNotSupportUSSD:                                "net not support USSD",
	ErrorUSSDNetNoReturn:                                      "net no return",
	ErrorUSSDNetOvertime:                                      "net overtime",
	ErrorUSSDTooLongContent:                                   "to long USSD content",
	ErrorUSSDXMLSpecialCharacterTransferFailed:                "XML special character transfer failed",
	ErrorWiFiPbcConnectFailed:                                 "pbc connect failed",
	ErrorWiFiStationConnectAPPasswordError:                    "station connect AP password error",
	ErrorWiFiStationConnectAPWISPRPasswordError:               "station connect AP WISPRP password error",
	ErrorWiFiWebPasswordOrDHCPOvertimeError:                   "web password or DHCP overtime error",
}

func ErrorMessage(code int64) string {
	if text, ok := errorText[code]; ok {
		return text
	}

	return ""
}
