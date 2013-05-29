package main

import (
    //"fmt"
    //"math"
    "math/rand"
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
    UpdateWindow (hWnd)
    
    var msg MSG
    for {
        if PeekMessage(&msg, HWND_TOP, 0, 0, PM_REMOVE) {
            if msg.Message == WM_QUIT {
                break
            }
            
            TranslateMessage(&msg)
            DispatchMessage(&msg)
        } else {
            DrawRectangle(hWnd)
        }
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

func RGB(r, g, b int32) COLORREF{    
    r = r & 0xff
    g = (g & 0xff) <<4
    b = (b & 0xff) <<8
    return COLORREF( r | g | b )
    
}

func DrawRectangle (hwnd HWND){
     
     if (cxClient == 0 || cyClient == 0){
          return 
     }
     
     var rect RECT
     SetRect (&rect, uint32(rand.Int31n(cxClient)) , uint32(rand.Int31n (cyClient)) ,
                     uint32(rand.Int31n(cxClient)) , uint32(rand.Int31n (cyClient)) ) 
     
     hBrush := CreateSolidBrush (RGB (rand.Int31n(255), rand.Int31n(255), rand.Int31n(255))) 

     hdc := GetDC (hwnd) 
     FillRect (hdc, &rect, hBrush) 
     ReleaseDC (hwnd, hdc) 
     DeleteObject(HGDIOBJ(hBrush))
}     

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    cxClient = int32(LOWORD(uint32(lParam)))
    cyClient = int32(HIWORD(uint32(lParam)))        
    return 0
}

func OnDestroy(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    PostQuitMessage(0)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr
var table map[uint32]EventHandler
func main() {    
    table = make(map[uint32]EventHandler)
    table[WM_SIZE] = OnSize
    table[WM_DESTROY] = OnDestroy
        
    initWindow("RandRect", "Random Rectangles", syscall.NewCallback(WndProc))
}
