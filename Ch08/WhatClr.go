package main

import (
    "fmt"
    //"math"
    . "github.com/cwchiu/go-winapi"
    "syscall"
    "unsafe"
)

const (
    ID_TIMER = 1
)

var (
    cr, crLast COLORREF
    hdcScreen HDC
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func FindWindowSize (pcxWindow *int32, pcyWindow *int32){
     
     var tm TEXTMETRIC
     
     hdcScreen := CreateIC (_T ("DISPLAY"), nil, nil, nil) 
     GetTextMetrics (hdcScreen, &tm) 
     DeleteDC (hdcScreen) 
     
     *pcxWindow = 2 * GetSystemMetrics (SM_CXBORDER)  + 
                        12 * tm.TmAveCharWidth 

     *pcyWindow = 2 * GetSystemMetrics (SM_CYBORDER)  +
                       GetSystemMetrics (SM_CYCAPTION) + 
                         2 * tm.TmHeight 
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
    wc.HbrBackground = HBRUSH(GetStockObject (WHITE_BRUSH))
    wc.LpszMenuName = nil
    wc.LpszClassName = szAppName

    if atom := RegisterClassEx(&wc); atom == 0 {
        panic("RegisterClassEx")
    }
    
    var cxWindow, cyWindow int32
    FindWindowSize (&cxWindow, &cyWindow) 
    
    hWnd := CreateWindowEx(
        0,
        szAppName,
        _T(title),
        WS_OVERLAPPEDWINDOW|WS_BORDER|WS_CAPTION|WS_SYSMENU|WS_MAXIMIZEBOX|WS_MINIMIZEBOX,
        CW_USEDEFAULT,
        CW_USEDEFAULT,
        cxWindow,
        cyWindow,
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
 DeleteDC (hdcScreen)
          KillTimer (hwnd, ID_TIMER) 
    PostQuitMessage(0)
    return 0
}

func OnCreate(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    hdcScreen = CreateDC (_T("DISPLAY"), nil, nil, nil)  
    SetTimer(hwnd, ID_TIMER, 1000, 0)
    return 0
}

func OnTimer(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var pt POINT
          GetCursorPos (&pt) 
          cr  = GetPixel (hdcScreen, pt.X, pt.Y) 
          
          if (cr != crLast){
               crLast = cr 
               InvalidateRect (hwnd, nil, false) 
          }
          
    return 0
}

func OnDisplayChange(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
          DeleteDC (hdcScreen) 
          hdcScreen = CreateDC (_T("DISPLAY"), nil, nil, nil) 
          return 0 
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)
    var rc RECT
          GetClientRect (hwnd, &rc) 
          
          szBuffer := fmt.Sprintf("  %02X %02X %02X  ",
                    GetRValue (cr), GetGValue (cr), GetBValue (cr)) 
          
          DrawText (hdc, _T(szBuffer), -1, &rc,
                    DT_SINGLELINE | DT_CENTER | DT_VCENTER) 
          
    EndPaint(hwnd, &ps)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr

var table map[uint32]EventHandler

func main() {
    table = make(map[uint32]EventHandler)
    table[WM_CREATE] = OnCreate
    table[WM_DISPLAYCHANGE] = OnDisplayChange
    table[WM_TIMER] = OnTimer
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy

    initWindow("WhatClr", "What Color", syscall.NewCallback(WndProc))
}
