/*

Copyright (c) 2015 Leon Baker
This projected is licensed under the terms of the MIT License.

gomilter
Go Bindings for libmilter

gomilter.go

*/

package gomilter

/*

#cgo LDFLAGS: -lmilter

#include <stdlib.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include "libmilter/mfapi.h"
#include "filter.h"

*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

type sockaddr_in struct {
	sin_family int8
	sin_port   uint8
	sin_addr   uint32
	sin_zero   uint64
}

// Return values for Callback functions
const (
	Continue = iota
	Reject
	Discard
	Accept
	Tempfail
	_
	_
	Noreply
	Skip
)

// flags
const (
	ADDHDRS     = 0x00000001 // 000000001
	CHGBODY     = 0x00000002 // 000000010
	ADDRCPT     = 0x00000004 // 000000100
	DELRCPT     = 0x00000008 // 000001000
	CHGHDRS     = 0x00000010 // 000010000
	QUARANTINE  = 0x00000020 // 000100000
	CHGFROM     = 0x00000040 // 001000000
	ADDRCPT_PAR = 0x00000080 // 010000000
	SETSYMLIST  = 0x00000100 // 100000000
)

// Interface that must be implemented in order to use gomilter
type Milter interface {
	GetFilterName() string
	GetDebug() bool
	GetFlags() int
	GetSocket() string
}

// An "empty" Milter with no callback functions
type MilterRaw struct {
	FilterName string
	Debug      bool
	Flags      int
	Socket     string
}

func (m *MilterRaw) GetFilterName() string {
	return m.FilterName
}

func (m *MilterRaw) GetDebug() bool {
	return m.Debug
}

func (m *MilterRaw) GetFlags() int {
	return m.Flags
}

func (m *MilterRaw) GetSocket() string {
	return m.Socket
}

// ********* Callback checking types *********
type checkForConnect interface {
	Connect(ctx uintptr, hostname string, ip net.IP) (sfsistat int8)
}

type checkForHelo interface {
	Helo(ctx uintptr, helohost string) (sfsistat int8)
}

type checkForEnvFrom interface {
	EnvFrom(ctx uintptr, myargv []string) (sfsistat int8)
}

type checkForEnvRcpt interface {
	EnvRcpt(ctx uintptr, myargv []string) (sfsistat int8)
}

type checkForHeader interface {
	Header(ctx uintptr, headerf, headerv string) (sfsistat int8)
}

type checkForEoh interface {
	Eoh(ctx uintptr) (sfsistat int8)
}

type checkForBody interface {
	Body(ctx uintptr, body []byte) (sfsistat int8)
}

type checkForEom interface {
	Eom(ctx uintptr) (sfsistat int8)
}

type checkForAbort interface {
	Abort(ctx uintptr) (sfsistat int8)
}

type checkForClose interface {
	Close(ctx uintptr) (sfsistat int8)
}

// ********* Global Milter variable *********
var milter Milter

// ********* Utility Functions (not exported) *********

type CtxPtr *C.struct_smfi_str

func int2ctx(ptr uintptr) CtxPtr {
	return CtxPtr(unsafe.Pointer(ptr))
}

func ctx2int(ctx CtxPtr) uintptr {
	return uintptr(unsafe.Pointer(ctx))
}

func cStringArrayToSlice(argv **C.char) (GoArgv []string) {
	// Build a slice of pointers with C array as a backing
	length := int(C.argv_len(argv))
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(argv)),
		Len:  length,
		Cap:  length,
	}
	argvSlice := *(*[]*C.char)(unsafe.Pointer(&hdr))

	// Build a string slice for Go friendly strings
	GoArgv = make([]string, length, length)
	for i := 0; i < length; i++ {
		GoArgv[i] = C.GoString(argvSlice[i])
	}
	return
}

func GobEncode(data interface{}) ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func GobDecode(buf []byte, data interface{}) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

// ********* Filter Callback functions *********
// These are registered with Milter
// They are only called if they get registered but need to be defined anyway

//export Go_xxfi_connect
func Go_xxfi_connect(ctx *C.SMFICTX, hostname *C.char, hostaddr *C._SOCK_ADDR) (sfsistat C.sfsistat) {
	defer func(sfsistat *C.sfsistat) {
		if r := recover(); r != nil {
			LoggerPrintf("Panic caught in Go_xxfi_connect(): %s", r)
			*sfsistat = 75 // tempfail
		}
	}(&sfsistat)

	ctxptr := ctx2int(ctx)
	var ip net.IP

	if hostaddr.sa_family == C.AF_INET {
		hostaddrin := (*sockaddr_in)(unsafe.Pointer(hostaddr))
		ip_addr := make([]byte, 4)
		binary.LittleEndian.PutUint32(ip_addr, hostaddrin.sin_addr)
		ip = net.IPv4(ip_addr[0], ip_addr[1], ip_addr[2], ip_addr[3])

	} else if hostaddr.sa_family == C.AF_INET6 {
		sa_in := (*C.struct_sockaddr_in6)(unsafe.Pointer(hostaddr))
		ip = net.IP(C.GoBytes(unsafe.Pointer(&sa_in.sin6_addr), 16))
	} else {
		if milter.GetDebug() {
			LoggerPrintln("hostaddr.sa_family value not implemented")
		}
		ip = net.ParseIP("::")
	}

	m := milter.(checkForConnect)
	code := m.Connect(ctxptr, C.GoString(hostname), ip)
	if milter.GetDebug() {
		LoggerPrintf("Connect callback returned: %d\n", code)
	}
	return C.sfsistat(code)
}

//export Go_xxfi_helo
func Go_xxfi_helo(ctx *C.SMFICTX, helohost *C.char) C.sfsistat {
	m := milter.(checkForHelo)
	code := m.Helo(ctx2int(ctx), C.GoString(helohost))
	if milter.GetDebug() {
		LoggerPrintf("Helo callback returned: %d\n", code)
	}
	return C.sfsistat(code)
}

//export Go_xxfi_envfrom
func Go_xxfi_envfrom(ctx *C.SMFICTX, argv **C.char) C.sfsistat {
	// Call our application's callback
	m := milter.(checkForEnvFrom)
	code := m.EnvFrom(ctx2int(ctx), cStringArrayToSlice(argv))
	if milter.GetDebug() {
		LoggerPrintf("EnvFrom callback returned: %d\n", code)
	}
	return C.sfsistat(code)
}

//export Go_xxfi_envrcpt
func Go_xxfi_envrcpt(ctx *C.SMFICTX, argv **C.char) C.sfsistat {
	// Call our application's callback
	m := milter.(checkForEnvRcpt)
	code := m.EnvRcpt(ctx2int(ctx), cStringArrayToSlice(argv))
	if milter.GetDebug() {
		LoggerPrintf("EnvRcpt callback returned: %d\n", code)
	}
	return C.sfsistat(code)
}

//export Go_xxfi_header
func Go_xxfi_header(ctx *C.SMFICTX, headerf, headerv *C.char) C.sfsistat {
	m := milter.(checkForHeader)
	code := m.Header(ctx2int(ctx), C.GoString(headerf), C.GoString(headerv))
	if milter.GetDebug() {
		LoggerPrintf("Header callback returned: %d\n", code)
	}
	return C.sfsistat(code)
}

//export Go_xxfi_eoh
func Go_xxfi_eoh(ctx *C.SMFICTX) C.sfsistat {
	// Call our application's callback
	m := milter.(checkForEoh)
	code := m.Eoh(ctx2int(ctx))
	if milter.GetDebug() {
		LoggerPrintf("Eoh callback returned: %d\n", code)
	}
	return C.sfsistat(code)
}

//export Go_xxfi_body
func Go_xxfi_body(ctx *C.SMFICTX, bodyp *C.uchar, bodylen C.size_t) C.sfsistat {
	// Create a byte slice from the body pointer.
	// The body pointer just points to a sequence of bytes
	// Our bit slice is backed by the original body. No copy is made here
	// Be aware that converting the byte slice to string will make a copy
	var b []byte
	b = (*[1 << 30]byte)(unsafe.Pointer(bodyp))[0:bodylen]
	// Call our application's callback
	m := milter.(checkForBody)
	code := m.Body(ctx2int(ctx), b)
	if milter.GetDebug() {
		LoggerPrintf("Body callback returned: %d\n", code)
	}
	return C.sfsistat(code)
}

//export Go_xxfi_eom
func Go_xxfi_eom(ctx *C.SMFICTX) C.sfsistat {
	// Call our application's callback
	m := milter.(checkForEom)
	code := m.Eom(ctx2int(ctx))
	if milter.GetDebug() {
		LoggerPrintf("Eom callback returned: %d\n", code)
	}
	return C.sfsistat(code)
}

//export Go_xxfi_abort
func Go_xxfi_abort(ctx *C.SMFICTX) C.sfsistat {
	// Call our application's callback
	m := milter.(checkForAbort)
	code := m.Abort(ctx2int(ctx))
	if milter.GetDebug() {
		LoggerPrintf("Abort callback returned: %d\n", code)
	}
	return C.sfsistat(code)
}

//export Go_xxfi_close
func Go_xxfi_close(ctx *C.SMFICTX) C.sfsistat {
	// Call our application's callback
	m := milter.(checkForClose)
	code := m.Close(ctx2int(ctx))
	if milter.GetDebug() {
		LoggerPrintf("Close callback returned: %d\n", code)
	}
	return C.sfsistat(code)
}

// ********* libmilter Data Access Functions *********

func GetSymVal(ctx uintptr, symname string) string {
	csymname := C.CString(symname)
	defer C.free(unsafe.Pointer(csymname))
	type CtxPtr *C.struct_smfi_str
	cval := C.smfi_getsymval(int2ctx(ctx), csymname)
	// Note: If you try to free the cval C string it will crash
	return C.GoString(cval)
}

// See also: http://bit.ly/1HVWA9I
func SetPriv(ctx uintptr, privatedata interface{}) int {
	// privatedata seems to work for any data type
	// Structs must have exported fields

	// Serialize Go privatedata into a byte slice
	bytedata, _ := GobEncode(privatedata)

	// length and size
	// length is a uint32 (usually 4 bytes)
	// the length will be stored in front of the byte sequence
	length := uint32(len(bytedata))
	lengthsize := C.size_t(unsafe.Sizeof(length))
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return -1
	}
	lengthbytes := buf.Bytes()

	// Allocate memory for the length and byte sequence
	CArray := (*C.uchar)(C.malloc(lengthsize + C.size_t(length)))

	var lenStart, seqStart uintptr
	lenStart = uintptr(unsafe.Pointer(CArray))
	seqStart = lenStart + uintptr(lengthsize)

	CArray = (*C.uchar)(unsafe.Pointer(lenStart))

	for i := uintptr(0); i < uintptr(lengthsize); i++ {
		CArray = (*C.uchar)(unsafe.Pointer(lenStart + i))
		*CArray = C.uchar(lengthbytes[i])
	}

	// Now copy the data bytes to the position after the length
	for i := uintptr(0); i < uintptr(length); i++ {
		CArray = (*C.uchar)(unsafe.Pointer(seqStart + i))
		*CArray = C.uchar(bytedata[i])
	}

	// Call libmilter smfi_setpriv
	type CtxPtr *C.struct_smfi_str
	return int(C.smfi_setpriv(int2ctx(ctx), unsafe.Pointer(lenStart)))
}

func GetPriv(ctx uintptr, privatedata interface{}) int {
	/*  Retrieve the private data stored by the milter
	    Retrieving the data will release the memory allocated for it
	    Don't try to retrieve it again unless you call SetPriv first
	*/

	// Call libmilter smfi_getpriv to get a pointer to our data
	CArray := (*byte)(C.smfi_getpriv(int2ctx(ctx)))

	// Make sure data has been set with a previous call to SetPriv
	if CArray == nil {
		return -1
	}

	// Read uint32 size bytes from the start of the pointer
	var length uint32
	lengthsize := unsafe.Sizeof(length)
	lengthbytes := make([]byte, lengthsize)

	lenStart := uintptr(unsafe.Pointer(CArray))
	seqStart := lenStart + uintptr(lengthsize)

	for i := uintptr(0); i < uintptr(lengthsize); i++ {
		CArray = (*byte)(unsafe.Pointer(lenStart + i))
		lengthbytes[i] = *CArray
	}

	// Binary decode the length bytes
	buf := bytes.NewBuffer(lengthbytes)
	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return -1
	}

	// Read byte sequence of data
	databytes := make([]byte, length)
	for i := uintptr(0); i < uintptr(length); i++ {
		CArray = (*byte)(unsafe.Pointer(seqStart + i))
		databytes[i] = *CArray
	}

	// Free the data malloc'ed by C
	C.free(unsafe.Pointer(lenStart))
	C.smfi_setpriv(int2ctx(ctx), nil)

	// Unserialize the data bytes back into a data structure
	err = GobDecode(databytes, privatedata)
	if err != nil {
		return -1
	}

	return 0
}

