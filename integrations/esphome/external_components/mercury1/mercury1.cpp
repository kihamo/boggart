#include "mercury1.h"
#include "esphome/core/log.h"

namespace esphome {
  namespace mercury1 {
    static const char *const TAG = "mercury1";

    float Mercury1::get_setup_priority() const {
      return setup_priority::DATA;
    }

    void Mercury1::setup() {
      this->clean_uart_buffer();

      this->packet_generate(read_power_counters_request_, Command::READ_POWER_COUNTERS);
      this->packet_generate(read_params_current_request_, Command::READ_PARAMS_CURRENT);
    }

    void Mercury1::loop() {

    }

    void Mercury1::read_from_uart() {
      // ESP_LOGV(TAG, "Read from UART start");

      memset(this->read_buffer_, 0, MERCURY1_READ_BUFFER_SIZE);
      int offset = 0;

      while (this->available()) {
        delay(10); // FIXME: задержка не портит буфер, без задержки байты читаются рандомно

        if(offset > MERCURY1_READ_BUFFER_SIZE) {
          ESP_LOGW(TAG, "Buffer overflow");
          this->clean_uart_buffer();
          break;
        }

        this->read_byte(&this->read_buffer_[offset]);
        offset++;
      }

      ESP_LOGV(TAG, "Response raw %s", hexencode(this->read_buffer_, offset).c_str());

      // ошибочное начало пакета
      if(this->read_buffer_[0] != 0x00) { // включительно, чтобы пропускать пакеты на отсылку команд
        ESP_LOGW(TAG, "Response first byte isn't 0x00, is %02X (raw %s)", this->read_buffer_[0], hexencode(this->read_buffer_, offset).c_str());
        return;
      }

      // игнорируем пакеты на отсылку команд самому счетчику
      if(offset <= MERCURY1_READ_REQUEST_SIZE) {
        ESP_LOGD(TAG, "Skip response with length %d", offset);
        return;
      }

      // игнорируем пакеты с некорректной контрольной суммой, так как в эфире бывает дичь из обрывков пакетов
      uint16_t computed_crc = this->crc16(this->read_buffer_, offset - 2);
      uint16_t remote_crc = uint16_t(this->read_buffer_[offset - 2]) | (uint16_t(this->read_buffer_[offset - 1]) << 8);

      if (computed_crc != remote_crc) {
        ESP_LOGW(TAG, "CRC Check failed! computed %02X != remote %02X", computed_crc, remote_crc);
        return;
      }

      // обработка данных с валидных пакетов
      switch (this->read_buffer_[4]) {
        case Command::READ_POWER_COUNTERS:
            // ADDR-CMD-count*4-CRC
            this->T1 = this->to_double<4>(&this->read_buffer_[5], 1);
            this->T2 = this->to_double<4>(&this->read_buffer_[9], 1);
            this->T3 = this->to_double<4>(&this->read_buffer_[13], 1);
            this->T4 = this->to_double<4>(&this->read_buffer_[17], 1);
            this->TTotal = this->T1 +this->T2 + this->T3 + this->T4;

            this->tariff1_sensor_->publish_state(this->T1);
            this->tariff2_sensor_->publish_state(this->T2);
            this->tariff3_sensor_->publish_state(this->T3);
            this->tariff4_sensor_->publish_state(this->T4);
            this->tariffs_total_sensor_->publish_state(this->TTotal);
          break;

        case Command::READ_PARAMS_CURRENT:
            // ADDR-CMD-V-I-P-CRC
            this->V = this->to_double(&this->read_buffer_[5], 10);
            this->A = this->to_double(&this->read_buffer_[7], 100);
            this->W = this->to_double<3>(&this->read_buffer_[9], 1);

            this->voltage_sensor_->publish_state(this->V);
            this->amperage_sensor_->publish_state(this->A);
            this->power_sensor_->publish_state(this->W);
          break;

        default:
            ESP_LOGW(TAG, "Unknown response command 0x%02X", this->read_buffer_[4]);
          break;
      }
    }

    void Mercury1::clean_uart_buffer() {
      while (this->available()) {
        this->read();
      }
    }

    void Mercury1::update() {
      ESP_LOGV(TAG, "Send READ_POWER_COUNTERS %s", hexencode(read_power_counters_request_, MERCURY1_READ_REQUEST_SIZE).c_str());
      this->write_array(read_power_counters_request_, MERCURY1_READ_REQUEST_SIZE);

      delay(MERCURY1_WAIT_AFTER_SEND_REQUEST);

      this->read_from_uart();

      delay(MERCURY1_WAIT_AFTER_READ_RESPONSE);

      ESP_LOGV(TAG, "Send READ_PARAMS_CURRENT %s", hexencode(read_params_current_request_, MERCURY1_READ_REQUEST_SIZE).c_str());
      this->write_array(read_params_current_request_, MERCURY1_READ_REQUEST_SIZE);

      delay(MERCURY1_WAIT_AFTER_SEND_REQUEST);

      this->read_from_uart();

      delay(MERCURY1_WAIT_AFTER_READ_RESPONSE);
    }

    void Mercury1::dump_config() {
      ESP_LOGCONFIG(TAG, "Mercury v1:");
      ESP_LOGCONFIG(TAG, "Address %d", this->address_);
      LOG_UPDATE_INTERVAL(this);
      LOG_SENSOR("", "Voltage", this->voltage_sensor_);
      LOG_SENSOR("", "Amperage", this->amperage_sensor_);
      LOG_SENSOR("", "Power", this->power_sensor_);
      LOG_SENSOR("", "Tariff 1", this->tariff1_sensor_);
      LOG_SENSOR("", "Tariff 2", this->tariff2_sensor_);
      LOG_SENSOR("", "Tariff 3", this->tariff3_sensor_);
      LOG_SENSOR("", "Tariff 4", this->tariff4_sensor_);
      LOG_SENSOR("", "Tariffs total", this->tariffs_total_sensor_);

      this->check_uart_settings(9600);
    }
  }  // namespace mercury1
}  // namespace esphome
