package main

import (
    "fmt"
    "syscall"
    "unsafe"

    //  "time"
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
        WS_OVERLAPPEDWINDOW|WS_BORDER|WS_CAPTION|WS_SYSMENU|WS_MAXIMIZEBOX|WS_MINIMIZEBOX|WS_VSCROLL|WS_HSCROLL,
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

    //time.Sleep(10000 * time.Millisecond)
    //fmt.Println(r)
    var msg MSG
    for GetMessage(&msg, HWND_TOP, 0, 0) == TRUE {
        TranslateMessage(&msg)
        DispatchMessage(&msg)
    }
}

type SysMetric struct {
    Index int32
    Label string
    Desc  string
}

var sysmetrics []SysMetric = []SysMetric{
    SysMetric{SM_CXSCREEN, "SM_CXSCREEN", "Screen width in pixels"},
    SysMetric{SM_CYSCREEN, "SM_CYSCREEN", "Screen height in pixels"},
    SysMetric{SM_CXVSCROLL, "SM_CXVSCROLL", "Vertical scroll width"},
    SysMetric{SM_CYHSCROLL, "SM_CYHSCROLL", "Horizontal scroll arrow height"},
    SysMetric{SM_CYCAPTION, "SM_CYCAPTION", "Caption bar height"},
    SysMetric{SM_CXBORDER, "SM_CXBORDER", "Window border width"},
    SysMetric{SM_CYBORDER, "SM_CYBORDER", "Window border height"},
    SysMetric{SM_CXDLGFRAME, "SM_CXDLGFRAME", "Dialog window frame width"},
    SysMetric{SM_CYDLGFRAME, "SM_CYDLGFRAME", "Dialog window frame height"},
    SysMetric{SM_CYVTHUMB, "SM_CYVTHUMB", "Vertical scroll thumb height"},
    SysMetric{SM_CXHTHUMB, "SM_CXHTHUMB", "Horizontal scroll thumb width"},
    SysMetric{SM_CXICON, "SM_CXICON", "Icon width"},
    SysMetric{SM_CYICON, "SM_CYICON", "Icon height"},
    SysMetric{SM_CXCURSOR, "SM_CXCURSOR", "Cursor width"},
    SysMetric{SM_CYCURSOR, "SM_CYCURSOR", "Cursor height"},
    SysMetric{SM_CYMENU, "SM_CYMENU", "Menu bar height"},
    SysMetric{SM_CXFULLSCREEN, "SM_CXFULLSCREEN", "Full screen client area width"},
    SysMetric{SM_CYFULLSCREEN, "SM_CYFULLSCREEN", "Full screen client area height"},
    SysMetric{SM_CYKANJIWINDOW, "SM_CYKANJIWINDOW", "Kanji window height"},
    SysMetric{SM_MOUSEPRESENT, "SM_MOUSEPRESENT", "Mouse present flag"},
    SysMetric{SM_CYVSCROLL, "SM_CYVSCROLL", "Vertical scroll arrow height"},
    SysMetric{SM_CXHSCROLL, "SM_CXHSCROLL", "Horizontal scroll arrow width"},
    SysMetric{SM_DEBUG, "SM_DEBUG", "Debug version flag"},
    SysMetric{SM_SWAPBUTTON, "SM_SWAPBUTTON", "Mouse buttons swapped flag"},
    SysMetric{SM_RESERVED1, "SM_RESERVED1", "Reserved"},
    SysMetric{SM_RESERVED2, "SM_RESERVED2", "Reserved"},
    SysMetric{SM_RESERVED3, "SM_RESERVED3", "Reserved"},
    SysMetric{SM_RESERVED4, "SM_RESERVED4", "Reserved"},
    SysMetric{SM_CXMIN, "SM_CXMIN", "Minimum window width"},
    SysMetric{SM_CYMIN, "SM_CYMIN", "Minimum window height"},
    SysMetric{SM_CXSIZE, "SM_CXSIZE", "Minimize/Maximize icon width"},
    SysMetric{SM_CYSIZE, "SM_CYSIZE", "Minimize/Maximize icon height"},
    SysMetric{SM_CXFRAME, "SM_CXFRAME", "Window frame width"},
    SysMetric{SM_CYFRAME, "SM_CYFRAME", "Window frame height"},
    SysMetric{SM_CXMINTRACK, "SM_CXMINTRACK", "Minimum window tracking width"},
    SysMetric{SM_CYMINTRACK, "SM_CYMINTRACK", "Minimum window tracking height"},
    SysMetric{SM_CXDOUBLECLK, "SM_CXDOUBLECLK", "Double click x tolerance"},
    SysMetric{SM_CYDOUBLECLK, "SM_CYDOUBLECLK", "Double click y tolerance"},
    SysMetric{SM_CXICONSPACING, "SM_CXICONSPACING", "Horizontal icon spacing"},
    SysMetric{SM_CYICONSPACING, "SM_CYICONSPACING", "Vertical icon spacing"},
    SysMetric{SM_MENUDROPALIGNMENT, "SM_MENUDROPALIGNMENT", "Left or right menu drop"},
    SysMetric{SM_PENWINDOWS, "SM_PENWINDOWS", "Pen extensions installed"},
    SysMetric{SM_DBCSENABLED, "SM_DBCSENABLED", "Double-Byte Char Set enabled"},
    SysMetric{SM_CMOUSEBUTTONS, "SM_CMOUSEBUTTONS", "Number of mouse buttons"},
    SysMetric{SM_SHOWSOUNDS, "SM_SHOWSOUNDS", "Present sounds visually"},
}

