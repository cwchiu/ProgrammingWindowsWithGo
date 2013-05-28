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

var aptFigure []POINT = []POINT{ 
    POINT{10,70}, 
    POINT{50,70}, 
    POINT{50,10}, 
    POINT{90,10}, 
    POINT{90,50},
    POINT{30,50}, 
    POINT{30,90}, 
    POINT{70,90}, 
    POINT{70,30}, 
    POINT{10,30},
}
                                     
func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)
    var NUM int32 = int32(len(aptFigure))
    var apt []POINT = make([]POINT, NUM)
    SelectObject (hdc, GetStockObject (GRAY_BRUSH))
    var i int32
    for i = 0 ; i < NUM ; i++{
       apt[i].X = cxClient * aptFigure[i].X / 200 
       apt[i].Y = cyClient * aptFigure[i].Y / 100        
    }

    SetPolyFillMode (hdc, ALTERNATE) 
    Polygon (hdc, &apt[0], NUM) 

    for i = 0 ; i < NUM ; i++{
       apt[i].X += cxClient / 2 
    }

    SetPolyFillMode (hdc, WINDING)
    Polygon (hdc, &apt[0], NUM)
    
    EndPaint(hwnd, &ps)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr
var table map[uint32]EventHandler
func main() {    
    table = make(map[uint32]EventHandler)
    table[WM_SIZE] = OnSize
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy
    
    initWindow("AltWind", "Alternate and Winding Fill Modes", syscall.NewCallback(WndProc))
}
