package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type KindsRepo struct {
    Kinds         map[string]models.Kind
    UpsertError   error
    TrimError     error
    FindError     error
    TrimArguments []interface{}
}

func NewKindsRepo() *KindsRepo {
    return &KindsRepo{
        Kinds:         make(map[string]models.Kind),
        TrimArguments: make([]interface{}, 0),
    }
}

func (fake *KindsRepo) Create(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
    key := kind.ID + kind.ClientID
    if _, ok := fake.Kinds[key]; ok {
        return kind, models.ErrDuplicateRecord{}
    }
    fake.Kinds[key] = kind
    return kind, nil
}

func (fake *KindsRepo) Update(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
    key := kind.ID + kind.ClientID
    fake.Kinds[key] = kind
    return kind, nil
}

func (fake *KindsRepo) Upsert(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
    key := kind.ID + kind.ClientID
    fake.Kinds[key] = kind
    return kind, fake.UpsertError
}

func (fake *KindsRepo) Find(conn models.ConnectionInterface, id, clientID string) (models.Kind, error) {
    key := id + clientID
    if kind, ok := fake.Kinds[key]; ok {
        return kind, fake.FindError
    }
    return models.Kind{}, models.ErrRecordNotFound{}
}

func (fake *KindsRepo) Trim(conn models.ConnectionInterface, clientID string, kindIDs []string) (int, error) {
    fake.TrimArguments = []interface{}{clientID, kindIDs}
    return 0, fake.TrimError
}
