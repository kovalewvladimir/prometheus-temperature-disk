# Prometheus Temperature Disk Exporter

Экспортер для Prometheus, который собирает данные о температуре дисков и их режиме питания с помощью smartctl.

## Требования

- smartmontools (для работы с S.M.A.R.T)
- Docker (опционально)

## Пример метрик

![Пример метрик в Prometheus](prometheus-metrics-example.png)

## Установка

### Без Docker

```bash
go build
./prometheus-temperature-disk
```

### С Docker

```bash
docker build -t prometheus-temperature-disk .
docker run --privileged -p 9586:9586 prometheus-temperature-disk
```

## Конфигурация

### Переменные окружения

- `EXCLUDED_DEVICES` - список исключаемых устройств через запятую (например: "sda,sdb")

## Метрики

- `disk_temperature_celsius` - температура диска в градусах Цельсия
- `disk_power_mode` - режим питания диска (1=ACTIVE, 0=STANDBY)

### Формат метрик

```
# HELP disk_temperature_celsius Current temperature of the disk
# TYPE disk_temperature_celsius gauge
disk_temperature_celsius{device="sda",path="/dev/sda"} 35

# HELP disk_power_mode Current power mode of the disk (1=ACTIVE, 0=STANDBY)
# TYPE disk_power_mode gauge
disk_power_mode{device="sda",path="/dev/sda"} 1
```

## Особенности работы

- Температура диска (`disk_temperature_celsius`) отображается только для дисков в активном режиме (ACTIVE)
- Режим питания (`disk_power_mode`) отображается для всех поддерживаемых дисков
- Экспортер не пробуждает диски в режиме ожидания (STANDBY), чтобы не влиять на энергопотребление и срок службы

## Prometheus конфигурация

```yaml
scrape_configs:
  - job_name: 'disk_temperature'
    static_configs:
      - targets: ['localhost:9586']
```

## Поддерживаемые устройства

- SATA диски (sd[a-y])
- NVMe диски (nvme*)
