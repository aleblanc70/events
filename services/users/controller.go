package users

import (
	"github.com/iris-contrib/middleware/pg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kataras/iris/v12"
)

type UserController struct {
}

func (c UserController) List(ctx iris.Context) {
	users := pg.Repository[User](ctx)
	userlist, _ := users.Select(ctx, "SELECT * FROM users")
	ctx.JSON(userlist)
}

func (c UserController) Create(ctx iris.Context) {
	var u User
	err := ctx.ReadJSON(&u)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("User deserialization failure").DetailErr(err))
		return
	}

	users := pg.Repository[User](ctx)
	errPg := users.InsertSingle(ctx, u, &u.Id)
	if errPg != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title(errPg.(*pgconn.PgError).Message).Detail(errPg.(*pgconn.PgError).Detail))
		return
	}
	ctx.StatusCode(iris.StatusCreated)
}

func (c UserController) Update(ctx iris.Context) {
	var u User
	var errPg error
	err := ctx.ReadJSON(&u)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("User Update failure").DetailErr(err))
		return
	}

	userRepo := pg.Repository[User](ctx)
	_, errPg = userRepo.Update(ctx, u)
	if errPg != nil {
		if pg.IsErrNoRows(errPg) {
			ctx.StopWithStatus(iris.StatusNotFound)
		} else {
			ctx.StopWithError(iris.StatusInternalServerError, errPg)
		}
		return
	}
	ctx.StatusCode(iris.StatusOK)
}

func InitRoutes(app *iris.Application, db iris.Handler) {
	users := UserController{}
	usersAPI := app.Party("/users", db)
	{
		// GET: http://localhost:8080/users
		usersAPI.Get("/", users.List)
		// POST: http://localhost:8080/users
		usersAPI.Post("/", users.Create)
		// PUT: http://localhost:8080/users
		usersAPI.Put("/", users.Update)
	}
}
