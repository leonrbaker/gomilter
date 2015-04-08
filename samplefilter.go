// +build ignore

/*
Copyright (c) 2015 Leon Baker
This projected is licensed under the terms of the LGPL-3.0 license.



Sample Milter



*/

package main

import (

  "."
  "fmt"

)


type Mymilter struct {
  gomilter.MilterRaw // Embed the basic functionality. No callbacks yet
}

// Define the callback functions we are going to use
func (m *Mymilter) Connect(ctx uintptr, hostname, ip string) (sfsistat int8) {
  fmt.Println("mymilter.connect was called")
  fmt.Printf("hostname: %s\n", hostname)
  fmt.Printf("ip: %s\n", ip)
  return gomilter.Reject
}

func (m *Mymilter) Helo(ctx uintptr, helohost string) (sfsistat int8) {
  fmt.Println("mymilter.helo was called")
  return
}

func main() {
  mymilter := new(Mymilter)
  mymilter.FilterName = "TestFilter"
  mymilter.Debug = true
  mymilter.Flags = gomilter.AddHdrs|gomilter.AddRcpt

  // Start Milter
  gomilter.Run(mymilter)
}