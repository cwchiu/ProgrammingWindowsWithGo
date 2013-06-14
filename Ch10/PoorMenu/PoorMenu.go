package main

import (
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

const (
    IDM_SYS_ABOUT  = 1
    IDM_SYS_HELP   = 2
    IDM_SYS_REMOVE = 3
)

func main() {
    app, _ := NewApp()

    app.On(WM_SYSCOMMAND, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        switch int32(LOWORD(uint32(wParam))) {
        case IDM_SYS_REMOVE:
            GetSystemMenu (hwnd, TRUE) 
            return 0

        case IDM_SYS_HELP:
            MessageBox(hwnd, _T("Help not yet implemented!"),
                _T("MenuDemo"), MB_ICONEXCLAMATION|MB_OK)
            return 0

        case IDM_SYS_ABOUT:
            MessageBox(hwnd, _T("A Poor-Person's Menu Program\n(c) Chui-Wen Chiu, 2013"),
                _T("MenuDemo"), MB_ICONINFORMATION|MB_OK)
            return 0
        }
        return DefWindowProc(hwnd, msg, wParam, lParam)
    })


    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })

    app.Init("PoorMenu", "The Poor-Person's Menu")
    
    hMenu := GetSystemMenu (app.HWnd, FALSE)
     
    AppendMenu (hMenu, MF_SEPARATOR, 0,           nil) 
    AppendMenu (hMenu, MF_STRING, IDM_SYS_ABOUT,  _T ("About...")) 
    AppendMenu (hMenu, MF_STRING, IDM_SYS_HELP,   _T ("Help...")) 
    AppendMenu (hMenu, MF_STRING, IDM_SYS_REMOVE, _T ("Remove Additions")) 
     
    app.Run()
}
