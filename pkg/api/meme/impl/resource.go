package impl

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"memessy-api/pkg/api"
	m "memessy-api/pkg/api/meme"
	"memessy-api/pkg/rest"
	"net/http"
)

type Resource struct {
	Service rest.Service
	Config  *Config
}

func (res *Resource) Get(w http.ResponseWriter, r *http.Request) {
	memeId := r.Context().Value("memeId").(string)
	if memeId == "" {
		log.Error().Msg("could get memeId from context")
		api.InternalError(w)
		return
	}
	meme, err := res.Service.Retrieve(r.Context(), memeId)
	if err != nil {
		api.InternalError(w)
		return
	}
	schema := m.CreateRetrieveSchema(meme)
	w.WriteHeader(http.StatusOK)
	err = api.RenderJson(schema, w)
	if err != nil {
		api.InternalError(w)
		log.Err(err).Send()
	}
}

func (res *Resource) Update(w http.ResponseWriter, r *http.Request) {
	memeId := r.Context().Value("memeId").(string)
	if memeId == "" {
		log.Error().Msg("could get memeId from context")
		api.InternalError(w)
		return
	}
	updateSchema := m.UpdateSchema{}
	err := json.NewDecoder(r.Body).Decode(&updateSchema)
	if err != nil {
		log.Err(err).Send()
		api.ValidationError(w, "Could not parse request.", nil)
		return
	}
	meme, err := res.Service.Retrieve(r.Context(), memeId)
	if err != nil {
		log.Err(err).Send()
		api.InternalError(w)
		return
	}
	meme.ParsedText = updateSchema.ParsedText
	meme.Description = updateSchema.Description
	meme.Categories = updateSchema.Categories
	meme, err = res.Service.Update(r.Context(), memeId, *meme)
	if err != nil {
		log.Err(err).Send()
		api.InternalError(w)
		return
	}
	schema := m.CreateRetrieveSchema(meme)
	w.WriteHeader(http.StatusOK)
	err = api.RenderJson(schema, w)
	if err != nil {
		api.InternalError(w)
		log.Err(err).Send()
	}
}

func (res *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	memeId := r.Context().Value("memeId").(string)
	if memeId == "" {
		log.Error().Msg("could get memeId from context")
		api.InternalError(w)
		return
	}
	err := res.Service.Delete(r.Context(), memeId)
	if err != nil {
		log.Err(err).Send()
		api.InternalError(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (res *Resource) List(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get(res.Config.SearchQueryKey)
	memes, err := res.Service.List(r.Context(), searchQuery)
	if err != nil {
		api.InternalError(w)
		log.Error().Err(err).Send()
		return
	}
	schemas := m.CreateListSchema(memes)
	w.Header().Add("content-type", "application/json")
	err = json.NewEncoder(w).Encode(schemas)
	if err != nil {
		api.InternalError(w)
		log.Error().Err(err).Send()
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (res *Resource) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(res.Config.FileMaxSize)
	if err != nil {
		api.ValidationError(
			w,
			"file size is too big",
			map[string]int64{"max_size": res.Config.FileMaxSize},
		)
		return
	}
	file, header, err := r.FormFile(res.Config.FileFormKey)
	if err != nil {
		api.ValidationError(
			w,
			"not found file in payload",
			map[string]string{"file_field": res.Config.FileFormKey},
		)
		return
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		api.ValidationError(w, "could not read file", nil)
		return
	}
	log.Debug().Msg(fmt.Sprintf("got file %s", header.Filename))

	meme, err := res.Service.Create(r.Context(), header.Filename, bytes)
	if err != nil {
		api.InternalError(w)
		log.Err(err).Send()
		return
	}
	schema := m.CreateRetrieveSchema(meme)
	w.WriteHeader(http.StatusCreated)
	err = api.RenderJson(schema, w)
	if err != nil {
		api.InternalError(w)
		log.Err(err).Send()
	}
}
