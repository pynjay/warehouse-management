# Warehouse management system

Запуск проекта (сброка образов, поднятие контейнеров в дев режиме)

```shell
make all
```

## Форматы запросов

### Создание продуктов

```shell
curl -X POST http://localhost:3000/api/v1/warehouses \
-H "Content-Type: application/json" \
--data-raw $'{"products": [{"name": "Jacket 001", "size": "xxl"}]}' -vvv
```

ответ:

[
  {
    "product_id": 1,
    "name": "Jacket 001",
    "size": "xxl",
    "sku": "SKU-b29869c3-b0b1-4af2-a529-5c1b7d21fbed"
  }
]


### Добавление складов

```shell
curl -X POST http://localhost:3000/api/v1/warehouses \
-H "Content-Type: application/json" \
--data-raw $'{"name": "my_warehouse", "is_available": false}' -vvv
```

ответ:

{
  "ID": 1,
  "Name": "my_warehouse",
  "IsAvailable": false,
  "CreatedAt": "2024-05-27T23:00:48.909191537Z",
  "UpdatedAt": "2024-05-27T23:00:48.909191537Z"
}

### Добавление продуктов в склад

```shell
curl -X POST http://localhost:3000/api/v1/products/1 \
-H "Content-Type: application/json" \
--data-raw $'{"quantity": 500, "warehouse_id": 1}' -vvv
```

ответ:

{
  "WarehouseId": 1,
  "ProductId": 1
}

### Получение общего количества продуктов на складе

```shell
curl -X GET http://localhost:3000/api/v1/warehouses/1/count \
-H "Content-Type: application/json"`
```

ответ:

{
  "count": 2010
}

### Резервация товаров

```shell
curl -X POST http://localhost:3000/api/v1/reservations \
-H "Content-Type: application/json" \
--data-raw $'{"order_id": 123125, "items": [
  {
    "product_sku": "SKU-dddd1e06-39d0-4474-bd8d-0457581107d2",
    "quantity": 100
  },
  {
    "product_sku": "SKU-82598cb7-bd2c-409d-91c8-1f7ecdc73c3f",
    "quantity": 100
  }
]}' -vvv
```

ответ:

[
  {
    "reservation_id": 1,
    "order_id": 123125,
    "quantity": 100,
    "status": "pending"
  },
  {
    "reservation_id": 2,
    "order_id": 123125,
    "quantity": 100,
    "status": "pending"
  }
]

### Обновление статуса резерваций

```shell
curl -X PATCH http://localhost:3000/api/v1/reservations/1 \
-H "Content-Type: application/json" \
--data-raw $'{"status": "fulfilled"}' -vvv```

ответ:

{
  "reservation_id": 1,
  "order_id": 123125,
  "quantity": 100,
  "status": "fulfilled"
}
