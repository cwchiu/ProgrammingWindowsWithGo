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
    cxChar, cyChar, cxCaps int32
) 
                    
          
          
          
func OnCreate(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    var tm TEXTMETRIC
    hdc := GetDC (hwnd) 
    SelectObject (hdc, GetStockObject (SYSTEM_FIXED_FONT)) 

    GetTextMetrics (hdc, &tm) 
    
    cxChar = tm.TmAveCharWidth 
    cyChar = tm.TmHeight + tm.TmExternalLeading 
    if tm.TmPitchAndFamily & 1 > 0 {
        cxCaps = 3
    }else{  
        cxCaps = 2
    }
    cxCaps = cxCaps * cxChar / 2 

    ReleaseDC (hwnd, hdc) 
    return 0
}

func OnDestroy(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    PostQuitMessage(0)
    return 0
}
    
type DeviceCaps struct{
     Index int32
     Label string
     Desc string
}

var devcaps []DeviceCaps = []DeviceCaps{
     DeviceCaps{HORZSIZE,      "HORZSIZE",     "Width in millimeters:"},
     DeviceCaps{VERTSIZE,      "VERTSIZE",     "Height in millimeters:"},
     DeviceCaps{HORZRES,       "HORZRES",      "Width in pixels:"},
     DeviceCaps{VERTRES,       "VERTRES",      "Height in raster lines:"},
     DeviceCaps{BITSPIXEL,     "BITSPIXEL",    "Color bits per pixel:"},
     DeviceCaps{PLANES,        "PLANES",       "Number of color planes:"},
     DeviceCaps{NUMBRUSHES,    "NUMBRUSHES",   "Number of device brushes:"},
     DeviceCaps{NUMPENS,       "NUMPENS",      "Number of device pens:"},
     DeviceCaps{NUMMARKERS,    "NUMMARKERS",   "Number of device markers:"},
     DeviceCaps{NUMFONTS,      "NUMFONTS",     "Number of device fonts:"},
     DeviceCaps{NUMCOLORS,     "NUMCOLORS",    "Number of device colors:"},
     DeviceCaps{PDEVICESIZE,   "PDEVICESIZE",  "Size of device structure:"},
     DeviceCaps{ASPECTX,       "ASPECTX",      "Relative width of pixel:"},
     DeviceCaps{ASPECTY,       "ASPECTY",      "Relative height of pixel:"},
     DeviceCaps{ASPECTXY,      "ASPECTXY",     "Relative diagonal of pixel:"},
     DeviceCaps{LOGPIXELSX,    "LOGPIXELSX",   "Horizontal dots per inch:"},
     DeviceCaps{LOGPIXELSY,    "LOGPIXELSY",   "Vertical dots per inch:"},
     DeviceCaps{SIZEPALETTE,   "SIZEPALETTE",  "Number of palette entries:"},
     DeviceCaps{NUMRESERVED,   "NUMRESERVED",  "Reserved palette entries:"},
     DeviceCaps{COLORRES,      "COLORRES",     "Actual color resolution:"},
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)
    SelectObject (hdc, GetStockObject (SYSTEM_FIXED_FONT)) 
    var NUMLINES int32 = int32(len(devcaps))
    var i int32
    var szBuffer string
    for i = 0 ; i <NUMLINES  ; i++ {
       TextOut (hdc, 0, cyChar * i, syscall.StringToUTF16Ptr(devcaps[i].Label), int32(len(devcaps[i].Label)))
       
       TextOut (hdc, 14 * cxCaps, cyChar * i, syscall.StringToUTF16Ptr(devcaps[i].Desc), int32(len(devcaps[i].Desc))) 
       
       SetTextAlign (hdc, TA_RIGHT | TA_TOP) 
       szBuffer = fmt.Sprintf("%5d", GetDeviceCaps (hdc, devcaps[i].Index))
       TextOut (hdc, 14 * cxCaps + 35 * cxChar, cyChar * i, syscall.StringToUTF16Ptr(szBuffer), int32(len(szBuffer)))       
       
       SetTextAlign (hdc, TA_LEFT | TA_TOP) ;
    }
    
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
    
    initWindow("DevCaps1", "Device Capabilities", syscall.NewCallback(WndProc))
}
