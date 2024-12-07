# Оценка производительности фреймворков БД

>
> Сравниваемые фреймворки
>
> sqlx + pgx
>
> gorm
>

## Результаты измерений
![](./res/benchmark/benchmark.svg)

## Утилизация ресурсов

### CPU
#### SQLX
![](./res/benchmark/sqlx/cpu_utilization.png)
#### GORM
![](./res/benchmark/gorm/cpu_utilization.png)

### DISK
#### SQLX
![](./res/benchmark/sqlx/disk_io.png)
#### GORM
![](./res/benchmark/gorm/disk_io.png)

### GC stops
#### SQLX
![](./res/benchmark/sqlx/gc_duration_quantile.png)
#### GORM
![](./res/benchmark/gorm/gc_duration_quantile.png)

### МБ на сценарий
#### SQLX
![](./res/benchmark/sqlx/mb_per_op_quantile.png)
#### GORM
![](./res/benchmark/gorm/mb_per_op_quantile.png)

### МС на сценарий
#### SQLX
![](./res/benchmark/sqlx/ms_per_op_quantile.png)
#### GORM
![](./res/benchmark/gorm/ms_per_op_quantile.png)

### Число выделений на сценарий
#### SQLX
![](./res/benchmark/sqlx/alloca_per_op_quantile.png)
#### GORM
![](./res/benchmark/gorm/alloca_per_op_quantile.png)

