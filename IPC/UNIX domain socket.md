
#### Socketpair API
![](../_resources/Pasted%20image%2020250106111820.png)
- Пара сокетов - двунаправленный канал связи
- Создание пары - системный вызов ***sockerpait***
	- ***sockerpair*** принимает конкретные файловые дескрипторы - по сути ***sockerpair*** это аналог пайпа с API сокетов
	- то есть ***sockerpair*** - это анонимный канал связи, он подойдет для связи близкородственных процессов
- На сокетах построено сетевое взаимодействие процессов
- Доменные сокеты - "локальные" сокеты, для взаимодействия процессов на локальной машине. API локальных сокетов совпадает с API сетевых


![](../_resources/Pasted%20image%2020250106112158.png)
- Протокол - соглашения о способе передачи данных и их обработке, интрепретации
	- Основные: TCP, UDP, IP etc

> Для обеспечения сетевого взаимодействия процессов обычно используются стандартный стек сетевых протоколов, реализованный в ядре.

![](../_resources/Pasted%20image%2020250106112306.png)
- Тип сокета - способ работы именно с сокетом
	- `SOCK_DGRAM` - передача данных пакетами/датаграммами
	- `SOCK_STREAM` - потоковая передача данных
	- ...
![](../_resources/Pasted%20image%2020250106112428.png)
- Домен - определяет семейство протоколов
	- `AF_UNIX/AF_LOCAL` - доменные сокеты
	- ...


```go
#define MAX_MSG_SIZE 1024 * 1024

void worker(int sock)
{
	char buffer[MAX_MSG_SIZE];
	ssize_t size;
	while ((size = read(sock, buffer, sizeof(buffer))) > 0)
		printf("Client received %d\n", (int) size);
}

int main(int argc, const char **argv)
{
        int socket_type = strcmp(argv[1], "stream") == 0 ?
                          SOCK_STREAM : SOCK_DGRAM;
	int sockets[2];
	socketpair(AF_UNIX, socket_type, 0, sockets);
	if (fork() == 0) {
		close(sockets[0]);
		worker(sockets[1]);
		return 0;
	}
	close(sockets[1]);
	char buffer[MAX_MSG_SIZE];
	int size;
	while (scanf("%d", &size) > 0) {
		printf("Server sent %d\n", size);
		write(sockets[0], buffer, size);
	}
	return 0;
}
```


![](../_resources/Pasted%20image%2020250106130345.png)
![](../_resources/Pasted%20image%2020250106130400.png)

#### Socket API

- API для создания именованных сокетов, которые смогут использовать любые  процессы для связи через канал по socket API
```go
int
socket(int domain, int type, int protocol); // создание анонимного сокета - получение дескриптора сокета

int
bind(int sockfd, const struct sockaddr *addr, 
     socklen_t addrlen); // привязка дескриптора сокета к адресу сокета 

int
listen(int sockfd, int backlog); // прослушивание дескриптора сокета(server side)

int
connect(int sockfd, const struct sockaddr *addr,
        socklen_t addrlen); // установка подключения до сокета с адресом remote_addr, sockfd - дескриптор клиентского сокета
// подключение устанавливается после bind() + listen() на стороне сервера

int
accept(int sockfd, struct sockaddr *addr,
       socklen_t *addrlen); // прием соединения, установленного клиентом через connect, на сервере(прием относительно дескриптора серверного сокета sockfd, возвращается дескриптор клиентского сокета)
```
- Результат `accept()` принятия `connect()` подключения - установка двунаправленного канала связи, по сути создание глобальной сокет-пары для связи разных процессов