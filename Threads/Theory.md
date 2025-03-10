
- Есть потребность в разбиении разных задач, напрямую не зависящих друг от друга, в распределении их на разные физические ядра и распараллеливании, либо в организации конкурентного выполнения

![](../_resources/Pasted%20image%2020250101210458.png) 
##### Решение 1 - исполнять такие задачи в разных однопоточных процессах
- Большая часть проблем связана с тем, что у ***процессов разные ресурсы*** - память, дескрипторы и проч --> из-за этого усложняется коммуникация между процессами, в принципе сложнее организовать работу разных процессов над одним ресурсом, потому что ***каждый процесс оперирует своим экземпляром ресурса***
- Поэтому появляются проблемы сериализации/десериализации данных сообщений для процессов(IPC), лишнего копирования данных и проч


![](../_resources/Pasted%20image%2020250101211104.png)

##### Решение 2 - исполнять такие задачи в рамках 1-ого многопоточного процесса
- Концепция потоков - единиц исполнения, разделяющих общие ресуры, такие как дескрипторы, память и проч., - устраняет проблемы первого подхода: мы бувально работаем с одной и той же памятью, с одними и теми же виртуальными адресами, которые неинвалидируются при передачи из 1-ого процесса в другой, напрямую вызываем функции, расположенные по общим адресам и проч


![](../_resources/Pasted%20image%2020250101211422.png)


![](../_resources/Pasted%20image%2020250101214419.png)
- Syscall для клонирования task struct - для создания потока нужно склонировать таблицу fd, директорию, pid, виртуальную память.
- thread id будут разными

![](../_resources/Pasted%20image%2020250101214336.png)
- Клонированный такс исполняется на своем стеке. Он заранее выделяется(на куче) и чистится после работы cloned таска(треда)



##### Volatile - отключение оптимизаций компилятора
Пример - https://slides.com/gerold103/sysprog_eng6#/12/0/5


##### Атомарность операций

![](../_resources/Pasted%20image%2020250101230157.png)
![](../_resources/Pasted%20image%2020250101230227.png)
- Проблема неатомарного инкремента - фактически инкремент это 3 операции, т.к. процессор умеет только загрузить в регистр(load), изменить значение в регистре, выгрузить значение из регистра в память(store)


![](../_resources/Pasted%20image%2020250101230457.png)
- Atomic API - инструменты для выполнения атомарных операций - ***такие функции будут компилироваться как неделимые атомарные ассемблерные инструкции*** 

> Атомарные операции - это userspace only операции, никаких syscalls не делается

![](../_resources/Pasted%20image%2020250101230628.png)
- Spin lock: атомарное CAS(лок перменной-флажка) - выполнение операций - cas(unlock флажка)


![](../_resources/Pasted%20image%2020250101230736.png)
- Проблема spin lock - это активное ожидание: вместо того чтобы отправить потоки спать, перепланировать их и тп, производится активная работа вхолостую

![](../_resources/Pasted%20image%2020250103163414.png)
![](../_resources/Pasted%20image%2020250103163241.png)
- Эту проблему призван решить futex(это и изменяемое значение, и название syscall) - системный вызов, отправляющий текущий поток спать в случае ***futex_wait***
- другой поток пробуждает спящие потоки, уведомляет их об изменении переменной futex
- в этом случае не производится активное ожидание, однако тут уже выполняется syscall

```C
int
futex_wait(int *futex, int val)
{
	return syscall(SYS_futex, futex, FUTEX_WAIT, val,
		       NULL, NULL, 0);
}

int
futex_wake(int *futex)
{
	return syscall(SYS_futex, futex, FUTEX_WAKE, 1,
		       NULL, NULL, 0);
}

void
futex_lock(int *futex)
{
	while (__sync_val_compare_and_swap(futex, 0, 1) != 0)
		futex_wait(futex, 1); // тут поток уйдет спать
}

void
futex_unlock(int *futex)
{
	__sync_bool_compare_and_swap(futex, 1, 0);
	futex_wake(futex); // уведомление спящих потоков
}
```
- ***Mutex*** на основе ***futex*** syscall: https://slides.com/gerold103/sysprog_eng6#/23/0/3



##### Переупорядочивание инструкций процессором, компилятором для ускорения работы

![](../_resources/Pasted%20image%2020250102120023.png)
- Проблема: процессор может выполнить инструкции, не связанные друг с другом логической последовательносстью, по данным и т.п., в произвольном порядке
- Пример: read a может пойти за данными по адресу, отсутствующему в CPU cache --> исполнение инструкции предпалагает поход в озу, в то время как write b обратиться к адресу, уже загруженному в кеш -->  это инструкция исполнится раньше read a

- Разные архитектуры процессора допускают разные упрорядочивания, ***x86_64*** наиболее строгая по числу возможных упорядочиваний

##### [Memory barriers](Memory%20barriers.md)

##### [pthread](pthread.md)


![](../_resources/Pasted%20image%2020250103170630.png)