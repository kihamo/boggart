Активные подключения
====================

Единственный вариант читать ```/ip arp```, но эта таблица кэшируется и надо ее чистить переодически от динамических записей.
Для этого нужно добавить в mikrotik скрипт

```
/system script add name="clear-arp"
/system script edit clear-arp source
```

вставить код

```
:log info "clearing arp table of dynamic entries"       
/ip arp remove [/ip arp find dynamic=yes]
```

можно проверить запустив вручную

```
/system script run clear-arp
```

добавить в планировщик с интервалом запуска 10 минут

```
/system scheduler add name=clea-arp-10m interval=10m on-event=clear-arp
```