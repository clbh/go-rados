/*
	These are librados bindings for the Go programming language.

	The current implementation is limited in scope and is aimed at
	exposing the minimum functionality required for creating the
	resources required for librbd to consume librados services.

	For Go librbd bindings, please visit:
		https://github.com/clbh/go-rbd

	To obtain the Ceph source from which librados can be built,
	please visit:
		https://github.com/ceph/ceph


	Authors:
		Benoit Page-Guitard (benoit@anchor.net.au)

	License:
		GNU General Public License v3
		http://www.gnu.org/licenses/gpl.html
*/

package gorados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
import "C"

import (
	"errors"
)

// Our bindings version
const VERSION string = "1.0.0"

// Exported types
type Conn struct {
	handle      C.rados_t
	isConnected bool
}

type Pool struct {
	handle   C.rados_ioctx_t
	name     string
	isActive bool
}

////
//   Library version querying
////

func Version() (major, minor, extra int) {
	var c_major, c_minor, c_extra C.int

	C.rados_version(&c_major, &c_minor, &c_extra)

	return int(c_major), int(c_minor), int(c_extra)
}

////
//   Connection creation/destruction functions
////

func Create(clientname string) (*Conn, error) {
	var handle C.rados_t

	if result := C.rados_create2(&handle, C.CString("ceph"), C.CString(clientname), 0); result < 0 {
		return nil, errors.New("Failed to create connection object")
	}

	return &Conn{
		handle:      handle,
		isConnected: false,
	}, nil
}

func (conn *Conn) Connect() error {
	if result := C.rados_connect(conn.handle); result < 0 {
		return errors.New("Failed to connect to cluster")
	}

	conn.isConnected = true
	return nil
}

func (conn *Conn) Close() {
	C.rados_shutdown(conn.handle)
	conn.isConnected = false
}

////
//   Connection configuration context functions
////

func (conn *Conn) ConfigSet(key string, value string) error {
	if result := C.rados_conf_set(conn.handle, C.CString(key), C.CString(value)); result < 0 {
		return errors.New("Failed to set configuration value")
	}

	return nil
}

func (conn *Conn) ConfigGet(key string) (value string, err error) {
	var buf [1024]C.char

	if result := C.rados_conf_get(conn.handle, C.CString(key), &buf[0], 1023); result < 0 {
		return "", errors.New("Failed to fetch configuration value")
	}

	return C.GoString(&buf[0]), nil
}

func (conn *Conn) ConfigReadFile(path string) error {
	if result := C.rados_conf_read_file(conn.handle, C.CString(path)); result < 0 {
		return errors.New("Failed to process external configuration file")
	}

	return nil
}

////
//   Pool management functions
////

func (conn *Conn) PoolOpen(name string) (*Pool, error) {
	var handle C.rados_ioctx_t
	if result := C.rados_ioctx_create(conn.handle, C.CString(name), &handle); result < 0 {
		return nil, errors.New("Failed to open pool")
	}

	return &Pool{
		handle:   handle,
		name:     name,
		isActive: true,
	}, nil
}

func (pool *Pool) Close() {
	C.rados_ioctx_destroy(pool.handle)
	pool.isActive = false
}

func (pool *Pool) Handle() C.rados_ioctx_t {
	return pool.handle
}

func (pool *Pool) Id() int64 {
	return int64(C.rados_ioctx_get_id(pool.handle))
}

func (pool *Pool) Name() string {
	return pool.name
}
