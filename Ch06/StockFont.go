package main

import (
    "fmt"
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
    iFont int32
) 

type StockFont struct{
    StockFontID int32
    StockFontName string
}

var stockfont []StockFont = []StockFont{ 
    StockFont{OEM_FIXED_FONT,      "OEM_FIXED_FONT"},
    StockFont{ANSI_FIXED_FONT,     "ANSI_FIXED_FONT"},    
    StockFont{ANSI_VAR_FONT,       "ANSI_VAR_FONT"},
    StockFont{SYSTEM_FONT,         "SYSTEM_FONT"},
    StockFont{DEVICE_DEFAULT_FONT, "DEVICE_DEFAULT_FONT"},
    StockFont{SYSTEM_FIXED_FONT,   "SYSTEM_FIXED_FONT"},
    StockFont{DEFAULT_GUI_FONT,    "DEFAULT_GUI_FONT"}, 
} 
                      
func OnCreate(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    SetScrollRange (hwnd, SB_VERT, 0, int32(len(stockfont) - 1), TRUE)
    return 0
}

func OnDestroy(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    PostQuitMessage(0)
    return 0
}
       
func OnDisplayChange(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    InvalidateRect (hwnd, nil, true)
    return 0
}

func Max(a, b int32) int32{
    if a>b {
        return a
    }
    
    return b
}

func Min(a, b int32)int32{
    if a>b {
        return b
    }
    
    return a
}

func OnVScroll(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    cFonts := int32(len(stockfont))
    switch (LOWORD(uint32(wParam))){
    case SB_TOP:            
        iFont = 0
        break 
    case SB_BOTTOM:
        iFont = cFonts - 1
        break
    case SB_LINEUP:
    case SB_PAGEUP:         
        iFont -= 1
        break 
    case SB_LINEDOWN:
    case SB_PAGEDOWN:       
        iFont += 1
        break
    case SB_THUMBPOSITION:  
        iFont = int32(HIWORD (uint32(wParam)) )
        break 
    }
    iFont = Max (0, Min (cFonts - 1, iFont)) 
    SetScrollPos (hwnd, SB_VERT, iFont, TRUE) 
    InvalidateRect (hwnd, nil, true) 
    return 0 ;
}

func OnKeyDown(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    switch (uint32(wParam)){
    case VK_HOME: 
        SendMessage (hwnd, WM_VSCROLL, SB_TOP, 0) 
        break 
    case VK_END:  
        SendMessage (hwnd, WM_VSCROLL, SB_BOTTOM, 0) 
        break 
    case VK_PRIOR:
    case VK_LEFT:
    case VK_UP:   
        SendMessage (hwnd, WM_VSCROLL, SB_LINEUP, 0) 
        break 
    case VK_NEXT: 
    case VK_RIGHT:
    case VK_DOWN: 
        SendMessage (hwnd, WM_VSCROLL, SB_PAGEDOWN, 0) 
        break 
    }
    return 0 ;
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)
    SelectObject (hdc, GetStockObject (stockfont[iFont].StockFontID)) 
    var buffer []byte = make([]byte, LF_FACESIZE)
    r := GetTextFaceA(hdc, LF_FACESIZE, (*uint16)(unsafe.Pointer(&buffer[0])) ) 
    szFaceName := string(buffer[0:r])
    //fmt.Println(szFaceName)
    var tm TEXTMETRIC
    GetTextMetrics (hdc, &tm) 
    cxGrid := Max (3 * tm.TmAveCharWidth, 2 * tm.TmMaxCharWidth) 
    cyGrid := tm.TmHeight + 3 ;
    var szBuffer string
    
    szBuffer = fmt.Sprintf(" %s: Face Name = %s, CharSet = %d", stockfont[iFont].StockFontName, 
                 szFaceName, tm.TmCharSet)
    TextOut (hdc, 0, 0, _T(szBuffer), int32(len(szBuffer)))

    SetTextAlign (hdc, TA_TOP | TA_CENTER) 

    // vertical and horizontal lines
    var i int32
    for i = 0 ; i < 17 ; i++ {
       MoveToEx (hdc, (i + 2) * cxGrid,  2 * cyGrid, nil) 
       LineTo   (hdc, (i + 2) * cxGrid, 19 * cyGrid) 

       MoveToEx (hdc,      cxGrid, (i + 3) * cyGrid, nil) 
       LineTo   (hdc, 18 * cxGrid, (i + 3) * cyGrid) 
    }
    
    // vertical and horizontal headings
    for i = 0 ; i < 16 ; i++ {
       szBuffer = fmt.Sprintf("%X-", i)
       TextOut (hdc, (2 * i + 5) * cxGrid / 2, 2 * cyGrid + 2, _T(szBuffer), int32(len(szBuffer)))
       
       szBuffer = fmt.Sprintf("-%X", i)
       TextOut (hdc, 3 * cxGrid / 2, (i + 3) * cyGrid + 2, _T(szBuffer), int32(len(szBuffer)))
    }
    
    // characters
    
    var x,y int32
    for y = 0 ; y < 16 ; y++{
        for x = 0 ; x < 16 ; x++ {    
            szBuffer = fmt.Sprintf("%c",byte(16 * x + y))            
            if szBuffer[0] == 0 {
                szBuffer = " "
            }
            
            TextOut (hdc, (2 * x + 5) * cxGrid / 2, 
                          (y + 3) * cyGrid + 2, 
                          _T(szBuffer), 
                          int32(len(szBuffer)))
            
        }
    }
    
    EndPaint(hwnd, &ps)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr
var table map[uint32]EventHandler
func main() {    
    table = make(map[uint32]EventHandler)
    table[WM_CREATE] = OnCreate
    table[WM_DISPLAYCHANGE] = OnDisplayChange
    table[WM_VSCROLL] = OnVScroll
    table[WM_KEYDOWN] = OnKeyDown
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy
    
    initWindow("StokFont", "Stock Fonts", syscall.NewCallback(WndProc))
}
