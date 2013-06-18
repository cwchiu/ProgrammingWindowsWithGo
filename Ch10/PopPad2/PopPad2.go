package main

import (    
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"    
)

const (
    ID_EDIT             = 1
    IDM_FILE_NEW        = 40001
    IDM_FILE_OPEN       = 40002
    IDM_FILE_SAVE       = 40003
    IDM_FILE_SAVE_AS    = 40004
    IDM_FILE_PRINT      = 40005
    IDM_APP_EXIT        = 40006
    IDM_EDIT_UNDO       = 40007
    IDM_EDIT_CUT        = 40008
    IDM_EDIT_COPY       = 40009
    IDM_EDIT_PASTE      = 40010
    IDM_EDIT_CLEAR      = 40011
    IDM_EDIT_SELECT_ALL = 40012
    IDM_HELP_HELP       = 40013
    IDM_APP_ABOUT       = 40014
    szAppName           = "PopPad2"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func AskConfirmation(hwnd HWND) int32 {
    return MessageBox(hwnd, _T("Really want to close PopPad2?"),
        _T(szAppName), MB_YESNO|MB_ICONQUESTION)
}

func main() {
    app, _ := NewApp()
    var hwndEdit HWND

    app.On(WM_CREATE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        hwndEdit = CreateWindowEx(0, _T("edit"), nil,
            WS_CHILD|WS_VISIBLE|WS_HSCROLL|WS_VSCROLL|
                WS_BORDER|ES_LEFT|ES_MULTILINE|
                ES_AUTOHSCROLL|ES_AUTOVSCROLL,
            0, 0, 0, 0, hwnd, HMENU(ID_EDIT),
            app.HInstance, nil)
        return 0
    })

    app.On(WM_COMMAND, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        cmd := LOWORD(uint32(wParam))
        if cmd == ID_EDIT {
            if HIWORD(uint32(wParam)) == EN_ERRSPACE ||
                HIWORD(uint32(wParam)) == EN_MAXTEXT {

                MessageBox(hwnd, _T("Edit control out of space."),
                    _T("PopPad1"), MB_OK|MB_ICONSTOP)

            }

        }

        switch cmd {
        case IDM_FILE_NEW, IDM_FILE_PRINT, IDM_FILE_OPEN, IDM_FILE_SAVE, IDM_FILE_SAVE_AS:
            MessageBeep(0)
            return 0

        case IDM_APP_EXIT:
            SendMessage(hwnd, WM_CLOSE, 0, 0)
            return 0

        case IDM_EDIT_UNDO:
            SendMessage(hwndEdit, WM_UNDO, 0, 0)
            return 0

        case IDM_EDIT_CUT:
            SendMessage(hwndEdit, WM_CUT, 0, 0)
            return 0

        case IDM_EDIT_COPY:
            SendMessage(hwndEdit, WM_COPY, 0, 0)
            return 0

        case IDM_EDIT_PASTE:
            SendMessage(hwndEdit, WM_PASTE, 0, 0)
            return 0

        case IDM_EDIT_CLEAR:
            SendMessage(hwndEdit, WM_CLEAR, 0, 0)
            return 0

        case IDM_EDIT_SELECT_ALL:
            var end int = -1
            SendMessage(hwndEdit, EM_SETSEL, 0, uintptr(end))
            return 0

        case IDM_HELP_HELP:
            MessageBox(hwnd, _T("Help not yet implemented!"),
                _T(szAppName), MB_OK|MB_ICONEXCLAMATION)
            return 0

        case IDM_APP_ABOUT:
            MessageBox(hwnd, _T("POPPAD2 (c) Chui-Wen Chiu, 2013"),
                _T(szAppName), MB_OK|MB_ICONINFORMATION)
            return 0

        }

        return DefWindowProc(hwnd, msg, wParam, lParam)
    })

    app.On(WM_INITMENUPOPUP, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        if int32(lParam) == 1 {
            var iUndo UINT
            if SendMessage(hwndEdit, EM_CANUNDO, 0, 0) > 0 {
                iUndo = MF_ENABLED
            } else {
                iUndo = MF_GRAYED
            }
            EnableMenuItem(HMENU(wParam), IDM_EDIT_UNDO, iUndo)

            var iPaste UINT
            if IsClipboardFormatAvailable(CF_TEXT) {
                iPaste = MF_ENABLED
            } else {
                iPaste = MF_GRAYED
            }

            EnableMenuItem(HMENU(wParam), IDM_EDIT_PASTE, iPaste)

            iSelect := SendMessage(hwndEdit, EM_GETSEL, 0, 0)

            var iEnable UINT
            if HIWORD(uint32(iSelect)) == LOWORD(uint32(iSelect)) {
                iEnable = MF_GRAYED
            } else {
                iEnable = MF_ENABLED
            }

            EnableMenuItem(HMENU(wParam), IDM_EDIT_CUT, iEnable)
            EnableMenuItem(HMENU(wParam), IDM_EDIT_COPY, iEnable)
            EnableMenuItem(HMENU(wParam), IDM_EDIT_CLEAR, iEnable)
            return 0
        }
        return DefWindowProc(hwnd, msg, wParam, lParam)
    })

    app.On(WM_SIZE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        MoveWindow(hwndEdit, 0, 0, int32(LOWORD(uint32(lParam))), int32(HIWORD(uint32(lParam))), true)

        return 0
    })

    app.On(WM_SETFOCUS, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        SetFocus(hwndEdit)
        return 0
    })

    app.On(WM_QUERYENDSESSION, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        if IDYES == AskConfirmation(hwnd) {
            return 1
        } else {
            return 0
        }
    })

    app.On(WM_CLOSE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        if IDYES == AskConfirmation(hwnd) {
            DestroyWindow(hwnd)
        }
        return 0
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })
    app.Icon = LoadIcon(app.HInstance, _T("POPPAD2"))
    app.MenuName = _T(szAppName)
    app.Init(szAppName, szAppName)
    app.Run()
}
