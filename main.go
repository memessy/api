package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"memessy-api/pkg/api"
	memeApi "memessy-api/pkg/api/meme/impl"
	"memessy-api/pkg/bus"
	busImpl "memessy-api/pkg/bus/memory"
	fileserver "memessy-api/pkg/fileserver/disk"
	"memessy-api/pkg/recognizer"
	recognizerImpl "memessy-api/pkg/recognizer/memessy"
	rest "memessy-api/pkg/rest/impl"
	storage "memessy-api/pkg/storage/mongo"
	"net/http"
	"os"
	"time"
)

type config struct {
	MongoUrl      string
	MongoDb       string
	MongoMemesCol string
	Port          string
	StaticDir     string
	StaticPrefix     string
	StaticBaseUrl string
	RecognizerUrl string
}

func parseConfig() config {
	return config{
		MongoUrl:      os.Getenv("MONGO_URL"),
		MongoDb:       os.Getenv("MONGO_DB"),
		MongoMemesCol: os.Getenv("MONGO_MEMES_COLLECTION"),
		Port:          os.Getenv("PORT"),
		StaticDir:     os.Getenv("STATIC_DIR"),
		StaticPrefix: os.Getenv("STATIC_PREFIX"),
		StaticBaseUrl: os.Getenv("STATIC_BASE_URL"),
		RecognizerUrl: os.Getenv("RECOGNIZER_URL"),
	}
}

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.SyncWriter(os.Stdout)

	conf := parseConfig()

	client, err := mongo.NewClient(options.Client().ApplyURI(conf.MongoUrl))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to mongo")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	memeStorage := storage.MemeStorage{Collection: client.Database(conf.MongoDb).Collection(conf.MongoMemesCol)}
	err = memeStorage.Init()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	fileServer := fileserver.FileServer{
		BaseUrl: conf.StaticBaseUrl,
		Dir: conf.StaticDir,
	}
	eventBus := busImpl.EventBus{}
	defer eventBus.Close()
	recognizerConsumer := recognizer.Consumer{
		Recognizer: &recognizerImpl.Recognizer{
			FileField: "file",
			Url:       conf.RecognizerUrl,
		},
		FileServer: &fileServer,
		Storage:    &memeStorage,
	}
	go bus.ConsumeCreated(&eventBus, recognizerConsumer.Consume)
	memeService := rest.NewService(&memeStorage, &fileServer, &eventBus)
	memeResource := memeApi.Resource{
		Service: memeService,
		Config: &memeApi.Config{
			FileMaxSize:    10 << 20,
			FileFormKey:    "file",
			SearchQueryKey: "q",
		},
	}

	apiApp := api.NewApi(
		api.Config{
			StaticDir:    http.Dir(conf.StaticDir),
			StaticPrefix: conf.StaticPrefix,
		},
		&memeResource,
	)
	log.Fatal().Err(http.ListenAndServe(":"+conf.Port, apiApp))
	defer client.Disconnect(ctx)
}
