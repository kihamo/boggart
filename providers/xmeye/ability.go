package xmeye

import (
	"context"
)

type ability string

const (
	AbilityEncodeCapability ability = "EncodeCapability"
	AbilityBlindCapability  ability = "BlindCapability"
	AbilityMotionArea       ability = "MotionArea"
	AbilityDDNSService      ability = "DDNSService"
	AbilityComProtocol      ability = "ComProtocol"
	AbilityPTZProtocol      ability = "PTZProtocol"
	AbilityTalkAudioFormat  ability = "TalkAudioFormat"
	AbilityMultiLanguage    ability = "MultiLanguage"
	AbilitySystemFunction   ability = "SystemFunction"
)

func (c *Client) Ability(ctx context.Context, name ability) (interface{}, error) {
	var result map[string]interface{}

	err := c.CmdWithResult(ctx, CmdAbilityGetRequest, string(name), &result)
	if err != nil {
		return nil, err
	}

	if ability, ok := result[string(name)]; ok {
		return ability, nil
	}

	return nil, err
}
