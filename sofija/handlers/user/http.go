package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"
	"github.com/Bloxico/exchange-gateway/sofija/core/ports"
	"github.com/Bloxico/exchange-gateway/sofija/server/auth"
	"github.com/Bloxico/exchange-gateway/sofija/server/params"
	"github.com/emicklei/go-restful/v3"
	"golang.org/x/crypto/bcrypt"
)

type EgwUserHttpHandler struct {
	userSvc ports.EgwUserUsecase
}

func NewEgwUserHandler(userSvc ports.EgwUserUsecase, wsCont *restful.Container) *EgwUserHttpHandler {
	httpHandler := &EgwUserHttpHandler{
		userSvc: userSvc,
	}

	ws := new(restful.WebService)

	ws.Path("/user").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/register").To(httpHandler.RegisterUser))
	ws.Route(ws.POST("/login").To(httpHandler.LoginUser))
	ws.Route(ws.PUT("").To(httpHandler.UpdateUser).Filter(auth.AuthJWT))

	wsCont.Add(ws)

	return httpHandler
}

func (e *EgwUserHttpHandler) UpdateUser(req *restful.Request, resp *restful.Response) {
	var a UpdateRequestData
	req.ReadEntity(&a)

	// get user ID for update query
	reqId, err := params.StringFrom(req.Request, auth.USER_ID_CTX_KEY)
	if err != nil || len(reqId) == 0 {
		resp.WriteError(http.StatusBadRequest, errors.New("no id found for user"))
		return
	}

	ctx := req.Request.Context()
	dataUser := &domain.EgwUser{ID: reqId, Name: a.Name, Surname: a.Surname}

	err = e.userSvc.Update(ctx, dataUser)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error updating user"))
		return
	}

	// return updated user as data
	var retUser *EgwUserModel = &EgwUserModel{}
	retUser.FromDomain(dataUser)
	resp.WriteAsJson(retUser)
}

// Performs login or register
func (e *EgwUserHttpHandler) RegisterUser(req *restful.Request, resp *restful.Response) {
	var reqData RegisterRequestData
	req.ReadEntity(&reqData)

	var egwUser *EgwUserModel = &EgwUserModel{}
	egwUser.Email = reqData.Email
	egwUser.Name = reqData.Name
	egwUser.Surname = reqData.Surname

	// todo: expand validation
	if len(egwUser.Email) == 0 {
		resp.WriteError(http.StatusBadRequest, errors.New("no email provided"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), 10)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error registering"))
		return
	}
	egwUser.PasswordHash = string(hashedPassword)

	err = e.registerUser(req.Request.Context(), egwUser)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error registering"))
		return
	}

	authToken, err := auth.CreateJWT(egwUser.Email, egwUser.ID)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error creating jwt"))
		return
	}

	// send user + token back
	respData := RegisterResponseData{AuthToken: authToken, User: *egwUser}

	resp.WriteAsJson(respData)
}

func (e *EgwUserHttpHandler) registerUser(ctx context.Context, egwUser *EgwUserModel) error {
	err := e.userSvc.RegisterUser(ctx, egwUser.ToDomain())
	if err != nil {
		return err
	}
	// retrieve their data from the DB to populate it (e.g. ID)
	userData, err := e.userSvc.FindByEmail(ctx, egwUser.Email)
	if err != nil {
		return err
	}
	egwUser.FromDomain(userData)
	return nil
}

func (e *EgwUserHttpHandler) LoginUser(req *restful.Request, resp *restful.Response) {
	var reqData LoginRequestData
	req.ReadEntity(&reqData)

	if len(reqData.Email) == 0 || len(reqData.Password) == 0 {
		resp.WriteError(http.StatusBadRequest, errors.New("bad login credentials"))
		return
	}

	userData, err := e.userSvc.FindByEmail(req.Request.Context(), reqData.Email)
	if err != nil {
		resp.WriteError(http.StatusForbidden, errors.New("unauthorized"))
		return
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(userData.PasswordHash), []byte(reqData.Password))
	if err != nil {
		resp.WriteError(http.StatusForbidden, errors.New("unauthorized"))
		return
	}

	authToken, err := auth.CreateJWT(userData.Email, userData.ID)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error creating jwt"))
		return
	}

	// send user + token back
	var egwUser *EgwUserModel = &EgwUserModel{}
	egwUser.FromDomain(userData)
	respData := RegisterResponseData{AuthToken: authToken, User: *egwUser}

	resp.WriteAsJson(respData)

}
