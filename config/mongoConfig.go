package config

import (
    "os"
)

type MongoConfig struct {
    Uri string
}

func LoadMongoConfig() *MongoConfig {
    return &MongoConfig {
        Uri: getEnv("SG_MONGO", "mongodb://softgen:QWAmTnsdBUaTL2z@118.139.167.125:27017/"),
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
