#include <windows.h>

/*
 * This program is used to refresh the windows environment.
 *
 * Compile using the following command(cl info: 19.15.26732.1 for x86):
 *
 *     cl refreshwin.c /link user32.lib /out:refreshwin.exe
 *
 */
int main() {
    SendMessage(HWND_BROADCAST, WM_SETTINGCHANGE, 0, (LPARAM)"Environment");
    return 0;
}

