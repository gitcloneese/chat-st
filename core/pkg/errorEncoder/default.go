package errorEncoder

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	nhttp "net/http"
	"strings"
)

const (
	baseContentType = "application"
)

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}

// DefaultErrorEncoder encodes the error to the HTTP response.
func DefaultErrorEncoder(w nhttp.ResponseWriter, r *nhttp.Request, err error) {
	se := FromError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(nhttp.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	w.WriteHeader(int(se.Code))
	_, _ = w.Write(body)
}

func FromError(err error) *errors.Error {
	if err == nil {
		return nil
	}
	if se := new(errors.Error); errors.As(err, &se) {
		return se
	}
	return errors.New(520, "", err.Error())
}
