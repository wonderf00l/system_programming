

![](../_resources/Pasted%20image%2020241204220114.png)

![](../_resources/Pasted%20image%2020241204220128.png)
![](../_resources/Pasted%20image%2020241204220139.png)

![](../_resources/Pasted%20image%2020241204220207.png)
- Что может подразумеваться под файловой системой?
	- сам набор файлов/папок
	- метод организации данных
		- нейминг
		- иерархия, формат хранения
		- методы доступа и тд
	- партиции диска


![](../_resources/Pasted%20image%2020241204220613.png)
![](../_resources/Pasted%20image%2020241204220553.png)

> С файловой системой взаимодействует именно ядро


#### Задачи файловой системы
- Система наименования
	- KV вида: Key - path, Value - file
- Формат хранения
- Интерфейс для взаимодействия с ФС
	- чтение
	- запись
	- разные режимы и проч
> Кастомная реализация ФС на СИ - это структура, инициализированная указателями на функции - то есть предоставляющая имплементацию API. Ядро вызывает разные имплементации в унифицированном виде, оно не знает про конкретные реализации ФС

 - Механизмы защиты, права доступа, принадлежность файлов и тд


![](../_resources/Pasted%20image%2020241204221710.png)

> Файловая система - это абстракция над физическими локальными/удаленным хранилищами данных, которые на логическом уровне представляют не что иное как массив байт. Файловая система - это по сути метод разметки этих байт, метод организации доступа к ним.

![](../_resources/Pasted%20image%2020241204223759.png)
- Причем можем монтировать разные файловые системы к разным видам хранилищ - tmpfs адаптирована для работы поверх ram, однако мы можем примонтировать ее к hdd/ssd

> ***proc*** - виртуальная файловая система, которая по сути является набором хендлеров в ядре, которые обрабатывают операции для работы с файловой системой. просто в userspace доступ к этому набору осуществляется через FS API

![](../_resources/Pasted%20image%2020241204224032.png)


![](../_resources/Pasted%20image%2020241204234448.png)
- Последовательный доступ
	- только запись(клавиатура - символы), только чтение(монитор - последовательность фреймов)
	- и эти <>-only операции производятся one-by-one
	- нет возможности обратиться к произвольному блоку данных на произвольном оффсете произвольного размера
	- в этом смысле сетевая карта - тоже sequential access устройство


#### Устройство HDD

![](../_resources/Pasted%20image%2020241205000725.png)
- Ядро посылает команды на контроллер: считать такие то байты, записать такие-то байты и тп, посылает посредством вызова функций драйвера
- Контроллер уже управляет считывающей головкой, байты считываются с секторов вращающихся дисков 

![](../_resources/Pasted%20image%2020241231154514.png)
![](../_resources/Pasted%20image%2020241231154531.png)
![](../_resources/Pasted%20image%2020241231154544.png)
![](../_resources/Pasted%20image%2020241231154644.png)



![](../_resources/Pasted%20image%2020241231164955.png)
- Диэлектрические пластины покрыты кобальтом/сталью, способными реагировать на изменение магнитного поля

![](../_resources/Pasted%20image%2020241231165039.png)
![](../_resources/Pasted%20image%2020241231165106.png)
- Это покрытие гранулярно, состоит из доменов, в каждом из которых происходит "наведение" зарядов(электронов) - ***при воздействии магнитного поля заряд заряженные частицы могут выстраиваться в той или иной части домена***

> Логическая абстракция над физическим расположением заряда:
> 	0 - одно направление
> 	1 - другое

![](../_resources/Pasted%20image%2020241231172337.png)
> ***Запись*** - намагничивание микроскопической области ферромагнитного покрытия
> 
> ***Чтение*** - считывание изменения полярности, считывание магнитного момента. Изменение полярности провоцирует появление тока и напряжения - что есть логическая единица. Отсутствие - логический ноль. Эта информация проходит через головки в предусилитель и далее на контроллер 

