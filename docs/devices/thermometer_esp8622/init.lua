local GPIO_SDA, GPIO_SCL = 2, 1

local WIFI_SSID = ""
local WIFI_PASSWORD = ""

local MQTT_HOSTNAME = ""
local MQTT_PORT = 1883
local MQTT_USER = ""
local MQTT_PASSWORD = ""
local MQTT_CLIENT_ID = wifi.sta.gethostname()
local MQTT_TOPIC = "esp8266/"..string.gsub(wifi.sta.getmac(),":","-").."/"

local BME280_ALTITUDE=180

-- program
local TMR_WIFI_STATUS = 0
local TMR_MQTT_CONNECT = 1
local TMR_SENSOR_READ = 2

if adc.force_init_mode(adc.INIT_VDD33)
then
    node.restart()
    return
end

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

            mqtt_connect()
        end
    end)
end

function mqtt_connect()
    m:connect(MQTT_HOSTNAME, MQTT_PORT, 0, 0,
        function(conn)
            print("MQTT connected")
            m:on("offline", function() mqtt_connect() end)
            tmr.alarm(TMR_SENSOR_READ, 10000, 1, function()
                sensor_read()
            end)
        end,
        function(conn, reason)
            print("MQTT connect failed "..reason)
            tmr.alarm(TMR_MQTT_CONNECT, 5000, tmr.ALARM_SINGLE, mqtt_connect())
        end
    )
end

function sensor_read()
    local totalAllocated, estimatedUsed = node.egc.meminfo()

    m:publish(MQTT_TOPIC.."egc/total-allocated", totalAllocated, 0, 1)
    m:publish(MQTT_TOPIC.."egc/estimated-used", estimatedUsed, 0, 1)
    m:publish(MQTT_TOPIC.."heap", node.heap(), 0, 1)
    m:publish(MQTT_TOPIC.."voltage", string.format("%.2f", adc.readvdd33(0) / 1000), 0, 1)

    local T, P, H, QNH = bme280.read(BME280_ALTITUDE)
    if T ~= nil then
        m:publish(MQTT_TOPIC.."sensor/temperature", string.format("%.2f", T / 100), 0, 1)
    end

    if P ~= nil then
        m:publish(MQTT_TOPIC.."sensor/pressure", string.format("%.2f", P / 1000), 0, 1)
    end

    if H ~= nil then
        m:publish(MQTT_TOPIC.."sensor/humidity", string.format("%.2f", H / 1000), 0, 0, nil)
    end
end

-- WiFi
wifi_connect()

-- I2C und BME280 init
i2c.setup(0, GPIO_SDA, GPIO_SCL, i2c.SLOW)
bme280.setup()

-- MQTT init
m = mqtt.Client(MQTT_CLIENT_ID, 120, MQTT_USER, MQTT_PASSWORD)