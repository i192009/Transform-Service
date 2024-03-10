#include "hello.h"
#ifdef __cplusplus
extern "C" {
#endif

void Hello(char *file) {
    printf("Process file %s!\n", file);
}

#ifdef __cplusplus
}
#endif