func SetReply(ctx uintptr, rcode, xcode, message string) int {
	type CtxPtr *C.struct_smfi_str
	crcode := C.CString(rcode)
	defer C.free(unsafe.Pointer(crcode))
	cxcode := C.CString(xcode)
	defer C.free(unsafe.Pointer(cxcode))
	cmessage := C.CString(message)
	defer C.free(unsafe.Pointer(cmessage))
	// Call libmilter setreply
	return int(C.smfi_setreply(int2ctx(ctx), crcode, cxcode, cmessage))
}

func SetMLReply(ctx uintptr, rcode, xcode string, message ...string) int {
	/*  Allows up to 32 lines of message
	 */

	crcode := C.CString(rcode)
	defer C.free(unsafe.Pointer(crcode))
	cxcode := C.CString(xcode)
	defer C.free(unsafe.Pointer(cxcode))

	// Build message C array
	// Get size of a C pointer on this system
	ptrSize := unsafe.Sizeof(crcode)

	// Allocate the char** array
	argv := C.malloc(C.size_t(len(message)) * C.size_t(ptrSize))
	defer C.free(argv)
	// Assign each line to its address offset
	for i := 0; i < len(message); i++ {
		linePtr := (**C.char)(unsafe.Pointer(uintptr(argv) + uintptr(i)*ptrSize))
		line := C.CString(message[i])
		defer C.free(unsafe.Pointer(line))
		*linePtr = line
	}

	// Call our C wrapper for setmlreply
	// The wrapper is needed as cgo doesn't support variadic C function calls
	return int(C.wrap_setmlreply(int2ctx(ctx), crcode, cxcode, C.int(len(message)), (**C.char)(argv)))
}

