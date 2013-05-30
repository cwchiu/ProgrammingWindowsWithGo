package main

import (
    //"fmt"
    "math"
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
    if table[msg] != nil {
       return table[msg](hwnd, msg, wParam, lParam)       
    }
    return DefWindowProc(hwnd, msg, wParam, lParam)
}

var (
    cxClient, cyClient int32
    hRgnClip HRGN
)

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    cxClient = int32(LOWORD(uint32(lParam)))
    cyClient = int32(HIWORD(uint32(lParam)))           

    hCursor := SetCursor(LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(IDC_WAIT)))))
    ShowCursor (TRUE) 

    if (hRgnClip == 0) {
       DeleteObject(HGDIOBJ(hRgnClip)) 
    }

    var hRgnTemp []HRGN = make([]HRGN, 6)
    hRgnTemp[0] = CreateEllipticRgn (0, cyClient / 3,
                                   cxClient / 2, 2 * cyClient / 3) 
    hRgnTemp[1] = CreateEllipticRgn (cxClient / 2, cyClient / 3,
                                   cxClient, 2 * cyClient / 3) 
    hRgnTemp[2] = CreateEllipticRgn (cxClient / 3, 0,
                                   2 * cxClient / 3, cyClient / 2) 
    hRgnTemp[3] = CreateEllipticRgn (cxClient / 3, cyClient / 2,
                                   2 * cxClient / 3, cyClient) 
    hRgnTemp[4] = CreateRectRgn (0, 0, 1, 1) 
    hRgnTemp[5] = CreateRectRgn (0, 0, 1, 1) 
    
    hRgnClip    = CreateRectRgn (0, 0, 1, 1) 

    CombineRgn (hRgnTemp[4], hRgnTemp[0], hRgnTemp[1], RGN_OR) 
    CombineRgn (hRgnTemp[5], hRgnTemp[2], hRgnTemp[3], RGN_OR) 
    CombineRgn (hRgnClip,    hRgnTemp[4], hRgnTemp[5], RGN_XOR) 

    var i int32
    for i = 0 ; i < 6 ; i++{
       DeleteObject (HGDIOBJ(hRgnTemp[i])) 
    }

    SetCursor (hCursor) 
    ShowCursor (FALSE) 
          
    return 0
}

func OnDestroy(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    PostQuitMessage(0)
    return 0
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    var ps PAINTSTRUCT
    const TWO_PI float64 = 2*3.1415926    
    
    hdc := BeginPaint(hwnd, &ps)
    
    SetViewportOrgEx (hdc, cxClient / 2, cyClient / 2, nil) 
    SelectClipRgn (hdc, hRgnClip) 

    fRadius := math.Hypot (float64(cxClient) / 2.0, float64(cyClient) / 2.0) 
    var fAngle float64
    for fAngle = 0.0 ; fAngle < TWO_PI ; fAngle += TWO_PI / 360{
        MoveToEx (hdc, 0, 0, nil) 
        LineTo (hdc, (int32) ( fRadius * math.Cos (fAngle) + 0.5),
                (int32) (-fRadius * math.Sin (fAngle) + 0.5)) 
    }
          
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
    
    initWindow("Clover", "Draw a Clover", syscall.NewCallback(WndProc))
}