![](../_resources/Pasted%20image%2020241231173500.png)
- Данные расположены "по кругу", считывание битов(магнитных моментов доменов) происходит по кругу
- Данные пишутся/считываются:
		1. в рамках одного из дисков
		2. по одному из колец диска
		3. в рамках одного из сеткоров кольца

![](../_resources/Pasted%20image%2020241231173703.png)
- ***LBA*** - логическая организация данных в виде блока последовательных байт - то есть как можно мыслить о байтах и как о них мыслит ядро ОС при запросе данных/записи данных
- e.g. записать m байт с оффсетом N от начала блока данных
- Контроллер уже занимается переводом такого запроса в контекст жесткого диска - правильно позиционирует головку i-ого диска над j-им сектором k-ого кольца 



#### Устройство SSD
![](../_resources/Pasted%20image%2020241231185335.png)
- Запись
	- подача на src плюса, на drain - минуса
	- подача плюса на control gate, чтобы электроны прошли через слои диэлектрика во floating gate под действием силы Кулона
- "Слив электронов" из floating gate
	- подача минуса на ctrl gate, плюса на drain
- Чтение
	- подача + на drain и ctrl gate. если суммарный заряд правой области 0, то электроны есть во floating gate


![](../_resources/Pasted%20image%2020241231185726.png)
- Со временем происходит "износ" диэлектрика: электронам все легче проходить через его слои, доходит до того, что уже не нужно подавать + на drain - такие области уже не пригодны к эксплупатации
- Контроллер ssd умеет отслеживать подобные зоны -- туда перестает подаваться заряд --> ***из-за этого со временем объем ssd уменьшается***
- Для продления срока службы используются разные механизмы:
	- резервные ячейки, которые не вводятся сразу в эксплуатацию
	- и т.п.
- У каждого ssd есть некое ограниченное гарантийное количество циклов чтения/записи

![](../_resources/Pasted%20image%2020241205004202.png)


![](../_resources/Pasted%20image%2020241205004223.png)
- MBR - содержит мета информацию о корневых партициях устройства хранения, которые могут разбиваться на подпартиции


![](../_resources/Pasted%20image%2020241205004819.png)
- В конечном счете каждая конечная партиция ассоциируется с той или иной файловой системой, коих не так уж и много - у каждой fs есть ""магическое число" - идентификатор
- Этот идентификатор и другая метадата хранится в заголовке файловой системы, который находится в начале партиции

- Партиция - либо очередное разбиение на подпартиции,либо уже файловая система, где лежат файлы и папки
- В начале каждой партиции также метаданные. Если партиция мапится на файловую систему, то там будет заголовок файловой системы
- В нем уже лежит magic num файловой системы, по которому понимаем как парсить байты партиции



#### FAT
![](../_resources/Pasted%20image%2020241207193652.png)

![](../_resources/Pasted%20image%2020241207193837.png)
- Принцип организации FAT
	- Вся партиция делится на блоки
	- файл - совокупность блоков, связанных между собой по принципу linked list
		- связь через FAT - таблицу, где номеру блока соответствует номер следующего
			- n - номер след блока
			- -1 - блок свободен
			- eof - блок не свободен и он последний, это конец файла, больше ничего нет
		- ***сами блоки могут быть расположены на устройстве хаотично - FAT не дает гарантий последовательного расположения блоков*** - файлы могут удаляться, перемещаться и тд, поэтому старые блоки могут менять содержимое, со временем связи между блоками будут перемешиваться
		- ***это плохо в случае hdd - где желательно последовательно считать линию сектора, random access много дороже***
	- Папки интерпретируются как directory entry structures
		- название папки
		- размер
		- номер 1-ого блока содержимого папки
			- это будет либо уже какой-то файл
			- либо подпапка(опять же интерпретируемая как файл)
> Чтение/запись файлов имеют сложность O(N), где N - размер файла. Поэтому такая организация крайне неэффективна для больших файлов - для чтения, добавления байтов нужно пройтись по всем блокам вплоть до конца и добавить новый блок в табличку


![](../_resources/Pasted%20image%2020241207195847.png)
- Так выглядит метадата файла/папки


