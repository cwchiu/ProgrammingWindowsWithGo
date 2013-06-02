package main

import (
    "fmt"
    //"math"
    . "github.com/cwchiu/go-winapi"
    "syscall"
    "unsafe"
)

const (
    NUM = 10
)

var (
    rect           RECT
    hwndButton     [NUM]HWND
    cxChar, cyChar int32
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr

var table map[uint32]EventHandler

func initWindow(appName string, title string, wndproc uintptr) {

    hInst := GetModuleHandle(nil)
    if hInst == 0 {
        panic("GetModuleHandle")
    }

    hIcon := LoadIcon(0, (*uint16)(unsafe.Pointer(uintptr(IDI_APPLICATION))))
    if hIcon == 0 {
        panic("LoadIcon")
    }

    hCursor := LoadCursor(0, MAKEINTRESOURCE(IDC_ARROW))
    if hCursor == 0 {
        panic("LoadCursor")
    }

    //hBrush := GetStockObject(WHITE_BRUSH)
    szAppName := _T(appName)

    var wc WNDCLASSEX
    wc.CbSize = uint32(unsafe.Sizeof(wc))
    wc.Style = CS_HREDRAW | CS_VREDRAW
    wc.LpfnWndProc = wndproc
    wc.HInstance = hInst
    wc.HIcon = hIcon
    wc.HCursor = hCursor
    wc.CbClsExtra = 0
    wc.CbWndExtra = 0
    wc.HbrBackground = HBRUSH(GetStockObject(WHITE_BRUSH))
    wc.LpszMenuName = nil
    wc.LpszClassName = szAppName

    if atom := RegisterClassEx(&wc); atom == 0 {
        panic("RegisterClassEx")
    }

    hWnd := CreateWindowEx(
        0,
        szAppName,
        _T(title),
        WS_OVERLAPPEDWINDOW|WS_BORDER|WS_CAPTION|WS_SYSMENU|WS_MAXIMIZEBOX|WS_MINIMIZEBOX,
        CW_USEDEFAULT,
        CW_USEDEFAULT,
        CW_USEDEFAULT,
        CW_USEDEFAULT,
        HWND_TOP,
        0,
        hInst,
        nil)

    if hWnd == 0 {
        panic("CreateWindowEx")
    }

    ShowWindow(hWnd, SW_NORMAL)
    UpdateWindow(hWnd)

    var msg MSG
    for GetMessage(&msg, HWND_TOP, 0, 0) == TRUE {
        TranslateMessage(&msg)
        DispatchMessage(&msg)
    }
}

func WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    if table[msg] != nil {
        return table[msg](hwnd, msg, wParam, lParam)
    }
    return DefWindowProc(hwnd, msg, wParam, lParam)
}

func OnDestroy(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    PostQuitMessage(0)
    return 0
}

type ButtonType struct {
    Style uint32
    Text  string
}

func OnCreate(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    cxChar = int32(LOWORD(uint32(GetDialogBaseUnits())))
    cyChar = int32(HIWORD(uint32(GetDialogBaseUnits())))
    var i int32

    var button []ButtonType = []ButtonType{
        ButtonType{BS_PUSHBUTTON, "PUSHBUTTON"},
        ButtonType{BS_DEFPUSHBUTTON, "DEFPUSHBUTTON"},
        ButtonType{BS_CHECKBOX, "CHECKBOX"},
        ButtonType{BS_AUTOCHECKBOX, "AUTOCHECKBOX"},
        ButtonType{BS_RADIOBUTTON, "RADIOBUTTON"},
        ButtonType{BS_3STATE, "3STATE"},
        ButtonType{BS_AUTO3STATE, "AUTO3STATE"},
        ButtonType{BS_GROUPBOX, "GROUPBOX"},
        ButtonType{BS_AUTORADIOBUTTON, "AUTORADIO"},
        ButtonType{BS_OWNERDRAW, "OWNERDRAW"},
    }

    for i = 0; i < NUM; i++ {
        hwndButton[i] = CreateWindowEx(WS_EX_WINDOWEDGE, _T("button"),
            _T(button[i].Text),
            WS_CHILD|WS_VISIBLE|button[i].Style,
            cxChar, cyChar*(1+2*i),
            20*cxChar, 7*cyChar/4,
            hwnd, HMENU(i),
            HINSTANCE(GetWindowLong(hwnd, GWL_HINSTANCE)), nil)
    }
    return 0
}

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    rect.Left = 24 * cxChar
    rect.Top = 2 * cyChar
    rect.Right = int32(LOWORD(uint32(lParam)))
    rect.Bottom = int32(HIWORD(uint32(lParam)))

    return 0
}

func OnCommand(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    ScrollWindow(hwnd, 0, -cyChar, &rect, &rect)

    hdc := GetDC(hwnd)
    SelectObject(hdc, GetStockObject(SYSTEM_FIXED_FONT))
    szFormat := "%-16s%04X-%04X    %04X-%04X"
    var msg_str string
    if msg == WM_DRAWITEM {
        msg_str = "WM_DRAWITEM"
    } else {
        msg_str = "WM_COMMAND"
    }
    szBuffer := fmt.Sprintf(szFormat,
        msg_str,
        HIWORD(uint32(wParam)), LOWORD(uint32(wParam)),
        HIWORD(uint32(lParam)), LOWORD(uint32(lParam)))
    TextOut(hdc, 24*cxChar, cyChar*(rect.Bottom/cyChar-1),
        _T(szBuffer), int32(len(szBuffer)))

    ReleaseDC(hwnd, hdc)
    ValidateRect(hwnd, &rect)

    return DefWindowProc(hwnd, msg, wParam, lParam)
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {

    InvalidateRect(hwnd, &rect, true)
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)

    SelectObject(hdc, GetStockObject(SYSTEM_FIXED_FONT))
    SetBkMode(hdc, TRANSPARENT)

    szTop := "message            wParam       lParam"
    szUnd := "_______            ______       ______"

    TextOut(hdc, 24*cxChar, cyChar, _T(szTop), int32(len(szTop)))
    TextOut(hdc, 24*cxChar, cyChar, _T(szUnd), int32(len(szUnd)))

    EndPaint(hwnd, &ps)
    return 0
}

func main() {
    table = make(map[uint32]EventHandler)
    table[WM_CREATE] = OnCreate
    table[WM_SIZE] = OnSize
    table[WM_DRAWITEM] = OnCommand
    table[WM_COMMAND] = OnCommand
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy

    initWindow("BtnLook", "Button Look", syscall.NewCallback(WndProc))
}
