package main

import (
    "syscall"
    wapi "github.com/cwchiu/go-winapi"
)


func main(){
    _T := syscall.StringToUTF16Ptr
    wapi.MessageBox(wapi.HWND_TOP, _T("Hello Windows"), _T("HelloMsg"), 0)
}