# Xor Filter

Xor Filter — это вероятностная структура данных, предназначенная для эффективной проверки принадлежности элемента к множеству.
Она является компактной альтернативой Bloom-фильтрам с более высокой производительностью и меньшим количеством ложных срабатываний.

Xor Filter требует предварительного построения на фиксированном наборе ключей. Алгоритм построения включает:
* Создание графа зависимостей между ключами
* Поиск решения системы уравнений, где каждый ключ связан с несколькими позициями в битовом массиве
* Заполнение массива таким образом, что XOR выбранных позиций дает хэш ключа

Добавление новых ключей потребует полного перестроения структуры, так как изменит все существующие зависимости.

## Особенности реализации

* Реализует XorBinaryFuse8
* Поддержка пользовательских хэш-функций
* Переиспользование фильтра через `sync.Pool` для минимизации аллокаций
* Отсутствие блокировок (lock-free)
* Компактное хранение данных

## Преимущества использования

1. **Меньший размер** по сравнению с Bloom-фильтрами
2. **Быстрые проверки** — только 3 обращения к памяти и 2 операции XOR
3. **Нет ложных отрицаний** — если ключ был добавлен, он всегда будет найден
4. **Низкая вероятность ложных срабатываний** — около 0.4%
5. **Параллелизм** — реализация не использует блокировки

## Математическое обоснование

Фильтр XOR строится в три этапа:

1. **Распределение ключей**:
   Для каждого ключа $x$ вычисляются позиции:

$$
h_1(x), h_2(x), h_3(x) \in \{0, \ldots, m-1\}
$$
   
   где $m$ - размер фильтра, обычно $m \approx 1.23 \cdot n$ для $n$ ключей.

2. **Вычисление fingerprint ключа**:
   Для ключа $x$ вычисляется:
   
$$
\text{fingerprint}(x) = \text{filter}[h_1(x)] \oplus \text{filter}[h_2(x)] \oplus \text{filter}[h_3(x)]
$$
   
   где $\text{fingerprint}(x)$ - 8-битный хэш ключа.

3. **Проверка вхождения ключа в множество**:
   Для ключа $y$ проверяется:

$$
\text{Contains}(y) = \left(\text{filter}[h_1(y)] \oplus \text{filter}[h_2(y)] \oplus \text{filter}[h_3(y)]\right) == \text{fingerprint}(y)
$$

при этом веротность ложноположительного результата составляет:

$$
\text{Pr}(\text{false positive}) \leq \frac{1}{2^8} = \frac{1}{256} \approx 0.39\%
$$

## Использование

```go
package main

import (
	"github.com/koykov/hash/xxhash"
	xor "github.com/koykov/pbtk/amq/xor_filter"
	"github.com/koykov/pbtk/metrics/prometheus"
)

func main() {
	hasher := xxhash.Hasher64[[]byte]{} // hash function
	config := xor.NewConfig(hasher).
		WithMetricsWriter(prometheus.NewAMQ("example_filter")) // cover with metrics
	f, err := xor.NewFilterWithKeys[[]byte](config, [][]byte{
		[]byte("foo"),
		[]byte("bar"),
	})
	_ = err

	println(f.Contains([]byte("foo"))) // true
	println(f.Contains([]byte("qwe"))) // false
}
```

## Примеры применения

1. **Кэширование** — быстрая проверка наличия данных в кэше перед дорогостоящим запросом
2. **Базы данных** — предфильтрация запросов для уменьшения обращений к диску
3. **Сетевые фильтры** — блокировка нежелательных IP-адресов или URL
4. **Уникальность элементов** — проверка на дубликаты в потоке данных
5. **Поисковые системы** — предварительный отбор документов для полнотекстового поиска

## Заключение

Xor Filter предоставляет эффективный и компактный способ проверки принадлежности элементов с минимальными аллокациями памяти и высокой параллельной производительностью.
Фильтр особенно полезен в сценариях, где важны скорость проверки и экономия памяти, а невозможность динамического добавления элементов не является критическим ограничением.
