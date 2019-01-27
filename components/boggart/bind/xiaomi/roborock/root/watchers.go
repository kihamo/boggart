package root

import (
	"bufio"
	"context"
	"os"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) runtimeConfigWatcher(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	sn := mqtt.NameReplace(b.SerialNumber())

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		match := reRuntimeConfigLine.FindSubmatch(scanner.Bytes())
		if len(match) == 0 {
			continue
		}

		key := string(match[1])
		current := string(match[2])

		b.cacheRuntimeConfigLock.Lock()
		prev, ok := b.cacheRuntimeConfig[key]

		if !ok || prev != current {
			b.cacheRuntimeConfig[key] = current

			// TODO
			_ = b.MQTTPublish(context.Background(), MQTTPublishTopicRuntimeConfig.Format(sn, key), 0, false, current)
		}

		b.cacheRuntimeConfigLock.Unlock()
	}

	return nil
}
