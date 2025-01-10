package globals

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var Db *mongo.Client
var DbName string
var Log *zap.SugaredLogger
var GraphDB neo4j.DriverWithContext
