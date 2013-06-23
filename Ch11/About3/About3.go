package main

import (
    "fmt"
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"
    "unsafe"
)

const (
    IDM_APP_ABOUT = 40001
    IDC_STATIC    = -1

    szAppName = "About3"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func main() {
    app, _ := NewApp()

    app.On(WM_COMMAND, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {

        switch int32((LOWORD(uint32(wParam)))) {
        case IDM_APP_ABOUT:
            DialogBox(app.HInstance, _T("AboutBox"), hwnd, syscall.NewCallback(func(hDlg HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
                switch msg {
                case WM_INITDIALOG:
                    return TRUE
                case WM_COMMAND:
                    switch int32(LOWORD(uint32(wParam))) {
                    case IDOK:
                        EndDialog(hDlg, TRUE)
                        return TRUE
                    }
                }
                return FALSE
            }))
        }

        return MSG_IGNORE
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })

    app.MenuName = _T(szAppName)
    app.Icon = LoadIcon(app.HInstance, _T(szAppName))
    err := app.Init(szAppName, "About Box Demo Program")
    if err != nil {
        fmt.Println(err)
        return
    }

    app.RegisterClass(_T("EllipPush"), func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {

        switch msg {
        case WM_PAINT:
            var rect RECT
            GetClientRect(hwnd, &rect)
            buf := make([]uint16, 40)
            size := GetWindowTextA(hwnd, uintptr(unsafe.Pointer(&buf[0])), 40)
            szText := syscall.UTF16ToString(buf[0:size])

            var ps PAINTSTRUCT
            hdc := BeginPaint(hwnd, &ps)

            hBrush := CreateSolidBrush(COLORREF(GetSysColor(COLOR_WINDOW)))
            hBrush = HBRUSH(SelectObject(hdc, HGDIOBJ(hBrush)))
            SetBkColor(hdc, COLORREF(GetSysColor(COLOR_WINDOW)))
            SetTextColor(hdc, COLORREF(GetSysColor(COLOR_WINDOWTEXT)))

            Ellipse(hdc, rect.Left, rect.Top, rect.Right, rect.Bottom)
            DrawTextA(hdc, _T(szText), -1, &rect,
                DT_SINGLELINE|DT_CENTER|DT_VCENTER)

            DeleteObject(HGDIOBJ(SelectObject(hdc, HGDIOBJ(hBrush))))

            EndPaint(hwnd, &ps)
            return FALSE
        case WM_KEYUP:
            if wParam != VK_SPACE {
                break
            }
            SendMessage(GetParent(hwnd), WM_COMMAND,
                uintptr(GetWindowLong(hwnd, GWL_ID)), uintptr(hwnd))
            return FALSE
        case WM_LBUTTONUP:
            SendMessage(GetParent(hwnd), WM_COMMAND,
                uintptr(GetWindowLong(hwnd, GWL_ID)), uintptr(hwnd))
            return FALSE
        }
        return DefWindowProc(hwnd, msg, wParam, lParam)
    }, nil, LoadCursor(0, MAKEINTRESOURCE(IDC_ARROW)), HBRUSH((COLOR_BTNFACE + 1)))

    app.Run()
}
