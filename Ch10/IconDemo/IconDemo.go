package main

import (
    "fmt"
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "os"
    "strings"
    "syscall"
    "unsafe"
)

const (
    IDI_ICON = "#101"
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

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
    var cxIcon, cyIcon int32
    var hIcon HICON 
    app.On(WM_CREATE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {          
          hIcon = LoadIcon (app.HInstance, _T(IDI_ICON)) 
          cxIcon = GetSystemMetrics (SM_CXICON) 
          cyIcon = GetSystemMetrics (SM_CYICON) 

        return 0
    })

    app.On(WM_SIZE, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        cxClient  = int32(LOWORD(uint32(lParam)))
        cyClient  = int32(HIWORD(uint32(lParam)))
        return 0
    })
    
    app.On(WM_PAINT, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        var ps PAINTSTRUCT  
         hdc := BeginPaint (hwnd, &ps) 
          var x,y int32
          for y = 0 ; y < cyClient ; y += cyIcon{
               for x = 0 ; x < cxClient ; x += cxIcon{
                    DrawIcon (hdc, x, y, hIcon) 
               }
               }
         EndPaint (hwnd, &ps)
        return 0
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })

    app.Init("IconDemo", "Icon Demo")
    app.Run()
}
