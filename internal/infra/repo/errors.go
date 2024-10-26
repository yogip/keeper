package repo

import "errors"

var ErrUniqConstrain = errors.New("object already exists")

// todo
var ErrNoMoney = errors.New("not enough money")
var ErrOrderAlreadyRegisteredByUser = errors.New("order already registered by user")
var ErrOrderAlreadyRegisteredByOther = errors.New("order already registered by other user")
