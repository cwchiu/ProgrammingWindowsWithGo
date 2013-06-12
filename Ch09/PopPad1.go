package main

import (
    //"fmt"
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"
    //"unsafe"
)

const ID_EDIT  = 1
var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func main() {
    app, _ := NewApp()
    var hwndEdit HWND
    app.On(WM_CREATE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
      hwndEdit = CreateWindowEx (0, _T("edit"), nil,
             WS_CHILD | WS_VISIBLE | WS_HSCROLL | WS_VSCROLL |
                       WS_BORDER | ES_LEFT | ES_MULTILINE |
                       ES_AUTOHSCROLL | ES_AUTOVSCROLL,
             0, 0, 0, 0, hwnd, HMENU(ID_EDIT),
             app.HInstance, nil) 
        return 0
    })

    app.On(WM_COMMAND, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        if (LOWORD (uint32(wParam)) == ID_EDIT){
               if (HIWORD (uint32(wParam)) == EN_ERRSPACE || 
                         HIWORD (uint32(wParam)) == EN_MAXTEXT){

                    MessageBox (hwnd, _T("Edit control out of space."),
                                _T("PopPad1"), MB_OK | MB_ICONSTOP) 
                                
                                }
                                
                                }
                                
        return 0
    })
    
    app.On(WM_SIZE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        MoveWindow (hwndEdit, 0, 0, int32(LOWORD (uint32(lParam))), int32(HIWORD (uint32(lParam))), true)

        return 0
    })

    app.On(WM_SETFOCUS, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        SetFocus (hwndEdit)
        return 0
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })

    app.Init("PopPad1", "PopPad1")
    app.Run()
}
