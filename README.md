# Backend 2025_2_Avrora
Команда Аврора делаем Cian

### Репозиторий backend TBD

### Наша команда(backend)

Степан @StephenMarkman
Руслан @Pasatechnik
Владислав П @vlpanichkin
Владислав Л @Vladislav_L_F

### Тесты

```bash
go test ./... -coverprofile=cover.out && go tool cover -html=cover.out -o=cover.html && open cover.html
```

### Создание БД

username=postgres
password=postgres
dbname="2025_2_Avrora"

```bash
psql -u postgres -h localhost
create database "2025_2_Avrora"
```

### Тестовая бд

```bash
psql -u postgres -h localhost
create database "2025_2_Avrora_test"
``` 