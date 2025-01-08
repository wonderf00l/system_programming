#### Коммутация
![](../_resources/Pasted%20image%2020250107114242.png)
- Пакетная коммутация - используется в современных сетях
	- узлы содержат информацию о других узлах сети
	- попадая на текущий узел, пакет направляется(маршрутизируется) на следующий узел, находящийся ближе к целевому
	- пакеты могут проходит разными маршрутами

#### Формироварние пакетов
![](../_resources/Pasted%20image%2020250107115535.png)
- В сокет пишем непосредственно полезную нагрузку
- Далее она проходит через стек сетевых протоколов в ядре, добавляются разные служебные пакеты


![](../_resources/Pasted%20image%2020250107115726.png)

#### OSI
![](../_resources/Pasted%20image%2020250107121806.png)
- прикладной уровень знает о смысле и назначении байт данных
- на этом уровне происходит обработка этих данных

![](../_resources/Pasted%20image%2020250107122416.png)
- Уровень представления не заботится о значении байт, ключевая задача этого уровня - упаковка даннных в нужный формат, предаине им структуры, выбор кодировки и т.п.

На последующих уровнях нас вообще не заботит содержимое байт, их формат, для нас это просто пласт данных для отправки

![](../_resources/Pasted%20image%2020250107122531.png)
- Уровень сессии оперирует понятиями аутентификации, авторизации, сессии и т.п.


![](../_resources/Pasted%20image%2020250108113953.png)
- Транспортный уровень
	- распределение данных по data streams
	- уровни гарантии доставки
	- целостность данных
	- менеджмент пакетов
	- содержимое и формат данных - не важны, точка назначения - определяется сетевым уровнем, вашен лишь порт

![](../_resources/Pasted%20image%2020250108114253.png)

![](../_resources/Pasted%20image%2020250108115406.png)
- Подсеть квартиры - дома - района - города - страны - материка - глобальная сеть


![](../_resources/Pasted%20image%2020250108124001.png)

![](../_resources/Pasted%20image%2020250108131703.png)
- Стек протоколов, реализуемый на практике - TCP/IP
- В ядре ОС обычно реализован сетевой стек - как раз стек TCP/IP


![](../_resources/Pasted%20image%2020250108131743.png)
![](../_resources/Pasted%20image%2020250108131818.png)
- При использовании интерфейса сокетов попадаем в ядро, где уже исполняется функционал сетевого стека ядра
- В ядре полезная нагрузка по стеку вызовов разных функций проходит через этапы добавление заголовков транспортного/сетевого/канального уровней, происходит прочая обработка
- Соответственно, на этом уровне реализованы механизмы установки соединения(если подразумеваются протоколом), проверки целостности данных, чек-суммы, проверка порядка сообщений, ретраи и прочая конкретика разных протоколов, у всех разная, предусматривает разные вещи


![](../_resources/Pasted%20image%2020250108131906.png)


![](../_resources/Pasted%20image%2020250108132034.png)
- Фактически тип сокета определяет основные методы обработки, упаковки и отправки данных, записываемых в сокет/читаемых из сокета(иначе говоря, определяет диапазон конкретных применяемых протоколов)
- **Why TCP (SOCK_STREAM) is Considered a Stream:** 
While TCP indeed operates on the concept of **packets (or segments)** at a lower level (the transport layer), **TCP provides a reliable byte stream** at the application layer. This means that for the user or application programming interface (API) level, the data is treated as a continuous, ordered flow of bytes, without worrying about the underlying packetization.  
- Here’s how it works:  
**Segmentation and Reassembly**: At the transport layer, TCP divides data into packets (segments) for transmission, and each segment contains a sequence number and other information. But when the data is received, **TCP reassembles these segments in the correct order** and presents the received data to the application as a continuous stream of bytes.  
**No Message Boundaries**: TCP doesn’t preserve message boundaries. This means that, from the application’s perspective, it just sends and receives an uninterrupted flow of data. The application doesn't know where one packet ends and the next one begins (those boundaries are not visible to the application).  
**Stream of Data**: When you send data over TCP, you write a series of bytes to the socket, and TCP handles breaking those bytes into segments, transmitting them, and reassembling them at the receiver. Similarly, when you read from a TCP socket, you read data as a continuous stream (not as discrete packets), and TCP will provide that data after reassembling any segmented parts

- **Contrast with SOCK_DGRAM (UDP):**  
On the other hand, **SOCK_DGRAM** is used for **UDP (User Datagram Protocol)**, which is connectionless and operates on discrete, independent packets (datagrams). UDP:  
Does **preserve message boundaries**: Each packet sent via UDP represents a separate message. The application knows exactly where one packet ends and the next begins.  
Does not guarantee reliability, ordering, or congestion control like TCP.  

Even though TCP uses packets (or segments), its key feature is to **present data as a continuous stream**, abstracting the packetization details from the application layer. This is why TCP is considered a "stream" protocol, while UDP is often referred to as a "message-oriented" or "datagram" protocol.  

- TCP - connection-oriented(виртуальный коннект), UDP - not connection oriented

> Распространенная практкика - реализация своего прикладного протокола в userspace под собственные нужды на основе протоколов L4. Но тут нужно учитывать возможные коллизии функционала своих протоколов поверх и протоколов L4(функционал обеспечения порядка доставки, гарантий доставки, проверки целостности данных etc)



![](../_resources/Pasted%20image%2020250108144000.png)
- ipv4 адреса кодируются в 4-ое число, порт просто представляется 2-ух байтовым числом
- общепринятый порядок байт для чисел  - bit endian(прямой порядок)
- на машинах с little endian - используются функции для разворота байтов