// ********* libmilter Message Modification Functions *********

func AddHeader(ctx uintptr, headerf, headerv string) int {
	/*  Add a header to the message.  SMFIF_ADDHDRS
	 */

	cheaderf := C.CString(headerf)
	defer C.free(unsafe.Pointer(cheaderf))
	cheaderv := C.CString(headerv)
	defer C.free(unsafe.Pointer(cheaderv))

	// Call smfi_addheader
	return int(C.smfi_addheader(int2ctx(ctx), cheaderf, cheaderv))
}

func ChgHeader(ctx uintptr, headerf string, hdridx int, headerv string) int {
	/*  Change or delete a header.  SMFIF_CHGHDRS
	 */

	cheaderf := C.CString(headerf)
	defer C.free(unsafe.Pointer(cheaderf))
	cheaderv := C.CString(headerv)
	defer C.free(unsafe.Pointer(cheaderv))

	// Call smfi_chgheader
	return int(C.smfi_chgheader(int2ctx(ctx), cheaderf, C.int(hdridx), cheaderv))
}

func InsHeader(ctx uintptr, hdridx int, headerf, headerv string) int {
	/*  Insert a header into the message. SMFIF_ADDHDRS
	 */

	cheaderf := C.CString(headerf)
	defer C.free(unsafe.Pointer(cheaderf))
	cheaderv := C.CString(headerv)
	defer C.free(unsafe.Pointer(cheaderv))

	// Call smfi_insheader
	return int(C.smfi_insheader(int2ctx(ctx), C.int(hdridx), cheaderf, cheaderv))
}

