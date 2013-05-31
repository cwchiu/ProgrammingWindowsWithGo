package main
import (
    "fmt"
    . "github.com/cwchiu/go-winapi"
    "syscall"
)


var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func MessageBoxPrintf (szCaption string, szFormat string, a ...interface{}){
    szBuffer := fmt.Sprintf(szFormat, a...)
    MessageBox (HWND(0), _T(szBuffer), _T(szCaption), 0)
}

func main() {
     cxScreen := GetSystemMetrics (SM_CXSCREEN)
     cyScreen := GetSystemMetrics (SM_CYSCREEN)

     MessageBoxPrintf ("ScrnSize",
                       "The screen is %d pixels wide by %d pixels high.",
                       cxScreen, cyScreen)
}