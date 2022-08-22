package util

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rotisserie/eris"
	"github.com/unrolled/render"
	"github.com/zufzuf/cake-store/libs/logger"
	"go.uber.org/zap"
)

type CTXValue string

const (
	CTXTrackerID = CTXValue("CTX.Tracker.ID")
)

func CTXTracker(ctx context.Context) string {
	v, _ := ctx.Value(CTXTrackerID).(string)
	return v
}

var (
	Render = render.New()
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Payload any    `json:"payload"`
	Err     any    `json:"error"`
}

type Error struct {
	TrackerID string `json:"tracker_id"`
	Cause     any    `json:"cause,omitempty"`
}

func HTTPResponse(rw http.ResponseWriter, code int, message string, payload any) error {
	return Render.JSON(rw, code, Response{
		Code:    code,
		Message: message,
		Payload: payload,
	})
}

func ErrHTTPResponse(ctx context.Context, rw http.ResponseWriter, err error) {
	var (
		trackerId = CTXTracker(ctx)
		code      = http.StatusInternalServerError
		unpack    = eris.Unpack(err)
	)

	logger.Log.With(
		zap.String("tracker_id", trackerId),
		zap.Any("error", eris.ToJSON(err, true)),
	).Error(unpack.ErrRoot.Msg)

	Render.JSON(rw, code, Response{
		Code:    code,
		Message: unpack.ErrRoot.Msg,
		Err: Error{
			TrackerID: trackerId,
		},
	})
}

func ErrorHTTPResponse(rw http.ResponseWriter, code int, message string, err any) error {
	return Render.JSON(rw, code, Response{
		Code:    code,
		Message: message,
		Err:     err,
	})
}

var (
	// Setup Validator
	Validate = validator.New()

	// Setup Validation Message Translation
	enTrans  = en.New()
	uni      = ut.New(enTrans, enTrans)
	Trans, _ = uni.GetTranslator("en")
)

type ValidationError struct {
	Key     string `json:"key"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func NewValidator() {
	// use the names which have been specified for JSON representations of structs, rather than normal Go field names
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Setup Error Message Translation
	en_translations.RegisterDefaultTranslations(Validate, Trans)
}

func Validation(data any) []ValidationError {
	err := Validate.Struct(data)
	if _, ok := err.(*validator.InvalidValidationError); ok {
		panic("error invalid validation error")
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	var (
		length     = len(errs)
		lengthName = 0
	)

	if length > 0 {
		rType := reflect.TypeOf(data)
		if rType.Kind() == reflect.Ptr {
			rType = rType.Elem()
		}
		lengthName = len(rType.Name())
	}

	val := make([]ValidationError, length)
	for i := 0; i < len(errs); i++ {
		// cut string namespace
		// note: add a `+1` to cut dot (.)
		key := errs[i].Namespace()[lengthName+1:]
		val[i] = ValidationError{
			Key:     key,
			Rule:    errs[i].Tag(),
			Message: errs[i].Translate(Trans),
		}
	}

	return val
}

func ResetSlice(slices ...any) {
	for i := 0; i < len(slices); i++ {
		v := reflect.ValueOf(slices[i])

		if v.Kind() != reflect.Ptr || v.IsNil() {
			continue
		} else {
			v = v.Elem()
		}

		if v.Kind() != reflect.Slice {
			continue
		}

		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	}
}

type Query struct {
	val  []string
	args []any
}

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) Where(query string, args ...any) *Query {
	q.val = append(q.val, query)
	q.args = append(q.args, args...)
	return q
}

func (q *Query) Build() (string, []any) {
	if len(q.val) <= 0 {
		return "", nil
	}
	return "WHERE " + strings.Join(q.val, " AND "), q.args
}
