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
    switch msg {
    case WM_SIZE:
        return OnSize(hwnd, msg, wParam, lParam)        
    case WM_DESTROY:
        return OnDestroy(hwnd, msg, wParam, lParam)        
    case WM_PAINT:
        return OnPaint(hwnd, msg, wParam, lParam)                
    case WM_LBUTTONDOWN:
    case WM_RBUTTONDOWN:
    case WM_MOUSEMOVE:
        return OnMouseMove(hwnd, msg, wParam, lParam)
    }
    return DefWindowProc(hwnd, msg, wParam, lParam)
}

var (
    cxClient, cyClient int32
)

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    cxClient = int32(LOWORD(uint32(lParam)))
    cyClient = int32(HIWORD(uint32(lParam)))        
    
    apt[0].X = cxClient / 4 
    apt[0].Y = cyClient / 2 
    
    apt[1].X = cxClient / 2 
    apt[1].Y = cyClient / 4 
    
    apt[2].X =     cxClient / 2 
    apt[2].Y = 3 * cyClient / 4 
    
    apt[3].X = 3 * cxClient / 4 
    apt[3].Y =     cyClient / 2 
    return 0
}

func OnMouseMove(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    hit_left_button := wParam & MK_LBUTTON>0
    hit_right_button := wParam & MK_RBUTTON >0
    if (hit_left_button || hit_right_button){
        hdc := GetDC (hwnd) 
        SelectObject (hdc, GetStockObject (WHITE_PEN))          
        DrawBezier (hdc, apt) 

        if (hit_left_button){
            apt[1].X = int32(LOWORD (uint32(lParam)))
            apt[1].Y = int32(HIWORD (uint32(lParam)))
        }

        if (hit_right_button){
            apt[2].X = int32(LOWORD(uint32(lParam))) 
            apt[2].Y = int32(HIWORD(uint32(lParam))) 
        }

        SelectObject (hdc, GetStockObject (BLACK_PEN)) 
        DrawBezier (hdc, apt) 
        ReleaseDC (hwnd, hdc) 
    }
    
    return 0
}

func OnDestroy(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    PostQuitMessage(0)
    return 0
}

var apt []POINT = make([]POINT, 4) ;
func DrawBezier (hdc HDC, apt[]POINT){
     PolyBezier (hdc, &apt[0], 4) 

     MoveToEx (hdc, apt[0].X, apt[0].Y, nil) 
     LineTo   (hdc, apt[1].X, apt[1].Y) 

     MoveToEx (hdc, apt[2].X, apt[2].Y, nil) 
     LineTo   (hdc, apt[3].X, apt[3].Y) 
}
     
func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    InvalidateRect (hwnd, nil, true)
    
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)
  
    DrawBezier (hdc, apt) 
    
    EndPaint(hwnd, &ps)
    return 0
}

func main() {
    initWindow("Bezier", "Bezier Splines", syscall.NewCallback(WndProc))
}
