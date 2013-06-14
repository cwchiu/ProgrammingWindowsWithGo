package main

import (
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

const (
    IDM_FILE_NEW     = 40001
    IDM_FILE_OPEN    = 40002
    IDM_FILE_SAVE    = 40003
    IDM_FILE_SAVE_AS = 40004
    IDM_APP_EXIT     = 40005
    IDM_EDIT_UNDO    = 40006
    IDM_EDIT_CUT     = 40007
    IDM_EDIT_COPY    = 40008
    IDM_EDIT_PASTE   = 40009
    IDM_EDIT_CLEAR   = 40010
    IDM_BKGND_WHITE  = 40011
    IDM_BKGND_LTGRAY = 40012
    IDM_BKGND_GRAY   = 40013
    IDM_BKGND_DKGRAY = 40014
    IDM_BKGND_BLACK  = 40015
    IDM_TIMER_START  = 40016
    IDM_TIMER_STOP   = 40017
    IDM_APP_HELP     = 40018
    IDM_APP_ABOUT    = 40019
    ID_MENUITEM40020 = 40020
    ID_TIMER         = 1
)

func main() {
    app, _ := NewApp()

    var hMenu HMENU
    var idColor [5]int = [5]int{WHITE_BRUSH, LTGRAY_BRUSH, GRAY_BRUSH,
        DKGRAY_BRUSH, BLACK_BRUSH}
    var iSelection int32 = IDM_BKGND_WHITE

    app.On(WM_COMMAND, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        hMenu = GetMenu(hwnd)
        switch int32(LOWORD(uint32(wParam))) {
        case IDM_FILE_SAVE_AS, IDM_FILE_NEW, IDM_FILE_OPEN, IDM_FILE_SAVE:
            MessageBeep(0)
            return 0

        case IDM_APP_EXIT:
            SendMessage(hwnd, WM_CLOSE, 0, 0)
            return 0

        case IDM_EDIT_CLEAR, IDM_EDIT_UNDO, IDM_EDIT_CUT, IDM_EDIT_COPY, IDM_EDIT_PASTE:
            MessageBeep(0)
            return 0

        case IDM_BKGND_BLACK, IDM_BKGND_WHITE, IDM_BKGND_LTGRAY, IDM_BKGND_GRAY, IDM_BKGND_DKGRAY:
            CheckMenuItem(hMenu, UINT(iSelection), MF_UNCHECKED)
            iSelection = int32(LOWORD(uint32(wParam)))
            CheckMenuItem(hMenu, UINT(iSelection), MF_CHECKED)

            SetClassLong(hwnd, GCL_HBRBACKGROUND, int32(GetStockObject(int32(idColor[int32(LOWORD(uint32(wParam)))-IDM_BKGND_WHITE]))))

            InvalidateRect(hwnd, nil, true)
            return 0

        case IDM_TIMER_START:
            if int32(SetTimer(hwnd, ID_TIMER, 1000, 0)) > 0 {
                EnableMenuItem(hMenu, IDM_TIMER_START, MF_GRAYED)
                EnableMenuItem(hMenu, IDM_TIMER_STOP, MF_ENABLED)
            }
            return 0

        case IDM_TIMER_STOP:
            KillTimer(hwnd, ID_TIMER)
            EnableMenuItem(hMenu, IDM_TIMER_START, MF_ENABLED)
            EnableMenuItem(hMenu, IDM_TIMER_STOP, MF_GRAYED)
            return 0

        case IDM_APP_HELP:
            MessageBox(hwnd, _T("Help not yet implemented!"),
                _T("MenuDemo"), MB_ICONEXCLAMATION|MB_OK)
            return 0

        case IDM_APP_ABOUT:
            MessageBox(hwnd, _T("Menu Demonstration Program\n(c) Charles Petzold, 1998"),
                _T("MenuDemo"), MB_ICONINFORMATION|MB_OK)
            return 0
        }
        return DefWindowProc(hwnd, msg, wParam, lParam)
    })

    app.On(WM_TIMER, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        MessageBeep(0)
        return 0
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })
    app.MenuName = _T("MenuDemo")
    app.Init("MenuDemo", "Menu Demonstration")
    app.Run()
}
