package main

import (
    "fmt"
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"
    "time"
    "unicode"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func ShowNumber(hwnd HWND, iNumber UINT) {
    szBuffer := fmt.Sprintf("%X", iNumber)
    SetDlgItemText(hwnd, VK_ESCAPE, _T(szBuffer))
}

func CalcIt(iFirstNum UINT, iOperation int32, iNum UINT) DWORD {
    switch iOperation {
    case '=':
        return DWORD(iNum)
    case '+':
        return DWORD(iFirstNum + iNum)
    case '-':
        return DWORD(iFirstNum - iNum)
    case '*':
        return DWORD(iFirstNum * iNum)
    case '&':
        return DWORD(iFirstNum & iNum)
    case '|':
        return DWORD(iFirstNum | iNum)
    case '^':
        return DWORD(iFirstNum ^ iNum)
    case '<':
        return DWORD(iFirstNum << iNum)
    case '>':
        return DWORD(iFirstNum >> iNum)
    case '/':
        if iNum > 0 {
            return DWORD(iFirstNum / iNum)
        } else {
            return MAXDWORD
        }
    case '%':
        if iNum > 0 {
            return DWORD(iFirstNum % iNum)
        } else {
            return MAXDWORD
        }
    default:
        return 0
    }
}

func main() {
    app, _ := NewApp()
    szAppName := "HexCalc"
    app.BackgroundBrush = HBRUSH(COLOR_BTNFACE + 1)
    app.Icon = LoadIcon (app.HInstance, _T(szAppName))
    var bNewNumber bool = true
    var iOperation int32 = '='
    var iNumber, iFirstNum UINT
    onCommand := func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        SetFocus(hwnd)
        wp_lo := LOWORD(uint32(wParam))

        if wp_lo == VK_BACK { // backspace
            iNumber = iNumber / 16
            ShowNumber(hwnd, iNumber)

        } else if wp_lo == VK_ESCAPE { // escape
            iNumber = 0
            ShowNumber(hwnd, iNumber)

        } else if unicode.Is(unicode.Hex_Digit, rune(wp_lo)) { // hex digit

            if bNewNumber {
                iFirstNum = iNumber
                iNumber = 0
            }
            bNewNumber = false

            var tmp UINT = UINT(MAXDWORD >> 4)
            if iNumber <= tmp {
                var v UINT
                if unicode.IsDigit(rune(wParam)) {
                    v = '0'
                } else {
                    v = 'A' - 10
                }

                iNumber = 16*iNumber + UINT(wParam) - v

                ShowNumber(hwnd, iNumber)
            } else {
                MessageBeep(0)
            }
        } else {
            // operation
            if !bNewNumber {
                iNumber = UINT(CalcIt(iFirstNum, iOperation, iNumber))
                ShowNumber(hwnd, iNumber)
            }
            bNewNumber = true
            iOperation = int32(LOWORD(uint32(wParam)))
        }
        return 0
    }

    onChar := func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        ch := int32(unicode.ToUpper(rune(wParam)))
        if ch == VK_RETURN {
            wParam = '='
        } else {
            wParam = uintptr(ch)
        }

        hButton := GetDlgItem(hwnd, int32(wParam))
        if hButton > 0 {
            SendMessage(hButton, BM_SETSTATE, 1, 0)
            time.Sleep(100)
            SendMessage(hButton, BM_SETSTATE, 0, 0)
            return onCommand(hwnd, msg, wParam, lParam)
        } else {
            MessageBeep(0)
            return MSG_IGNORE
        }

    }

    app.On(WM_COMMAND, onCommand)
    app.On(WM_CHAR, onChar)

    app.On(WM_KEYDOWN, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        if wParam != VK_LEFT {
            return MSG_IGNORE
        }

        wParam = VK_BACK
        return onChar(hwnd, msg, wParam, lParam)
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })
    app.InitWithDialog(szAppName)

    app.Run()
}
