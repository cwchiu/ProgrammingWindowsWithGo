package main

import (
    //"fmt"
    //"math"
    . "github.com/cwchiu/go-winapi"
    "syscall"
    "unsafe"
)

const (
    ID_TIMER = 1
)

var (
    f24Hour, fSuppress bool
    hBrushRed HBRUSH 
    cxClient, cyClient int32
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func DisplayDigit (hdc HDC, iNumber uint16){

     var fSevenSegment [10][7]BOOL = [10][7]BOOL{
                         {TRUE,  TRUE,  TRUE, FALSE,  TRUE,  TRUE,  TRUE},   // 0
                        {FALSE, FALSE,  TRUE, FALSE, FALSE,  TRUE, FALSE},   // 1
                        { TRUE, FALSE,  TRUE,  TRUE,  TRUE, FALSE,  TRUE},   // 2
                        { TRUE, FALSE,  TRUE,  TRUE, FALSE,  TRUE,  TRUE},   // 3
                        {FALSE,  TRUE,  TRUE,  TRUE, FALSE,  TRUE, FALSE},   // 4
                        { TRUE,  TRUE, FALSE,  TRUE, FALSE,  TRUE,  TRUE},   // 5
                        { TRUE,  TRUE, FALSE,  TRUE,  TRUE,  TRUE,  TRUE},   // 6
                        { TRUE, FALSE,  TRUE, FALSE, FALSE,  TRUE, FALSE},   // 7
                        { TRUE,  TRUE,  TRUE,  TRUE,  TRUE,  TRUE,  TRUE},   // 8
                        { TRUE,  TRUE,  TRUE,  TRUE, FALSE,  TRUE,  TRUE}  } // 9
         
                         
                         
                         
                         
                         
                         
                         
                         
                         
     var ptSegment [7][6]POINT = [7][6]POINT{
                         {POINT{ 7,  6},  POINT{11,  2},  POINT{31,  2},  POINT{35,  6},  POINT{31, 10},  POINT{11, 10}},
                         {POINT{ 6,  7},  POINT{10, 11},  POINT{10, 31},  POINT{ 6, 35},  POINT{ 2, 31},  POINT{ 2, 11}},
                         {POINT{36,  7},  POINT{40, 11},  POINT{40, 31},  POINT{36, 35},  POINT{32, 31},  POINT{32, 11}},
                         {POINT{ 7, 36},  POINT{11, 32},  POINT{31, 32},  POINT{35, 36},  POINT{31, 40},  POINT{11, 40}},
                         {POINT{ 6, 37},  POINT{10, 41},  POINT{10, 61},  POINT{ 6, 65},  POINT{ 2, 61},  POINT{ 2, 41}},
                         {POINT{36, 37},  POINT{40, 41},  POINT{40, 61},  POINT{36, 65},  POINT{32, 61},  POINT{32, 41}},
                         {POINT{ 7, 66},  POINT{11, 62},  POINT{31, 62},  POINT{35, 66},  POINT{31, 70},  POINT{11, 70}} } ;
     var iSeg  int32
     
     for iSeg = 0 ; iSeg < 7 ; iSeg++{
          if (fSevenSegment [iNumber][iSeg] == TRUE) {
               Polygon (hdc, &ptSegment [iSeg][0], 6) 
               }
               }
}

func DisplayTwoDigits (hdc HDC, iNumber uint16, fSuppress bool){
     if (!fSuppress || (iNumber / 10 != 0)){
          DisplayDigit (hdc, iNumber / 10) 
    }
    
     OffsetWindowOrgEx (hdc, -42, 0, nil) 
     DisplayDigit (hdc, iNumber % 10) 
     OffsetWindowOrgEx (hdc, -42, 0, nil) 
}

func DisplayColon (hdc HDC){
     var ptColon [2][4]POINT = [2][4]POINT{ 
        {POINT{2,  21}, POINT{6,  17}, POINT{10, 21}, POINT{6,  25}},
        {POINT{2,  51}, POINT{6,  47}, POINT{10, 51}, POINT{6,  55}}} 

     Polygon (hdc, &ptColon [0][0], 4) 
     Polygon (hdc, &ptColon [1][0], 4) 

     OffsetWindowOrgEx (hdc, -12, 0, nil) 
}

func DisplayTime (hdc HDC,f24Hour, fSuppress bool){
     var st SYSTEMTIME

     GetLocalTime (&st) 

     if (f24Hour){
          DisplayTwoDigits (hdc, st.WHour, fSuppress) 
     }else{
        var h uint16
            if st.WHour % 12 == 0 {
                h = 12
            } else {
                h = st.WHour
            }
          DisplayTwoDigits (hdc, h, fSuppress) 
     }

     DisplayColon (hdc) 
     DisplayTwoDigits (hdc, st.WMinute, false) 
     DisplayColon (hdc) 
     DisplayTwoDigits (hdc, st.WSecond, false) 
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
          KillTimer (hwnd, ID_TIMER) 
          DeleteObject (HGDIOBJ(hBrushRed))
    PostQuitMessage(0)
    return 0
}

func OnCreate(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    hBrushRed = CreateSolidBrush (RGB (255, 0, 0))
    SetTimer(hwnd, ID_TIMER, 1000, 0)
    return OnSettingChange(hwnd, msg, wParam, lParam)
}

func OnTimer(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    InvalidateRect(hwnd, nil, true)
    return 0
}

func OnSettingChange(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var szBuffer [2]byte
    
          GetLocaleInfo (LOCALE_USER_DEFAULT, LOCALE_ITIME, (*uint16)(unsafe.Pointer(&szBuffer[0])), 2) 
          f24Hour = (szBuffer[0] == '1') 

          GetLocaleInfo (LOCALE_USER_DEFAULT, LOCALE_ITLZERO, (*uint16)(unsafe.Pointer(&szBuffer[0])), 2) 
          fSuppress = (szBuffer[0] == '0') 

          InvalidateRect (hwnd, nil, true) 
          return 0 ;
}

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    cxClient = int32(LOWORD(uint32(lParam)))
    cyClient = int32(HIWORD(uint32(lParam)))   
    
    return 0
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)

          SetMapMode (hdc, MM_ISOTROPIC) 
          SetWindowExtEx (hdc, 276, 72, nil) 
          SetViewportExtEx (hdc, cxClient, cyClient, nil) 

          SetWindowOrgEx (hdc, 138, 36, nil) 
          SetViewportOrgEx (hdc, cxClient / 2, cyClient / 2, nil) 

          SelectObject (hdc, GetStockObject (NULL_PEN)) 
          SelectObject (hdc, HGDIOBJ(hBrushRed) )

          DisplayTime (hdc, f24Hour, fSuppress) 
          
    EndPaint(hwnd, &ps)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr

var table map[uint32]EventHandler

func main() {
    table = make(map[uint32]EventHandler)
    table[WM_CREATE] = OnCreate
    table[WM_SETTINGCHANGE] = OnSettingChange
    table[WM_SIZE] = OnSize
    table[WM_TIMER] = OnTimer
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy

    initWindow("Beeper1", "Beeper1 Timer Demo", syscall.NewCallback(WndProc))
}
