#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/wait.h>

int main()
{
    int to_child[2];
    pipe(to_child);
    char buf[16];
    dup2(to_child[0], 0); // отвязываем стандарнтый дескриптор stdin от struct file, описывающую файл потока ввода в kernespace(
    // привязываем в файлу, в который будут приходить данные в резульатте записи в пайпа, из него можем только читать (по сути перетерли указатель в таблице файловых дескрипторов)
    // *struct file{...} = 0x123(stdin console file pointer) --> *struct{...0} 0x345(read end of the pipe)
    if (fork() == 0)
    {
        close(to_child[1]);
        return execlp("python3", "python3", "-i", NULL); // вызываем другую программу, читающую из пайпа --> в итоге связываем две программы
    }
    close(to_child[0]);
    const char cmd[] = "print(100 + 200)";
    write(to_child[1], cmd, sizeof(cmd));
    close(to_child[1]);
    wait(NULL);
    return 0;
}