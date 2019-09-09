package hilink

const (
	ErrorBusy                                                         int64 = 100004
	ErrorCheckSimCardCanUnuseable                                     int64 = 101004
	ErrorCheckSimCardPINLock                                          int64 = 101002
	ErrorCheckSimCardPUNLock                                          int64 = 101003
	ErrorCompressLogFileFailed                                        int64 = 103102
	ErrorCradleCodingFailed                                           int64 = 118005
	ErrorCradleGetCrurrentConnectedUserIPFailed                       int64 = 118001
	ErrorCradleGetCrurrentConnectedUserMACFailed                      int64 = 118002
	error_cradle_get_wan_information_failed                           int64 = 118004
	error_cradle_set_mac_failed                                       int64 = 118003
	error_cradle_update_profile_failed                                int64 = 118006
	error_default                                                     int64 = -1
	error_device_at_execute_failed                                    int64 = 103001
	ErrorDeviceCompressLogFileFailed                                  int64 = 103015
	Error_device_get_api_version_failed                               int64 = 103006
	Error_device_get_autorun_version_failed                           int64 = 103005
	Error_device_get_log_informaton_level_failed                      int64 = 103014
	Error_device_get_pc_aissst_information_failed                     int64 = 103012
	Error_device_get_product_informaton_failed                        int64 = 103007
	Error_device_not_support_remote_operate                           int64 = 103010
	Error_device_pin_modiffy_failed                                   int64 = 103003
	Error_device_pin_validate_failed                                  int64 = 103002
	Error_device_puk_dead_lock                                        int64 = 103011
	Error_device_puk_modiffy_failed                                   int64 = 103004
	Error_device_restore_file_decrypt_failed                          int64 = 103016
	Error_device_restore_file_failed                                  int64 = 103018
	Error_device_restore_file_version_match_failed                    int64 = 103017
	Error_device_set_log_informaton_level_failed                      int64 = 103013
	Error_device_set_time_failed                                      int64 = 103101
	Error_device_sim_card_busy                                        int64 = 103008
	Error_device_sim_lock_input_error                                 int64 = 103009
	Error_dhcp_error                                                  int64 = 104001
	Error_dialup_add_prorile_error                                    int64 = 107724
	Error_dialup_dialup_managment_parse_error                         int64 = 107722
	Error_dialup_get_auto_apn_match_error                             int64 = 107728
	Error_dialup_get_connect_file_error                               int64 = 107720
	Error_dialup_get_prorile_list_error                               int64 = 107727
	Error_dialup_modify_prorile_error                                 int64 = 107725
	Error_dialup_set_auto_apn_match_error                             int64 = 107729
	Error_dialup_set_connect_file_error                               int64 = 107721
	Error_dialup_set_default_prorile_error                            int64 = 107726
	Error_disable_auto_pin_failed                                     int64 = 101008
	Error_disable_pin_failed                                          int64 = 101006
	Error_enable_auto_pin_failed                                      int64 = 101009
	Error_enable_pin_failed                                           int64 = 101005
	Error_first_send                                                  int64 = 1
	Error_format_error                                                int64 = 100005
	Error_get_config_file_error                                       int64 = 100008
	Error_get_connect_status_failed                                   int64 = 102004
	Error_get_net_type_failed                                         int64 = 102001
	Error_get_roam_status_failed                                      int64 = 102003
	Error_get_service_status_failedint64                                    = 102002
	Error_language_get_failed                                         int64 = 109001
	Error_language_set_failed                                         int64 = 109002
	Error_login_already_logined                                       int64 = 108003
	Error_login_modify_password_failed                                int64 = 108004
	Error_login_no_exist_user                                         int64 = 108001
	Error_login_password_error                                        int64 = 108002
	Error_login_too_many_times                                        int64 = 108007
	Error_login_too_many_users_logined                                int64 = 108005
	Error_login_username_or_password_error                            int64 = 108006
	Error_net_current_net_mode_not_support                            int64 = 112007
	Error_net_memory_alloc_failed                                     int64 = 112009
	Error_net_net_connected_order_not_match                           int64 = 112006
	Error_net_register_net_failed                                     int64 = 112005
	Error_net_sim_card_not_ready_status                               int64 = 112008
	Error_not_support                                                 int64 = 100002
	Error_no_device                                                   int64 = -2
	Error_no_right                                                    int64 = 100003
	Error_no_sim_card_or_invalid_sim_card                             int64 = 101001
	Error_online_update_already_booted                                int64 = 110002
	Error_online_update_cancel_downloding                             int64 = 110007
	Error_online_update_connect_error                                 int64 = 110009
	Error_online_update_get_device_information_failed                 int64 = 110003
	Error_online_update_get_local_group_commponent_information_failed int64 = 110004
	Error_online_update_invalid_url_list                              int64 = 110021
	Error_online_update_low_battery                                   int64 = 110024
	Error_online_update_need_reconnect_server                         int64 = 110006
	Error_online_update_not_boot                                      int64 = 110023
	Error_online_update_not_find_file_on_server                       int64 = 110005
	Error_online_update_not_support_url_list                          int64 = 110022
	Error_online_update_same_file_list                                int64 = 110008
	Error_online_update_server_not_accessed                           int64 = 110001
	Error_parameter_error                                             int64 = 100006
	Error_pb_call_system_fucntion_error                               int64 = 115003
	Error_pb_local_telephone_full_error                               int64 = 115199
	Error_pb_null_argument_or_illegal_argument                        int64 = 115001
	Error_pb_overtime                                                 int64 = 115002
	Error_pb_read_file_error                                          int64 = 115005
	Error_pb_write_file_error                                         int64 = 115004
	Error_safe_error                                                  int64 = 106001
	Error_save_config_file_errorint64                                       = 100007
	Error_sd_directory_exist                                          int64 = 114002
	Error_sd_file_exist                                               int64 = 114001
	Error_sd_file_is_uploading                                        int64 = 114007
	Error_sd_file_name_too_long                                       int64 = 114005
	Error_sd_file_or_directory_not_exist                              int64 = 114004
	Error_sd_is_operted_by_other_user                                 int64 = 114004
	Error_sd_no_right                                                 int64 = 114006
	Error_set_net_mode_and_band_failed                                int64 = 112003
	Error_set_net_mode_and_band_when_dailup_failed                    int64 = 112001
	Error_set_net_search_mode_failed                                  int64 = 112004
	Error_set_net_search_mode_when_dailup_failed                      int64 = 112002
	Error_sms_delete_sms_failed                                       int64 = 113036
	Error_sms_local_space_not_enough                                  int64 = 113053
	Error_sms_null_argument_or_illegal_argument                       int64 = 113017
	Error_sms_overtime                                                int64 = 113018
	Error_sms_query_sms_index_list_error                              int64 = 113020
	Error_sms_save_config_file_failed                                 int64 = 113047
	Error_sms_set_sms_center_number_failed                            int64 = 113031
	Error_sms_telephone_number_too_long                               int64 = 113054
	Error_stk_call_system_fucntion_error                              int64 = 116003
	Error_stk_null_argument_or_illegal_argument                       int64 = 116001
	Error_stk_overtime                                                int64 = 116002
	Error_stk_read_file_error                                         int64 = 116005
	Error_stk_write_file_error                                        int64 = 116004
	Error_unknown                                                     int64 = 100001
	Error_unlock_pin_failed                                           int64 = 101007
	Error_ussd_at_send_failed                                         int64 = 111018
	Error_ussd_coding_error                                           int64 = 111017
	Error_ussd_empty_command                                          int64 = 111016
	Error_ussd_error                                                  int64 = 111001
	Error_ussd_fucntion_return_error                                  int64 = 111012
	Error_ussd_in_ussd_session                                        int64 = 111013
	Error_ussd_net_not_support_ussd                                   int64 = 111022
	Error_ussd_net_no_return                                          int64 = 11019
	Error_ussd_net_overtimeint64                                      int64 = 111020
	Error_ussd_too_long_content                                       int64 = 111014
	Error_ussd_xml_special_character_transfer_failed                  int64 = 111021
	Error_wifi_pbc_connect_failed                                     int64 = 117003
	Error_wifi_station_connect_ap_password_error                      int64 = 117001
	Error_wifi_station_connect_ap_wispr_password_error                int64 = 117004
	Error_wifi_web_password_or_dhcp_overtime_error                    int64 = 117002
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
