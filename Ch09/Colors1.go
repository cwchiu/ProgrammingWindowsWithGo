package main

import (
    "fmt"
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"
    //"unsafe"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

type ScrollProc func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr)

func main() {
    app, _ := NewApp()
    app.BackgroundBrush = HBRUSH(CreateSolidBrush(0))

    var cyChar uint32
    var cxClient, cyClient uint32
    var color [3]int32
    var hBrushStatic HBRUSH
    var hBrush [3]HBRUSH
    var rcColor RECT
    var hwndScroll, hwndLabel, hwndValue [3]HWND
    var OldScroll [3]int32
    var hwndRect HWND
    var idFocus int32

    var crPrim [3]COLORREF = [3]COLORREF{
        COLORREF(0xff),
        COLORREF(0xff00),
        COLORREF(0xff0000),
    }

    app.On(WM_CREATE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        hwndRect = CreateWindowEx(WS_EX_WINDOWEDGE, _T("static"), nil,
            WS_CHILD|WS_VISIBLE|SS_WHITERECT,
            0, 0, 0, 0,
            hwnd, HMENU(9), app.HInstance, nil)

        var proc ScrollProc = func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
            id := GetWindowLong(hwnd, GWL_ID)
            switch msg {
            case WM_KEYDOWN:
                if wParam == VK_TAB {
                    var next_id int32
                    if GetKeyState(VK_SHIFT) < 0 {
                        next_id = (id + 2) % 3
                    } else {
                        next_id = (id + 1) % 3
                    }
                    SetFocus(GetDlgItem(GetParent(hwnd), next_id))
                }
                break
            case WM_SETFOCUS:
                idFocus = id
                break
            }
            return CallWindowProc(uintptr(OldScroll[id]), hwnd, msg, wParam, lParam)
        }

        var szColorLabel []string = []string{"Red", "Green", "Blue"}
        var i uint32
        for i = 0; i < 3; i++ {
            // The three scroll bars have IDs 0, 1, and 2, with
            // scroll bar ranges from 0 through 255.
            hwndScroll[i] = CreateWindowEx(0, _T("scrollbar"), nil,
                WS_CHILD|WS_VISIBLE|
                    WS_TABSTOP|SBS_VERT,
                0, 0, 0, 0,
                hwnd, HMENU(i), app.HInstance, nil)

            SetScrollRange(hwndScroll[i], SB_CTL, 0, 255, FALSE)
            SetScrollPos(hwndScroll[i], SB_CTL, 0, FALSE)

            // The three color-name labels have IDs 3, 4, and 5,
            // and text strings "Red", "Green", and "Blue".
            hwndLabel[i] = CreateWindowEx(WS_EX_WINDOWEDGE, _T("static"), _T(szColorLabel[i]),
                WS_CHILD|WS_VISIBLE|SS_CENTER,
                0, 0, 0, 0,
                hwnd, HMENU(i+3),
                app.HInstance, nil)

            // The three color-value text fields have IDs 6, 7,
            // and 8, and initial text strings of "0".
            hwndValue[i] = CreateWindowEx(WS_EX_WINDOWEDGE, _T("static"), _T("0"),
                WS_CHILD|WS_VISIBLE|SS_CENTER,
                0, 0, 0, 0,
                hwnd, HMENU(i+6),
                app.HInstance, nil)

            OldScroll[i] = SetWindowLong(hwndScroll[i], GWL_WNDPROC, int32(syscall.NewCallback(proc)))

            hBrush[i] = CreateSolidBrush(crPrim[i])
        }

        hBrushStatic = CreateSolidBrush(COLORREF(GetSysColor(COLOR_BTNHIGHLIGHT)))

        cyChar = uint32(HIWORD(uint32(GetDialogBaseUnits())))
        return 0
    })

    app.On(WM_SIZE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        cxClient = uint32(LOWORD(uint32(lParam)))
        cyClient = uint32(HIWORD(uint32(lParam)))

        SetRect(&rcColor, cxClient/2, 0, cxClient, cyClient)

        MoveWindow(hwndRect, 0, 0, int32(cxClient/2), int32(cyClient), true)

        var i uint32
        for i = 0; i < 3; i++ {
            MoveWindow(hwndScroll[i],
                int32((2*i+1)*cxClient/14),
                int32(2*cyChar),
                int32(cxClient/14),
                int32(cyClient-4*cyChar),
                true)

            MoveWindow(hwndLabel[i],
                int32((4*i+1)*cxClient/28),
                int32(cyChar/2),
                int32(cxClient/7),
                int32(cyChar),
                true)

            MoveWindow(hwndValue[i],
                int32((4*i+1)*cxClient/28),
                int32(cyClient-3*cyChar/2),
                int32(cxClient/7),
                int32(cyChar),
                true)
        }
        SetFocus(hwnd)

        return 0
    })

    app.On(WM_CTLCOLORSTATIC, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        i := GetWindowLong(HWND(lParam), GWL_ID)

        if i >= 3 && i <= 8 { // static text controls
            SetTextColor(HDC(wParam), crPrim[i%3])
            SetBkColor(HDC(wParam), COLORREF(GetSysColor(COLOR_BTNHIGHLIGHT)))
            return uintptr(hBrushStatic)
        }

        return DefWindowProc(hwnd, msg, wParam, lParam)
    })

    app.On(WM_CTLCOLORSCROLLBAR, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        i := GetWindowLong(HWND(lParam), GWL_ID)
        return uintptr(hBrush[i])
    })

    app.On(WM_SYSCOLORCHANGE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        DeleteObject(HGDIOBJ(hBrushStatic))
        hBrushStatic = CreateSolidBrush(COLORREF(GetSysColor(COLOR_BTNHIGHLIGHT)))
        return 0
    })

    app.On(WM_SETFOCUS, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        SetFocus(hwndScroll[idFocus])
        return 0
    })

    app.On(WM_VSCROLL, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        i := GetWindowLong(HWND(lParam), GWL_ID)

        switch LOWORD(uint32(wParam)) {
        case SB_PAGEDOWN:
            color[i] += 15 // fall through
            color[i] = Min(255, color[i]+1)
        case SB_LINEDOWN:
            color[i] = Min(255, color[i]+1)
        case SB_PAGEUP:
            color[i] -= 15 // fall through
            color[i] = Max(0, color[i]-1)
        case SB_LINEUP:
            color[i] = Max(0, color[i]-1)
        case SB_TOP:
            color[i] = 0
        case SB_BOTTOM:
            color[i] = 255
        case SB_THUMBTRACK, SB_THUMBPOSITION:
            color[i] = int32(HIWORD(uint32(wParam)))
        }

        SetScrollPos(hwndScroll[i], SB_CTL, color[i], TRUE)
        SetWindowText(hwndValue[i], _T(fmt.Sprintf("%d", color[i])))

        DeleteObject(HGDIOBJ(SetClassLong(hwnd, GCL_HBRBACKGROUND, int32(CreateSolidBrush(RGB(color[0], color[1], color[2]))))))

        InvalidateRect(hwnd, &rcColor, true)

        return 0
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        var i uint32
        for i = 0; i < 3; i++ {
            DeleteObject(HGDIOBJ(hBrush[i]))
        }
        DeleteObject(HGDIOBJ(hBrushStatic))

        PostQuitMessage(0)
        return 0
    })

    app.Init("Colors1", "Color Scroll")
    app.Run()
}
