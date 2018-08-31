package internal

import (
	"bufio"
	"os"
	"regexp"
	"sync"

	"github.com/kihamo/boggart/components/roborock"
)

/*
Example

RoboController :
{
  runtime :
  {
    bin_in_time = 168;
    time_slept = 0;
    sound_volume = 38;
    clean_id = 304;
    carpet_mode_enabled = 1;
    carpet_mode_curr_integral = 450;
    carpet_mode_curr_highwater = 500;
    carpet_mode_curr_lowwater = 400;
    carpet_mode_stall_time = 10;
    fan_power = 38;
    temp_fan_power = -1;
  };
};
*/

var (
	reRuntimeConfigLine    = regexp.MustCompile(`(?m)\s*([[:alnum:]_]+)\s*=\s*([^;]+);`)
	cacheRuntimeConfig     = make(map[string]string, 11)
	cacheRuntimeConfigLock sync.Mutex
)

func (c *Component) runtimeConfigWatcher(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		match := reRuntimeConfigLine.FindSubmatch(scanner.Bytes())
		if len(match) == 0 {
			continue
		}

		key := string(match[1])
		value := string(match[2])

		cacheRuntimeConfigLock.Lock()
		prevValue, ok := cacheRuntimeConfig[key]

		if !ok || prevValue != value {
			cacheRuntimeConfig[key] = value
			c.mqtt.Publish(roborock.MQTTTopicPrefix+"runtime/"+key, 0, false, value)
		}

		cacheRuntimeConfigLock.Unlock()
	}

	return nil
}