func ChgFrom(ctx uintptr, mail, args string) int {
	/*  Change the envelope sender address. SMFIF_CHGFROM
	 */

	cmail := C.CString(mail)
	defer C.free(unsafe.Pointer(cmail))
	cargs := C.CString(args)
	defer C.free(unsafe.Pointer(cargs))

	return int(C.smfi_chgfrom(int2ctx(ctx), cmail, cargs))
}

func AddRcpt(ctx uintptr, rcpt string) int {
	/*  Add a recipient to the envelope.  SMFIF_ADDRCPT
	 */
	crcpt := C.CString(rcpt)
	defer C.free(unsafe.Pointer(crcpt))
	// Call smfi_addrcpt
	return int(C.smfi_addrcpt(int2ctx(ctx), crcpt))
}

func AddRcpt_Par(ctx uintptr, rcpt, args string) int {
	/*  Add a recipient including ESMTP parameter to the envelope.  SMFIF_ADDRCPT_PAR
	 */
	crcpt := C.CString(rcpt)
	defer C.free(unsafe.Pointer(crcpt))
	cargs := C.CString(args)
	defer C.free(unsafe.Pointer(cargs))
	// Call smfi_addrcpt
	return int(C.smfi_addrcpt_par(int2ctx(ctx), crcpt, cargs))
}

