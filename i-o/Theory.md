![](../_resources/Pasted%20image%2020250108224053.png)

```C
void
make_fd_nonblocking(int fd)
{
	int old_flags = fcntl(fd, F_GETFL);
	fcntl(fd, F_SETFL, old_flags | O_NONBLOCK);
}
```

- `fcntl()` - системный вызов, предоставляет различные операции над открытыми файлами, включая управление блокировками файлов, изменение флагов открытия файла и получение информации о файле



![](../_resources/Pasted%20image%2020250108232609.png)
- Процесс может брать лок на файл. Shared или exclusive
- Попытка залочить уже взятый лок - блокирующая операция(если не указан флаг с _NB) 
- lock only the entire file
- advisory lock - it means the lock does not stop from writing/reading the file

![](../_resources/Pasted%20image%2020250108232805.png)
- API с гибкими настройками для взятия лока на часть файла, на file range

![](../_resources/Pasted%20image%2020250108233038.png)


![](../_resources/Pasted%20image%2020250108233707.png)
- ***Все виды file локов - рекомендательные*** - то есть они работают, покуда процессы знают о локах и явно вызывают lock api. А так любой процесс,не знающий о локах, может спокойно записать даже в залоченный файл


