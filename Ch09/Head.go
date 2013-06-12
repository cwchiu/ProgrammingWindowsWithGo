package main

import (
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "os"
    "regexp"
    "syscall"
    "unsafe"
)

const (
    ID_LIST = 1
    ID_TEXT = 2
    MAXREAD = 8192
    DIRATTR = (DDL_READWRITE | DDL_READONLY | DDL_HIDDEN | DDL_SYSTEM | DDL_DIRECTORY | DDL_ARCHIVE | DDL_DRIVES)
    DTFLAGS = (DT_WORDBREAK | DT_EXPANDTABS | DT_NOCLIP | DT_NOPREFIX)
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func GetCurDir() string {
    buf := make([]uint16, MAX_PATH+1)
    GetCurrentDirectory(MAX_PATH+1, uintptr(unsafe.Pointer(&buf[0])))
    szBuffer := syscall.UTF16ToString(buf)
    return szBuffer
}

func main() {

    app, _ := NewApp()
    app.BackgroundBrush = HBRUSH(COLOR_BTNFACE + 1)
    var hwndList, hwndText HWND
    var cxChar, cyChar int32
    var OldList int32
    var rect RECT
    var bValidFile bool
    var szFile string
    app.On(WM_CREATE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        dlg_base_uint := uint32(GetDialogBaseUnits())
        cxChar = int32(LOWORD(dlg_base_uint))
        cyChar = int32(HIWORD(dlg_base_uint))

        rect.Left = 20 * cxChar
        rect.Top = 3 * cyChar

        hwndList = CreateWindowEx(0, _T("listbox"), nil,
            WS_CHILDWINDOW|WS_VISIBLE|LBS_STANDARD,
            cxChar, cyChar*3,
            cxChar*13+GetSystemMetrics(SM_CXVSCROLL),
            cyChar*10,
            hwnd, HMENU(ID_LIST),
            app.HInstance,
            nil)

        szBuffer := GetCurDir()

        hwndText = CreateWindowEx(0, _T("static"), _T(szBuffer),
            WS_CHILDWINDOW|WS_VISIBLE|SS_LEFT,
            cxChar, cyChar, cxChar*MAX_PATH, cyChar,
            hwnd, HMENU(ID_TEXT),
            app.HInstance,
            nil)

        OldList = SetWindowLong(hwndList, GWL_WNDPROC, int32(syscall.NewCallback(func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
            if msg == WM_KEYDOWN && wParam == VK_RETURN {
                //SendMessage (GetParent (hwnd), WM_COMMAND,
                //             MAKELONG (1, LBN_DBLCLK), (LPARAM) hwnd)
            }

            return CallWindowProc(uintptr(OldList), hwnd, msg, wParam, lParam)
        })))

        SendMessage(hwndList, LB_DIR, DIRATTR, uintptr(unsafe.Pointer(_T("*.*"))))

        return 0
    })

    app.On(WM_SIZE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        rect.Right = int32(LOWORD(uint32(lParam)))
        rect.Bottom = int32(HIWORD(uint32(lParam)))
        return 0
    })

    app.On(WM_COMMAND, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        if LOWORD(uint32(wParam)) == ID_LIST &&
            HIWORD(uint32(wParam)) == LBN_DBLCLK {
            i := int32(SendMessage(hwndList, LB_GETCURSEL, 0, 0))
            if LB_ERR == i {
                return 0
            }

            iLength := int32(SendMessage(hwndList, LB_GETTEXTLEN, uintptr(i), 0)) + 1
            buf := make([]uint16, iLength)
            SendMessage(hwndList, LB_GETTEXT, uintptr(i), uintptr(unsafe.Pointer(&buf[0])))
            szBuffer := syscall.UTF16ToString(buf)

            fi, err := os.Stat(szBuffer)
            if (err == nil || os.IsExist(err)) && !fi.IsDir() {
                bValidFile = true
                szFile = szBuffer
                szBuffer = GetCurDir()
                if szBuffer[len(szBuffer)-1] != '\\' {
                    szBuffer = szBuffer + "\\"
                }

                SetWindowText(hwndText, _T(szBuffer+szFile))
            } else {
                bValidFile = false
                re_drive := regexp.MustCompile("\\[\\-(\\w)\\-\\]")
                drive := re_drive.FindStringSubmatch(szBuffer)
                if len(drive) == 0 || SetCurrentDirectory(_T(drive[1]+":")) == FALSE {
                    re_folder := regexp.MustCompile("\\[(.*)\\]")
                    folder := re_folder.FindStringSubmatch(szBuffer)
                    if len(folder) == 0 || SetCurrentDirectory(_T(folder[1])) == FALSE {
                        SetCurrentDirectory(_T(szBuffer[0:2]))
                    }
                }

                // Get the new directory name and fill the list box.
                szBuffer = GetCurDir()
                SetWindowText(hwndText, _T(szBuffer))
                SendMessage(hwndList, LB_RESETCONTENT, 0, 0)
                SendMessage(hwndList, LB_DIR, DIRATTR, uintptr(unsafe.Pointer(_T("*.*"))))
            }

            InvalidateRect(hwnd, nil, true)
        }

        return 0
    })

    app.On(WM_PAINT, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        if !bValidFile {
            return 0
        }

        //fmt.Println(szFile)

        f, err := os.Open(szFile)

        if err != nil {
            bValidFile = false
            return 0
        }

        defer func() { f.Close() }()
        var buf []byte = make([]byte, MAXREAD)
        n, err := f.Read(buf)
        if err != nil || n == 0 {
            return 0
        }

        // i now equals the number of bytes in buffer.
        // Commence getting a device context for displaying text.
        var ps PAINTSTRUCT
        hdc := BeginPaint(hwnd, &ps)
        SelectObject(hdc, GetStockObject(SYSTEM_FIXED_FONT))
        SetTextColor(hdc, COLORREF(GetSysColor(COLOR_BTNTEXT)))
        SetBkColor(hdc, COLORREF(GetSysColor(COLOR_BTNFACE)))
        // Assume the file is ASCII
        DrawTextA(hdc, (*uint16)(unsafe.Pointer((&buf[0]))), int32(n), &rect, DTFLAGS)
        EndPaint(hwnd, &ps)
        return 0
    })

    app.On(WM_SETFOCUS, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        SetFocus(hwndList)
        return 0
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })

    app.Init("head", "head")
    app.Run()
}