func DelRcpt(ctx uintptr, rcpt string) int {
	/*  Delete a recipient from the envelope. SMFIF_DELRCPT
	 */
	crcpt := C.CString(rcpt)
	defer C.free(unsafe.Pointer(crcpt))
	// Call smfi_addrcpt
	return int(C.smfi_delrcpt(int2ctx(ctx), crcpt))
}

func ReplaceBody(ctx uintptr, body []byte) int {
	/*  Replace the body of the message.  SMFIF_CHGBODY
	 */

	// Allocate memory
	length := len(body)

	// Allocate memory for the length and byte sequence
	cbody := (*C.uchar)(C.malloc(C.size_t(length)))
	start := uintptr(unsafe.Pointer(cbody))

	for i := uintptr(0); i < uintptr(length); i++ {
		cbody = (*C.uchar)(unsafe.Pointer(start + i))
		*cbody = C.uchar(body[i])
	}
	cbody = (*C.uchar)(unsafe.Pointer(start))

	// Call smfi_replacebody
	return int(C.smfi_replacebody(int2ctx(ctx), cbody, C.int(length)))
}

// ********* Other Message Handling Functions *********

func progress(ctx uintptr) int {
	/*  Report operation in progress.
	 */

	// Call smfi_progress
	return int(C.smfi_progress(int2ctx(ctx)))
}

// ********* Run the milter *********

func Stop() {
	C.smfi_stop()
}

