package main

import (
    "fmt"
    //"math"
    "syscall"
    "unsafe"
    . "github.com/cwchiu/go-winapi"
)

var _T func (s string) *uint16 = syscall.StringToUTF16Ptr

func initWindow(appName string, title string, wndproc uintptr) {
    

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

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    cxClient = int32(LOWORD(uint32(lParam)))
    cyClient = int32(HIWORD(uint32(lParam)))           

    rectScroll.Left   = 0 
    rectScroll.Right  = cxClient 
    rectScroll.Top    = cyChar 
    rectScroll.Bottom = cyChar * (cyClient / cyChar) 

    InvalidateRect (hwnd, nil, true) 
          
    return 0
}

func OnDestroy(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    PostQuitMessage(0)
    return 0
}

var (
     cxClientMax, cyClientMax, cxClient, cyClient, cxChar, cyChar int32
     cLinesMax, cLines int32
     pmsg []MSG
     rectScroll RECT
)


var szMessage []string = []string{ 
    "WM_KEYDOWN",    "WM_KEYUP", 
    "WM_CHAR",       "WM_DEADCHAR", 
    "WM_SYSKEYDOWN", "WM_SYSKEYUP", 
    "WM_SYSCHAR",    "WM_SYSDEADCHAR",
}

func Min(a, b int32) int32{
    if a<b {
        return a
    } else {
        return b
    }
}

func OnDisplayChange(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    cxClientMax = GetSystemMetrics (SM_CXMAXIMIZED) 
    cyClientMax = GetSystemMetrics (SM_CYMAXIMIZED) 

      // Get character size for fixed-pitch font

    hdc := GetDC (hwnd)

    SelectObject (hdc, GetStockObject (SYSTEM_FIXED_FONT))
    
    var tm TEXTMETRIC
    GetTextMetrics (hdc, &tm)
    cxChar = tm.TmAveCharWidth
    cyChar = tm.TmHeight

    ReleaseDC (hwnd, hdc) 
    
    cLinesMax = cyClientMax / cyChar 
    pmsg = make([]MSG, cLinesMax)
    cLines = 0 
    
    rectScroll.Left   = 0 
    rectScroll.Right  = cxClient 
    rectScroll.Top    = cyChar 
    rectScroll.Bottom = cyChar * (cyClient / cyChar) 

    InvalidateRect (hwnd, nil, true)
    
    return 0
}

func OnKeyHandler(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    // Rearrange storage array
    var i int32
    for i = cLinesMax - 1 ; i > 0 ; i--{
       pmsg[i] = pmsg[i - 1] ;
    }
    
    // Store new message
    pmsg[0].HWnd = hwnd 
    pmsg[0].Message = msg 
    pmsg[0].WParam = wParam 
    pmsg[0].LParam = lParam 

    cLines = Min (cLines + 1, cLinesMax) 

    // Scroll up the display
    ScrollWindow (hwnd, 0, -cyChar, &rectScroll, &rectScroll) 
    
    return DefWindowProc (hwnd, msg, wParam, lParam)
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)    
    
    SelectObject (hdc, GetStockObject (SYSTEM_FIXED_FONT)) 
    SetBkMode (hdc, TRANSPARENT) 
    
    szTop := "Message        Key       Char     Repeat Scan Ext ALT Prev Tran"                                    
    szUnd := "_______        ___       ____     ______ ____ ___ ___ ____ ____"    

    TextOut (hdc, 0, 0, _T(szTop), int32(len (szTop))) 
    TextOut (hdc, 0, 0, _T(szUnd), int32(len (szUnd))) 

    
    szFormat := []string{           
        "%-13s %3d %-15s%s%6d %4d %3s %3s %4s %4s",
        "%-13s            0x%04X%1s%s %6d %4d %3s %3s %4s %4s",
    }
               
    szYes  := "Yes"
    szNo   := "No"
    szDown := "Down"
    szUp   := "Up"

    var szBuffer, szKeyName string
    var iType int32
    var char string
    var szExt, szAlt, szPrev, szTran string
    var i int32
    var bKeyName []byte = make([]byte, 32)
    //fmt.Printf("%d, %d, %d, %d\n", cyClient, cyChar, cLines, cyClient / cyChar - 1)
    for i = 0 ; i < Min (cLines, cyClient / cyChar - 1) ; i++{
       if( pmsg[i].Message == WM_CHAR ||
           pmsg[i].Message == WM_SYSCHAR ||
           pmsg[i].Message == WM_DEADCHAR ||
           pmsg[i].Message == WM_SYSDEADCHAR ){
           iType = 1
           szKeyName = " "           
           char = fmt.Sprintf("%c", uint32(pmsg[i].WParam))
       }else{
           iType = 0
           r := GetKeyNameTextA(pmsg[i].LParam, (*uint16)(unsafe.Pointer(&bKeyName[0])), int32(len(bKeyName)))            
           fmt.Println(r)
           if r > 0 {                
                szKeyName = string(bKeyName[0:r])
                fmt.Println(szKeyName)
           }
           
           char = " "
       }

       if 0x01000000 & pmsg[i].LParam > 0 {
            szExt = szYes
       } else {
            szExt = szNo
       } 
       
       if 0x20000000 & pmsg[i].LParam > 0 {
            szAlt = szYes
       } else {
            szAlt = szNo
       } 
       
       if 0x40000000 & pmsg[i].LParam > 0 {
            szPrev = szDown
       } else {
            szPrev = szUp
       } 
       
       if 0x80000000 & pmsg[i].LParam > 0 {
            szTran = szUp
       } else {
            szTran = szDown
       } 
       
       szBuffer = fmt.Sprintf(szFormat [iType],
                     szMessage [pmsg[i].Message - WM_KEYFIRST],                   
                     pmsg[i].WParam,
                     szKeyName,
                     char,
                     LOWORD (uint32(pmsg[i].LParam)),
                     HIWORD (uint32(pmsg[i].LParam)) & 0xFF,
                     szExt,
                     szAlt,
                     szPrev,
                     szTran)
                     
        fmt.Println(szBuffer)
       TextOut (hdc, 0, (cyClient / cyChar - 1 - i) * cyChar, _T(szBuffer),
                int32(len(szBuffer))) 
    }
          
    EndPaint(hwnd, &ps)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr
var table map[uint32]EventHandler
func main() {
    table = make(map[uint32]EventHandler)
    table[WM_CREATE] = OnDisplayChange
    table[WM_DISPLAYCHANGE] = OnDisplayChange
    table[WM_SIZE] = OnSize
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy
    table[WM_KEYDOWN] = OnKeyHandler
    table[WM_KEYUP] = OnKeyHandler
    table[WM_CHAR] = OnKeyHandler
    table[WM_DEADCHAR] = OnKeyHandler
    table[WM_SYSKEYDOWN] = OnKeyHandler
    table[WM_SYSKEYUP] = OnKeyHandler
    table[WM_SYSCHAR] = OnKeyHandler
    table[WM_SYSDEADCHAR] = OnKeyHandler
    
    initWindow("KeyView1", "Keyboard Message Viewer #1", syscall.NewCallback(WndProc))
}
