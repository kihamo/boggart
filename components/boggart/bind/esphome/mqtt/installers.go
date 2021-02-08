package mqtt

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	components := b.Components()
	sort.SliceStable(components, func(i, j int) bool {
		return components[i].ID() < components[j].ID()
	})

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	channels := make([]*openhab.Channel, 0, len(components))

	for _, component := range components {
		id := openhab.IDNormalizeCamelCase(component.ID())

		switch component.Type() {
		case ComponentTypeBinarySensor:
			var (
				channel *openhab.Channel
				item    *openhab.Item
			)

			if cmp, ok := component.(*ComponentBinarySensor); ok && cmp.PayloadOn() != "" && cmp.PayloadOff() != "" {
				channel = openhab.NewChannel(id, openhab.ChannelTypeContact).
					WithOn(cmp.PayloadOn()).
					WithOff(cmp.PayloadOff())

				item = openhab.NewItem(itemPrefix+id, openhab.ItemTypeContact).
					WithIcon(OpenHabIconConverter(component.Icon(), "contact"))
			} else if component.CommandTopic() == "" {
				// если switch не управляемый переводим его в режим readonly
				channel = openhab.NewChannel(id, openhab.ChannelTypeContact).
					WithOn("ON").
					WithOff("OFF")

				item = openhab.NewItem(itemPrefix+id, openhab.ItemTypeContact).
					WithIcon(OpenHabIconConverter(component.Icon(), "contact"))
			} else {
				channel = openhab.NewChannel(id, openhab.ChannelTypeSwitch).
					WithOn("ON").
					WithOff("OFF")

				item = openhab.NewItem(itemPrefix+id, openhab.ItemTypeSwitch).
					WithIcon(OpenHabIconConverter(component.Icon(), "contact"))
			}

			channels = append(channels,
				channel.
					WithStateTopic(component.StateTopic()).
					WithCommandTopic(component.CommandTopic()).
					AddItems(
						item.WithLabel(strings.Title(component.Name())),
					),
			)

		case ComponentTypeSensor:
			label := strings.Title(component.Name())

			if cmp, ok := component.(*ComponentSensor); ok {
				label += " [%." + strconv.FormatUint(cmp.AccuracyDecimals(), 10) + "f"

				if unit := cmp.UnitOfMeasurement(); unit != "" {
					label += " " + strings.ReplaceAll(unit, "%", "%%")
				}

				label += "]"
			}

			channels = append(channels,
				openhab.NewChannel(id, openhab.ChannelTypeNumber).
					WithStateTopic(component.StateTopic()).
					WithCommandTopic(component.CommandTopic()).
					AddItems(
						openhab.NewItem(itemPrefix+id, openhab.ItemTypeNumber).
							WithLabel(label).
							WithIcon(OpenHabIconConverter(component.Icon(), "")),
					),
			)

		case ComponentTypeSwitch:
			var (
				channel *openhab.Channel
				item    *openhab.Item
			)

			// если switch не управляемый переводим его в режим readonly
			if component.CommandTopic() == "" {
				channel = openhab.NewChannel(id, openhab.ChannelTypeContact)

				item = openhab.NewItem(itemPrefix+id, openhab.ItemTypeContact).
					WithIcon(OpenHabIconConverter(component.Icon(), "contact"))
			} else {
				channel = openhab.NewChannel(id, openhab.ChannelTypeSwitch).
					WithCommandTopic(component.CommandTopic())

				item = openhab.NewItem(itemPrefix+id, openhab.ItemTypeSwitch).
					WithIcon(OpenHabIconConverter(component.Icon(), "switch"))
			}

			channels = append(channels,
				channel.
					WithStateTopic(component.StateTopic()).
					WithOn("ON").
					WithOff("OFF").
					AddItems(
						item.WithLabel(strings.Title(component.Name())),
					),
			)

		case ComponentTypeLight:
			channels = append(channels,
				openhab.NewChannel(id, openhab.ChannelTypeSwitch).
					WithStateTopic(component.StateTopic()).
					WithCommandTopic(component.CommandTopic()).
					WithOn("ON").
					WithOff("OFF").
					WithTransformationPattern("JSONPATH:$.state").
					WithFormatBeforePublish(`{\"state\":%s}`).
					AddItems(
						openhab.NewItem(itemPrefix+id, openhab.ItemTypeSwitch).
							WithLabel(strings.Title(component.Name())).
							WithIcon(OpenHabIconConverter(component.Icon(), "light")),
					),
			)

			if cmp, ok := component.(*ComponentLight); ok {
				if cmp.Brightness() {
					const idPostfix = "Brightness"

					channels = append(channels,
						openhab.NewChannel(id+idPostfix, openhab.ChannelTypeDimmer).
							WithStateTopic(component.StateTopic()).
							WithCommandTopic(component.CommandTopic()).
							WithMin(0).
							WithMax(255).
							WithStep(1).
							WithTransformationPattern("JSONPATH:$.brightness").
							WithFormatBeforePublish(`{\"brightness\":%s}`).
							AddItems(
								openhab.NewItem(itemPrefix+id+idPostfix, openhab.ItemTypeDimmer).
									WithLabel(strings.Title(component.Name())+" brightness [%.0f %%]").
									WithIcon(OpenHabIconConverter(component.Icon(), "heating-40")),
							),
					)
				}

				if cmp.Effect() {
					const idPostfix = "Effect"

					channels = append(channels,
						openhab.NewChannel(id+idPostfix, openhab.ChannelTypeString).
							WithStateTopic(component.StateTopic()).
							WithCommandTopic(component.CommandTopic()).
							WithTransformationPattern("JSONPATH:$.effect").
							WithFormatBeforePublish(`{\"effect\":%s}`).
							AddItems(
								openhab.NewItem(itemPrefix+id+idPostfix, openhab.ItemTypeString).
									WithLabel(strings.Title(component.Name())+" effect").
									WithIcon(OpenHabIconConverter(component.Icon(), "rgb")),
							),
					)
				}
			}

		default:
			channels = append(channels,
				openhab.NewChannel(id, openhab.ChannelTypeString).
					WithStateTopic(component.StateTopic()).
					WithCommandTopic(component.CommandTopic()).
					AddItems(
						openhab.NewItem(itemPrefix+id, openhab.ItemTypeString).
							WithLabel(strings.Title(component.Name())).
							WithIcon(OpenHabIconConverter(component.Icon(), "")),
					),
			)
		}
	}

	return openhab.StepsByBind(b, nil, channels...)
}

