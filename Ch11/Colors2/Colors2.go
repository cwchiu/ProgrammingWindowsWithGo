package main

import (
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

type ScrollProc func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr)

func main() {
    app, _ := NewApp()
    app.BackgroundBrush = HBRUSH(CreateSolidBrush(0))
    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        DeleteObject(HGDIOBJ(SetClassLong(hwnd, GCL_HBRBACKGROUND,
            int32(GetStockObject(WHITE_BRUSH)))))

        PostQuitMessage(0)
        return 0
    })

    app.BackgroundBrush = CreateSolidBrush(0)
    app.Init("Colors2", "Color Scroll")

    var iColor [3]int32
    app.AddModelessDialog(CreateDialog(
        app.HInstance,
        _T("ColorScrDlg"),
        app.HWnd,
        syscall.NewCallback(func(hDlg HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
            switch msg {
            case WM_INITDIALOG:
                for iCtrlID := 10; iCtrlID < 13; iCtrlID++ {
                    hCtrl := GetDlgItem(hDlg, int32(iCtrlID))
                    SetScrollRange(hCtrl, SB_CTL, 0, 255, FALSE)
                    SetScrollPos(hCtrl, SB_CTL, 0, FALSE)
                }
                return TRUE

            case WM_VSCROLL:
                hCtrl := HWND(lParam)
                iCtrlID := GetWindowLong(hCtrl, GWL_ID)
                iIndex := iCtrlID - 10
                hwndParent := GetParent(hDlg)

                switch LOWORD(uint32(wParam)) {
                case SB_PAGEDOWN:
                    iColor[iIndex] += 15 // fall through
                    iColor[iIndex] = Min(255, iColor[iIndex]+1)
                case SB_LINEDOWN:
                    iColor[iIndex] = Min(255, iColor[iIndex]+1)
                case SB_PAGEUP:
                    iColor[iIndex] -= 15 // fall through
                    iColor[iIndex] = Max(0, iColor[iIndex]-1)
                case SB_LINEUP:
                    iColor[iIndex] = Max(0, iColor[iIndex]-1)
                case SB_TOP:
                    iColor[iIndex] = 0
                case SB_BOTTOM:
                    iColor[iIndex] = 255
                case SB_THUMBTRACK, SB_THUMBPOSITION:
                    iColor[iIndex] = int32(HIWORD(uint32(wParam)))
                default:
                    return FALSE
                }

                SetScrollPos(hCtrl, SB_CTL, iColor[iIndex], TRUE)
                SetDlgItemInt(hDlg, iCtrlID+3, UINT(iColor[iIndex]), FALSE)

                DeleteObject(HGDIOBJ(SetClassLong(hwndParent, GCL_HBRBACKGROUND, int32(CreateSolidBrush(RGB(iColor[0], iColor[1], iColor[2]))))))

                InvalidateRect(hwndParent, nil, true)
                return TRUE
            }

            return FALSE
        },
    )))

    app.Run()
}
