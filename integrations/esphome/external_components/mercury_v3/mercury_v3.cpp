#include "mercury_v3.h"
#include "esphome/core/log.h"

namespace esphome {
  namespace mercury_v3 {
    static const char *const TAG = "mercury_v3";

    float MercuryV3::get_setup_priority() const {
      return setup_priority::DATA;
    }

    void MercuryV3::setup() {
      this->last_open_channel_ = 0;

      channel_open_request_[0] = this->address_;
      channel_open_request_[1] = 0x1; // code open channel
      channel_open_request_[2] = 0x1; // level
      channel_open_request_[3] = 0x1; // password 1
      channel_open_request_[4] = 0x1; // password 2
      channel_open_request_[5] = 0x1; // password 3
      channel_open_request_[6] = 0x1; // password 4
      channel_open_request_[7] = 0x1; // password 5
      channel_open_request_[8] = 0x1; // password 6
      auto crc = this->crc16(channel_open_request_, 9);
      channel_open_request_[9] = crc >> 0;
      channel_open_request_[10] = crc >> 8;

      this->packet_generate(read_voltage_request_, 0x8, 0x16, 0x1<<4 | 0x1);
      this->packet_generate(read_current_request_, 0x8, 0x16, 0x2<<4 | 0x1);
      this->packet_generate(read_power_request_, 0x8, 0x16, 0x0<<4 | 0x0 << 2 | 0x0);
      this->packet_generate(read_tariff_1_request_, 0x5, 0x0, 0x1);
    }

    void MercuryV3::loop() {

    }

    void MercuryV3::update() {
      this->update_voltage();
      this->update_current();
      this->update_power();
      this->update_tariffs();
    }

    void MercuryV3::dump_config() {
      ESP_LOGCONFIG(TAG, "Mercury v3:");
      LOG_UPDATE_INTERVAL(this);

      LOG_SENSOR("  ", "Voltage A", this->phase_[0].voltage_sensor_);
      LOG_SENSOR("  ", "Current A", this->phase_[0].current_sensor_);
      LOG_SENSOR("  ", "Power A", this->phase_[0].power_sensor_);
      LOG_SENSOR("  ", "Voltage B", this->phase_[1].voltage_sensor_);
      LOG_SENSOR("  ", "Current B", this->phase_[1].current_sensor_);
      LOG_SENSOR("  ", "Power B", this->phase_[1].power_sensor_);
      LOG_SENSOR("  ", "Voltage C", this->phase_[2].voltage_sensor_);
      LOG_SENSOR("  ", "Current C", this->phase_[2].current_sensor_);
      LOG_SENSOR("  ", "Power C", this->phase_[2].power_sensor_);
      LOG_SENSOR("  ", "Tariff 1", this->tariff1_sensor_);
      LOG_SENSOR("  ", "Power", this->power_sensor_);

      this->check_uart_settings(9600);
    }

    void MercuryV3::update_voltage() {
        if (this->phase_[0].voltage_sensor_ == nullptr && this->phase_[1].voltage_sensor_ == nullptr && this->phase_[2].voltage_sensor_ == nullptr) {
            return;
        }

        if (!this->invoke(read_voltage_request_, sizeof(read_voltage_request_))) {
            return;
        }

        if (this->payload_.size() != 9) {
            ESP_LOGW(TAG, "Wrong size of voltage response packet have %d want 9", this->payload_.size());
            return;
        }

        if (this->phase_[0].voltage_sensor_ != nullptr) {
            this->phase_[0].voltage_sensor_->publish_state((float) this->to_long<3>(this->payload_, 0) / 100);
        }
        if (this->phase_[1].voltage_sensor_ != nullptr) {
            this->phase_[1].voltage_sensor_->publish_state((float) this->to_long<3>(this->payload_, 3) / 100);
        }
        if (this->phase_[2].voltage_sensor_ != nullptr) {
            this->phase_[2].voltage_sensor_->publish_state((float) this->to_long<3>(this->payload_, 6) / 100);
        }
    }

    void MercuryV3::update_current() {
        if (this->phase_[0].current_sensor_ == nullptr && this->phase_[1].current_sensor_ == nullptr && this->phase_[2].current_sensor_ == nullptr) {
            return;
        }

        if (!this->invoke(read_current_request_, sizeof(read_current_request_))) {
            return;
        }

        if (this->payload_.size() != 9) {
            ESP_LOGW(TAG, "Wrong size of current response packet have %d want 9", this->payload_.size());
            return;
        }

        if (this->phase_[0].current_sensor_ != nullptr) {
            this->phase_[0].current_sensor_->publish_state((float) this->to_long<3>(this->payload_, 0) / 1000);
        }
        if (this->phase_[1].current_sensor_ != nullptr) {
            this->phase_[1].current_sensor_->publish_state((float) this->to_long<3>(this->payload_, 3) / 1000);
        }
        if (this->phase_[2].current_sensor_ != nullptr) {
            this->phase_[2].current_sensor_->publish_state((float) this->to_long<3>(this->payload_, 6) / 1000);
        }
    }

    void MercuryV3::update_power() {
        if (this->phase_[0].power_sensor_ == nullptr && this->phase_[1].power_sensor_ == nullptr && this->phase_[2].power_sensor_ == nullptr) {
            return;
        }

        if (!this->invoke(read_power_request_, sizeof(read_power_request_))) {
            return;
        }

        if (this->payload_.size() != 12) {
            ESP_LOGW(TAG, "Wrong size of power response packet have %d want 12", this->payload_.size());
            return;
        }

        if (this->power_sensor_ != nullptr) {
            this->power_sensor_->publish_state((float) this->to_long<2>(this->payload_, 1) / 100);
        }
        if (this->phase_[0].power_sensor_ != nullptr) {
            this->phase_[0].power_sensor_->publish_state((float) this->to_long<2>(this->payload_, 4) / 100);
        }
        if (this->phase_[1].power_sensor_ != nullptr) {
            this->phase_[1].power_sensor_->publish_state((float) this->to_long<2>(this->payload_, 7) / 100);
        }
        if (this->phase_[2].power_sensor_ != nullptr) {
            this->phase_[2].power_sensor_->publish_state((float) this->to_long<2>(this->payload_, 10) / 100);
        }
    }

