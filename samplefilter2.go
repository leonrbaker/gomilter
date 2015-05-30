// +build ignore

/*
Copyright (c) 2015 Leon Baker
*/

package main

import (
	"fmt"
	m "github.com/leonrbaker/gomilter"
)

type Mymilter struct {
	m.MilterRaw // Embed the basic functionality. No callbacks defined yet
}

// Data type I want use for private data. Can be any type
type T struct {
	A uint8
	B string
	C string
}

// Define the callback functions we are going to use
func (milter *Mymilter) Connect(ctx uintptr, hostname, ip string) (sfsistat int8) {
	fmt.Println("mymilter.Connect was called")
	fmt.Printf("hostname: %s\n", hostname)
	fmt.Printf("ip: %s\n", ip)

	t := T{1, hostname, "Test"}
	m.SetPriv(ctx, &t)

	return m.Continue
}

func (milter *Mymilter) Helo(ctx uintptr, helohost string) (sfsistat int8) {
	fmt.Println("mymilter.Helo was called")
	fmt.Printf("helohost: %s\n", helohost)
	return
}

func (milter *Mymilter) EnvFrom(ctx uintptr, myargv []string) (sfsistat int8) {
	fmt.Println("mymilter.EnvFrom was called")
	fmt.Printf("myargv: %s\n", myargv)
	// Show the value of a symbol
	fmt.Printf("{mail_addr}: %v\n", m.GetSymVal(ctx, "{mail_addr}"))
	return
}

func (milter *Mymilter) EnvRcpt(ctx uintptr, myargv []string) (sfsistat int8) {
	fmt.Println("mymilter.EnvRcpt was called")
	fmt.Printf("myargv: %s\n", myargv)
	// Show the value of a symbol
	fmt.Printf("{rcpt_addr}: %v\n", m.GetSymVal(ctx, "{rcpt_addr}"))
	return
}

func (milter *Mymilter) Header(ctx uintptr, headerf, headerv string) (sfsistat int8) {
	fmt.Println("mymilter.Header was called")
	fmt.Printf("header field: %s\n", headerf)
	fmt.Printf("header value: %s\n", headerv)
	return
}

func (milter *Mymilter) Eoh(ctx uintptr) (sfsistat int8) {
	fmt.Println("mymilter.Eoh was called")
	return
}

func (milter *Mymilter) Body(ctx uintptr, body []byte) (sfsistat int8) {
	// Be careful as a conversion of body to string will make a copy of body
	fmt.Println("mymilter.Body was called")
	fmt.Println(string(body))
	fmt.Printf("Body Length: %d\n", len(body))
	return
}

func (milter *Mymilter) Eom(ctx uintptr) (sfsistat int8) {
	fmt.Println("mymilter.Eom was called")

	var t T
	fmt.Println("m.GetPri:", m.GetPriv(ctx, &t))
	fmt.Println("t:", t)

	m.AddHeader(ctx, "LEONUX-Mailer",
		"test server;\n\ttest1=\"foobar\"")

	newBody := []byte("This is a new body")
	m.ReplaceBody(ctx, newBody)
	return
}

func (milter *Mymilter) Abort(ctx uintptr) (sfsistat int8) {
	fmt.Println("mymilter.Abort was called")
	return
}

func (milter *Mymilter) Close(ctx uintptr) (sfsistat int8) {
	fmt.Println("mymilter.Close was called")
	return
}

func main() {
	mymilter := new(Mymilter)
	mymilter.FilterName = "TestFilter"
	mymilter.Debug = true
	mymilter.Flags = m.ADDHDRS | m.ADDRCPT | m.CHGFROM | m.CHGBODY
	mymilter.Socket = "unix:/var/milterattachcheck/socket"

	// Start Milter
	m.Run(mymilter)
}
