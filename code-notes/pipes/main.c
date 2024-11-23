#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/wait.h>

int main()
{
    int to_parent[2];
    int to_child[2];
    pipe(to_child);  // создание однонаправленного канала записи: будут созданы 2 дескриптора в случае успеха pipe() -> один для чтения(to_child[0]), другой - для записи(to_child[1])
    pipe(to_parent); // создание еще одного однонаправленного канала -> двунаправленная коммуникация
    char buf[16];
    if (fork() == 0)
    {
        close(to_parent[0]); // важно закрывать задублированные дескрипторы в ребенке: 1. они просто не нужны, дубли;
        close(to_child[1]);  // 2. пайпы устроены так, что если одна из сторон закрыла дескриптор, другая возвращает ошибку --> не закрыв эти дубли мы повиснем на чтении(записи), даже если запиь(чтение) на другой стороне прекратились
        read(to_child[0], buf, sizeof(buf));
        printf("%d: read %s\n", (int)getpid(), buf);
        write(to_parent[1], "hello2", sizeof("hello2"));
        return 0;
    }
    close(to_parent[1]);
    close(to_child[0]);
    write(to_child[1], "hello1", sizeof("hello"));
    read(to_parent[0], buf, sizeof(buf));
    printf("%d: read %s\n", (int)getpid(), buf);
    wait(NULL);
    return 0;
}