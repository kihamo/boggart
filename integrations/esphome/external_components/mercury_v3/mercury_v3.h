#pragma once

#include "esphome/core/component.h"
#include "esphome/components/uart/uart.h"
#include "esphome/components/sensor/sensor.h"

#define MERCURY_V3_ADDRESS_UNIVERSAL 0x0

#define MERCURY_V3_WAIT_AFTER_SEND_REQUEST 30
#define MERCURY_V3_WAIT_AFTER_READ_RESPONSE 100

#define MERCURY_V3_CHANNEL_OPEN_TIMEOUT 200000

namespace esphome {
  namespace mercury_v3 {
    class MercuryV3 : public PollingComponent, public uart::UARTDevice {
      public:
        MercuryV3() = default;

        void loop() override;
        void update() override;
        void setup() override;
        void dump_config() override;

        float get_setup_priority() const override;

        void set_voltage_sensor(int phase, sensor::Sensor *s) { this->phase_[phase].voltage_sensor_ = s; }
        void set_current_sensor(int phase, sensor::Sensor *s) { this->phase_[phase].current_sensor_ = s; }
        void set_power_sensor(int phase, sensor::Sensor *s) { this->phase_[phase].power_sensor_ = s; }
        void set_tariff1_sensor(sensor::Sensor *s) { tariff1_sensor_ = s; }

      protected:
        sensor::Sensor *tariff1_sensor_;

        struct MercuryV3Phase {
            sensor::Sensor *voltage_sensor_{nullptr};
            sensor::Sensor *current_sensor_{nullptr};
            sensor::Sensor *power_sensor_{nullptr};
          } phase_[3];

        uint8_t address_ = MERCURY_V3_ADDRESS_UNIVERSAL; // по-умолчанию устанавливаем универсальный адрес, этого достаточно, если в сети всего один счетчик

        unsigned char channel_open_request_[11];
        unsigned char read_voltage_request_[6];
        unsigned char read_current_request_[6];
        unsigned char read_power_request_[6];
        unsigned char read_tariff_1_request_[6];

        std::vector<uint8_t> response_;
        std::vector<uint8_t> payload_;
        uint32_t last_open_channel_;

        void update_voltage();
        void update_current();
        void update_power();
        void update_tariffs();

        void uart_tx(const uint8_t *data, size_t len);
        bool uart_rx();
        bool invoke(const uint8_t *data, size_t len);

        void packet_generate(unsigned char* packet, unsigned char code, unsigned char parameterCode, unsigned char parameterExtension) {
          packet[0] = this->address_;
          packet[1] = code;
          packet[2] = parameterCode;
          packet[3] = parameterExtension;

          auto crc = this->crc16(packet, 4);
          packet[4] = crc >> 0;
          packet[5] = crc >> 8;
        }

        uint16_t crc16(const uint8_t *data, uint8_t len) {
          uint16_t crc = 0xFFFF;
          while (len--) {
            crc ^= *data++;
            for (uint8_t i = 0; i < 8; i++) {
              if ((crc & 0x01) != 0) {
                crc >>= 1;
                crc ^= 0xA001;
              } else {
                crc >>= 1;
              }
            }
          }
          return crc;
        }

        template <size_t N = 3>
        long to_long(std::vector<uint8_t> payload, std::size_t start) {
          long out = 0;

          if (N < 2 || N > 4 || payload.size() < start + N) {
            return out;
          }

          if (N == 2) { // 16
            out |= payload[start+1] << 8;
            out |= payload[start];
          } else if (N == 3) { // 32
            out |= payload[start]   << 16;
            out |= payload[start+2] << 8;
            out |= payload[start+1];
          } else if (N == 4) { // 32
            out |= payload[start+1] << 24;
            out |= payload[start]   << 16;
            out |= payload[start+3] << 8;
            out |= payload[start+2];
          }

          return out;
        }
    };
  }  // namespace mercury_v3
}  // namespace esphome