var (
    cxChar, cyChar                                  int32
    cxCaps                                          int32
    iVscrollPos, iVscrollMax, iVscrollInc int32
    iHscrollMax, iHscrollPos              int32
    iMaxWidth                                       int32
)

func Max(a, b int32) int32 {
    if a > b {
        return a
    } else {
        return b
    }
}

func Min(a, b int32) int32 {
    if a < b {
        return a
    } else {
        return b
    }
}

var si SCROLLINFO

func WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    NUMLINES := int32(len(sysmetrics))
    switch msg {
    case WM_CREATE:
        hdc := GetDC(hwnd)
        var tm TEXTMETRIC
        GetTextMetrics(hdc, &tm)
        cxChar = tm.TmAveCharWidth
        cyChar = tm.TmHeight + tm.TmExternalLeading
        if tm.TmPitchAndFamily&1 > 0 {
            cxCaps = 3
        } else {
            cxCaps = 2
        }
        cxCaps = cxCaps * cxChar / 2
        ReleaseDC(hwnd, hdc)
        iMaxWidth = 40*cxChar + 22*cxCaps

        return 0
    case WM_SIZE:
        cxClient := int32(LOWORD(uint32(lParam)))
        cyClient := int32(HIWORD(uint32(lParam)))
        
        si.CbSize = uint32(unsafe.Sizeof(si))
        si.Mask = SIF_RANGE | SIF_PAGE
        si.Min = 0
        si.Max = NUMLINES - 1
        si.Page = uint32(cyClient / cyChar)
        SetScrollInfo(hwnd, SB_VERT, &si, TRUE)

        si.CbSize = uint32(unsafe.Sizeof(si))
        si.Mask = SIF_RANGE | SIF_PAGE
        si.Min = 0
        si.Max = 2 + iMaxWidth/cxChar
        si.Page = uint32(cxClient / cxChar)
        SetScrollInfo(hwnd, SB_HORZ, &si, TRUE)
        return 0

    case WM_VSCROLL:
        si.CbSize = uint32(unsafe.Sizeof(si))
        si.Mask = SIF_ALL
        GetScrollInfo(hwnd, SB_VERT, &si)
        iVertPos := si.Pos

        v := LOWORD(uint32(wParam))
        if v == SB_LINEUP {
            si.Pos -= 1
        } else if v == SB_LINEDOWN {
            si.Pos += 1
        } else if v == SB_TOP {
            si.Pos = si.Min
        } else if v == SB_BOTTOM {
            si.Pos = si.Max
        } else if v == SB_PAGEUP {
            si.Pos -= int32(si.Page)
        } else if v == SB_PAGEDOWN {
            si.Pos += int32(si.Page)
        } else if v == SB_THUMBPOSITION {
            si.Pos = si.TrackPos
        }
        si.Mask = SIF_POS
        SetScrollInfo(hwnd, SB_VERT, &si, TRUE)
        GetScrollInfo(hwnd, SB_VERT, &si)

        if si.Pos != iVertPos {
            ScrollWindow(hwnd, 0, cyChar*(iVertPos-si.Pos), nil, nil)
            UpdateWindow(hwnd)
        }
        
        return 0
    case WM_HSCROLL:
        si.CbSize = uint32(unsafe.Sizeof(si))
        si.Mask = SIF_ALL
        GetScrollInfo(hwnd, SB_HORZ, &si)
        iHorzPos := si.Pos

        v := LOWORD(uint32(wParam))
    
        if v == SB_LINELEFT {
            si.Pos -= 1
        } else if v == SB_LINERIGHT {
            si.Pos += 1
        } else if v == SB_PAGELEFT {
            si.Pos -= int32(si.Page)
        } else if v == SB_PAGERIGHT {
            si.Pos += int32(si.Page)
        } else if v == SB_THUMBPOSITION {
            si.Pos = si.TrackPos
        }

        si.Mask = SIF_POS
        SetScrollInfo(hwnd, SB_HORZ, &si, TRUE)
        GetScrollInfo(hwnd, SB_HORZ, &si)

        if si.Pos != iHorzPos {
            ScrollWindow(hwnd, cxChar*(iHorzPos-si.Pos), 0, nil, nil)
            UpdateWindow(hwnd)
        }
        return 0
    case WM_DESTROY:
        PostQuitMessage(0)
        return 0
    case WM_PAINT:
        var ps PAINTSTRUCT

        hdc := BeginPaint(hwnd, &ps)
        si.CbSize = uint32(unsafe.Sizeof(si))
        si.Mask = SIF_POS
        GetScrollInfo(hwnd, SB_VERT, &si)
        iVertPos := si.Pos

        GetScrollInfo(hwnd, SB_HORZ, &si)
        iHorzPos := si.Pos

        iPaintBeg := Max(0, iVertPos+ps.RcPaint.Top/cyChar)
        iPaintEnd := Min(NUMLINES-1,
            iVertPos+ps.RcPaint.Bottom/cyChar)

        var i int32
        for i = iPaintBeg; i <= iPaintEnd; i++ {
            x := cxChar * (1 - iHorzPos)
            y := cyChar * (i - iVertPos)

            TextOut(hdc, x, y, syscall.StringToUTF16Ptr(sysmetrics[i].Label), int32(len(sysmetrics[i].Label)))
            TextOut(hdc, x+22*cxCaps, y, syscall.StringToUTF16Ptr(sysmetrics[i].Desc), int32(len(sysmetrics[i].Desc)))
            val := fmt.Sprintf("%5d", GetSystemMetrics(sysmetrics[i].Index))
            SetTextAlign(hdc, TA_RIGHT|TA_TOP)
            TextOut(hdc, x+22*cxCaps+40*cxChar, y, syscall.StringToUTF16Ptr(val), int32(len(val)))
            SetTextAlign(hdc, TA_LEFT|TA_TOP)
        }
        EndPaint(hwnd, &ps)
        return 0
    case WM_KEYDOWN:
        switch (uint32(wParam)){
          case VK_HOME:
               SendMessage (hwnd, WM_VSCROLL, SB_TOP, 0) 
               break 
               
          case VK_END:
               SendMessage (hwnd, WM_VSCROLL, SB_BOTTOM, 0) 
               break 
               
          case VK_PRIOR:
               SendMessage (hwnd, WM_VSCROLL, SB_PAGEUP, 0) 
               break 
               
          case VK_NEXT:
               SendMessage (hwnd, WM_VSCROLL, SB_PAGEDOWN, 0) 
               break 
               
          case VK_UP:
               SendMessage (hwnd, WM_VSCROLL, SB_LINEUP, 0) 
               break 
               
          case VK_DOWN:
               SendMessage (hwnd, WM_VSCROLL, SB_LINEDOWN, 0) 
               break 
               
          case VK_LEFT:
               SendMessage (hwnd, WM_HSCROLL, SB_PAGEUP, 0) 
               break 
               
          case VK_RIGHT:
               SendMessage (hwnd, WM_HSCROLL, SB_PAGEDOWN, 0) 
               break 
          }
          return 0 
    }
    return DefWindowProc(hwnd, msg, wParam, lParam)
}

func main() {
    initWindow("System1", "Get System Metrics No. 3", syscall.NewCallback(WndProc))
}
