package hilink

const (
	ErrorNoDevice                                             int64 = -2
	ErrorDefault                                              int64 = -1
	ErrorFirstSend                                            int64 = 1
	ErrorUSSDNetNoReturn                                      int64 = 11019
	ErrorUnknown                                              int64 = 100001
	ErrorNotSupport                                           int64 = 100002
	ErrorNoRight                                              int64 = 100003
	ErrorBusy                                                 int64 = 100004
	ErrorFormatError                                          int64 = 100005
	ErrorParameterError                                       int64 = 100006
	ErrorSaveConfigFileError                                  int64 = 100007
	ErrorGetConfigFileError                                   int64 = 100008
	ErrorNoSimCardOrInvalidSimCard                            int64 = 101001
	ErrorCheckSimCardPINLock                                  int64 = 101002
	ErrorCheckSimCardPUNLock                                  int64 = 101003
	ErrorCheckSimCardCanUnuseable                             int64 = 101004
	ErrorEnablePINFailed                                      int64 = 101005
	ErrorDisablePINFailed                                     int64 = 101006
	ErrorUnlockPINFailed                                      int64 = 101007
	ErrorDisableAutoPINFailed                                 int64 = 101008
	ErrorEnableAutoPINFailed                                  int64 = 101009
	ErrorGetNetTypeFailed                                     int64 = 102001
	ErrorGetServiceStatusFailed                               int64 = 102002
	ErrorGetRoamStatusFailed                                  int64 = 102003
	ErrorGetConnectStatusFailed                               int64 = 102004
	ErrorDeviceATExecuteFailed                                int64 = 103001
	ErrorDevicePINValidateFailed                              int64 = 103002
	ErrorDevicePINModifyFailed                                int64 = 103003
	ErrorDevicePUKModifyFailed                                int64 = 103004
	ErrorDeviceGetAutorunVersionFailed                        int64 = 103005
	ErrorDeviceGetAPIVersionFailed                            int64 = 103006
	ErrorDeviceGetProductInformatonFailed                     int64 = 103007
	ErrorDeviceSimCardBusy                                    int64 = 103008
	ErrorDeviceSimLockInputError                              int64 = 103009
	ErrorDeviceNotSupportRemoteOperate                        int64 = 103010
	ErrorDevicePUKDeadLock                                    int64 = 103011
	ErrorDeviceGetPCAissstInformationFailed                   int64 = 103012
	ErrorDeviceSetLogInformatonLevelFailed                    int64 = 103013
	ErrorDeviceGetLogInformatonLevelFailed                    int64 = 103014
	ErrorDeviceCompressLogFileFailed                          int64 = 103015
	ErrorDeviceRestoreFileDecryptFailed                       int64 = 103016
	ErrorDeviceRestoreFileVersionMatchFailed                  int64 = 103017
	ErrorDeviceRestoreFileFailed                              int64 = 103018
	ErrorDeviceSetTimeFailed                                  int64 = 103101
	ErrorCompressLogFileFailed                                int64 = 103102
	ErrorDHCPError                                            int64 = 104001
	ErrorSafeError                                            int64 = 106001
	ErrorDialupGetConnectFileError                            int64 = 107720
	ErrorDialupSetConnectFileError                            int64 = 107721
	ErrorDialupDialupManagmentParseError                      int64 = 107722
	ErrorDialupAddProfileError                                int64 = 107724
	ErrorDialupModifyProfileError                             int64 = 107725
	ErrorDialupSetDefaultProfileError                         int64 = 107726
	ErrorDialupGetProfileListError                            int64 = 107727
	ErrorDialupGetAutoAPNMatchError                           int64 = 107728
	ErrorDialupSetAutoAPNMatchError                           int64 = 107729
	ErrorLoginUserNoExist                                     int64 = 108001
	ErrorLoginPasswordError                                   int64 = 108002
	ErrorLoginAlreadyLogined                                  int64 = 108003
	ErrorLoginModifyPasswordFailed                            int64 = 108004
	ErrorLoginTooManyUsersLogined                             int64 = 108005
	ErrorLoginUsernameOrPasswordError                         int64 = 108006
	ErrorLoginTooManyTimes                                    int64 = 108007
	ErrorGetLanguageFailed                                    int64 = 109001
	ErrorSetLanguageFailed                                    int64 = 109002
	ErrorOnlineUpdateServerNotAccessed                        int64 = 110001
	ErrorOnlineUpdateAlreadyBooted                            int64 = 110002
	ErrorOnlineUpdateGetDeviceInformationFailed               int64 = 110003
	ErrorOnlineUpdateGetLocalGroupCommponentInformationFailed int64 = 110004
	ErrorOnlineUpdateNotFindFileOnServer                      int64 = 110005
	ErrorOnlineUpdateNeedReconnectServer                      int64 = 110006
	ErrorOnlineUpdateCancelDownloding                         int64 = 110007
	ErrorOnlineUpdateSameFileList                             int64 = 110008
	ErrorOnlineUpdateConnectError                             int64 = 110009
	ErrorOnlineUpdateInvalidURLList                           int64 = 110021
	ErrorOnlineUpdateNotSupportURLList                        int64 = 110022
	ErrorOnlineUpdateNotBoot                                  int64 = 110023
	ErrorOnlineUpdateLowBattery                               int64 = 110024
	ErrorUSSDError                                            int64 = 111001
	ErrorUSSDFucntionReturnError                              int64 = 111012
	ErrorUSSDInUSSDSession                                    int64 = 111013
	ErrorUSSDTooLongContent                                   int64 = 111014
	ErrorUSSDEmptyCommand                                     int64 = 111016
	ErrorUSSDCodingError                                      int64 = 111017
	ErrorUSSDATSendFailed                                     int64 = 111018
	ErrorUSSDNetOvertime                                      int64 = 111020
	ErrorUSSDXMLSpecialCharacterTransferFailed                int64 = 111021
	ErrorUSSDNetNotSupportUSSD                                int64 = 111022
	ErrorSetNetModeAndBandWhenDailupFailed                    int64 = 112001
	ErrorSetNetSearchModeWhenDailupFailed                     int64 = 112002
	ErrorSetNetModeAndBandFailed                              int64 = 112003
	ErrorSetNetSearchModeFailed                               int64 = 112004
	ErrorNetRegisterNetFailed                                 int64 = 112005
	ErrorNetConnectedOrderNotMatch                            int64 = 112006
	ErrorNetCurrentNetModeNotSupport                          int64 = 112007
	ErrorNetSimCardNotReadyStatus                             int64 = 112008
	ErrorNetMemoryAllocFailed                                 int64 = 112009
	ErrorSMSNullArgumentOrIllegalArgument                     int64 = 113017
	ErrorSMSOvertime                                          int64 = 113018
	ErrorSMSQuerySMSIndexListError                            int64 = 113020
	ErrorSMSSetSMSCenterNumberFailed                          int64 = 113031
	ErrorSMSDeleteSMSFailed                                   int64 = 113036
	ErrorSMSSaveConfigFileFailed                              int64 = 113047
	ErrorSMSLocalSpaceNotEnough                               int64 = 113053
	ErrorSMSTelephoneNumberTooLong                            int64 = 113054
	ErrorSDFileExist                                          int64 = 114001
	ErrorSDDirectoryExist                                     int64 = 114002
	ErrorSDFileOrDirectoryNotExist                            int64 = 114004
	ErrorSDFileNameTooLong                                    int64 = 114005
	ErrorSDNoRight                                            int64 = 114006
	ErrorSDFileIsUploading                                    int64 = 114007
	ErrorPbNullArgumentOrIllegalArgument                      int64 = 115001
	ErrorPbOvertime                                           int64 = 115002
	ErrorPbCallSystemFucntionError                            int64 = 115003
	ErrorPbWriteFileError                                     int64 = 115004
	ErrorPbReadFileError                                      int64 = 115005
	ErrorPbLocalTelephoneFullError                            int64 = 115199
	ErrorSTKNullArgumentOrIllegalArgument                     int64 = 116001
	ErrorSTKOvertime                                          int64 = 116002
	ErrorSTKCallSystemFucntionError                           int64 = 116003
	ErrorSTKWriteFileError                                    int64 = 116004
	ErrorSTKReadFileError                                     int64 = 116005
	ErrorWiFiStationConnectAPPasswordError                    int64 = 117001
	ErrorWiFiWebPasswordOrDHCPOvertimeError                   int64 = 117002
	ErrorWiFiPbcConnectFailed                                 int64 = 117003
	ErrorWiFiStationConnectAPWISPRPasswordError               int64 = 117004
	ErrorCradleGetCurrentConnectedUserIPFailed                int64 = 118001
	ErrorCradleGetCurrentConnectedUserMACFailed               int64 = 118002
	ErrorCradleSetMACFailed                                   int64 = 118003
	ErrorCradleGetWANInformationFailed                        int64 = 118004
	ErrorCradleCodingFailed                                   int64 = 118005
	ErrorCradleUpdateProfileFailed                            int64 = 118006
	ErrorTokenWrong                                           int64 = 125001
	ErrorSessionWrong                                         int64 = 125002
	ErrorSessionTokenWrong                                    int64 = 125003
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
