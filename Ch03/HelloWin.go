package main

import (
    "syscall"
    "unsafe"
    "fmt"
  //  "time"
    . "github.com/cwchiu/go-winapi"
)

var WndProcPtr uintptr = syscall.NewCallback(WndProc)

func WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
    switch(msg){
    case WM_CREATE:
        PlaySound(syscall.StringToUTF16Ptr("hellowin.wav"), HWND(0), SND_FILENAME|SND_ASYNC )   
        //r := PlaySound(syscall.StringToUTF16Ptr("SystemStart"), HWND(0), SND_ALIAS )   
        return 0;
    case WM_DESTROY:
        PostQuitMessage(0)
        return 0
    case WM_PAINT:
        var ps PAINTSTRUCT
        hdc := BeginPaint(hwnd, &ps)
        //fmt.Println(hdc)
        
        var rect RECT
        GetClientRect(hwnd, &rect)
        //fmt.Println(rect)
        var params DRAWTEXTPARAMS
        params.CbSize = uint32(unsafe.Sizeof(params)) 
        DrawTextEx(
            hdc, syscall.StringToUTF16Ptr("Hello, Window"), -1, &rect, 
            DT_SINGLELINE | DT_CENTER | DT_VCENTER, &params)
        
        EndPaint(hwnd, &ps)
        return 0
    }
    return DefWindowProc(hwnd, msg, wParam, lParam)
}

func main(){
    _T := syscall.StringToUTF16Ptr
    
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
    szAppName := _T("HelloWin")
    
    var wc WNDCLASSEX
    wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.Style = CS_HREDRAW | CS_VREDRAW
	wc.LpfnWndProc = WndProcPtr
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
		_T("The Hello Program"),
		WS_OVERLAPPEDWINDOW  | WS_BORDER | WS_CAPTION | WS_SYSMENU | WS_MAXIMIZEBOX | WS_MINIMIZEBOX,
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

    //time.Sleep(10000 * time.Millisecond)
    //fmt.Println(r)
    var msg MSG
    for GetMessage(&msg, HWND_TOP, 0, 0) == TRUE{
        TranslateMessage(&msg)
        DispatchMessage(&msg)
    }
}