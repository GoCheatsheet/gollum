// Copyright 2015-2019 trivago N.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build cgo,!unit

package librdkafka

// #cgo CFLAGS: -I/usr/local/include -std=c99 -Wno-deprecated-declarations
// #cgo LDFLAGS: -L/usr/local/lib -L/usr/local/opt/librdkafka/lib -lrdkafka
// #include "wrapper.h"
// #include <string.h>
import "C"

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"unsafe"
)

var (
	// Log is the standard logger used for non-message related errors
	Log = log.New(os.Stderr, "librdkafka: ", log.Lshortfile)
)

//export goErrorHandler
func goErrorHandler(code C.int, reason *C.char) {
	reasonHeader := reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(reason)),
		Len:  int(C.strlen(reason)),
	}
	reasonString := (*string)(unsafe.Pointer(&reasonHeader))
	Log.Printf("%s -- %s", codeToString(int(code)), *reasonString)
}

//export goLogHandler
func goLogHandler(level C.int, facility *C.char, message *C.char) {
	facHeader := reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(facility)),
		Len:  int(C.strlen(facility)),
	}
	facString := (*string)(unsafe.Pointer(&facHeader))

	msgHeader := reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(message)),
		Len:  int(C.strlen(message)),
	}
	msgString := (*string)(unsafe.Pointer(&msgHeader))

	Log.Printf("[%s] %s", *facString, *msgString)
}

func codeToString(code int) string {
	nativeString := C.rd_kafka_err2str(C.rd_kafka_resp_err_t(code))
	if nativeString == nil {
		return "Unknown error"
	}

	textHeader := reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(nativeString)),
		Len:  int(C.strlen(nativeString)),
	}
	text := (*string)(unsafe.Pointer(&textHeader))
	return fmt.Sprintf("(%d) %s", code, *text)
}

// ErrorHandle is a convenience wrapper for handling librdkafka native errors.
// This struct fulfills the standard golang error interface.
type ErrorHandle struct {
	errBuffer [512]byte
}

func (l *ErrorHandle) buffer() *C.char {
	return (*C.char)(unsafe.Pointer(&l.errBuffer[0]))
}

func (l *ErrorHandle) len() C.size_t {
	return C.size_t(len(l.errBuffer))
}

func (l *ErrorHandle) Error() string {
	for i := 0; i < len(l.errBuffer); i++ {
		if l.errBuffer[i] == 0 {
			return string(l.errBuffer[:i])
		}
	}
	return string(l.errBuffer[:len(l.errBuffer)])
}

// ResponseError is used as a wrapper for errors generated by the batch
// producer. The Code member wraps directly to the librdkafka error
// number. The original message is attached to allow backtracking.
type ResponseError struct {
	Userdata []byte
	Code     int
}

func (r ResponseError) Error() string {
	return codeToString(r.Code)
}