func Run(amilter Milter) int {
	milter = amilter
	if milter.GetDebug() {
		LoggerPrintf("Debugging enabled")
	}

	// Declare an empty smfiDesc structure
	var smfilter C.smfiDesc_str

	// Set filter name
	fname := C.CString(milter.GetFilterName())
	defer C.free(unsafe.Pointer(fname))
	smfilter.xxfi_name = fname
	if milter.GetDebug() {
		LoggerPrintf("Filter Name: %s\n", C.GoString(smfilter.xxfi_name))
	}

	// Set version code
	smfilter.xxfi_version = C.SMFI_VERSION

	// Set Flags
	smfilter.xxfi_flags = C.ulong(milter.GetFlags())
	if milter.GetDebug() {
		LoggerPrintf("Flags: 0x%b\n", smfilter.xxfi_flags)
	}

	// Set Callbacks if they are implemented

	// Check if Connect method was implemented
	if _, ok := milter.(checkForConnect); ok {
		if milter.GetDebug() {
			LoggerPrintln("Connect callback implemented")
		}
		C.setConnect(&smfilter)
	} else {
		if milter.GetDebug() {
			LoggerPrintln("Connect callback not implemented")
		}
	}
	// Check if Helo method was implemented
	if _, ok := milter.(checkForHelo); ok {
		if milter.GetDebug() {
			LoggerPrintln("Helo callback implemented")
		}
		C.setHelo(&smfilter)
	} else {
		if milter.GetDebug() {
			LoggerPrintln("Helo callback not implemented")
		}
	}
	// Check if EnvFrom method was implemented
	if _, ok := milter.(checkForEnvFrom); ok {
		if milter.GetDebug() {
			LoggerPrintln("EnvFrom callback implemented")
		}
		C.setEnvFrom(&smfilter)
	} else {
		if milter.GetDebug() {
			LoggerPrintln("EnvFrom callback not implemented")
		}
	}
	// Check if EnvRcpt method was implemented
	if _, ok := milter.(checkForEnvRcpt); ok {
		if milter.GetDebug() {
			LoggerPrintln("EnvRcpt callback implemented")
		}
		C.setEnvRcpt(&smfilter)
	} else {
		if milter.GetDebug() {
			LoggerPrintln("EnvRcpt callback not implemented")
		}
	}
	// Check if Header method was implemented
	if _, ok := milter.(checkForHeader); ok {
		if milter.GetDebug() {
			LoggerPrintln("Header callback implemented")
		}
		C.setHeader(&smfilter)
	} else {
		if milter.GetDebug() {
			LoggerPrintln("Header callback not implemented")
		}
	}
	// Check if Eoh method was implemented
	if _, ok := milter.(checkForEoh); ok {
		if milter.GetDebug() {
			LoggerPrintln("Eoh callback implemented")
		}
		C.setEoh(&smfilter)
	} else {
		if milter.GetDebug() {
			LoggerPrintln("Eoh callback not implemented")
		}
	}
	// Check if Body method was implemented
	if _, ok := milter.(checkForBody); ok {
		if milter.GetDebug() {
			LoggerPrintln("Body callback implemented")
		}
		C.setBody(&smfilter)
	} else {
		if milter.GetDebug() {
			LoggerPrintln("Body callback not implemented")
		}
	}
	// Check if Eom method was implemented
	if _, ok := milter.(checkForEom); ok {
		if milter.GetDebug() {
			LoggerPrintln("Eom callback implemented")
		}
		C.setEom(&smfilter)
	} else {
		if milter.GetDebug() {
			LoggerPrintln("Eom callback not implemented")
		}
	}
	// Check if Abort method was implemented
	if _, ok := milter.(checkForAbort); ok {
		if milter.GetDebug() {
			LoggerPrintln("Abort callback implemented")
		}
		C.setAbort(&smfilter)
	} else {
		if milter.GetDebug() {
			LoggerPrintln("Abort callback not implemented")
		}
	}
	// Check if Close method was implemented
	if _, ok := milter.(checkForClose); ok {
		if milter.GetDebug() {
			LoggerPrintln("Close callback implemented")
		}
		C.setClose(&smfilter)
	} else {
		if milter.GetDebug() {
			LoggerPrintln("Close callback not implemented")
		}
	}

	if milter.GetDebug() {
		LoggerPrintln("smfilter:")
		LoggerPrintln(fmt.Sprint(smfilter))
	}

	// Setup socket connection
	socket := milter.GetSocket()
	if socket == "" {
		panic("No socket name. Set MilterRaw.Socket")
	}
	// Try to delete old socket if it exists
	socketparts := strings.Split(socket, ":")
	if len(socketparts) == 2 {
		os.Remove(socketparts[1])
	}

	csocket := C.CString(socket)
	defer C.free(unsafe.Pointer(csocket))
	if code := C.smfi_setconn(csocket); code != 0 {
		LoggerPrintf("smfi_setconn failed: %d\n", code)
	}

	// Register the filter
	if code := C.smfi_register(smfilter); code == C.MI_FAILURE {
		LoggerPrintf("smfi_register failed: %d\n", code)
	}

	// Hand control to libmilter
	if milter.GetDebug() {
		LoggerPrintln("Handing over to libmilter")
	}
	result := C.smfi_main()
	if milter.GetDebug() {
		LoggerPrintf("smfi_main returned: %v\n", result)
	}
	return int(result)
}

var LoggerPrintln func(...interface{})
var LoggerPrintf func(string, ...interface{})

func init() {
	LoggerPrintln = func(i ...interface{}) { fmt.Println(i...) }
	LoggerPrintf = func(i string, j ...interface{}) { fmt.Printf(i, j...) }
}
