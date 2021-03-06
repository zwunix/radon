/*
 * Radon
 *
 * Copyright 2018 The Radon Authors.
 * Code is licensed under the GPLv3.
 *
 */

package v1

import (
	"errors"
	"testing"

	"proxy"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
	querypb "github.com/xelabs/go-mysqlstack/sqlparser/depends/query"
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/sqltypes"
	"github.com/xelabs/go-mysqlstack/xlog"
)

func TestCtlV1CreateUser(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	// fakedbs.
	{
		fakedbs.AddQuery("GRANT SELECT ON *.* TO 'mock'@'localhost' IDENTIFIED BY 'pwd'", &sqltypes.Result{})
	}

	// server
	api := rest.NewApi()
	router, _ := rest.MakeRouter(
		rest.Post("/v1/user/add", CreateUserHandler(log, proxy)),
	)
	api.SetApp(router)
	handler := api.MakeHandler()

	{
		p := &userParams{
			User:     "mock",
			Password: "pwd",
		}
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/user/add", p))
		recorded.CodeIs(200)
	}
}

func TestCtlV1CreateUserError(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	// fakedbs.
	{
		fakedbs.AddQueryError("GRANT SELECT ON *.* TO 'mock'@'localhost' IDENTIFIED BY 'pwd'", errors.New("mock.create.user.error"))
	}

	// server
	api := rest.NewApi()
	router, _ := rest.MakeRouter(
		rest.Post("/v1/user/add", CreateUserHandler(log, proxy)),
	)
	api.SetApp(router)
	handler := api.MakeHandler()

	{
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/user/add", nil))
		recorded.CodeIs(500)
	}

	{
		p := &userParams{
			User:     "mock",
			Password: "pwd",
		}
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/user/add", p))
		recorded.CodeIs(503)
	}
}

func TestCtlV1AlterUser(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	// fakedbs.
	{
		fakedbs.AddQuery("ALTER USER 'mock'@'localhost' IDENTIFIED BY 'pwd'", &sqltypes.Result{})
	}

	// server
	api := rest.NewApi()
	router, _ := rest.MakeRouter(
		rest.Post("/v1/user/update", AlterUserHandler(log, proxy)),
	)
	api.SetApp(router)
	handler := api.MakeHandler()

	{
		p := &userParams{
			User:     "mock",
			Password: "pwd",
		}
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/user/update", p))
		recorded.CodeIs(200)
	}
}

func TestCtlV1AlterUserError(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	// fakedbs.
	{
		fakedbs.AddQueryError("ALTER USER 'mock'@'localhost' IDENTIFIED BY 'pwd'", errors.New("mock.alter.user.error"))
	}

	// server
	api := rest.NewApi()
	router, _ := rest.MakeRouter(
		rest.Post("/v1/user/update", AlterUserHandler(log, proxy)),
	)
	api.SetApp(router)
	handler := api.MakeHandler()

	// 500.
	{
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/user/update", nil))
		recorded.CodeIs(500)
	}

	// 503.
	{
		p := &userParams{
			User:     "mock",
			Password: "pwd",
		}
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/user/update", p))
		recorded.CodeIs(503)
	}
}

func TestCtlV1DropUser(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	// fakedbs.
	{
		fakedbs.AddQueryPattern("DROP USER 'mock'@'localhost'", &sqltypes.Result{})
	}

	// server
	api := rest.NewApi()
	router, _ := rest.MakeRouter(
		rest.Post("/v1/user/remove", DropUserHandler(log, proxy)),
	)
	api.SetApp(router)
	handler := api.MakeHandler()

	{
		p := &userParams{
			User: "mock",
		}
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/user/remove", p))
		recorded.CodeIs(200)
	}
}

func TestCtlV1DropError(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	// fakedbs.
	{
		fakedbs.AddQueryErrorPattern("DROP .*", errors.New("mock.drop.user.error"))
	}

	// server
	api := rest.NewApi()
	router, _ := rest.MakeRouter(
		rest.Post("/v1/user/remove", DropUserHandler(log, proxy)),
	)
	api.SetApp(router)
	handler := api.MakeHandler()

	// 503.
	{
		p := &userParams{
			User: "mock",
		}
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/user/remove", p))
		recorded.CodeIs(503)
	}
}

func TestCtlV1Userz(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	selectResult := &sqltypes.Result{
		Fields: []*querypb.Field{
			{
				Name: "User",
				Type: querypb.Type_VARCHAR,
			},
			{
				Name: "Host",
				Type: querypb.Type_VARCHAR,
			},
		},
		Rows: [][]sqltypes.Value{
			{
				sqltypes.MakeTrusted(querypb.Type_VARCHAR, []byte("test1")),
				sqltypes.MakeTrusted(querypb.Type_VARCHAR, []byte("localhost")),
			},
			{
				sqltypes.MakeTrusted(querypb.Type_VARCHAR, []byte("test2")),
				sqltypes.MakeTrusted(querypb.Type_VARCHAR, []byte("localhost")),
			},
		},
	}

	//fakedbs.
	{
		fakedbs.AddQueryPattern("select .*", selectResult)
	}

	// server
	api := rest.NewApi()
	router, _ := rest.MakeRouter(
		rest.Get("/v1/user/userz", UserzHandler(log, proxy)),
	)
	api.SetApp(router)
	handler := api.MakeHandler()

	{
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/v1/user/userz", nil))
		recorded.CodeIs(200)

		want := "[{\"User\":\"test1\",\"Host\":\"localhost\"},{\"User\":\"test2\",\"Host\":\"localhost\"}]"
		got := recorded.Recorder.Body.String()
		log.Debug(got)
		assert.Equal(t, want, got)
	}
}

func TestCtlV1UserzError(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	// fakedbs.
	{
		fakedbs.AddQueryErrorPattern("select .*", errors.New("api.v1.userz.get.mysql.user.error"))
	}

	// server
	api := rest.NewApi()
	router, _ := rest.MakeRouter(
		rest.Get("/v1/user/userz", UserzHandler(log, proxy)),
	)
	api.SetApp(router)
	handler := api.MakeHandler()

	{
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/v1/user/userz", nil))
		recorded.CodeIs(503)
	}
}
