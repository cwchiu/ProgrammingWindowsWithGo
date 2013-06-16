package main

import (
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    //"os"
    //"strings"
    "syscall"
    //"unsafe"
)

const (
    IDM_APP_ABOUT = 40001
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func main() {

    app, _ := NewApp()
    app.On(WM_COMMAND, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {          
        switch int32((LOWORD (uint32(wParam))))          {
          case IDM_APP_ABOUT :
               DialogBox(app.HInstance, _T("AboutBox"), hwnd, syscall.NewCallback(func (hDlg HWND, msg uint32, wParam, lParam uintptr) (result uintptr){
                    switch (msg)                 {
                 case WM_INITDIALOG :
                      return TRUE 
                      
                 case WM_COMMAND :
                      switch (int32(LOWORD (uint32(wParam))))                      {
                      case IDCANCEL, IDOK :
                           EndDialog (hDlg, 0) 
                           return TRUE 
                      }
                 }               
                    return 0
               }))
          }
        return 0
    })

    app.On(WM_DESTROY, func(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })    
    
    app.MenuName = _T("About1")
    app.Init("About1", "About Box Demo Program")
    app.Run()
}
