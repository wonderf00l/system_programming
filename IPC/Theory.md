
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
