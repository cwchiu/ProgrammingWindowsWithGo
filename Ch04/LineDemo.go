package main

import (
    //"fmt"
    //"math"
    "syscall"
    "unsafe"
    . "github.com/cwchiu/go-winapi"
)

func initWindow(appName string, title string, wndproc uintptr) {
    _T := syscall.StringToUTF16Ptr

    hInst := GetModuleHandle(nil)
    if hInst == 0 {
        panic("GetModuleHandle")
    }

    hIcon := LoadIcon(0, (*uint16)(unsafe.Pointer(uintptr(IDI_APPLICATION))))
    if hIcon == 0 {
        panic("LoadIcon")
    }

    hCursor := LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(IDC_ARROW))))
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
    // UpdateWindow
    ShowWindow(hWnd, SW_NORMAL)

    var msg MSG
    for GetMessage(&msg, HWND_TOP, 0, 0) == TRUE {
        TranslateMessage(&msg)
        DispatchMessage(&msg)
    }
}


func WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    switch msg {
    case WM_SIZE:
        return OnSize(hwnd, msg, wParam, lParam)        
    case WM_DESTROY:
        return OnDestroy(hwnd, msg, wParam, lParam)        
    case WM_PAINT:
        return OnPaint(hwnd, msg, wParam, lParam)                
    }
    return DefWindowProc(hwnd, msg, wParam, lParam)
}

var (
    cxClient, cyClient int32
)

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    cxClient = int32(LOWORD(uint32(lParam)))
    cyClient = int32(HIWORD(uint32(lParam)))        
    return 0
}

func OnDestroy(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    PostQuitMessage(0)
    return 0
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)
  
    Rectangle (hdc, cxClient / 8,cyClient / 8, 7 * cxClient / 8, 7 * cyClient / 8) 

    MoveToEx  (hdc, 0, 0, nil) 
    LineTo    (hdc, cxClient, cyClient)

    MoveToEx  (hdc,        0, cyClient, nil) 
    LineTo    (hdc, cxClient,        0) 

    Ellipse   (hdc,     cxClient / 8,     cyClient / 8,
                    7 * cxClient / 8, 7 * cyClient / 8) 

    RoundRect (hdc,     cxClient / 4,     cyClient / 4,
                    3 * cxClient / 4, 3 * cyClient / 4,
                        cxClient / 4,     cyClient / 4) 
    EndPaint(hwnd, &ps)
    return 0
}

func main() {
    initWindow("LineDemo", "Line Demonstration", syscall.NewCallback(WndProc))
}
