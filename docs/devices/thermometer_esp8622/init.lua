dofile("config.lua")

local MQTT_CLIENT_ID = wifi.sta.gethostname()
local MQTT_TOPIC = "/esp8266/"..string.gsub(wifi.sta.getmac(),":","-").."/"

-- program
local TMR_WIFI_STATUS = 0
local TMR_MQTT_CONNECT = 1
local TMR_SENSOR_READ = 2

tmr.stop(TMR_WIFI_STATUS)
tmr.stop(TMR_MQTT_CONNECT)
tmr.stop(TMR_SENSOR_READ)

function wifi_connect()
    local cfg={}
    cfg.ssid=WIFI_SSID
    cfg.pwd=WIFI_PASSWORD
    cfg.save=true

    wifi.setmode(wifi.STATION)
    wifi.nullmodesleep(false)

    wifi.sta.config(cfg)
    wifi.sta.autoconnect(1)
    wifi.sta.sethostname("esp8266-"..node.chipid())
    wifi.sta.connect()

    tmr.alarm(TMR_WIFI_STATUS, 1000, 1, function()
        if (wifi.sta.status() == 5) then
            tmr.stop(TMR_WIFI_STATUS)

            print("WiFi connected")
            print(wifi.sta.getip())

            mqtt_connect()
        end
    end)
end

function mqtt_connect()
    if wifi.sta.status() ~= 5 then
        print("No net connect")
        return
    end

    m:lwt('/lwt/' .. MQTT_CLIENT_ID, "died", 0, 0)

    m:connect(MQTT_HOSTNAME, MQTT_PORT, 0, 0,
        function()
            tmr.stop(TMR_MQTT_CONNECT)
            m:on("offline", function()
                print("MQTT is OFFLINE")

                mqtt_connect()
            end)

            print("MQTT connected")

            tmr.alarm(TMR_SENSOR_READ, 3000, tmr.ALARM_SINGLE, function()
                sensor_read()
            end)
        end,
        function(_, reason)
            if reason == -5 then print("MQTT connect failed CONN_FAIL_SERVER_NOT_FOUND")
            elseif reason == -4 then print("MQTT connect failed CONN_FAIL_NOT_A_CONNACK_MSG")
            elseif reason == -3 then print("MQTT connect failed CONN_FAIL_DNS")
            elseif reason == -2 then print("MQTT connect failed CONN_FAIL_TIMEOUT_RECEIVING")
            elseif reason == -1 then print("MQTT connect failed CONN_FAIL_TIMEOUT_SENDING")
            elseif reason == 0 then print("MQTT connect failed CONNACK_ACCEPTED")
            elseif reason == 1 then print("MQTT connect failed CONNACK_REFUSED_PROTOCOL_VER")
            elseif reason == 2 then print("MQTT connect failed CONNACK_REFUSED_ID_REJECTED")
            elseif reason == 3 then print("MQTT connect failed CONNACK_REFUSED_SERVER_UNAVAILABLE")
            elseif reason == 4 then print("MQTT connect failed CONNACK_REFUSED_BAD_USER_OR_PASS")
            elseif reason == 5 then print("MQTT connect failed CONNACK_REFUSED_NOT_AUTHORIZED")
            end

            tmr.alarm(TMR_MQTT_CONNECT, 5000, tmr.ALARM_SINGLE, function()
                print("MQTT atempt connect")

                mqtt_connect()
            end)
        end
    )
end

function sensor_read()
    print("Sensor read started")

    -- local totalAllocated, estimatedUsed = node.egc.meminfo()
    -- m:publish(MQTT_TOPIC.."egc/total-allocated", totalAllocated, 1, 1)
    -- m:publish(MQTT_TOPIC.."egc/estimated-used", estimatedUsed, 0, 1)
    -- m:publish(MQTT_TOPIC.."heap", node.heap(), 0, 1)
    m:publish(MQTT_TOPIC.."voltage", string.format("%.2f", adc.readvdd33(0) / 1000), 0, 1)

    local T, P, H, _ = bme280.read(BME280_ALTITUDE)

    if T ~= nil then
        m:publish(MQTT_TOPIC.."sensor/temperature", string.format("%.2f", T / 100), 0, 1)
    end

    if P ~= nil then
        -- из hPa в mmHg
        m:publish(MQTT_TOPIC.."sensor/pressure", string.format("%.2f", P / 1000 * 0.75), 0, 1)
    end

    if H ~= nil then
        m:publish(MQTT_TOPIC.."sensor/humidity", string.format("%.2f", H / 1000), 0, 1)
    end

    if T == nil or P == nil or H == nil then
        return
    end

    print("Sensor read done")
    tmr.stop(TMR_SENSOR_READ)

    -- задержка в секунду, что бы publish отработал нормально
    tmr.alarm(TMR_SENSOR_READ, 1000, tmr.ALARM_SINGLE, function()
        print("Sleep mode on")
        node.dsleep(SLEEP_INTERVAL * 1000000, 2)
    end)
end

if adc.force_init_mode(adc.INIT_VDD33)
then
    node.restart()
    return
end

reset_code, reset_reason = node.bootreason()
if reset_code == 1 then print("Boot code: "..reset_code.." power-on")
elseif reset_code == 2 then print("Boot code: "..reset_code.." reset")
elseif reset_code == 3 then print("Boot code: "..reset_code.." hardware reset via reset pin")
elseif reset_code == 4 then print("Boot code: "..reset_code.." WDT reset")
end

if reset_reason == 0 then print("Boot reason: "..reset_reason.." power-on")
elseif reset_reason == 1 then print("Boot reason: "..reset_reason.." hardware watchdog reset")
elseif reset_reason == 2 then print("Boot reason: "..reset_reason.." exception reset")
elseif reset_reason == 3 then print("Boot reason: "..reset_reason.." software watchdog reset")
elseif reset_reason == 4 then print("Boot reason: "..reset_reason.." software restart")
elseif reset_reason == 5 then print("Boot reason: "..reset_reason.." wake from deep sleep")
elseif reset_reason == 6 then print("Boot reason: "..reset_reason.." external reset")
end

-- I2C und BME280 init
i2c.setup(0, GPIO_SDA, GPIO_SCL, i2c.SLOW)
bme280.setup(nil, nil, nil, 0)

-- MQTT init
m = mqtt.Client(MQTT_CLIENT_ID, 120, MQTT_USER, MQTT_PASSWORD)

-- WiFi
wifi_connect()
