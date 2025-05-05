# Similarity Estimation

Similarity estimation (оценка схожести) - это методы измерения степени похожести между двумя наборами данных.
Эти алгоритмы решают задачи:
* Поиска дубликатов или почти дубликатов
* Кластеризации похожих документов
* Рекомендательных систем
* Обнаружения плагиата
* Удаления похожих записей в данных

## Реализованные алгоритмы
* **Hamming Distance** - Измеряет количество отличающихся бит между двумя векторами.
  Эффективен для сравнения бинарных данных или хэшей фиксированной длины.
* **Cosine Similarity** - Оценивает схожесть по углу между векторами в многомерном пространстве.
  Широко используется для текстовых данных, представленных как вектора признаков.
* **Jaccard Distance** - Вычисляет меру различия между множествами как долю несовпадающих элементов.
  Хорошо подходит для сравнения наборов слов или шинглов.

## Особенности реализации

* **Высокая производительность**: Минимизированы аллокации памяти, использованы эффективные структуры данных
* **Единый интерфейс**: все алгоритмы реализуют интерфейс `Estimator`
* **Взаимозаменяемость**: алгоритмы можно менять без изменения основного кода
* **Интеграция с LSH**: работают с векторами, полученными из LSH-алгоритма

## Использование

```go
package main

import (
	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh/minhash"
	"github.com/koykov/pbtk/shingle"
	"github.com/koykov/pbtk/similarity/jaccard"
)

func main() {
	hasher := xxhash.Hasher64[[]byte]{}
	shingler := shingle.NewChar[[]byte](3, "") // 3-gram
	lsh, _ := minhash.NewHasher[[]byte](minhash.NewConfig[[]byte](hasher, 50, shingler))
	est, err := jaccard.NewEstimator[[]byte](jaccard.NewConfig[[]byte](lsh))
	_ = err

	e, _ := est.Estimate([]byte("Four children are doing backbends in the gym"), []byte("Four children are doing backbends in the park"))
	println(e) // 0.8478260869565217 (high similarity)

	est.Reset()
	e, _ = est.Estimate([]byte("A man is sitting near a bike and is writing a note"), []byte("A man is standing near a bike and is writing on a piece paper"))
	println(e) // 0.532258064516129 (medium similarity)

	est.Reset()
	e, _ = est.Estimate([]byte("One white dog and one black one are sitting side by side on the grass"), []byte("A black and a white dog are joyfully running on the grass"))
	println(e) // 0.44155844155844154 (low similarity)
}
```

## Примеры применения

1. **Поиск похожих документов** в большой коллекции текстов
2. **Выявление дубликатов** товаров в каталоге интернет-магазина
3. **Рекомендации контента** на основе схожести с ранее просмотренным
4. **Обнаружение плагиата** в академических работах
5. **Кластеризация новостей** на одну тему из разных источников

## Заключение

Реализованные алгоритмы предоставляют гибкий инструментарий для решения широкого круга задач сравнения текстовых данных.
Благодаря единому интерфейсу и оптимизациям производительности, они могут быть легко интегрированы в существующие системы обработки данных.
