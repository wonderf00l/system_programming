
![](../_resources/Pasted%20image%2020250103232354.png)
- Простые средства IPC

![](../_resources/Pasted%20image%2020250103233336.png)

https://slides.com/gerold103/sysprog_eng7#/8/1
```C
void
write_to_shared_mem(char *mem, const char *src, int size)
{
	volatile char *is_readable = &mem[IS_READABLE];
	volatile char *is_writable = &mem[IS_WRITABLE];
	mem += MEM_META_SIZE;
	int mem_size = MEM_SIZE - MEM_META_SIZE;
	int saved_size = mem_size;
	char *saved_mem = mem;
	while (1) {
	// в данном случае важно установить барьеры записи при помощи atomic операции(нет проблем именно многопоточного конкурентного доступа, только барьеры), т.к. компилятор/процессор не знают об IPC и могут переставить инструкции 
		while (! __atomic_load_n(is_writable, __ATOMIC_ACQUIRE)) // инструкции записи в shared memory должны быть выполнены строго после установки is_writable в 1 другим процессом 
			sched_yield();
		int to_copy = mem_size > size ? size : mem_size;
		memcpy(mem, src, to_copy);
		size -= to_copy;
		mem_size -= to_copy;
		mem += to_copy;
		src += to_copy;

		__atomic_store_n(is_writable, 0, __ATOMIC_RELEASE); // БАРЬЕР RELEASE НЕОБХОДИМ, ЧТОБЫ ИНСТРУКЦИИ ЗАПИСИ БАЙТ ФАЙЛА В SHARED MEMORY ГАРАТНИРОВАННО БЫЛИ ВЫПОЛНЕНЫ ДО УСТАНОВКИ IS_WRITABLE в 0, на случай, если другой процесс будет как-то ориентироваться на этот флаг. 
		__atomic_store_n(is_readable, 1, __ATOMIC_RELEASE); // is_readable должно стать 1 строго после установки is_writable в 0, иначе теоритически возможна ситуация, когда инструкция is_readable = 1 исполниться раньше, процесс прочитает данные и установит is_writable байт shared memory в 0 одновременно с пишущим процессом
		if (size == 0)
			break;
		mem = saved_mem;
		mem_size = saved_size;
	}
}
```

![](../_resources/Pasted%20image%2020250105191957.png)
- ***pipe, mmap*** можно использовать для связи родственных процессов, однако для связи разнородных процессов нужны иные средства IPC

#### FIFO
![](../_resources/Pasted%20image%2020250105192234.png)
- FIFO - именованный pipe
- при создании fifo через mkfifo в файловой системе создается файл с типом pipe
- если в случае pipe данные читателя/писателя отправляются в некий промежуточнйы буфер, то есть операции неблокирущие, если есть куда писать/есть что читать
- в случае fifo передача данных идет напрямую через буферы писателя/читателя; размер файла fifo 1 байт - по сути это небуферизированный канал для записи/чтения сообщений

##### server
```go
int main()
{
	mkfifo("/tmp/fifo_server", S_IRWXU | S_IRWXO);
	int fd = open("/tmp/fifo_server", O_RDONLY);
	while (1) {
		pid_t new_client;
		if (read(fd, &new_client, sizeof(pid_t)) > 0)
			printf("new client %d\n", (int) new_client);
		else
			sched_yield();
	}
}
```

##### client
```go
int main()
{
	int fd = open("/tmp/fifo_server", O_WRONLY);
	pid_t pid = getpid();
	write(fd, &pid, sizeof(pid));
	printf("my pid %d is sent to server\n", (int) pid);
	close(fd);
	return 0;
}
```

![](../_resources/Pasted%20image%2020250105192901.png)

#### XSI IPC

![](../_resources/Pasted%20image%2020250105192926.png)
- Стандарт XSI - определяет станларт работы ряда средств IPC для разных ОС
- Все средства XSI IPC работает по следующим принципам:
	- Создаются по ключу, если не были созданы. В дальнейшем к ним производится подключение, так же по ключу
	- Живут и после завершения процесса, если не были удалены вручную 


![](../_resources/Pasted%20image%2020250105193108.png)
- Очередь сообщений - создается с нуля/либо к ней осуществляется подключение
- Живет и после завершения процесса - нужно чистить вручную при необходимости
- Сообщения разделяются по ключам


![](../_resources/Pasted%20image%2020250105195833.png)
- Получить список существующих IPC


#### Semaphore
![](../_resources/Pasted%20image%2020250105195912.png)
- XSI Semaphore - межпроцессный семафор

- Межпоточный семафор в рамках 1-ого процесса
```go
struct semaphore {
	int counter;
	pthread_mutex_t mutex;
	pthread_cond_t cond;
};

static inline void
semaphore_get(struct semaphore *sem)
{
	pthread_mutex_lock(&sem->mutex);
	while (sem->counter == 0)
		pthread_cond_wait(&sem->cond, &sem->mutex);
	sem->counter--;
	pthread_mutex_unlock(&sem->mutex);
}

static inline void
semaphore_put(struct semaphore *sem)
{
	pthread_mutex_lock(&sem->mutex);
	sem->counter++;
	pthread_cond_signal(&sem->cond);
	pthread_mutex_unlock(&sem->mutex);
}
```


#### Shared memory(аналог mmap)
![](../_resources/Pasted%20image%2020250106093654.png)
- ***shmat*** - мапинг id/дескриптора, полученного в shmget, на кусок памяти в системе