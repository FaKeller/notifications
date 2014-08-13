package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/dgrijalva/jwt-go"
)

type NotifyUser struct {
    courier     postal.CourierInterface
    errorWriter ErrorWriterInterface
    clientsRepo models.ClientsRepoInterface
    kindsRepo   models.KindsRepoInterface
}

func NewNotifyUser(courier postal.CourierInterface, errorWriter ErrorWriterInterface,
    clientsRepo models.ClientsRepoInterface, kindsRepo models.KindsRepoInterface) NotifyUser {

    return NotifyUser{
        courier:     courier,
        errorWriter: errorWriter,
        clientsRepo: clientsRepo,
        kindsRepo:   kindsRepo,
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    userGUID := strings.TrimPrefix(req.URL.Path, "/users/")

    params, err := NewNotifyParams(req.Body)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    if !params.Validate() {
        handler.errorWriter.Write(w, ParamsValidationError(params.Errors))
        return
    }

    rawToken := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

    token, err := jwt.Parse(rawToken, func(token *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    clientID := token.Claims["client_id"].(string)
    client, err := handler.FindClient(clientID)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    kind, err := handler.FindKind(params.KindID, clientID)
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    responses, err := handler.courier.Dispatch(rawToken, postal.UserGUID(userGUID), params.ToOptions(client, kind))
    if err != nil {
        handler.errorWriter.Write(w, err)
        return
    }

    output, err := json.Marshal(responses)
    if err != nil {
        panic(err)
    }

    w.WriteHeader(http.StatusOK)
    w.Write(output)
}

func (handler NotifyUser) FindClient(clientID string) (models.Client, error) {
    client, err := handler.clientsRepo.Find(models.Database().Connection, clientID)
    if err != nil {
        if _, ok := err.(models.ErrRecordNotFound); ok {
            return models.Client{}, nil
        } else {
            return models.Client{}, err
        }
    }
    return client, nil
}

func (handler NotifyUser) FindKind(kindID, clientID string) (models.Kind, error) {
    kind, err := handler.kindsRepo.Find(models.Database().Connection, kindID, clientID)
    if err != nil {
        if _, ok := err.(models.ErrRecordNotFound); ok {
            return models.Kind{}, nil
        } else {
            return models.Kind{}, err
        }
    }
    return kind, nil
}
