package main

import (
    //"fmt"
    "math"
    . "github.com/cwchiu/go-winapi"
    "syscall"
    "unsafe"
)

const (
    ID_TIMER = 1
    TWOPI    =   2 * 3.14159
)

var (
    cxClient, cyClient int32
    stPrevious SYSTEMTIME
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func SetIsotropic (hdc HDC, cxClient, cyClient int32){
     SetMapMode (hdc, MM_ISOTROPIC) 
     SetWindowExtEx (hdc, 1000, 1000, nil) 
     SetViewportExtEx (hdc, cxClient / 2, -cyClient / 2, nil) 
     SetViewportOrgEx (hdc, cxClient / 2,  cyClient / 2, nil) 
}

func RotatePoint (pt []POINT , iNum, iAngle int32){
     var i int32
     var ptTemp  POINT
     
     for i = 0 ; i < iNum ; i++     {
          ptTemp.X = int32 (float64(pt[i].X) * math.Cos (TWOPI * float64(iAngle) / float64(360)) +
               float64(pt[i].Y) * math.Sin (TWOPI * float64(iAngle) / 360.0)) 
          
          ptTemp.Y = int32 (float64(pt[i].Y) * math.Cos (TWOPI * float64(iAngle) / float64(360)) -
               float64(pt[i].X) * math.Sin (TWOPI * float64(iAngle) / float64(360))) 
          
          pt[i] = ptTemp 
     }
}

func DrawClock (hdc HDC){
     var   iAngle int32
     var pt [3]POINT 
     
     for iAngle = 0 ; iAngle < 360 ; iAngle += 6      {
          pt[0].X =   0 
          pt[0].Y = 900 
          
          RotatePoint (pt[0:], 1, iAngle) 
          
          if iAngle % 5 == 0 {
            pt[2].Y = 100 
          }else{
            pt[2].Y = 33
          }
          pt[2].X = pt[2].Y
          
          pt[0].X -= pt[2].X / 2 
          pt[0].Y -= pt[2].Y / 2 
          
          pt[1].X  = pt[0].X + pt[2].X
          pt[1].Y  = pt[0].Y + pt[2].Y 
          
          SelectObject (hdc, GetStockObject (BLACK_BRUSH)) 
          
          Ellipse (hdc, pt[0].X, pt[0].Y, pt[1].X, pt[1].Y) 
     }
}

func DrawHands (hdc HDC, pst *SYSTEMTIME, fChange bool){
     var pt [3][5]POINT = [3][5]POINT{ 
       {POINT{0, -150}, POINT{100, 0}, POINT{0, 600}, POINT{-100, 0}, POINT{0, -150}},
       {POINT{0,  -200},  POINT{50, 0}, POINT{0, 800},  POINT{-50, 0}, POINT{0, -200}},
       {POINT{0,     0},   POINT{0, 0}, POINT{0,   0},    POINT{0, 0}, POINT{0,  800}}, 
     } 
     var i int32
     var iAngle[3]int32
     //var ptTemp[3][5]POINT 
     
     iAngle[0] = int32((pst.WHour * 30) % 360 + pst.WMinute / 2) 
     iAngle[1] = int32( pst.WMinute  *  6 )
     iAngle[2] = int32( pst.WSecond  *  6 )
     
     //memcpy (ptTemp, pt, sizeof (pt)) ;
     
     
     if fChange {
        i = 0
     } else {
        i = 2
     }
     for  ; i < 3 ; i++      {
          RotatePoint (pt[i][0:], 5, iAngle[i]) 
          
          Polyline (hdc, &pt[i][0], 5) 
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
          KillTimer (hwnd, ID_TIMER) 
    PostQuitMessage(0)
    return 0
}

func OnCreate(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    SetTimer(hwnd, ID_TIMER, 1000, 0)
    var st SYSTEMTIME        
          GetLocalTime (&st) 
          stPrevious = st ;
    return 0
}

func OnTimer(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
var st SYSTEMTIME        
          GetLocalTime (&st) 
                    
          fChange := st.WHour   != stPrevious.WHour ||
                     st.WMinute != stPrevious.WMinute 
          
          hdc := GetDC (hwnd) 
          
          SetIsotropic (hdc, cxClient, cyClient) 
          
          SelectObject (hdc, GetStockObject (WHITE_PEN)) 
          DrawHands (hdc, &stPrevious, fChange) 
          
          SelectObject (hdc, GetStockObject (BLACK_PEN)) 
          DrawHands (hdc, &st, true) 
          
          ReleaseDC (hwnd, hdc) 
          
          stPrevious = st 
          
    return 0
}

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    cxClient = int32(LOWORD(uint32(lParam)))
    cyClient = int32(HIWORD(uint32(lParam)))   
    
    return 0
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)
    
          SetIsotropic (hdc, cxClient, cyClient) 
          DrawClock    (hdc) 
          DrawHands    (hdc, &stPrevious, true) 
          
    EndPaint(hwnd, &ps)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr

var table map[uint32]EventHandler

func main() {
    table = make(map[uint32]EventHandler)
    table[WM_CREATE] = OnCreate
    table[WM_SIZE] = OnSize
    table[WM_TIMER] = OnTimer
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy

    initWindow("Clock", "Analog Clock", syscall.NewCallback(WndProc))
}
