package root

import (
	"bufio"
	"context"
	"os"
)

func (b *Bind) runtimeConfigWatcher(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	sn := b.Meta().SerialNumber()

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
			_ = b.MQTT().Publish(context.Background(), b.config.TopicRuntimeConfig.Format(sn, key), current)
		}

		b.cacheRuntimeConfigLock.Unlock()
	}

	return nil
}
