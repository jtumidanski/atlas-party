package character

import (
	"atlas-party/json"
	"atlas-party/party"
	"atlas-party/rest"
	"atlas-party/rest/response"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	getCharacterParty = "get_character_party"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	r := router.PathPrefix("/characters").Subrouter()
	r.HandleFunc("/{id}/party", registerGetCharacterParty(l)).Methods(http.MethodGet)
}

func registerGetCharacterParty(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(getCharacterParty, func(span opentracing.Span) http.HandlerFunc {
		return parseCharacterId(l, func(characterId uint32) http.HandlerFunc {
			return handleGetCharacterParty(l)(span)(characterId)
		})
	})
}

type characterIdHandler func(characterId uint32) http.HandlerFunc

func parseCharacterId(l logrus.FieldLogger, next characterIdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		characterId, err := strconv.Atoi(vars["id"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing characterId as uint32")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(characterId))(w, r)
	}
}

func handleGetCharacterParty(l logrus.FieldLogger) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				p, err := party.GetByMember(l, span)(characterId)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				result := response.NewDataContainer(true)
				result.AddData(p.Id(), "parties", party.MakeAttribute(p), party.MakePartyRelationships(p))

				err = json.ToJSON(result, w)
				if err != nil {
					l.WithError(err).Errorf("Encoding response")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}
