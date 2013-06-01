package main

import (
    //"fmt"
    //"math"
    . "github.com/cwchiu/go-winapi"
    "syscall"
    "unsafe"
)

const (
    MAXPOINTS = 1000
)

var (
    iCount int32
    pt     []POINT = make([]POINT, MAXPOINTS)
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

func OnLeftButtonDown(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    iCount = 0
    InvalidateRect(hwnd, nil, true)
    return 0
}

func OnLeftButtonUp(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    InvalidateRect(hwnd, nil, false)
    return 0
}

func OnMouseMove(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    if (wParam&MK_LBUTTON) > 0 && iCount < 1000 {
        pt[iCount].X = int32(LOWORD(uint32(lParam)))
        iCount++
        pt[iCount].Y = int32(HIWORD(uint32(lParam)))

        hdc := GetDC(hwnd)
        SetPixel(hdc, int32(LOWORD(uint32(lParam))), int32(HIWORD(uint32(lParam))), 0)
        ReleaseDC(hwnd, hdc)
    }
    return 0
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)

    SetCursor(LoadCursor(HINSTANCE(0), MAKEINTRESOURCE(IDC_WAIT)))
    ShowCursor(TRUE)
    var i, j int32
    for i = 0; i < iCount-1; i++ {
        for j = i + 1; j < iCount; j++ {
            MoveToEx(hdc, pt[i].X, pt[i].Y, nil)
            LineTo(hdc, pt[j].X, pt[j].Y)
        }
    }
    ShowCursor(FALSE)
    SetCursor(LoadCursor(HINSTANCE(0), MAKEINTRESOURCE(IDC_ARROW)))

    EndPaint(hwnd, &ps)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr

var table map[uint32]EventHandler

func main() {
    table = make(map[uint32]EventHandler)
    table[WM_LBUTTONDOWN] = OnLeftButtonDown
    table[WM_LBUTTONUP] = OnLeftButtonUp
    table[WM_MOUSEMOVE] = OnMouseMove
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy

    initWindow("Connect", "Connect-the-Points Mouse Demo", syscall.NewCallback(WndProc))
}
