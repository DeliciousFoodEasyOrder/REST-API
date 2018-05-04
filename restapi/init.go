package restapi

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Resp is the general response of this api
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Err will be sent as data when error occurs
type Err struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

var (
	formatter     *render.Render
	jwtMiddleware *jwtmiddleware.JWTMiddleware
)

func init() {
	formatter = render.New(render.Options{
		IndentJSON: true,
	})

	jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}

// NewResp returns a new instance of Struct Resp
func NewResp(code int, msg string, data interface{}) *Resp {
	return &Resp{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// NewErr returns a new instance of Struct Err
func NewErr(err, errDes string) *Err {
	return &Err{
		Error:            err,
		ErrorDescription: errDes,
	}
}

// NewServer returns a new Negroni Server
func NewServer() *negroni.Negroni {
	n := negroni.Classic()
	router := mux.NewRouter()

	initRouter(router)

	n.UseHandler(router)

	return n
}

func initRouter(router *mux.Router) {

	// ## JWT Collection
	routeJWTCollection(router)

	// ## Merchant Collection
	// routeMerchantCollection(router);

	// ## Order Collection
	// routeOrderCollection(router);

	// ## Seat Collection
	// routeSeatCollection(router);

	// ## Food Collection
	// routeFoodCollection(router)

}
