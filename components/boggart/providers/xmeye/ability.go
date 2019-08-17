package xmeye

const (
	AbilityEncodeCapability = "EncodeCapability"
	AbilityBlindCapability  = "BlindCapability"
	AbilityMotionArea       = "MotionArea"
	AbilityDDNSService      = "DDNSService"
	AbilityComProtocol      = "ComProtocol"
	AbilityPTZProtocol      = "PTZProtocol"
	AbilityTalkAudioFormat  = "TalkAudioFormat"
	AbilityMultiLanguage    = "MultiLanguage"
	AbilitySystemFunction   = "SystemFunction"
)

func (c *Client) Ability(name string) (interface{}, error) {
	var result map[string]interface{}

	err := c.CmdWithResult(CmdAbilityGetRequest, name, &result)
	if err != nil {
		return nil, err
	}

	if ability, ok := result[name]; ok {
		return ability, nil
	}

	return nil, err
}
