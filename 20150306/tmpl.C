#include <stdlib.h>
#include <stdio.h>

int main(int argc, char **argv) {
    if (argc != 2) {
        return 127;
    }
    unsigned int nargs = atoi(argv[1]);
    unsigned long long ns = 1;
    printf("%llu\n", ns);
    return 0;
}
