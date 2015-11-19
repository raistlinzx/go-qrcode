package main

import (
  qrcode "github.com/skip2/go-qrcode"
  colorful "github.com/lucasb-eyer/go-colorful"
  // "time"
  "bytes"
  "strings"
  "image/color"
  "image"
  "image/draw"
  "net/http"
  "fmt"
  "github.com/nfnt/resize"
  "image/png"
  "strconv"
)

func Generate(url, fcolor, bcolor, size, logo string) []byte {
  
  var fc, bc color.Color
  var bcFmt = []string{"#","","","","","",""}
  var qrSize = 300
  if size != "" {
    qrSize,_ = strconv.Atoi(size)
  }

  if fcolor != "" {
    // println(strings.Join([]string{"#","","","","","",""}, "a"))
    fc, _ = colorful.Hex(strings.Join(bcFmt, fcolor))
  } else {
    fc = color.Black
  }
  
  if bcolor != "" {
    println(bcolor)
    bc, _ = colorful.Hex(strings.Join(bcFmt, bcolor))
  } else {
    bc = color.White
  }

  q, e := qrcode.New(url, qrcode.High)
  if e != nil {
      println(e.Error())
      return nil
  }
  
  // println(fc)
  // println(q)
  q.ForegroundColor = fc
  q.BackgroundColor = bc
  
  qrPng, _ := q.PNG(qrSize)
  qrPng = DrawLogo(qrPng, logo, bc)
  return qrPng
}

func DrawLogo(qrPng []byte, logo string, logoBgColor color.Color) []byte {
  
  if logo != "" {
    // println(logo)
    resp, err := http.Get(logo)
    
    if err!=nil {
      fmt.Printf("Logo Error=%s\n", err.Error())
      // println(resp)
    } else {
      
      logoImg , _ , _ := image.Decode(resp.Body)
      qrImg , _ , _ := image.Decode(bytes.NewReader(qrPng))
      pb := qrImg.Bounds()
      logoImg = resize.Resize(uint(pb.Dx()/4), uint(pb.Dy()/4), logoImg, resize.Lanczos3)
      logoBgImg := image.NewRGBA(image.Rect(0, 0, logoImg.Bounds().Dx()+2, logoImg.Bounds().Dy()+2))
      draw.Draw(logoBgImg, logoBgImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
      draw.Draw(logoBgImg, logoImg.Bounds(), logoImg, image.Point{-1,-1}, draw.Over)
      

      offsetX := (pb.Dx()-logoBgImg.Bounds().Dx())/2
      newImg := image.NewRGBA(pb)
      draw.Draw(newImg, pb, qrImg, image.Point{}, draw.Src)
      draw.Draw(newImg, newImg.Bounds(), logoBgImg, image.Point{-offsetX,-offsetX}, draw.Over)
      // println(logoImg)
      buf := new(bytes.Buffer)
      err = png.Encode(buf, newImg)
      // send_s3 := buf.Bytes()
      return buf.Bytes()
      
    }
  }
  return qrPng
}
