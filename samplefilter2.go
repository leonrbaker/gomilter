// +build ignore

/*
Copyright (c) 2015 Leon Baker
*/

package main

import (
	m "github.com/leonrbaker/gomilter"
	"log"
	"os"
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

var logger *log.Logger

// Define the callback functions we are going to use
func (milter *Mymilter) Connect(ctx uintptr, hostname, ip string) (sfsistat int8) {
	logger.Println("mymilter.Connect was called")
	logger.Printf("hostname: %s\n", hostname)
	logger.Printf("ip: %s\n", ip)

	t := T{1, hostname, "Test"}
	m.SetPriv(ctx, &t)

	return m.Continue
}

func (milter *Mymilter) Helo(ctx uintptr, helohost string) (sfsistat int8) {
	logger.Println("mymilter.Helo was called")
	logger.Printf("helohost: %s\n", helohost)
	return
}

func (milter *Mymilter) EnvFrom(ctx uintptr, myargv []string) (sfsistat int8) {
	logger.Println("mymilter.EnvFrom was called")
	logger.Printf("myargv: %s\n", myargv)
	// Show the value of a symbol
	logger.Printf("{mail_addr}: %v\n", m.GetSymVal(ctx, "{mail_addr}"))
	return
}

func (milter *Mymilter) EnvRcpt(ctx uintptr, myargv []string) (sfsistat int8) {
	logger.Println("mymilter.EnvRcpt was called")
	logger.Printf("myargv: %s\n", myargv)
	// Show the value of a symbol
	logger.Printf("{rcpt_addr}: %v\n", m.GetSymVal(ctx, "{rcpt_addr}"))
	return
}

func (milter *Mymilter) Header(ctx uintptr, headerf, headerv string) (sfsistat int8) {
	logger.Println("mymilter.Header was called")
	logger.Printf("header field: %s\n", headerf)
	logger.Printf("header value: %s\n", headerv)
	return
}

func (milter *Mymilter) Eoh(ctx uintptr) (sfsistat int8) {
	logger.Println("mymilter.Eoh was called")
	return
}

func (milter *Mymilter) Body(ctx uintptr, body []byte) (sfsistat int8) {
	// Be careful as a conversion of body to string will make a copy of body
	logger.Println("mymilter.Body was called")
	logger.Println(string(body))
	logger.Printf("Body Length: %d\n", len(body))
	return
}

func (milter *Mymilter) Eom(ctx uintptr) (sfsistat int8) {
	logger.Println("mymilter.Eom was called")

	var t T
	logger.Println("m.GetPri:", m.GetPriv(ctx, &t))
	logger.Println("t:", t)

	m.AddHeader(ctx, "LEONUX-Mailer",
		"test server;\n\ttest1=\"foobar\"")

	newBody := []byte("This is a new body")
	m.ReplaceBody(ctx, newBody)
	return
}

func (milter *Mymilter) Abort(ctx uintptr) (sfsistat int8) {
	logger.Println("mymilter.Abort was called")
	return
}

func (milter *Mymilter) Close(ctx uintptr) (sfsistat int8) {
	logger.Println("mymilter.Close was called")
	return
}

func main() {
	mymilter := new(Mymilter)
	mymilter.FilterName = "TestFilter"

	logger = log.New(os.Stdout, "", log.LstdFlags)
	mymilter.Logger = logger
	mymilter.Debug = true

	mymilter.Flags = m.ADDHDRS | m.ADDRCPT | m.CHGFROM | m.CHGBODY
	mymilter.Socket = "unix:/var/milterattachcheck/socket"

	// Start Milter
	m.Run(mymilter)
}
