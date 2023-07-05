package events

import (
	"github.com/iris-contrib/middleware/pg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kataras/iris/v12"
)

type EventController struct {
}

func (c *EventController) List(ctx iris.Context) {
	eventRepo := pg.Repository[Event](ctx)
	eventlist, _ := eventRepo.Select(ctx, "SELECT * FROM events")
	ctx.JSON(eventlist)
}

func (c *EventController) Create(ctx iris.Context) {
	var e Event
	err := ctx.ReadJSON(&e)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Create Event deserialization failure").DetailErr(err))
		return
	}

	eventRepo := pg.Repository[Event](ctx)
	errPg := eventRepo.InsertSingle(ctx, e, &e.Id)
	if errPg != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title(errPg.(*pgconn.PgError).Message).Detail(errPg.(*pgconn.PgError).Detail))
		return
	}
	ctx.StatusCode(iris.StatusCreated)
}

func (c *EventController) Update(ctx iris.Context) {
	var e Event
	var errPg error

	err := ctx.ReadJSON(&e)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Update Event deserialization failure").DetailErr(err))
		return
	}
	userRepo := pg.Repository[Event](ctx)
	_, errPg = userRepo.Update(ctx, e)
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
	events := EventController{}

	eventsAPI := app.Party("/events", db)
	{
		// GET: http://localhost:8080/events
		eventsAPI.Get("/", events.List)
		// POST: http://localhost:8080/events
		eventsAPI.Post("/", events.Create)
		// PUT: http://localhost:8080/events
		eventsAPI.Put("/", events.Update)
	}
}
