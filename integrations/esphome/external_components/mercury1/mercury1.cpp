#include "mercury1.h"
#include "esphome/core/log.h"

namespace esphome {
  namespace mercury1 {
    static const char *const TAG = "mercury1";

    void Mercury1::setup() {
      // Clear UART buffer
      while (this->available())
        this->read();

      // TODO: generate commands for update
    }

    void Mercury1::loop() {
      if (!this->available())
        return;

      int len;
      memset(this->read_buffer_, 0, MERCURY1_READ_BUFFER_SIZE);
      this->read_index_ = 0;

      while (this->available()) {
        if(this->read_index_ > MERCURY1_READ_BUFFER_SIZE) {
          ESP_LOGW(TAG, "Buffer overflow. Total length %d", this->read_index_);

          // Clear UART buffer
          while (this->available())
            this->read();

          break;
        }

        uint8_t byte;
        this->read_byte(&byte);

        this->read_buffer_[this->read_index_] = byte;
        this->read_index_++;
      }

      delay(100);

      /*
      while((len = this->available()) > 0) {
        len = std::min(len, MERCURY1_READ_BUFFER_SIZE);
        this->read_index_ += len;

        if(this->read_index_ > MERCURY1_READ_BUFFER_SIZE) {
          ESP_LOGW(TAG, "Buffer overflow. Total length %d", this->read_index_);

          while (this->available())
            this->read();

          return;
        }

        this->read_array(this->read_buffer_, len);
      }
      */

      ESP_LOGV(TAG, "Response raw %s", hexencode(this->read_buffer_, this->read_index_).c_str());

      // ошибочное начало пакета
      if(this->read_buffer_[0] != 0x00) { // включительно, чтобы пропускать пакеты на отсылку команд
        ESP_LOGW(TAG, "Response first byte isn't 0x00, is %02X (raw %s)", this->read_buffer_[0], hexencode(this->read_buffer_, this->read_index_).c_str());
        return;
      }

      if(this->read_index_ <= 7) { // включительно, чтобы пропускать пакеты на отсылку команд
        ESP_LOGD(TAG, "Skip response with length %d", this->read_index_);
        return;
      }

      ESP_LOGD(TAG, "Response command %02X", this->read_buffer_[4]);

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
            this->W = this->to_double<3>(&this->read_buffer_[9], 100);

            this->voltage_sensor_->publish_state(this->V);
            this->amperage_sensor_->publish_state(this->A);
            this->power_sensor_->publish_state(this->W);
          break;

        default:
            ESP_LOGW(TAG, "Unknown response command 0x%02X", this->read_buffer_[4]);
          break;
      }
    }

    void Mercury1::update() {
        // TODO: Send update commands
    }

    void Mercury1::dump_config() {
      ESP_LOGCONFIG(TAG, "Mercury v1:");
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