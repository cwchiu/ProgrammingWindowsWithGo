package main

import (
    "fmt"
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "os"
    "regexp"
    //"math"
    "strings"
    "syscall"
    "unsafe"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

const (
    IDS_APPNAME = 1
    IDS_CAPTION = 2
    IDS_ERRMSG  = 3
)

// ref: env_windows.go
func FillListBox(hwndList HWND) {
    for _, v := range os.Environ() {
        part := strings.Split(v, "=")
        fmt.Printf("%s=%s\n", part[0], part[1])
        if len(part[0]) == 0 {
            continue
        }

        SendMessage(hwndList, LB_ADDSTRING, 0, uintptr(unsafe.Pointer(_T(part[0]))))
    }
}

func main() {
    app, _ := NewApp()

    var cxClient, cyClient int32
    var cxChar, cyChar int32
    var xScroll, iPosition int32
    var hResource HGLOBAL
    var pText []byte
    var hScroll HWND
    var rect RECT
    var iNumLines int32
    app.On(WM_CREATE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        hdc := GetDC(hwnd)

        var tm TEXTMETRIC
        GetTextMetrics(hdc, &tm)
        cxChar = tm.TmAveCharWidth
        cyChar = tm.TmHeight + tm.TmExternalLeading
        ReleaseDC(hwnd, hdc)

        xScroll = GetSystemMetrics(SM_CXVSCROLL)

        hScroll = CreateWindowEx(0, _T("scrollbar"), nil,
            WS_CHILD|WS_VISIBLE|SBS_VERT,
            0, 0, 0, 0,
            hwnd, HMENU(1), app.HInstance, nil)

        res := FindResource(HMODULE(app.HInstance), _T("AnnabelLee"),
            _T("TEXT"))
        
        hResource = LoadResource(HMODULE(app.HInstance), res)

        size := SizeofResource(HMODULE(app.HInstance), res)
        //fmt.Println(size)

        //shift := math.Ceil(math.Log2(float64(size)))
        b := (*[1 << 11]byte)(unsafe.Pointer(LockResource(hResource)))
        //fmt.Println(b)
        pText = (*b)[0:size]
        re := regexp.MustCompile("\n")
        iNumLines = int32(len(re.FindAllString(string(pText), -1)))

        SetScrollRange(hScroll, SB_CTL, 0, int32(iNumLines), FALSE)
        SetScrollPos(hScroll, SB_CTL, 0, FALSE)

        return 0
    })

    app.On(WM_SIZE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        cxClient = int32(LOWORD(uint32(lParam)))
        cyClient = int32(HIWORD(uint32(lParam)))

        MoveWindow(hScroll, int32(LOWORD(uint32(lParam)))-xScroll, 0,
            xScroll, cyClient, true)
        SetFocus(hwnd)
        return 0
    })

    app.On(WM_VSCROLL, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        v := LOWORD(uint32(wParam))
        if v == SB_LINEUP {
            iPosition -= 1
        } else if v == SB_BOTTOM {
            iPosition = iNumLines
        } else if v == SB_TOP {
            iPosition = 0
        } else if v == SB_LINEDOWN {
            iPosition += 1
        } else if v == SB_PAGEUP {
            iPosition -= cyClient / cyChar
        } else if v == SB_PAGEDOWN {
            iPosition += cyClient / cyChar
        } else if v == SB_THUMBPOSITION {
            iPosition = int32(HIWORD(uint32(wParam)))
        }
        iPosition = Max(0, Min(iPosition, iNumLines))
        if iPosition != GetScrollPos(hScroll, SB_CTL) {
            SetScrollPos(hScroll, SB_CTL, iPosition, TRUE)
            InvalidateRect(hwnd, nil, true)
        }
        return 0
    })

    app.On(WM_SETFOCUS, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        SetFocus(hScroll)
        return 0
    })
    app.On(WM_PAINT, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        var ps PAINTSTRUCT
        hdc := BeginPaint(hwnd, &ps)
        //pText = (char *) LockResource (hResource) ;

        GetClientRect(hwnd, &rect)
        rect.Left += cxChar
        rect.Top += cyChar * (1 - iPosition)
        DrawTextA(hdc, (*uint16)(unsafe.Pointer((&pText[0]))), int32(len(pText)), &rect, DT_EXTERNALLEADING)

        EndPaint(hwnd, &ps)
        return 0
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        FreeResource(hResource)
        PostQuitMessage(0)
        return 0
    })

    var bufAppName [16]byte
    n1 := LoadStringA(app.HInstance, IDS_APPNAME, uintptr(unsafe.Pointer(&bufAppName[0])), 16)
    szAppName := string(bufAppName[0:n1])

    var bufCaption [64]byte
    n2 := LoadStringA(app.HInstance, IDS_CAPTION, uintptr(unsafe.Pointer(&bufCaption[0])), 64)
    szCaption := string(bufCaption[0:n2])

    app.Icon = LoadIcon(app.HInstance, _T(strings.ToUpper(szAppName)))

    app.Init(szAppName, szCaption)
    app.Run()
}