#### Ext2 filesystem
![](../_resources/Pasted%20image%2020241207204832.png)
![](../_resources/Pasted%20image%2020241207204858.png)


![](../_resources/Pasted%20image%2020241207210805.png)
![](../_resources/Pasted%20image%2020241207211133.png)
- Так выглядит организация массива байтов устройства хранения в случае ext2 файловой системы(опять же, это все software представление)


> Фрагментация данных устройства хранения - ситуация, когда массивы данных, которые логически должны быть расположены рядом последовательно(данные 1-ого файла, файлы одной папки и тд), физически разбросаны по разным местам устройства - это замедляет операции с устройством


#### Виртуальная файловая система(единый интерфейс - множество реализаций) 
![](../_resources/Pasted%20image%2020241207224358.png)
- Конкретная файловая система - реализация интерфейса super_operations
- В коде ядра это буквально fat.h, fat.c; ext2.h, ext2.c etc
- Ядро дергает конкретные реализации методов интерфейса по переданному адресу

![](../_resources/Pasted%20image%2020241207224448.png)
 1. Вызываем функцию write стандартной библиотеки
 2. Под капотом она делает системный вызов write
 userspace
выполнение инструкции jump процессором - переход к инструкции обработчика системного вызова(текст инструкции - код ядра - в сегменте .text адресного пространства ядра) 

 ---
kernel space

исполнение инструкций обработчика системного вызова
 1.   поиск *struct file по полученному из userspace файловому дескриптору(аргументы fd. buf, size) - берутся из стека, регистров
 2. вызов у struct inode структуры struct file метода интерфейса super_operations write_inode
 3. исполнение имплементации метода конкретной файловой системы
	 1. по сути определение того, куда и как запишутся переданные байты буфера размера size 
	 2. алгоритм этого определения и отличает разные файловые системы
 4. конкретные файловые системы уже дергают методы драйвера устройства, к которому примонтированы(драйвер - такой же код в ядре)
	 1. к одному устройству может быть примонтировано несколько файловых систем - все они дернут единый метод драйвера для записи байтов в диск по смещению
 5. код драйвера уже посылает команды на устройство
	 1. вероятнее всего, при исполнении инструкций драйвера работа задействуются конкретные пины/контакты процессора, посылаются байты команды и данные для записи на само устройство через шину данных
 6. процессор исполняет следующие инструкции
 7. когда устройство будет готово, контроллер сгенерирует прерывание, которое попадет в PIC и будет обработано процессором в порядке очереди


> Можно подписываться на события, происходящие с файлом - через inode watcher


![](../_resources/Pasted%20image%2020241207234719.png)
- Задача i/o шедулера - минимизировать количество физических операций с устройством
- Шедулер выполняет merge, sort разных логических операций
- Минимизирует кол-во физических обращений к устройству, если с ним взаимодействует множество процессов/потоков
![](../_resources/Pasted%20image%2020241207235011.png)


![](../_resources/Pasted%20image%2020241207235543.png)
- Еще один способ минимизировать число физических обращений - не делать их

> На самом деле системные вызовы write, read взаимодействуют с page cache(page - т.к. данные читаются страницами) - это структура в ядре, которая кеширует страницы устройства

- Закешированные страницы хранятся в ram, попадают в cpu cache

![](../_resources/Pasted%20image%2020241207235804.png)

> В случае успешного завершения write нет никаких гарантий, что данные действительно попали на устройство
> 	Данные пойдут в кеш, а будут вытеснены на диск только спустя некоторое время
> 	***Поэтому там, где важны гарантии записи и возвращения ответа только после полноценной операции - нужно синковать данные - сбрасывать записанные байты из page cache на устройство(системный вызов fsync())*** 


> То есть проходим через кеш, если все же нужно идти на диск, то операции проходят через i/o шедулер - только после этого фактически выполняется операция - все это в kernel space

![](../_resources/Pasted%20image%2020241208000427.png)
- Также есть буферизация ввода/вывода и в userspace