    void MercuryV3::update_tariffs() {
        if (this->tariff1_sensor_ == nullptr) {
            return;
        }

        if (!this->invoke(read_tariff_1_request_, sizeof(read_tariff_1_request_))) {
            return;
        }

        if (this->payload_.size() != 16) {
            ESP_LOGW(TAG, "Wrong size of tariff #1 response packet have %d want 16", this->payload_.size());
            return;
        }

        // активная прямая (А+)
        this->tariff1_sensor_->publish_state(this->to_long<4>(this->payload_, 0));
    }

    void MercuryV3::uart_tx(const uint8_t *data, size_t len) {
        ESP_LOGV(TAG, "UART TX size %d payload %s", len, format_hex_pretty(data, len).c_str());

        // clean uart buffer
        while (this->available()) {
            this->read();
        }

        this->write_array(data, len);

        this->flush();
        delay(MERCURY_V3_WAIT_AFTER_SEND_REQUEST);
    }

    bool MercuryV3::uart_rx() {
        this->response_.clear();
        while (this->available()) {
          uint8_t byte = this->read();
          this->response_.push_back(byte);
        }

         ESP_LOGV(TAG, "UART RX size %d payload %s", this->response_.size(), format_hex_pretty(this->payload_.data(), this->payload_.size()).c_str());

        if (this->payload_.size() < MERCURY_V3_FIELD_ADDRESS_LENGTH + MERCURY_V3_FIELD_COMMAND_LENGTH + MERCURY_V3_FIELD_CRC_LENGTH) {
            ESP_LOGW(TAG, "Payload is shot have %d want %d", this->response_.size(), MERCURY_V3_FIELD_ADDRESS_LENGTH + MERCURY_V3_FIELD_COMMAND_LENGTH + MERCURY_V3_FIELD_CRC_LENGTH);
            return false;
        }

        this->payload_.assign(this->response_.begin() + MERCURY_V3_FIELD_ADDRESS_LENGTH, this->response_.end() - MERCURY_V3_FIELD_CRC_LENGTH);

        // check address byte
        auto address = *this->response_.begin();
        if (address_ != this->address_) {
            ESP_LOGW(TAG, "Wrong address of response packet have %d want %d", address_, this->address_);
            return false;
        }

        // check checksum
        auto payload_crc = std::vector<uint8_t>(this->response_.begin(), this->response_.end() - MERCURY_V3_FIELD_CRC_LENGTH);
        uint16_t computed_crc = this->crc16(payload_crc.data(), payload_crc.size());
        uint16_t remote_crc = uint16_t(this->response_[this->response_.size() - MERCURY_V3_FIELD_CRC_LENGTH]) | (uint16_t(this->response_[this->response_.size() - MERCURY_V3_FIELD_ADDRESS_LENGTH]) << 8);

        if (computed_crc != remote_crc) {
            ESP_LOGW(TAG, "CRC16 check failed! computed %02X != remote %02X", computed_crc, remote_crc);
            return false;
        }

        // check error in response
        if (this->payload_.size() == 1) {
            switch (*this->payload_.begin()) {
                case 0x1: // Недопустимая команда или параметр
                    ESP_LOGE(TAG, "Invalid command or parameter");
                    return false;
                case 0x2: // Внутренняя ошибка счётчика
                    ESP_LOGE(TAG, "Internal meter error");
                    return false;
                case 0x3: // Недостаточен уровень для удовлетворения запроса
                    ESP_LOGE(TAG, "Insufficient level to satisfy the request");
                    return false;
                case 0x4: // Внутренние часы счётчика уже корректировались в течение текущих суток
                    ESP_LOGE(TAG, "Internal clock of the meter has already been corrected during the current day");
                    return false;
                case 0x5: // Не открыт канал связи
                    ESP_LOGE(TAG, "Communication channel not open");
                    return false;
            }
        }

        delay(MERCURY_V3_WAIT_AFTER_READ_RESPONSE);

        return true;
    }

    bool MercuryV3::invoke(const uint8_t *data, size_t len) {
        const uint32_t now = millis();

        /*
        Если значение байта состояния обмена в последовательности ответа равно нулю, то разрешается доступ к данным в
        течение 240 секунд, т.е. счётчик, будет отвечать на запросы в соответствии с уровнем доступа, определяемым
        введённым паролем. Каждый сле- дующий корректный запрос к счётчику переустанавливает таймер открытого канала в
        исходное состояние, т.е. на 240 секунд. Если к счётчику не было запросов в течение 240 секунд, то канал
        автоматически закрывается.
        */
        if (this->last_open_channel_ == 0 || (now >= MERCURY_V3_CHANNEL_OPEN_TIMEOUT && now - this->last_open_channel_ >= MERCURY_V3_CHANNEL_OPEN_TIMEOUT)) {
            this->uart_tx(channel_open_request_, sizeof(channel_open_request_));
            if (!this->uart_rx()) {
                return false;
            }

            this->last_open_channel_ = now;
        }

        this->uart_tx(data, len);

        if (!this->uart_rx()) {
            return false;
        }

        this->last_open_channel_ = now;

        return true;
    }
  }  // namespace mercury_v3
}  // namespace esphome
