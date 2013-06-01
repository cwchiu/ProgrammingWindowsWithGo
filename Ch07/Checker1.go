package main

import (
    //"fmt"
    //"math"
    . "github.com/cwchiu/go-winapi"
    "syscall"
    "unsafe"
)

const (
    DIVISIONS = 5
)

var (
    fState           [DIVISIONS][DIVISIONS]BOOL
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
    x := int32(LOWORD(uint32(lParam))) / cxBlock
    y := int32(HIWORD(uint32(lParam))) / cyBlock

    if x < DIVISIONS && y < DIVISIONS {
        fState[x][y] ^= 1
        rect := RECT{x * cxBlock, y * cyBlock, (x + 1) * cxBlock, (y + 1) * cyBlock}

        InvalidateRect(hwnd, &rect, false)
    } else {
        MessageBeep(0)
    }

    return 0
}

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    cxBlock = int32(LOWORD(uint32(lParam))) / DIVISIONS
    cyBlock = int32(HIWORD(uint32(lParam))) / DIVISIONS
    return 0
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)
    var x, y int32
    for x = 0; x < DIVISIONS; x++ {
        for y = 0; y < DIVISIONS; y++ {
            Rectangle(hdc, x*cxBlock, y*cyBlock,
                (x+1)*cxBlock, (y+1)*cyBlock)

            if fState[x][y] == TRUE {
                MoveToEx(hdc, x*cxBlock, y*cyBlock, nil)
                LineTo(hdc, (x+1)*cxBlock, (y+1)*cyBlock)
                MoveToEx(hdc, x*cxBlock, (y+1)*cyBlock, nil)
                LineTo(hdc, (x+1)*cxBlock, y*cyBlock)
            }
        }
    }

    EndPaint(hwnd, &ps)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr

var table map[uint32]EventHandler

func main() {
    table = make(map[uint32]EventHandler)
    table[WM_LBUTTONDOWN] = OnLeftButtonDown
    table[WM_SIZE] = OnSize
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy

    initWindow("Checker1", "Checker1 Mouse Hit-Test Demo", syscall.NewCallback(WndProc))
}
