# Probabilistic toolkit

Коллекция вероятностных структур и алгоритмов. Позволяет решать задачи:

* Membership testing или AMQ (Approximate Membership Query)
* Cardinality estimation
* Frequency estimation
* Similarity estimation
* Symmetric difference
* LSH (Local Sensitive Hashing)
* Heavy hitters

Все решения написаны с расчётом на использование в highload окружении и предлагают:

* компактные структуры с минимальным потреблением памяти
* нулевые или минимальные аллокации памяти
* отсутствие блокировок за счёт использования atomic операций
* поддержка конкурентного режима для использования в многопоточных средах
* использование SIMD оптимизаций
* гибкая инициализация (все вспомогательные структуры абстрагированы, например во всех структурах можно задать нужный алгоритм хэширования)
* каждая структура реализует единый (в рамках своей задачи) интерфейс, что позволяет легко их переключать между собой
* коробочное покрытие метриками

Полное дерево решений:

* [AMQ](amq/readme.ru.md)
  * [Bloom filter](amq/bloom_filter/readme.ru.md)
  * [Counting bloom filter](amq/bloom_filter/readme.ru.md)
  * [Cuckoo filter](amq/cuckoo_filter/readme.ru.md)
  * [Quotient filter](amq/quotient_filter)
  * [Xor filter](amq/xor_filter/readme.ru.md)
* [Cardinality estimation](cardinality/readme.ru.md)
  * [LogLog](cardinality/loglog)
  * [HyperLogLog](cardinality/hyperloglog)
  * [HyperBitBit](cardinality/hyperbitbit)
  * [Linear counting](cardinality/linear_counting)
* [Frequency estimation](frequency/readme.ru.md)
  * [Count-Min Sketch](frequency/cmsketch)
  * [Conservative Update Sketch](frequency/cusketch)
  * [Count Sketch](frequency/countsketch)
  * [TinyLFU](frequency/tinylfu)
  * [TinyLFU (EWMA)](frequency/tinylfu_ewma/readme.ru.md)
* [Similarity estimation](similarity/readme.ru.md)
  * [Cosine similarity](similarity/cosine)
  * [Jaccard similarity](similarity/jaccard)
  * [Hamming similarity](similarity/hamming)
* [Symmetric difference](symmetric/readme.ru.md)
  * [Odd Sketch](symmetric/oddsketch)
* [LSH](lsh/readme.ru.md)
  * [SimHash](lsh/simhash)
  * [MinHash](lsh/minhash)
  * [b-Bit MinHash](lsh/bbitminhash)
* [Shingle](shingle/readme.ru.md)
  * [Char](shingle/char.go)
  * [Word](shingle/word.go)
* [Heavy hitters](heavy/readme.ru.md)
  * [Space-Saving](heavy/spacesaving/readme.ru.md)
  * [Misra-Gries](heavy/misragries/readme.ru.md)
  * [Lossy Counting](heavy/lossy/readme.ru.md)

Ниже есть краткое описание каждой задачи. Описание конкретных алгоритмов можно найти в соответствующих разделах.

## AMQ (Approximate Membership Query)

AMQ структуры решают задачу membership testing - определение принадлежности ключа к множеству. Обычно для этой задачи
используются хэш-таблицы, но они допустимы только для небольших множеств. AMQ же позволяет хранить очень большие
множества, жертвуя взамен точностью - возможны ложноположительные результаты, но ложноотрицательные - нет.

[Подробное описание](amq/readme.ru.md).

## Cardinality estimation

Cardinality estimation структуры решают задачу определения количества уникальных ключей в множестве. Эту задачу также
можно решить с помощью хэш-таблиц, но на больших множествах расход памяти слишком велик. Cardinality estimation структуры
позволяют уменьшить расход памяти до минимума, но взамен выдают приблизительный результат.

[Подробное описание](cardinality/readme.ru.md).

## Frequency estimation

Frequency estimation структуры решают задачу определения частоты ключей в множестве. Подобно прочим вероятностным структурам
они потребляют минимум памяти за счёт уменьшения точности результатов.

[Подробное описание](frequency/readme.ru.md).

## Similarity estimation

Similarity estimation структуры решают задачу определения схожести двух множеств. В этом пакете в качестве множеств
выступют только строковые типы данных, таким образом они решают задачу нечёткого сравнения двух текстов. Для работы
алгоритмы требуют указания вспомогательной LSH структуры и все этапы, включая шинглирование, хэширование и векторизацию
выполняются неявно для пользователя.

[Подробное описание](similarity/readme.ru.md).

## LSH (Local Sensitive Hashing)

В теории LSH это метод приближенного поиска соседей, который хэширует входные данные таким образом, что похожие объекты
с высокой вероятностью попадают в один бакет. В рамках этого пакета LSH работает только со строковыми типами данных и
занимается только векторизацией заранее шинглированных текстов. Полученные векторы далее передаются в similarity estimation
структуру для определения схожести двух текстов.
На практике LSH инициализируется вспомогательным Shingle алгоритмом и занимается шинглированием неявно для пользователя.

[Подробное описание](lsh/readme.ru.md).

## Shingle

Шинглирование это способ токенизации текста. В рамках этого пакета реализованы алгоритмы шинглирования по символам и по словам.
Далее предполагается отправка шинглированных текстов в LSH структуру для векторизации и затем отправка векторов в similarity estimation
структуру для определения схожести двух текстов.

[Подробное описание](shingle/readme.ru.md).

## Symmetric difference

Symmetric difference структуры решают задачу определения симметрической разности двух множеств. В рамках данного пакета
реализованы алгоритмы для определения симметрической разности двух текстов. Практически симметрическая разность
это операция обратная к similarity estimation - чем более похожи тексты, тем меньше их симметрическая разность.
Аналогично similarity estimation структурам, symmetric difference требует для работы вспомогательной LSH структуры.

[Подробное описание](symmetric/readme.ru.md).

## Heavy hitters

Heavy hitters алгоритмы решают задачу идентификации наиболее часто встречающихся элементов в потоке данных. Они подходят
для случаев когда поток данных настолько большой, что традиционные методы (такие как хэш-таблицы) оказываются неэффективными
по расходу расурсов. Подобно прочим вероятностным структурам они дают приблизительные результаты, но потребляют при этом
минимум ресурсов.

> [!IMPORTANT]
> Приведённые реализации не являются lock-free структурами и тем самым являются исключением из правил для этого репозитория.
> Они по умолчанию работают в режиме конкурентного доступа.

[Подробное описание](heavy/readme.ru.md)

## Заключение

Реализованные структуры позволяют проводить анализ больших данных или потоков данных в реальном времени с минимальным
потреблением ресурсов и оптимальной производительностью. Абстракции позволяют легко переключаться между различными
алгоритмами и выбрать оптимальный для конкретной задачи. Конкурентный режим позволит работать в многопоточной среде без
блокировок. А метрики помогут оценить оптимально ли настроена структура для конкретной задачи и донастроить при необходимости.
