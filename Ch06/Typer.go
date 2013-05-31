package main

import (
    "fmt"
    //"math"
    . "github.com/cwchiu/go-winapi"
    "syscall"
    "unsafe"
)

var (
    cxChar, cyChar     int32
    cxClient, cyClient int32
    cxBuffer, cyBuffer int32
    xCaret, yCaret int32
    pBuffer []byte
    dwCharSet          uint32 = DEFAULT_CHARSET
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func Min(a, b int32) int32 {
    if a < b {
        return a
    } else {
        return b
    }
}

func Max(a, b int32) int32 {
    if a < b {
        return b
    } else {
        return a
    }
}

func XY2Pos(x,y int32) int32{
    return y * cxBuffer + x
}

func setBuffer(x, y int32, v byte) {
    pBuffer[ XY2Pos(x,y) ] = v
}

func getBuffer(x, y int32) byte {
    return pBuffer[ XY2Pos(x,y) ]
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

func OnCreate(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var tm TEXTMETRIC
    hdc := GetDC(hwnd)
    SelectObject(hdc, HGDIOBJ(CreateFont(0, 0, 0, 0, 0, 0, 0, 0,
        dwCharSet, 0, 0, 0, FIXED_PITCH, nil)))

    GetTextMetrics(hdc, &tm)
    cxChar = tm.TmAveCharWidth
    cyChar = tm.TmHeight

    DeleteObject(SelectObject(hdc, GetStockObject(SYSTEM_FONT)))
    ReleaseDC(hwnd, hdc)
    return OnSize(hwnd, msg, wParam, lParam)
}

func OnSize(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    // obtain window size in pixels

    if msg == WM_SIZE {
        cxClient = int32(LOWORD(uint32(lParam)))
        cyClient = int32(HIWORD(uint32(lParam)))
    }
    // calculate window size in characters

    cxBuffer = Max(1, cxClient/cxChar)
    cyBuffer = Max(1, cyClient/cyChar)

    var y, x int32    
    // allocate memory for buffer and clear it
    pBuffer = make([]byte, cxBuffer * cyBuffer)    
    for y = 0; y < cyBuffer; y++ {
        for x = 0; x < cxBuffer; x++ {
            setBuffer(x, y, byte(0x20) )
        }
    }

    // set caret to upper left corner

    xCaret = 0
    yCaret = 0

    if hwnd == GetFocus() {
        SetCaretPos(xCaret*cxChar, yCaret*cyChar)
    }

    InvalidateRect(hwnd, nil, true)
    return 0
}

func OnDestroy(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    PostQuitMessage(0)
    return 0
}

func OnChar(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var i int32
    var x,y int32
    
    for i = 0; i < int32(LOWORD(uint32(lParam))); i++ {
        switch byte(wParam) {
        case '\b': // backspace
            if xCaret > 0 {
                xCaret--
                SendMessage(hwnd, WM_KEYDOWN, VK_DELETE, 1)
            }
            break

        case '\t': // tab
            for {
                SendMessage(hwnd, WM_CHAR, ' ', 1)
                
                if(xCaret%8 == 0){
                    break
                }
            }
            
            break

        case '\n': // line feed
            yCaret++
            if yCaret == cyBuffer {
                yCaret = 0
            }
            break

        case '\r': // carriage return
            xCaret = 0
            yCaret++
            if yCaret == cyBuffer {
                yCaret = 0
            }
            break

        case '\x1B': // escape            
            for y = 0; y < cyBuffer; y++ {
                for x = 0; x < cxBuffer; x++ {
                    setBuffer(x, y, byte(' '))
                }
            }
            xCaret = 0
            yCaret = 0

            InvalidateRect(hwnd, nil, false)
            break

        default: // character codes
            setBuffer(xCaret, yCaret, byte(wParam))

            HideCaret(hwnd)
            hdc := GetDC(hwnd)

            SelectObject(hdc, HGDIOBJ(CreateFont(0, 0, 0, 0, 0, 0, 0, 0,
                dwCharSet, 0, 0, 0, FIXED_PITCH, nil)))

            TextOut(hdc, xCaret*cxChar, yCaret*cyChar,
                _T(string(getBuffer(xCaret, yCaret))), 1)

            DeleteObject(
                SelectObject(hdc, HGDIOBJ(GetStockObject(SYSTEM_FONT))))
            ReleaseDC(hwnd, hdc)
            ShowCaret(hwnd)
            xCaret++
            if xCaret == cxBuffer {
                xCaret = 0
                yCaret++
                if yCaret == cyBuffer {
                    yCaret = 0
                }
            }
            break
        }
    }

    SetCaretPos(xCaret*cxChar, yCaret*cyChar)
    return 0
}

func OnKeyDown(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    switch (wParam) {
    case VK_HOME:
        xCaret = 0
        break

    case VK_END:
        xCaret = cxBuffer - 1
        break

    case VK_PRIOR:
        yCaret = 0
        break

    case VK_NEXT:
        yCaret = cyBuffer - 1
        break

    case VK_LEFT:
        xCaret = Max(xCaret-1, 0)
        break

    case VK_RIGHT:
        xCaret = Min(xCaret+1, cxBuffer-1)
        break

    case VK_UP:
        yCaret = Max(yCaret-1, 0)
        break

    case VK_DOWN:
        yCaret = Min(yCaret+1, cyBuffer-1)
        break

    case VK_DELETE:
        var x int32
        for x = xCaret; x < cxBuffer-1; x++ {
            setBuffer(x, yCaret, getBuffer(x+1, yCaret))
        }

        setBuffer(cxBuffer-1, yCaret, byte(' '))

        HideCaret(hwnd)
        hdc := GetDC(hwnd)

        SelectObject(hdc, HGDIOBJ(CreateFont(0, 0, 0, 0, 0, 0, 0, 0,
            dwCharSet, 0, 0, 0, FIXED_PITCH, nil)))

        begin :=  yCaret * cxBuffer + xCaret          
        size := cxBuffer-xCaret
        //fmt.Println(string(pBuffer[begin:begin+size])
        TextOut(hdc, xCaret*cxChar, yCaret*cyChar,
            _T(string(pBuffer[begin:begin+size])),
            size)

        DeleteObject(SelectObject(hdc, GetStockObject(SYSTEM_FONT)))
        ReleaseDC(hwnd, hdc)
        ShowCaret(hwnd)
        break
    }
    SetCaretPos(xCaret*cxChar, yCaret*cyChar)

    return 0
}

func OnSetFocus(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    // create and show the caret

    CreateCaret(hwnd, HBITMAP(0), cxChar, cyChar)
    SetCaretPos(xCaret*cxChar, yCaret*cyChar)
    ShowCaret(hwnd)
    return 0

}

func OnKillFocus(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    HideCaret(hwnd)
    DestroyCaret()
    return 0
}
func OnInputLangChange(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    dwCharSet = uint32(wParam)
    return OnCreate(hwnd, msg, wParam, lParam)
}

func OnPaint(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    var ps PAINTSTRUCT
    hdc := BeginPaint(hwnd, &ps)

    SelectObject(hdc, HGDIOBJ(CreateFont(0, 0, 0, 0, 0, 0, 0, 0,
        dwCharSet, 0, 0, 0, FIXED_PITCH, nil)))
    var y int32
    for y = 0; y < cyBuffer; y++ {
        begin := XY2Pos(0,y)
        fmt.Println(pBuffer[begin:begin+cxBuffer])
        TextOut(hdc, 0, y*cyChar, _T(string(pBuffer[begin:begin+cxBuffer])), cxBuffer)
    }
    DeleteObject(SelectObject(hdc, GetStockObject(SYSTEM_FONT)))

    EndPaint(hwnd, &ps)
    return 0
}

type EventHandler func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr

var table map[uint32]EventHandler

func main() {
    table = make(map[uint32]EventHandler)
    table[WM_INPUTLANGCHANGE] = OnInputLangChange
    table[WM_CREATE] = OnCreate
    table[WM_SIZE] = OnSize
    table[WM_SETFOCUS] = OnSetFocus
    table[WM_KILLFOCUS] = OnKillFocus
    table[WM_KEYDOWN] = OnKeyDown
    table[WM_CHAR] = OnChar
    table[WM_PAINT] = OnPaint
    table[WM_DESTROY] = OnDestroy

    initWindow("Typer", "Typing Program", syscall.NewCallback(WndProc))
}
