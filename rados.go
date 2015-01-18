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

type Rados struct {
	rados *C.rados_t
}

type RadosConfig struct {
	rados_config *C.rados_config_t
}

type RadosIoctx struct {
	rados_ioctx *C.rados_ioctx_t
}

const VERSION string = "0.1"

func RadosVersion() (major, minor, extra int) {
	var c_major, c_minor, c_extra C.int

	C.rados_version(&c_major, &c_minor, &c_extra)

	return int(c_major), int(c_minor), int(c_extra)
}
