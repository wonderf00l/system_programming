#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/wait.h>
#include <string.h>

int main(int argc, const char **argv)
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

    int p = pipes1(argc, argv); // внутри создается новый процесс: стек вызовов сохранится, дочерний процесс вернет 123 из pipe1 и напечатает в свой stdout дескриптор, т.к. они дублируются - видим вывод 2-ух процессов в одну и ту же консоль
    printf("pipes1 res: %d\n", p);

    return 0;
}

int pipes1(int argc, const char **argv) // будем блокироваться, если не закроем duplicated дескрипторы
{
    int channel[2];
    int need_close = argc > 1 && strcmp(argv[1], "close") == 0;
    pipe(channel);
    if (fork() == 0)
    {
        char buf[16];
        printf("child: started\n");
        if (need_close)
        {
            printf("child: close output channel\n");
            close(channel[1]);
        }
        while (read(channel[0], buf, sizeof(buf)) > 0)
            printf("child: read %s\n", buf);
        printf("child: EOF\n");
        return 123;
    }
    write(channel[1], "100", 3);
    printf("parent: written 100\n");
    if (need_close)
    {
        printf("parent: close output channel\n");
        close(channel[1]);
    }
    printf("parent: waiting for child termination ...\n");
    wait(NULL);
    return 0;
}