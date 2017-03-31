package main

import (

	// "github.com/go-redis/redis"
	"fmt"

	"github.com/attilasatan/model"
	"github.com/mediocregopher/radix.v2/redis"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
)

func main() {
	// Receives optional iris.Configuration{}, see ./configuration.go
	// for more.
	app := iris.New()

	// Order doesn't matter,
	// You can split it to different .Adapt calls.
	// See ./adaptors folder for more.
	app.Adapt(
		// adapt a logger which prints all errors to the os.Stdout
		iris.DevLogger(),
		// adapt the adaptors/httprouter or adaptors/gorillamux
		httprouter.New(),
		// 5 template engines are supported out-of-the-box:
		//
		// - standard html/template
		// - amber
		// - django
		// - handlebars
		// - pug(jade)
		//
		// Use the html standard engine for all files inside "./views" folder with extension ".html"
		view.HTML("./views", ".html").Reload(true),
		// Cors wrapper to the entire application, allow all origins.
		cors.New(cors.Options{AllowedOrigins: []string{"*"}}))

	// http://localhost:6300
	// Method: "GET"
	// Render ./views/index.html
	app.Get("/", func(ctx *iris.Context) {
		ctx.Render("index.html", iris.Map{"Title": "Page Title"}, iris.RenderOptions{"gzip": true})
	})

	app.Get("/foobar", CheckRedix)

	// Group routes, optionally: share middleware, template layout and custom http errors.
	userAPI := app.Party("/users", userAPIMiddleware).
		Layout("layouts/userLayout.html")
	{
		// Fire userNotFoundHandler when Not Found
		// inside http://localhost:6300/users/*anything
		userAPI.OnError(404, userNotFoundHandler)

		// http://localhost:6300/users
		// Method: "GET"
		userAPI.Get("/", getAllHandler)

		// http://localhost:6300/users/42
		// Method: "GET"
		userAPI.Get("/:id", getByIDHandler)

		// http://localhost:6300/users
		// Method: "POST"
		userAPI.Post("/", saveUserHandler)
	}

	// Start the server at 127.0.0.1:6300
	app.Listen(":6300")
}

func userAPIMiddleware(ctx *iris.Context) {
	// your code here...
	println("Request: " + ctx.Path())
	ctx.Next() // go to the next handler(s)
}

func userNotFoundHandler(ctx *iris.Context) {
	// your code here...
	ctx.HTML(iris.StatusNotFound, "<h1> User page not found </h1>")
}

func getAllHandler(ctx *iris.Context) {
	// your code here...
}

func getByIDHandler(ctx *iris.Context) {
	// take the :id from the path, parse to integer
	// and set it to the new userID local variable.
	userID, _ := ctx.ParamInt("id")
	type User struct {
		Name string
		Job  string
		ID   int
	}

	user := User{"Çiğdem", "Aşk", userID}

	err := model.Save(user)

	if err != nil {
		panic(err)
	}

	// send back a response to the client,
	// .JSON: content type as application/json; charset="utf-8"
	// iris.StatusOK: with 200 http status code.
	//
	// send user as it is or make use of any json valid golang type,
	// like the iris.Map{"username" : user.Username}.

	ctx.JSON(iris.StatusOK, user)
}

func saveUserHandler(ctx *iris.Context) {
	// your code here...
}

// func CheckRedis() {
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	pong, err := client.Ping().Result()
// 	fmt.Println(pong, err)
// 	// Output: PONG <nil>
// }

func CheckRedix(ctx *iris.Context) {
	client, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		// handle err
	}
	rSetFoo := client.Cmd("SET", "foo", "foo val")
	if rSetFoo.Err != nil {
		panic(rSetFoo.Err)
	}

	rSetBar := client.Cmd("SET", "bar", "bar val")
	if rSetBar.Err != nil {
		panic(rSetBar.Err)
	}

	rSetBaz := client.Cmd("SET", "baz", "baz val")
	if rSetBaz.Err != nil {
		panic(rSetBaz.Err)
	}
	mGetR := client.Cmd("MGET", "foo", "bar", "baz")

	if mGetR.Err != nil {
		// handle error
	}

	// This:
	l, _ := mGetR.List()
	for _, elemStr := range l {
		fmt.Println(elemStr)
	}

	// is equivalent to this:
	elems, err := mGetR.Array()
	for i := range elems {
		elemStr, _ := elems[i].Str()
		fmt.Println(elemStr)
	}
}
