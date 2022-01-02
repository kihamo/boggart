#pragma once

#include "esphome/core/component.h"
#include "esphome/components/uart/uart.h"
#include "esphome/components/sensor/sensor.h"

#define MERCURY_V1_READ_BUFFER_SIZE 40
#define MERCURY_V1_READ_REQUEST_SIZE 7
#define MERCURY_V1_WAIT_AFTER_SEND_REQUEST 30
#define MERCURY_V1_WAIT_AFTER_READ_RESPONSE 100

#define MERCURY_V1_FIELD_ADDRESS_LENGTH 4
#define MERCURY_V1_FIELD_COMMAND_LENGTH 1
#define MERCURY_V1_FIELD_CRC_LENGTH 2

namespace esphome {
  namespace mercury_v1 {
    class MercuryV1 : public PollingComponent, public uart::UARTDevice {
      public:
        MercuryV1() = default;

        void loop() override;
        void update() override;
        void setup() override;
        void dump_config() override;

        float get_setup_priority() const override;

        void set_voltage_sensor(sensor::Sensor *voltage_sensor) { voltage_sensor_ = voltage_sensor; }
        void set_amperage_sensor(sensor::Sensor *amperage_sensor) { amperage_sensor_ = amperage_sensor; }
        void set_power_sensor(sensor::Sensor *power_sensor) { power_sensor_ = power_sensor; }
        void set_tariff1_sensor(sensor::Sensor *tariff1_sensor) { tariff1_sensor_ = tariff1_sensor; }
        void set_tariff2_sensor(sensor::Sensor *tariff2_sensor) { tariff2_sensor_ = tariff2_sensor; }
        void set_tariff3_sensor(sensor::Sensor *tariff3_sensor) { tariff3_sensor_ = tariff3_sensor; }
        void set_tariff4_sensor(sensor::Sensor *tariff4_sensor) { tariff4_sensor_ = tariff4_sensor; }
        void set_tariffs_total_sensor(sensor::Sensor *tariffs_total_sensor) { tariffs_total_sensor_ = tariffs_total_sensor; }

        void set_address(uint32_t address) {
            address_[0] = address >> 24; // TODO: 0x00 для модели 200, если вбили не правильное число
            address_[1] = address >> 16;
            address_[2] = address >> 8;
            address_[3] = address;
        }

      protected:
        sensor::Sensor *voltage_sensor_;
        sensor::Sensor *amperage_sensor_;
        sensor::Sensor *power_sensor_;
        sensor::Sensor *tariff1_sensor_;
        sensor::Sensor *tariff2_sensor_;
        sensor::Sensor *tariff3_sensor_;
        sensor::Sensor *tariff4_sensor_;
        sensor::Sensor *tariffs_total_sensor_;

        uint8_t address_[MERCURY_V1_FIELD_ADDRESS_LENGTH];
        uint8_t read_buffer_[MERCURY_V1_READ_BUFFER_SIZE]{};
        uint8_t packet_buffer_[MERCURY_V1_READ_BUFFER_SIZE]{};

        double V, A, W;
        double T1, T2, T3, T4, TTotal;

        enum Command : uint8_t {
          READ_POWER_COUNTERS = 0x27,
          READ_PARAMS_CURRENT = 0x63,
        };

        unsigned char read_power_counters_request_[MERCURY_V1_READ_REQUEST_SIZE];
        unsigned char read_params_current_request_[MERCURY_V1_READ_REQUEST_SIZE];

        void read_from_uart();
        void clean_uart_buffer();

        void packet_generate(unsigned char* packet, unsigned char cmd) {
          memcpy(packet, this->address_, MERCURY_V1_FIELD_ADDRESS_LENGTH);
          packet[4] = cmd;
          auto crc = this->crc16(packet, MERCURY_V1_FIELD_ADDRESS_LENGTH + MERCURY_V1_FIELD_COMMAND_LENGTH);
          packet[5] = crc >> 0;
          packet[6] = crc >> 8;
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

        long pow(long a, int s) {
          long out = 1;
          for (int i = 0; i < s; i++) out *= a;
          return out;
        }

        template <size_t N = 2>
        long to_long(unsigned char *inp) {
          long out = 0;

          for (int i = 0; i < N; i++) {
            unsigned char v = inp[i];
            int p = this->pow(10, ((N - 1) - i) * 2);
            out += (((v >> 4) & 15) * 10 + (v & 15)) * p;
          }

          return out;
        }
    };
  }  // namespace mercury_v1
}  // namespace esphome
