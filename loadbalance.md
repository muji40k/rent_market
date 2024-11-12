# Результаты нагрузочного тестирования

```
Утилита нагузочного тестирования: apache benchmark
URL: /api/v1/products?search=Lorem%20ipsum%20Fermentum&category=923297f7-3dc3-4d02-a4ef-f82dd4ce1fb9
Время выполнения рассчитывается как средние из 10 измерений
```

## 1 Пользователь
![](res/loadbalancing/1.svg)

## 100 Пользователей
![](res/loadbalancing/100.svg)

## 1000 Пользователей
![](res/loadbalancing/1000.svg)

# Попытка 2: Увеличено число worker'ов

## 1 Пользователь
![](res/loadbalancing2/1.svg)

## 100 Пользователей
![](res/loadbalancing2/100.svg)

## 1000 Пользователей
![](res/loadbalancing2/1000.png)


