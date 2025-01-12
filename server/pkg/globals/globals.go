package globals

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var Db *mongo.Client
var DbName string
var Log *zap.SugaredLogger
var RootPassword string
var TusUploadDir string
