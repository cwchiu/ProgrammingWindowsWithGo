package main

import (
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"
)

const (
    IDM_APP_ABOUT = 40001
    IDC_BLACK     = 1000
    IDC_BLUE      = 1001
    IDC_GREEN     = 1002
    IDC_CYAN      = 1003
    IDC_RED       = 1004
    IDC_MAGENTA   = 1005
    IDC_YELLOW    = 1006
    IDC_WHITE     = 1007
    IDC_RECT      = 1008
    IDC_ELLIPSE   = 1009
    IDC_PAINT     = 1010
    szAppName     = "About2"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

var (
    iCurrentColor  int32 = IDC_BLACK
    iCurrentFigure int32 = IDC_RECT
)

func PaintWindow(hwnd HWND, iColor, iFigure int32) {
    var crColor [8]COLORREF = [8]COLORREF{
        RGB(0, 0, 0), RGB(0, 0, 255),
        RGB(0, 255, 0), RGB(0, 255, 255),
        RGB(255, 0, 0), RGB(255, 0, 255),
        RGB(255, 255, 0), RGB(255, 255, 255),
    }

    hdc := GetDC(hwnd)
    var rect RECT
    GetClientRect(hwnd, &rect)
    hBrush := CreateSolidBrush(crColor[iColor-IDC_BLACK])
    hBrush = HBRUSH(SelectObject(hdc, HGDIOBJ(hBrush)))

    if iFigure == IDC_RECT {
        Rectangle(hdc, rect.Left, rect.Top, rect.Right, rect.Bottom)
    } else {
        Ellipse(hdc, rect.Left, rect.Top, rect.Right, rect.Bottom)
    }
    DeleteObject(SelectObject(hdc, HGDIOBJ(hBrush)))
    ReleaseDC(hwnd, hdc)
}

func PaintTheBlock(hCtrl HWND, iColor, iFigure int32) {
    InvalidateRect(hCtrl, nil, true)
    UpdateWindow(hCtrl)
    PaintWindow(hCtrl, iColor, iFigure)
}

func main() {
    var iColor, iFigure int32
    var hCtrlBlock HWND

    app, _ := NewApp()
    app.On(WM_COMMAND, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {

        switch int32((LOWORD(uint32(wParam)))) {
        case IDM_APP_ABOUT:
            ret := DialogBox(app.HInstance, _T("AboutBox"), hwnd, syscall.NewCallback(func(hDlg HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
                switch msg {
                case WM_INITDIALOG:
                    iColor = iCurrentColor
                    iFigure = iCurrentFigure

                    CheckRadioButton(hDlg, IDC_BLACK, IDC_WHITE, iColor)
                    CheckRadioButton(hDlg, IDC_RECT, IDC_ELLIPSE, iFigure)

                    hCtrlBlock = GetDlgItem(hDlg, IDC_PAINT)

                    SetFocus(GetDlgItem(hDlg, iColor))
                case WM_PAINT:
                    PaintTheBlock(hCtrlBlock, iColor, iFigure)
                case WM_COMMAND:
                    switch int32(LOWORD(uint32(wParam))) {
                    case IDCANCEL:
                        EndDialog(hDlg, FALSE)
                        return TRUE
                    case IDOK:
                        iCurrentColor = iColor
                        iCurrentFigure = iFigure
                        EndDialog(hDlg, TRUE)
                        return TRUE
                    case IDC_BLACK, IDC_WHITE, IDC_CYAN, IDC_MAGENTA, IDC_BLUE, IDC_YELLOW, IDC_GREEN, IDC_RED:
                        iColor = int32(LOWORD(uint32(wParam)))
                        CheckRadioButton(hDlg, IDC_BLACK, IDC_WHITE, int32(LOWORD(uint32(wParam))))
                        PaintTheBlock(hCtrlBlock, iColor, iFigure)
                        return TRUE
                    case IDC_RECT, IDC_ELLIPSE:
                        iFigure = int32(LOWORD(uint32(wParam)))
                        CheckRadioButton(hDlg, IDC_RECT, IDC_ELLIPSE, int32(LOWORD(uint32(wParam))))
                        PaintTheBlock(hCtrlBlock, iColor, iFigure)
                        return TRUE

                    }
                }
                return 0
            }))

            if ret > 0 {
                InvalidateRect(hwnd, nil, true)
            }
        }

        return MSG_IGNORE
    })

    app.On(WM_PAINT, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        var ps PAINTSTRUCT
        BeginPaint(hwnd, &ps)
        PaintWindow(hwnd, iCurrentColor, iCurrentFigure)
        EndPaint(hwnd, &ps)
        return 0
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })

    app.MenuName = _T(szAppName)
    app.Icon = LoadIcon(app.HInstance, _T(szAppName))
    app.Init(szAppName, "About Box Demo Program")
    app.Run()
}
