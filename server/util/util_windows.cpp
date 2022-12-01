// util_windows.cpp provide win32 api calls from go lang to restart/shutdown, and change UEFI firmware BootNext variable
// code adopted from: https://serverfault.com/questions/813695/how-do-i-stop-windows-10-install-from-modifying-bios-boot-settings

#include "util_windows.h"
#include <memory>
// Define Global UEFI GUID for PC
const LPCWSTR globalGUID = L"{8BE4DF61-93CA-11D2-AA0D-00E098032B8C}";
const LPCWSTR BootOrderStr = L"BootOrder";
const LPCWSTR BootNextStr= L"BootNext";

struct CloseHandleHelper
{
  void operator()(void *p) const { CloseHandle(p); }
};


/** Function to obtain required priviledges to issue shutdown or restart**/
BOOL MySetPrivilege(HANDLE process, LPCWSTR name, BOOL on)
{
  HANDLE token;
  if (!OpenProcessToken(process, TOKEN_ADJUST_PRIVILEGES, &token))
    return FALSE;
  std::unique_ptr<void, CloseHandleHelper> tokenLifetime(token);
  TOKEN_PRIVILEGES tp;
  tp.PrivilegeCount = 1;
  if (!LookupPrivilegeValueW(NULL, name, &tp.Privileges[0].Luid))
    return FALSE;
  tp.Privileges[0].Attributes = on ? SE_PRIVILEGE_ENABLED : 0;
  return AdjustTokenPrivileges(token, FALSE, &tp, sizeof(tp), NULL, NULL);
}

/**Shutdown function**/
void shutdown(uint16_t mode)
{
  //MODE 1 - shutdown 0 - restart
  if (mode == 1)
    InitiateSystemShutdownEx(NULL, NULL, 0, FALSE, FALSE, 0);
  else
    InitiateSystemShutdownEx(NULL, NULL, 2, FALSE, TRUE, 0);
}
/**Change UEFI boot function**/
void changeBoot(uint16_t data, uint16_t mode)
{
  // DATA: BootID
  // MODE:  0 - change BootNext ( temporary next boot change)
  //        1 - change BootOrder (permanent EFI boot order change)

  const int bootOrderBytes = 2;
  uint16_t bootOrderData = data;
  const LPCWSTR *bootOrderName = (mode == 0) ? &BootNextStr : &BootOrderStr;
  SetFirmwareEnvironmentVariableW(*bootOrderName, globalGUID, &bootOrderData, bootOrderBytes);
  
}


/**Suspend function**/  
bool SystemSuspend()
{
  MySetPrivilege(GetCurrentProcess(), (LPCWSTR)SE_SHUTDOWN_NAME, TRUE);
  SetSuspendState(FALSE, FALSE, FALSE);
  return true;
}

/**Shutdown/restart function**/  
bool SystemShutdown(uint16_t mode)
{
  //MODE 1 - shutdown 0 - restart

  // Get priviledge to shutdown/restart
  MySetPrivilege(GetCurrentProcess(), (LPCWSTR)SE_SHUTDOWN_NAME, TRUE);
  // we are just doign a normal shutdown
  shutdown(mode);
  // shutdown was successful
  return true;
}
/**Change UEFI boot function**/
bool SystemChangeBoot(uint16_t data, uint16_t mode)
{
  // DATA: BootID
  // MODE:  0 - change BootNext ( temporary next boot change)
  //        1 - change BootOrder (permanent EFI boot order change) 

  // get priviledge to change UEFI variables   
  MySetPrivilege(GetCurrentProcess(), (LPCWSTR)SE_SYSTEM_ENVIRONMENT_NAME, TRUE);
  
  // data is boot integer id
  // Mode 0: BootNext, 1: BootOrder
  changeBoot(data,mode);
  return true;
}
