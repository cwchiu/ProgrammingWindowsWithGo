package main

import (
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

const (
    IDM_FILE         = 40001
    IDM_EDIT         = 40002
    IDM_FILE_NEW     = 40003
    IDM_FILE_OPEN    = 40004
    IDM_FILE_SAVE    = 40005
    IDM_FILE_SAVE_AS = 40006
    IDM_MAIN         = 40007
    IDM_EDIT_UNDO    = 40008
    IDM_EDIT_CUT     = 40009
    IDM_EDIT_COPY    = 40010
    IDM_EDIT_PASTE   = 40011
    IDM_EDIT_CLEAR   = 40012
)

func main() {
    app, _ := NewApp()
    var hMenuMain, hMenuEdit, hMenuFile HMENU
    app.On(WM_CREATE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        hMenuMain = LoadMenu(app.HInstance, _T("MenuMain"))
        hMenuFile = LoadMenu(app.HInstance, _T("MenuFile"))
        hMenuEdit = LoadMenu(app.HInstance, _T("MenuEdit"))

        SetMenu(hwnd, hMenuMain)
        return 0
    })

    app.On(WM_COMMAND, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        switch int32(LOWORD(uint32(wParam))) {
        case IDM_FILE_NEW, IDM_FILE_OPEN, IDM_FILE_SAVE, IDM_FILE_SAVE_AS,
            IDM_EDIT_UNDO, IDM_EDIT_CUT, IDM_EDIT_COPY, IDM_EDIT_PASTE,
            IDM_EDIT_CLEAR:
            MessageBeep(0)
            return 0
        case IDM_MAIN:
            SetMenu(hwnd, hMenuMain)
            return 0

        case IDM_FILE:
            SetMenu(hwnd, hMenuFile)
            return 0

        case IDM_EDIT:
            SetMenu(hwnd, hMenuEdit)
            return 0
        }
        return DefWindowProc(hwnd, msg, wParam, lParam)
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        SetMenu(hwnd, hMenuMain)
        DestroyMenu(hMenuFile)
        DestroyMenu(hMenuEdit)

        PostQuitMessage(0)
        return 0
    })

    app.Init("NoPopUps", "No-Popup Nested Menu Demonstration")

    app.Run()
}
