package router

import (
	"cloudace/progcon/test/gae/handler/app"
	"net/http"
	"prompter/registry/util"

	"github.com/gorilla/mux"
)

// Route ルーター構造体
type Route struct {
	// Name 名前
	Name string
	// Method HTTPメソッド
	Method string
	// Pattern URLパターン
	Pattern string
	// HandlerFunc 実行される関数
	HandlerFunc http.HandlerFunc
}

// Routes ルーター
type Routes []Route

var routes = Routes{
	// ヘルスチェック用API
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	// test
	Route{
		"test",
		"GET",
		"/test",
		app.TestHandle,
	},

	// test
	Route{
		"testput",
		"PUT",
		"/test",
		app.TestPut,
	},

	Route{
		"testlog",
		"GET",
		"/test/log",
		app.Logger,
	},

	Route{
		"testMirror",
		"POST",
		"/api/v1/search",
		app.TestMirror,
	},

	Route{
		"testCheck",
		"GET",
		"/check",
		app.CheckRequest,
	},
}

// NewRouter コンストラクタ　routesで定義したroute構造体配列を用いて、muxのルーターを作成
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = util.Logger(route.HandlerFunc, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

// Index ルート
func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type stackdriver struct {
	Severity string `json:"severity"`
	Msg      string `json:"msg"`
}
