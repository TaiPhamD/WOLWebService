#include <windows.h>
#include <powrprof.h>
#include <stdint.h>
#include <stdbool.h>
bool SystemSuspend();
bool SystemShutdown(uint16_t mode);
bool SystemChangeBoot(uint16_t data, uint16_t mode);
