package main

import (
    "fmt"
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
    wc.HbrBackground = HBRUSH(GetStockObject (WHITE_BRUSH))
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
    if table[msg] != nil {
       return table[msg](hwnd, msg, wParam, lParam)       
    }
    return DefWindowProc(hwnd, msg, wParam, lParam)
}

var (
    cxChar, cyChar int32
) 

func OnCreate(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    var tm TEXTMETRIC
    hdc := GetDC (hwnd) 
    SelectObject (hdc, GetStockObject (SYSTEM_FIXED_FONT)) 

    GetTextMetrics (hdc, &tm) 
    cxChar = tm.TmAveCharWidth 
    cyChar = tm.TmHeight + tm.TmExternalLeading 

    ReleaseDC (hwnd, hdc) 
    return 0
}

func OnDestroy(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    PostQuitMessage(0)
    return 0
}
     
func Show (hwnd HWND, hdc HDC, xText, yText, iMapMode int32, szMapMode string){
     var  rect RECT
     
     SaveDC (hdc) 
     
     SetMapMode (hdc, iMapMode) 
     GetClientRect (hwnd, &rect) 
     DPtoLP (hdc, (*POINT)(unsafe.Pointer(&rect)), 2) 
     
     RestoreDC (hdc, -1) 
     szBuffer := fmt.Sprintf("%-20s %7d %7d %7d %7d", szMapMode,
              rect.Left, rect.Right, rect.Top, rect.Bottom)
     TextOut (hdc, xText, yText, syscall.StringToUTF16Ptr(szBuffer), int32(len(szBuffer))) 
} 
     
func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    szHeading := "Mapping Mode            Left   Right     Top  Bottom"
    szUndLine := "------------            ----   -----     ---  ------"
     
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)
    SelectObject (hdc, GetStockObject (SYSTEM_FIXED_FONT)) 

    SetMapMode (hdc, MM_ANISOTROPIC) 
    SetWindowExtEx (hdc, 1, 1, nil) 
    SetViewportExtEx (hdc, cxChar, cyChar, nil) 

    TextOut (hdc, 1, 1, syscall.StringToUTF16Ptr(szHeading), int32(len(szHeading))) 
    TextOut (hdc, 1, 2, syscall.StringToUTF16Ptr(szUndLine), int32(len(szUndLine))) 

    Show (hwnd, hdc, 1, 3, MM_TEXT,      "TEXT (pixels)")
    Show (hwnd, hdc, 1, 4, MM_LOMETRIC,  "LOMETRIC (.1 mm)")
    Show (hwnd, hdc, 1, 5, MM_HIMETRIC,  "HIMETRIC (.01 mm)")
    Show (hwnd, hdc, 1, 6, MM_LOENGLISH, "LOENGLISH (.01 in)")
    Show (hwnd, hdc, 1, 7, MM_HIENGLISH, "HIENGLISH (.001 in)")
    Show (hwnd, hdc, 1, 8, MM_TWIPS,     "TWIPS (1/1440 in)")
    
    EndPaint(hwnd, &ps)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr
var table map[uint32]EventHandler
func main() {    
    table = make(map[uint32]EventHandler)
    table[WM_CREATE] = OnCreate
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy
    
    initWindow("WhatSize", "What Size is the Window?", syscall.NewCallback(WndProc))
}
