
![](../_resources/Pasted%20image%2020241201230933.png)
![](../_resources/Pasted%20image%2020241201230945.png)
- Аппаратные прерывания - некоторые сигналы от периферии, которые обрабатываются процесором
- Программные - для обработки ситуаций, происходящих на самом процессоре, а не на периферии, генерируются и обрабатываются процессором, но причина их появления - пользовательский код
	- zero division
	- дебаг - breakpoint - это инструкция, генерирующая прерывание в точке останова, обработка которого приводит к остановке на этой точке
	- page fault - при TLB miss, write to COW memory segment etc


![](../_resources/Pasted%20image%2020241201231521.png)
Flow чтения с диска:
- В коде делается read syscall
- Через спец инструкцию syscall попадаем сразу в ядро
- Идем в файловую систему в ядре, находим нужный файл и соответствующее ему устройство
- Посылаем в драйвер устройства команду для исполнения, байты данных
- Это команда/данные ставятся на шину данных
- Все устройства ее видят, но по шине она доходит до диска, исполняется



> Главная проблема - все периферийные устройства много медленнее процессора. Поэтому нужен механизм, который позволит процессору отправить команду устройствам асинхронно и пойти дальше выполнять инструкции. А когда устройство завершит работу, то даст об этом знать 

![](../_resources/Pasted%20image%2020241201232756.png)
![](../_resources/Pasted%20image%2020241201232823.png)
![](../_resources/Pasted%20image%2020241201232829.png)
![](../_resources/Pasted%20image%2020241201232843.png)![](../_resources/Pasted%20image%2020241201232853.png)
- Если устройствам есть что сказать cpu, то они генерируют прерывание
- Прерывание отправляется не непосредственно на cpu, а в PIC, где они копятся в порядке очереди.
- У PIC есть регистр, в который пишется номер устройства, идентифицирующий прерывание
- У процессора есть элемент, сигнализирующий о готовности обработать прерывание
- PIC это видит и посылает по шине прерывание
- Процессор его получает, считывает через шину номер устройства, обрабатывает прерывание
- И так в порядке очереди, значение устройства в регистре PIC поочередно меняется


![](../_resources/Pasted%20image%2020241201233421.png)
- Программные же прерывания возникают на самом процессоре, обычно это разные ошибки
- Выглядит это так, что процессору встречается инструкция, которую нельзя выполнить. Тут генерится прерывание и вызывается обработчик в ядре. Итогом обработки станет либо последующая обработка, либо exit с кодом ошибки

![](../_resources/Pasted%20image%2020241201234242.png)
- Errors(exceptions) - прерывания - ошибки(не всегда критичные)
	- для таких прерываний вызовется обработчик
	- после работы хендлера попытаемся заретраить ***инструкцию, вызвавшую прерывание*** - будем ретраиться, пока инструкция не выполниться, либо пока хендлер не завершит исполнение инструкций принудительно
	- ретрай пройдет, например, если при выполнении инструкции получили page fault, или запись в cow сегмент --> хендлер прерывания загрузит мапинг в TLB/выделит память в пространстве дочернего процесса --> при ретрае корректно исполним инструкцию 
- Traps - после обработки прерываний, вызванных инструкцией, перейдем сразу к следующей инструкции, без ретрая
	- Системные вызовы - перепрыгиваем в ядро, исполняем syscall, возвращаемся. незачем это ретраить

> Прерывание - остановка исполнения текущего кода, сохранение состояния, начало исполнения другого кода, возврат состояния, продолжение исполнения


> 1. Как и где найти хендлер текущего прерывания? 
> 2. Как и где сохранить контекст текущего исполнения, исполнить хендлер и вернуться бесшовно для текущего хода исполнения?

![](../_resources/Pasted%20image%2020241202002806.png)
- По сути контекст исполнения - это набор регистров
	1. Текущее состояние стека характеризуется значением в регистре SP
	2. Адрес текущей инструкции, которая сгенерила прерывание, то есть откуда продолжить исполнение - хранится в регистре
	3. Аргументы функций
		1. Либо передавались явно в регистрах
		2. Либо они на стеке - см. п.1
- То есть прихраним значения этих регистров в спец стеке в ОЗУ, потом восстановим по ним контекст


- После сохранения контекста посмотрим в регистр idtr, который хранит адрес начала массива с хендлерами, где индекс массива - номер прерывания -- ***IDT***(хранится в ОЗУ, инициализируется на старте системы)
- IDT - это массив адресов, по которым находятся машинные инструкции
- Извлечем нужный адрес, процессор сделает ***jump*** на этот адрес и исполнит код хендлера

> Указатели на функции - ни что иное как адреса инструкций, которые нужно выполнить

- После исполнения хендлера возвращаемся обратно, восстанавливаем контекст, раскладываем его по регистрам - DONE! продолжаем исполнение как ни в чем не бывало



![](../_resources/Pasted%20image%2020241204195706.png)
- 32-255 - отдаются на откуп ядру



![](../_resources/Pasted%20image%2020241204200909.png)
- Хендлеры прерываний отрабатывают на отдельных стеках, которые создаются в ядре для каждого физического ядра процессора
	- более безопасно: например, в случае stackoverflow обработчик не может исполняться на том же стеке, что и битая программа
- ***Top half - быстрая первичная обработка одного из множества прерываний в критическом, горячем контексте - происходит последовательно one-by-one в блокирующем режиме***
	- найти устройство-причину прерывания(если надо)
	- сохранить первичный срочные данные
	- зашедулить исполнение основной тяжеловесной части ***bottom half***
- ***Bottom half - основная тяжеловесная обработка***
	- Задачи разбираются воркер пулом, исполняются параллельно
	- Требования к скорости обработки не столь критичные
	- Примеры тяжелых задач:
		- доставить сигнал в userspace
		- перевести спящие потоки из состояния waiting в ready to run, добавить в список потоков для шедулинга и последующего исполнения

![](../_resources/Pasted%20image%2020241204203339.png)
- Регистрация обработчика прерывания для какого-либо устройства


![](../_resources/Pasted%20image%2020241204203419.png)
- Зарегистрированные обработчики прерываний и кол-во их исполнений

> ***/proc*** - это не файловая система, это интерфейс для взаимодействия с ядром и получения различной служебной информации в userspace. Каждый read таких "файлов" перехватывается ядром и обрабатыватся


![](../_resources/Pasted%20image%2020241204210809.png)
- Когда процессу посылается сигнал, он перехватывается процессом(userspace), вызывается обработчик сигнала - предопределенный или наш кастомный


![](../_resources/Pasted%20image%2020241204210922.png)

