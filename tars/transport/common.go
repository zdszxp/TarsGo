package transport

import(
	"net"
	"errors"
	"context"
	"strings"
) 

func isNoDataError(err error) bool {
	netErr, ok := err.(net.Error)
	if ok && netErr.Timeout() && netErr.Temporary() {
		return true
	}
	return false
}

func IsTimeoutError(err error) bool {
	netErr, ok := err.(net.Error)
	if ok && netErr.Timeout() && netErr.Temporary() {
		return true
	}
	return false
}

// ErrNetClosing is returned when a network descriptor is used after
// it has been closed. equal to "internal/poll"
var ErrNetClosing = errors.New("use of closed network connection")

// IsClosedConnError reports whether err is an error from use of a closed
// network connection.
func IsClosedConnError(err error) bool {
	if err == nil {
		return false
	}

	if oe, ok := err.(*net.OpError); ok && oe.Op == "read" {
		if oe.Err.Error() == ErrNetClosing.Error() {
			return true
		}
	}

	// TODO: remove this string search and be more like the Windows
	// case below. That might involve modifying the standard library
	// to return better error types.
	if strings.Contains(err.Error(), "use of closed network connection") {
		return true
	}

	// if runtime.GOOS == "windows" {
	// 	if oe, ok := err.(*net.OpError); ok && oe.Op == "read" {
	// 		if se, ok := oe.Err.(*os.SyscallError); ok && se.Syscall == "wsarecv" {
	// 			const WSAECONNABORTED = 10053
	// 			const WSAECONNRESET = 10054
	// 			if n := errno(se.Err); n == WSAECONNRESET || n == WSAECONNABORTED {
	// 				return true
	// 			}
	// 		}
	// 	}
	// }
	
	return false
}

var connKey = struct{}{}

func GetNetConnFromContext(ctx context.Context) (net.Conn, error) {
	value := ctx.Value(connKey)
	if value == nil {
		return nil, errors.New("ctx has not set connKey")
	}

	netConn, ok := value.(net.Conn)
	if !ok {
		return nil, errors.New("type of connKey is not net.Conn")
	}

	return netConn, nil
}

func contextWithNetConn(ctx context.Context, conn net.Conn) context.Context {
	return context.WithValue(ctx, connKey, conn)
}