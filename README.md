# gomilter
Go Bindings for Sendmail's libmilter

Tested on Linux and FreeBSD

## Installation

The Sendmail header file libmilter/mfapi.h is required. For Redhat/CentOS, install the sendmail-devel package:
```sh
yum install sendmail-devel
```

Install the gomilter package:

```sh
go get github.com/leonrbaker/gomilter
```

##Usage

The milter is implemented in a struct. Start by defining your own struct type and embeding the gomilter MilterRaw struct.

```go
type Mymilter struct {
	gomilter.MilterRaw // Embed the basic functionality.
}
```

Milter callbacks are added by implementing methods for the struct with matching predefined names.

### Callbacks

* Connect
* Helo
* EnvFrom
* EnvRcpt
* Data
* Header
* Eoh
* Body
* Eom
* Abort
* Close

Not all the callbacks need to be defined. The callbacks are explained on the milter.org site. Unfortunately the milter.org site has been shut down but it is still on [web.archive.org](http://web.archive.org/web/20150510034154/https://www.milter.org/developers/api/index)

### Message Modification Functions

* AddHeader
* ChgHeader
* InsHeader
* ChgFrom
* AddRcpt
* AddRcpt_Par
* DelRcpt
* ReplaceBody

### Other Message Handling Functions
* progress

### Startup

The Socket field of the milter struct must be set. For example:
```go
mymilter.Socket = "unix:/var/gomilter/socket"
```

Control is handed over to the libmilter smfi_main function by calling the Run method and passing it a pointer to your milter struct
```go
gomilter.Run(mymilter)
```

The milter has a Stop method which calls the libmilter smfi_stop function.

### Private Data

libmilter is able to store private data for a connection. This data can be accessed from other functions and callbacks for the same connection. You can pass a pointer to any data structure to SetPriv. The data is retrieved with GetPriv
```go
t := T{1, 2, 3}
m.SetPriv(ctx, &t)
```
Retrieve the data with
```go
var t T
m.GetPriv(ctx, &t))
```

GetPriv should only be called once. If the private data is needed in another function or callback then call SetPriv again.

## Sample Programs
There are two sample programs included, samplefilter.go and samplefilter2.go

##Other Libraries

A usefull MIME parsing library is [go.enmime](https://github.com/jhillyerd/go.enmime)

