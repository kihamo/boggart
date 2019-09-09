package hilink

const (
	ErrorBusy                                                        int64 = 100004
	ErrorCheckSimCardCanUnuseable                                    int64 = 101004
	ErrorCheckSimCardPINLock                                         int64 = 101002
	ErrorCheckSimCardPUNLock                                         int64 = 101003
	ErrorCompressLogFileFailed                                       int64 = 103102
	ErrorCradleCodingFailed                                          int64 = 118005
	ErrorCradleGetCrurrentConnectedUserIPFailed                      int64 = 118001
	ErrorCradleGetCrurrentConnectedUserMACFailed                     int64 = 118002
	errorCradleGetWanInformationFailed                               int64 = 118004
	errorCradleSetMACFailed                                          int64 = 118003
	errorCradleUpdateProfileFailed                                   int64 = 118006
	errorDefault                                                     int64 = -1
	errorDeviceATExecuteFailed                                       int64 = 103001
	ErrorDeviceCompressLogFileFailed                                 int64 = 103015
	ErrorDeviceGetAPIVersionFailed                                   int64 = 103006
	ErrorDeviceGetAutorunVersionFailed                               int64 = 103005
	ErrorDeviceGetLogInformatonLevelFailed                           int64 = 103014
	ErrorDeviceGetPCAissstInformationFailed                          int64 = 103012
	ErrorDeviceGetProductInformatonFailed                            int64 = 103007
	ErrorDeviceNotSupportRemoteOperate                               int64 = 103010
	ErrorDevicePINModiffyFailed                                      int64 = 103003
	ErrorDevicePINValidateFailed                                     int64 = 103002
	ErrorDevicePUKDeadLock                                           int64 = 103011
	ErrorDevicePUKModiffyFailed                                      int64 = 103004
	ErrorDevice_restore_file_decrypt_failed                          int64 = 103016
	ErrorDevice_restore_file_failed                                  int64 = 103018
	ErrorDevice_restore_file_version_match_failed                    int64 = 103017
	ErrorDevice_set_log_informaton_level_failed                      int64 = 103013
	ErrorDevice_set_time_failed                                      int64 = 103101
	ErrorDevice_sim_card_busy                                        int64 = 103008
	ErrorDevice_sim_lock_input_error                                 int64 = 103009
	ErrorDHCP_error                                                  int64 = 104001
	ErrorDialup_add_prorile_error                                    int64 = 107724
	ErrorDialup_dialup_managment_parse_error                         int64 = 107722
	ErrorDialup_get_auto_apn_match_error                             int64 = 107728
	ErrorDialup_get_connect_file_error                               int64 = 107720
	ErrorDialup_get_prorile_list_error                               int64 = 107727
	ErrorDialup_modify_prorile_error                                 int64 = 107725
	ErrorDialup_set_auto_apn_match_error                             int64 = 107729
	ErrorDialup_set_connect_file_error                               int64 = 107721
	ErrorDialup_set_default_prorile_error                            int64 = 107726
	ErrorDisable_auto_pin_failed                                     int64 = 101008
	ErrorDisable_pin_failed                                          int64 = 101006
	ErrorEnable_auto_pin_failed                                      int64 = 101009
	ErrorEnable_pin_failed                                           int64 = 101005
	ErrorFirst_send                                                  int64 = 1
	ErrorFormat_error                                                int64 = 100005
	ErrorGet_config_file_error                                       int64 = 100008
	ErrorGet_connect_status_failed                                   int64 = 102004
	ErrorGet_net_type_failed                                         int64 = 102001
	ErrorGet_roam_status_failed                                      int64 = 102003
	ErrorGet_service_status_failedint64                                    = 102002
	ErrorLanguage_get_failed                                         int64 = 109001
	ErrorLanguage_set_failed                                         int64 = 109002
	ErrorLogin_already_logined                                       int64 = 108003
	ErrorLogin_modify_password_failed                                int64 = 108004
	ErrorLogin_no_exist_user                                         int64 = 108001
	ErrorLogin_password_error                                        int64 = 108002
	ErrorLogin_too_many_times                                        int64 = 108007
	ErrorLogin_too_many_users_logined                                int64 = 108005
	ErrorLogin_username_or_password_error                            int64 = 108006
	ErrorNet_current_net_mode_not_support                            int64 = 112007
	ErrorNet_memory_alloc_failed                                     int64 = 112009
	ErrorNet_net_connected_order_not_match                           int64 = 112006
	ErrorNet_register_net_failed                                     int64 = 112005
	ErrorNet_sim_card_not_ready_status                               int64 = 112008
	ErrorNot_support                                                 int64 = 100002
	ErrorNo_device                                                   int64 = -2
	ErrorNo_right                                                    int64 = 100003
	ErrorNo_sim_card_or_invalid_sim_card                             int64 = 101001
	ErrorOnline_update_already_booted                                int64 = 110002
	ErrorOnline_update_cancel_downloding                             int64 = 110007
	ErrorOnline_update_connect_error                                 int64 = 110009
	ErrorOnline_update_get_device_information_failed                 int64 = 110003
	ErrorOnline_update_get_local_group_commponent_information_failed int64 = 110004
	ErrorOnline_update_invalid_url_list                              int64 = 110021
	ErrorOnline_update_low_battery                                   int64 = 110024
	ErrorOnline_update_need_reconnect_server                         int64 = 110006
	ErrorOnline_update_not_boot                                      int64 = 110023
	ErrorOnline_update_not_find_file_on_server                       int64 = 110005
	ErrorOnline_update_not_support_url_list                          int64 = 110022
	ErrorOnline_update_same_file_list                                int64 = 110008
	ErrorOnline_update_server_not_accessed                           int64 = 110001
	ErrorParameter_error                                             int64 = 100006
	ErrorPb_call_system_fucntion_error                               int64 = 115003
	ErrorPb_local_telephone_full_error                               int64 = 115199
	ErrorPb_null_argument_or_illegal_argument                        int64 = 115001
	ErrorPb_overtime                                                 int64 = 115002
	ErrorPb_read_file_error                                          int64 = 115005
	ErrorPb_write_file_error                                         int64 = 115004
	ErrorSafe_error                                                  int64 = 106001
	ErrorSave_config_file_errorint64                                       = 100007
	ErrorSD_directory_exist                                          int64 = 114002
	ErrorSD_file_exist                                               int64 = 114001
	ErrorSD_file_is_uploading                                        int64 = 114007
	ErrorSD_file_name_too_long                                       int64 = 114005
	ErrorSD_file_or_directory_not_exist                              int64 = 114004
	ErrorSD_is_operted_by_other_user                                 int64 = 114004
	ErrorSD_no_right                                                 int64 = 114006
	ErrorSet_net_mode_and_band_failed                                int64 = 112003
	ErrorSet_net_mode_and_band_when_dailup_failed                    int64 = 112001
	ErrorSet_net_search_mode_failed                                  int64 = 112004
	ErrorSet_net_search_mode_when_dailup_failed                      int64 = 112002
	ErrorSMS_delete_sms_failed                                       int64 = 113036
	ErrorSMS_local_space_not_enough                                  int64 = 113053
	ErrorSMS_null_argument_or_illegal_argument                       int64 = 113017
	ErrorSMS_overtime                                                int64 = 113018
	ErrorSMS_query_sms_index_list_error                              int64 = 113020
	ErrorSMS_save_config_file_failed                                 int64 = 113047
	ErrorSMS_set_sms_center_number_failed                            int64 = 113031
	ErrorSMS_telephone_number_too_long                               int64 = 113054
	ErrorSTK_call_system_fucntion_error                              int64 = 116003
	ErrorSTK_null_argument_or_illegal_argument                       int64 = 116001
	ErrorSTK_overtime                                                int64 = 116002
	ErrorSTK_read_file_error                                         int64 = 116005
	ErrorSTK_write_file_error                                        int64 = 116004
	ErrorUnknown                                                     int64 = 100001
	ErrorUnlock_pin_failed                                           int64 = 101007
	ErrorUSSD_at_send_failed                                         int64 = 111018
	ErrorUSSD_coding_error                                           int64 = 111017
	ErrorUSSD_empty_command                                          int64 = 111016
	ErrorUSSD_error                                                  int64 = 111001
	ErrorUSSD_fucntion_return_error                                  int64 = 111012
	ErrorUSSD_in_ussd_session                                        int64 = 111013
	ErrorUSSD_net_not_support_ussd                                   int64 = 111022
	ErrorUSSD_net_no_return                                          int64 = 11019
	ErrorUSSD_net_overtimeint64                                      int64 = 111020
	ErrorUSSD_too_long_content                                       int64 = 111014
	ErrorUSSD_xml_special_character_transfer_failed                  int64 = 111021
	ErrorWiFi_pbc_connect_failed                                     int64 = 117003
	ErrorWiFi_station_connect_ap_password_error                      int64 = 117001
	ErrorWiFi_station_connect_ap_wispr_password_error                int64 = 117004
	ErrorWiFi_web_password_or_dhcp_overtime_error                    int64 = 117002
)

var errorText = map[int64]string{
	ErrorBusy:                        "Busy",
	ErrorDeviceCompressLogFileFailed: "Compress log file failed",
}

func ErrorMessage(code int64) string {
	if text, ok := errorText[code]; ok {
		return text
	}

	return ""
}
