package main

import (
    //"fmt"
    //"math"
    . "github.com/cwchiu/go-winapi"
    "syscall"
    "unsafe"
)

const ()

var (
    fBlocking, fValidBox             BOOL
    ptBeg, ptEnd, ptBoxBeg, ptBoxEnd POINT
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

func DrawBoxOutline(hwnd HWND, ptBeg, ptEnd POINT) {
    hdc := GetDC(hwnd)

    SetROP2(hdc, R2_NOT)
    SelectObject(hdc, GetStockObject(NULL_BRUSH))
    Rectangle(hdc, ptBeg.X, ptBeg.Y, ptEnd.X, ptEnd.Y)

    ReleaseDC(hwnd, hdc)
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

func OnMouseMove(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    if fBlocking == TRUE {
        SetCursor(LoadCursor(HINSTANCE(0), MAKEINTRESOURCE(IDC_CROSS)))

        DrawBoxOutline(hwnd, ptBeg, ptEnd)

        ptEnd.X = int32(LOWORD(uint32(lParam)))
        ptEnd.Y = int32(HIWORD(uint32(lParam)))

        DrawBoxOutline(hwnd, ptBeg, ptEnd)
    }
    return 0
}

func OnLeftButtonUp(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    if fBlocking == TRUE {
        DrawBoxOutline(hwnd, ptBeg, ptEnd)

        ptBoxBeg = ptBeg
        ptBoxEnd.X = int32(LOWORD(uint32(lParam)))
        ptBoxEnd.Y = int32(HIWORD(uint32(lParam)))

        ReleaseCapture () 
        SetCursor(LoadCursor(HINSTANCE(0), MAKEINTRESOURCE(IDC_ARROW)))

        fBlocking = FALSE
        fValidBox = TRUE

        InvalidateRect(hwnd, nil, true)
    }
    return 0
}

func OnChar(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    if fBlocking == TRUE && (wParam == '\x1B') { // i.e., Escape

        DrawBoxOutline(hwnd, ptBeg, ptEnd)

        ReleaseCapture ()
        SetCursor(LoadCursor(HINSTANCE(0), MAKEINTRESOURCE(IDC_ARROW)))

        fBlocking = FALSE
    }
    return 0
}

func OnLeftButtonDown(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    ptEnd.X = int32(LOWORD(uint32(lParam)))
    ptBeg.X = ptEnd.X
    ptEnd.Y = int32(HIWORD(uint32(lParam)))
    ptBeg.Y = ptEnd.Y

    DrawBoxOutline(hwnd, ptBeg, ptEnd)

    SetCapture (hwnd)
    SetCursor(LoadCursor(HINSTANCE(0), MAKEINTRESOURCE(IDC_CROSS)))

    fBlocking = TRUE

    return 0
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)
    if fValidBox == TRUE {
        SelectObject(hdc, GetStockObject(BLACK_BRUSH))
        Rectangle(hdc, ptBoxBeg.X, ptBoxBeg.Y,
            ptBoxEnd.X, ptBoxEnd.Y)
    }

    if fBlocking == TRUE {
        SetROP2(hdc, R2_NOT)
        SelectObject(hdc, GetStockObject(NULL_BRUSH))
        Rectangle(hdc, ptBeg.X, ptBeg.Y, ptEnd.X, ptEnd.Y)
    }

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
    table[WM_CHAR] = OnChar
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy

    initWindow("BlokOut2", "Mouse Button & Capture Demo", syscall.NewCallback(WndProc))
}
