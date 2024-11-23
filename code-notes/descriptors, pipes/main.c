#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>

int main()
{
    int fd = open("tmp.txt", O_RDWR | O_CREAT);
    int fd2 = dup(fd);
    printf("fd=%d, fd2=%d\n", fd, fd2);
    dprintf(fd2, "1 ");
    close(fd2);
    dprintf(fd, "2 ");
    close(fd);
    return 0;
}