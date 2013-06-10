package main

import (
    . "github.com/cwchiu/go-winapi"
    . "github.com/cwchiu/winclass"
    "syscall"
    "unsafe"
)

const (
    ID_SMALLER  =    1
    ID_LARGER   =    2
)

var _T func(s string) *uint16 = syscall.StringToUTF16Ptr

func Triangle (hdc HDC, pt []POINT){
     SelectObject (hdc, GetStockObject (BLACK_BRUSH)) 
     Polygon (hdc, &pt[0], 3) 
     SelectObject (hdc, GetStockObject (WHITE_BRUSH)) 
}

func main() {
    app, _ := NewApp()  
    var cxChar, cyChar uint32
    var cxClient, cyClient uint32
    var hwndSmaller, hwndLarger HWND
    app.On(WM_CREATE , func (hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        cxChar = uint32(LOWORD (uint32(GetDialogBaseUnits ()))) 
        cyChar = uint32(HIWORD (uint32(GetDialogBaseUnits ())))
        
        BTN_WIDTH := int32(8 * cxChar)
        BTN_HEIGHT := int32(4 * cyChar)

        // Create the owner-draw pushbuttons
        hwndSmaller = CreateWindowEx (WS_EX_WINDOWEDGE, _T("button"), _T(""),
                                  WS_CHILD | WS_VISIBLE | BS_OWNERDRAW,
                                  0, 0, BTN_WIDTH, BTN_HEIGHT,
                                  hwnd, HMENU(ID_SMALLER), app.HInstance, nil) 

        hwndLarger  = CreateWindowEx (WS_EX_WINDOWEDGE, _T("button"), _T(""),
                                  WS_CHILD | WS_VISIBLE | BS_OWNERDRAW,
                                  0, 0, BTN_WIDTH, BTN_HEIGHT,
                                  hwnd, HMENU(ID_LARGER), app.HInstance, nil) 
                                  
        return 0
    })
    
    app.On(WM_SIZE, func (hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        cxClient = uint32(LOWORD (uint32(lParam)))
        cyClient = uint32(HIWORD (uint32(lParam)))

        BTN_WIDTH  := 8 * cxChar
        BTN_HEIGHT := 4 * cyChar

        // Move the buttons to the new center          
        MoveWindow (hwndSmaller, int32(cxClient / 2 - 3 * BTN_WIDTH  / 2),
                               int32(cyClient / 2 -     BTN_HEIGHT / 2),
                  int32(BTN_WIDTH), int32(BTN_HEIGHT), true)

        MoveWindow (hwndLarger,  int32(cxClient / 2 +     BTN_WIDTH  / 2),
                               int32(cyClient / 2 -     BTN_HEIGHT / 2),
                  int32(BTN_WIDTH), int32(BTN_HEIGHT), true) 
        return 0 
    })
    
    app.On(WM_COMMAND , func (hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        var rc RECT
        GetWindowRect (hwnd, &rc) 
          
               // Make the window 10% smaller or larger
          
          switch (int32(wParam)){
          case ID_SMALLER :
               rc.Left   += int32(cxClient / 20) 
               rc.Right  -= int32(cxClient / 20) 
               rc.Top    += int32(cyClient / 20) 
               rc.Bottom -= int32(cyClient / 20) 
               break                
          case ID_LARGER :
               rc.Left   -= int32(cxClient / 20) 
               rc.Right  += int32(cxClient / 20) 
               rc.Top    -= int32(cyClient / 20) 
               rc.Bottom += int32(cyClient / 20) 
               break 
          }
          
          MoveWindow (hwnd, rc.Left, rc.Top, rc.Right  - rc.Left,
                            rc.Bottom - rc.Top, true) 
          return 0 
    })
    
    app.On(WM_DRAWITEM , func (hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        pdis := (*DRAWITEMSTRUCT)(unsafe.Pointer(lParam))
          
        // Fill area with white and frame it black          
        FillRect (pdis.HDC, &pdis.RcItem,
                    HBRUSH(GetStockObject (WHITE_BRUSH))) 
          
        FrameRect (pdis.HDC, &pdis.RcItem,
                     HBRUSH(GetStockObject (BLACK_BRUSH))) 
          
        // Draw inward and outward black triangles          
        cx := pdis.RcItem.Right  - pdis.RcItem.Left 
        cy := pdis.RcItem.Bottom - pdis.RcItem.Top  
          
          var pt [3]POINT
          
          switch (pdis.CtlID){
          case ID_SMALLER :
               pt[0].X = 3 * cx / 8  
               pt[0].Y = 1 * cy / 8 
               
               pt[1].X = 5 * cx / 8 
               pt[1].Y = 1 * cy / 8 
               
               pt[2].X = 4 * cx / 8 
               pt[2].Y = 3 * cy / 8 
               
               Triangle (pdis.HDC, pt[0:]) 
               
               pt[0].X = 7 * cx / 8  
               pt[0].Y = 3 * cy / 8 
               
               pt[1].X = 7 * cx / 8 
               pt[1].Y = 5 * cy / 8 
               
               pt[2].X = 5 * cx / 8  
               pt[2].Y = 4 * cy / 8 
               
               Triangle (pdis.HDC, pt[0:]) 
               
               pt[0].X = 5 * cx / 8
               pt[0].Y = 7 * cy / 8 
               
               pt[1].X = 3 * cx / 8
               pt[1].Y = 7 * cy / 8 
               
               pt[2].X = 4 * cx / 8
               pt[2].Y = 5 * cy / 8 
               
               Triangle (pdis.HDC, pt[0:]) 
               
               pt[0].X = 1 * cx / 8   
               pt[0].Y = 5 * cy / 8 
               
               pt[1].X = 1 * cx / 8   
               pt[1].Y = 3 * cy / 8 
               
               pt[2].X = 3 * cx / 8   
               pt[2].Y = 4 * cy / 8 
               
               Triangle (pdis.HDC, pt[0:])
               break
               
          case ID_LARGER :
               pt[0].X = 5 * cx / 8   
               pt[0].Y = 3 * cy / 8 
               pt[1].X = 3 * cx / 8   
               pt[1].Y = 3 * cy / 8 
               pt[2].X = 4 * cx / 8   
               pt[2].Y = 1 * cy / 8 
               
               Triangle (pdis.HDC, pt[0:]) 
               
               pt[0].X = 5 * cx / 8   
               pt[0].Y = 5 * cy / 8 
               pt[1].X = 5 * cx / 8   
               pt[1].Y = 3 * cy / 8 
               pt[2].X = 7 * cx / 8   
               pt[2].Y = 4 * cy / 8 
               
               Triangle (pdis.HDC, pt[0:])
               
               pt[0].X = 3 * cx / 8   
               pt[0].Y = 5 * cy / 8 
               
               pt[1].X = 5 * cx / 8   
               pt[1].Y = 5 * cy / 8 
               
               pt[2].X = 4 * cx / 8   
               pt[2].Y = 7 * cy / 8 
               
               Triangle (pdis.HDC, pt[0:]) 
               
               pt[0].X = 3 * cx / 8   
               pt[0].Y = 3 * cy / 8 
               
               pt[1].X = 3 * cx / 8   
               pt[1].Y = 5 * cy / 8 
               
               pt[2].X = 1 * cx / 8   
               pt[2].Y = 4 * cy / 8 
               
               Triangle (pdis.HDC, pt[0:])
               break
          }
          
          // Invert the rectangle if the button is selected          
          if (pdis.ItemState & ODS_SELECTED > 0){
               InvertRect (pdis.HDC, &pdis.RcItem) 
          }
          
          // Draw a focus rectangle if the button has the focus          
          if (pdis.ItemState & ODS_FOCUS > 0){
               pdis.RcItem.Left   += cx / 16 
               pdis.RcItem.Top    += cy / 16 
               pdis.RcItem.Right  -= cx / 16 
               pdis.RcItem.Bottom -= cy / 16 
               
               DrawFocusRect (pdis.HDC, &pdis.RcItem) 
          }
          return 0 
    })
    
    app.On(WM_DESTROY, func (hwnd HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
        PostQuitMessage(0)
        return 0
    })
    app.Init("OwnDraw", "Owner-Draw Button Demo")
    app.Run()
}