func OpenHabIconConverter(icon, def string) string {
	switch icon {
	//case "mdi:axis-arrow":
	//case "mdi:axis-x-arrow":
	//case "mdi:axis-y-arrow":
	//case "mdi:axis-z-arrow":
	//case "mdi:arrow-expand-vertical":
	case "mdi:battery":
		return "battery"
	case "mdi:briefcase-download":
		return "returnpipe"
	case "mdi:bug":
		return "status"
		//case "mdi:check-circle-outline":
		//case "mdi:chemical-weapon":
		//case "mdi:counter":
	case "mdi:current-ac":
		return "line"
	case "mdi:flash":
		return "energy"
	case "mdi:flask-outline":
		return "sewerage"
	case "mdi:flower":
		return "lawnmower"
	case "mdi:gas-cylinder":
		return "gas"
	case "mdi:gauge":
		return "pressure"
	case "mdi:brightness-5", "mdi:lightbulb":
		return "light"
		//case "mdi:magnet":
	case "mdi:molecule-co2":
		return "carbondioxide"
	case "mdi:motion-sensor":
		return "motion"
		// case "mdi:new-box":
		// case "mdi:percent":
	case "mdi:restart", "mdi:power":
		return "switch"
		// case "mdi:pulse":
	case "mdi:radiator":
		return "radiator"
		// case "mdi:rotate-right":
	case "mdi:ruler", "mdi:scale":
		return "niveau"
	case "mdi:screen-rotation":
		return "screen"
		// case "mdi:sign-direction":
	case "mdi:signal-distance-variant", "mdi:signal":
		return "qualityofservice"
	case "mdi:thermometer":
		return "temperature"
	case "mdi:timelapse", "mdi:timer-outline":
		return "time"
	case "mdi:water-percent":
		return "humidity"
	case "mdi:weather-sunset", "mdi:weather-sunset-down", "mdi:weather-sunset-up":
		return "sunset"
	case "mdi:weather-windy":
		return "wind"
	case "mdi:wifi":
		return "network"

		// OpenHab default icons
		// https://github.com/eclipse-archived/smarthome/tree/master/extensions/ui/iconset/org.eclipse.smarthome.ui.iconset.classic/icons
	case "alarm", "attic",

		"baby_1", "baby_2", "baby_3", "baby_4", "baby_5", "baby_6", "bath", "battery", "battery-0", "battery-10",
		"battery-100", "battery-20", "battery-30", "battery-40", "battery-50", "battery-60", "battery-70", "battery-80",
		"battery-90", "batterylevel", "batterylevel-0", "batterylevel-10", "batterylevel-100", "batterylevel-20",
		"batterylevel-30", "batterylevel-40", "batterylevel-50", "batterylevel-60", "batterylevel-70", "batterylevel-80",
		"batterylevel-90", "bedroom", "bedroom_blue", "bedroom_orange", "bedroom_red", "blinds", "blinds-0", "blinds-10",
		"blinds-100", "blinds-20", "blinds-30", "blinds-40", "blinds-50", "blinds-60", "blinds-70", "blinds-80", "blinds-90",
		"bluetooth", "boy_1", "boy_2", "boy_3", "boy_4", "boy_5", "boy_6",

		"calendar", "camera", "carbondioxide", "cellar", "chart", "cinema", "cinemascreen", "cinemascreen-0", "cinemascreen-10",
		"cinemascreen-100", "cinemascreen-20", "cinemascreen-30", "cinemascreen-40", "cinemascreen-50", "cinemascreen-60",
		"cinemascreen-70", "cinemascreen-80", "cinemascreen-90", "cistern", "cistern-0", "cistern-10", "cistern-100", "cistern-20",
		"cistern-30", "cistern-40", "cistern-50", "cistern-60", "cistern-70", "cistern-80", "cistern-90", "climate", "climate-on",
		"colorlight", "colorpicker", "colorwheel", "contact", "contact-ajar", "contact-closed", "contact-open", "corridor",

		"door", "door-closed", "door-open", "dryer ", "dryer-0", "dryer-1", "dryer-2", "dryer-3", "dryer-4", "dryer-5",

		"energy", "error",

		"fan", "fan_box", "fan_ceiling", "faucet", "fire", "fire-off", "fire-on", "firstfloor", "flow", "flowpipe", "frontdoor",
		"frontdoor-closed", "frontdoor-open",

		"garage", "garage_detached", "garage_detached_selected", "garagedoor", "garagedoor-0", "garagedoor-10", "garagedoor-100",
		"garagedoor-20", "garagedoor-30", "garagedoor-40", "garagedoor-50", "garagedoor-60", "garagedoor-70", "garagedoor-80",
		"garagedoor-90", "garagedoor-ajar", "garagedoor-closed", "garagedoor-open", "garden", "gas", "girl_1", "girl_2",
		"girl_3", "girl_4", "girl_5", "girl_6", "greenhouse", "groundfloor", "group",

		"heating", "heating-0", "heating-100", "heating-20", "heating-40", "heating-60", "heating-80", "heating-off", "heating-on",
		"house", "humidity", "humidity-0", "humidity-10", "humidity-100", "humidity-20", "humidity-30", "humidity-40", "humidity-50",
		"humidity-60", "humidity-70", "humidity-80", "humidity-90",

		"incline",

		"keyring", "kitchen",

		"lawnmower", "light", "light-0", "light-10", "light-100", "light-20", "light-30", "light-40", "light-50", "light-60",
		"light-70", "light-80", "light-90", "light-off", "light-on", "lightbulb", "line", "line-decline", "line-incline",
		"line-stagnation", "lock", "lock-closed", "lock-open", "lowbattery", "lowbattery-off", "lowbattery-on",

		"man_1", "man_2", "man_3", "man_4", "man_5", "man_6", "mediacontrol", "microphone", "moon", "motion", "movecontrol",

		"network", "network-off", "network-on", "niveau", "none",

		"office", "oil", "outdoorlight",

		"pantry", "parents-off", "parents_1_1", "parents_1_2", "parents_1_3", "parents_1_4", "parents_1_5", "parents_1_6",
		"parents_2_1", "parents_2_2", "parents_2_3", "parents_2_4", "parents_2_5", "parents_2_6", "parents_3_1", "parents_3_2",
		"parents_3_3", "parents_3_4", "parents_3_5", "parents_3_6", "parents_4_1", "parents_4_2", "parents_4_3", "parents_4_4",
		"parents_4_5", "parents_4_6", "parents_5_1", "parents_5_2", "parents_5_3", "parents_5_4", "parents_5_5", "parents_5_6",
		"parents_6_1", "parents_6_2", "parents_6_3", "parents_6_4", "parents_6_5", "parents_6_6", "party", "pie", "piggybank",
		"player", "poweroutlet", "poweroutlet-off", "poweroutlet-on", "poweroutlet_au", "poweroutlet_eu", "poweroutlet_uk",
		"poweroutlet_us", "presence", "presence-off", "presence-on", "pressure", "price", "projector", "pump",

		"qualityofservice", "qualityofservice-0", "qualityofservice-1", "qualityofservice-2", "qualityofservice-3", "qualityofservice-4",

		"radiator", "rain", "receiver", "receiver-off", "receiver-on", "recorder", "returnpipe", "rgb", "rollershutter",
		"rollershutter-0", "rollershutter-10", "rollershutter-100", "rollershutter-20", "rollershutter-30", "rollershutter-40",
		"rollershutter-50", "rollershutter-60", "rollershutter-70", "rollershutter-80", "rollershutter-90",

		"screen", "screen-off", "screen-on", "settings", "sewerage", "sewerage-0", "sewerage-10", "sewerage-100", "sewerage-20",
		"sewerage-30", "sewerage-40", "sewerage-50", "sewerage-60", "sewerage-70", "sewerage-80", "sewerage-90", "shield",
		"shield-0", "shield-1", "siren", "siren-off", "siren-on", "slider", "slider-0", "slider-10", "slider-100", "slider-20",
		"slider-30", "slider-40", "slider-50", "slider-60", "slider-70", "slider-80", "slider-90", "smiley", "smoke", "snow",
		"sofa", "softener", "solarplant", "soundvolume", "soundvolume-0", "soundvolume-100", "soundvolume-33", "soundvolume-66",
		"soundvolume_mute", "status", "suitcase", "sun", "sun_clouds", "sunrise", "sunset", "switch", "switch-off", "switch-on",

		"temperature", "temperature_cold", "temperature_hot", "terrace", "text", "time", "time-on", "toilet",

		"vacation", "video",

		"wallswitch", "wallswitch-off", "wallswitch-on", "wardrobe", "washingmachine", "washingmachine_2", "washingmachine_2-0",
		"washingmachine_2-1", "washingmachine_2-2", "washingmachine_2-3", "water", "whitegood", "wind", "window", "window-ajar",
		"window-closed", "window-open", "woman_1", "woman_2", "woman_3", "woman_4", "woman_5", "woman_6",

		"zoom":
		return icon
	}

	return def
}
