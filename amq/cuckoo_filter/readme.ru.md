# Cuckoo Filter

Cuckoo Filter - это вероятностная структура данных для проверки принадлежности элемента множеству (membership testing),
которая сочетает в себе преимущества Cuckoo Hashing и Bloom фильтров. В отличие от Bloom фильтра,
Cuckoo Filter поддерживает удаление элементов без потери точности.

## Особенности реализации

* Поддержка пользовательских хэш-функций
* Подключаемый режим конкуренции - безопасные асинхронные чтение и запись
* Отсутствие блокировок за счёт atomic операций
* Использование SIMD инструкций для ускорения операций
* Фиксированный размер бакета (4 элемента)
* Автоматический расчёт размера таблицы на основе ожидаемого количества элементов и kicks limit

## Принцип работы Cuckoo Hashing

* Каждый элемент может находиться в одной из двух позиций, определяемых хэш-функциями
* Каждая позиция (бакет) содержит несколько отпечатков (fingerprints) элемента
* При коллизии происходит вытеснение (kicking) существующего отпечатка в его альтернативную позицию

## Математические основы

### Формула оптимального размера таблицы

Функция `optimalM` вычисляет оптимальный размер таблицы (количество бакетов) для заданного максимального количества элементов `n`:

$$
m = \frac{2^{\lceil \log_2(n) \rceil}}{b}
$$

где:
- `b` - размер бакета (в данной реализации жестко задан как 4)
- `⌈log₂(n)⌉` - округление вверх до ближайшей степени двойки

Эта формула обеспечивает:
* Размер таблицы равный степени двойки, что позволяет использовать быстрые битовые операции вместо дорогих делений
* Нагрузочный коэффициент около 95% при размере бакета 4
* Минимизацию вероятности бесконечных циклов при вытеснениях

### Вероятность ложноположительного срабатывания

Вероятность ложного срабатывания для Cuckoo Filter вычисляется как:

$$
\epsilon \approx \frac{2b}{2^f}
$$

где:
- `f` - длина отпечатка в битах
- `b` - размер бакета

## Пример рассчёта размера

Для `n = 1000`:
* Находим ближайшую степень двойки: `2^10 = 1024`
* Делим на размер бакета (4): `1024 / 4 = 256`
* Итоговый размер таблицы: 256 бакетов

## Использование

```go
package main

import (
	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/amq/cuckoo_filter"
	"github.com/koykov/pbtk/metrics/prometheus"
)

const N = 1e7

func main() {
	hasher := xxhash.Hasher64[[]byte]{} // hash function
	config := cuckoo.NewConfig(N, hasher).
		WithKicksLimit(10).                                    // limit for cuckoo kicks to avoid infinite loop
		WithConcurrency().                                     // switch to race protected buckets array (atomic based)
		WithMetricsWriter(prometheus.NewAMQ("example_filter")) // cover with metrics
	f, err := cuckoo.NewFilter[string](config)
	_ = err
	_ = f.Set("foobar")
	print(f.Contains("foobar")) // true
	print(f.Contains("qwerty")) // false
}

```

## Области применения

* Кэширование - быстрая проверка наличия элемента в кэше
* Сетевые приложения - фильтрация дубликатов пакетов
* Базы данных - ускорение запросов за счёт отсечения бесперспективных запросов
* Системы мониторинга - отслеживание уникальных событий
* Распределенные системы - дедупликация данных

## Заключение

Эта реализация предоставляет эффективный инструмент для вероятностного хранения данных с поддержкой конкурентного доступа.
Использование атомарных операций и SIMD инструкций позволяет достичь высокой производительности, а математически
обоснованный выбор параметров обеспечивает оптимальное использование памяти.

Фильтр особенно полезен в сценариях, где требуется частое обновление данных при ограниченных ресурсах памяти и высокой нагрузке.
