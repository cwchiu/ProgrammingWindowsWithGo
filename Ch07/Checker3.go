package main

import (
    //"fmt"
    //"math"
    . "github.com/cwchiu/go-winapi"
    "syscall"
    "unsafe"
)

const (
    DIVISIONS    = 5
    szChildClass = "Checker3_Child"
)

var (
    fState           [DIVISIONS][DIVISIONS]BOOL
    hwndChild        [DIVISIONS][DIVISIONS]HWND
    cxBlock, cyBlock int32
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func Min(a, b int32) int32 {
    if a < b {
        return a
    } else {
        return b
    }
}

func Max(a, b int32) int32 {
    if a < b {
        return b
    } else {
        return a
    }
}

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
    wc.HbrBackground = COLOR_BTNFACE + 1 //HANDLE(hBrush)
    wc.LpszMenuName = nil
    wc.LpszClassName = szAppName

    if atom := RegisterClassEx(&wc); atom == 0 {
        panic("RegisterClassEx")
    }

    // Child Window
    wc.LpfnWndProc = syscall.NewCallback(ChildWndProc)
    wc.CbWndExtra = 64
    wc.HIcon = HICON(0)
    wc.LpszClassName = _T(szChildClass)

    RegisterClassEx(&wc)

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

func ChildWndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {

    switch int(msg) {
    case WM_CREATE:
        SetWindowLong(hwnd, 0, 0) // on/off flag
        return 0

    case WM_LBUTTONDOWN:
        SetWindowLong(hwnd, 0, 1^GetWindowLong(hwnd, 0))
        InvalidateRect(hwnd, nil, false)
        return 0

    case WM_PAINT:
        var ps PAINTSTRUCT
        hdc := BeginPaint(hwnd, &ps)
        var rect RECT
        GetClientRect(hwnd, &rect)
        Rectangle(hdc, 0, 0, rect.Right, rect.Bottom)

        if GetWindowLong(hwnd, 0) > 0 {
            MoveToEx(hdc, 0, 0, nil)
            LineTo(hdc, rect.Right, rect.Bottom)
            MoveToEx(hdc, 0, rect.Bottom, nil)
            LineTo(hdc, rect.Right, 0)
        }

        EndPaint(hwnd, &ps)
        return 0
    }
    return DefWindowProc(hwnd, msg, wParam, lParam)
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

func OnLeftButtonDown(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    MessageBeep(0)
    return 0
}

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    cxBlock = int32(LOWORD(uint32(lParam))) / DIVISIONS
    cyBlock = int32(HIWORD(uint32(lParam))) / DIVISIONS
    var x, y int32
    for x = 0; x < DIVISIONS; x++ {
        for y = 0; y < DIVISIONS; y++ {
            MoveWindow(hwndChild[x][y],
                x*cxBlock, y*cyBlock,
                cxBlock, cyBlock, true)
        }
    }
    return 0
}

func OnCreate(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var x, y int32

    for x = 0; x < DIVISIONS; x++ {
        for y = 0; y < DIVISIONS; y++ {
            hwndChild[x][y] = CreateWindowEx(WS_EX_WINDOWEDGE,
                _T(szChildClass), nil,
                WS_CHILDWINDOW|WS_VISIBLE,
                0, 0, 0, 0,
                hwnd, HMENU(y<<8|x),
                HINSTANCE(GetWindowLong(hwnd, GWL_HINSTANCE)),
                nil)
        }
    }
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr

var table map[uint32]EventHandler

func main() {
    table = make(map[uint32]EventHandler)
    table[WM_CREATE] = OnCreate
    table[WM_SIZE] = OnSize
    table[WM_LBUTTONDOWN] = OnLeftButtonDown
    table[WM_DESTROY] = OnDestroy

    initWindow("Checker3", "Checker3 Mouse Hit-Test Demo", syscall.NewCallback(WndProc))
}
