package party

import (
	"atlas-party/json"
	"atlas-party/party/member"
	"atlas-party/rest"
	"atlas-party/rest/resource"
	"atlas-party/rest/response"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

const (
	CreateParty   = "create_party"
	GetAllParties = "get_all_parties"
	GetParty      = "get_party"
	GetMembers    = "get_members"
	JoinParty     = "join_party"
	LeaveParty    = "leave_party"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	r := router.PathPrefix("/parties").Subrouter()
	r.HandleFunc("", registerGetAllParties(l)).Methods(http.MethodGet).Queries("include", "{include}")
	r.HandleFunc("", registerGetAllParties(l)).Methods(http.MethodGet)
	r.HandleFunc("", registerCreateParty(l)).Methods(http.MethodPost)
	r.HandleFunc("/{id}", registerGetParty(l)).Methods(http.MethodGet).Queries("include", "{include}")
	r.HandleFunc("/{id}", registerGetParty(l)).Methods(http.MethodGet)
	r.HandleFunc("/{id}/members", registerGetMembers(l)).Methods(http.MethodGet)
	r.HandleFunc("/{id}/members", registerJoinParty(l)).Methods(http.MethodPut)
	r.HandleFunc("/{id}/members/{memberId}", registerLeaveParty(l)).Methods(http.MethodDelete)
}

func registerLeaveParty(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(LeaveParty, func(span opentracing.Span) http.HandlerFunc {
		return ParseId(l, func(partyId uint32) http.HandlerFunc {
			return ParseMemberId(l, func(memberId uint32) http.HandlerFunc {
				return ParseCreateInput(l, func(input *inputDataContainer) http.HandlerFunc {
					return handleLeaveParty(l)(span)(partyId)(memberId)(input)
				})
			})
		})
	})
}

func handleLeaveParty(l logrus.FieldLogger) func(span opentracing.Span) func(partyId uint32) func(memberId uint32) func(input *inputDataContainer) http.HandlerFunc {
	return func(span opentracing.Span) func(partyId uint32) func(memberId uint32) func(input *inputDataContainer) http.HandlerFunc {
		return func(partyId uint32) func(memberId uint32) func(input *inputDataContainer) http.HandlerFunc {
			return func(memberId uint32) func(input *inputDataContainer) http.HandlerFunc {
				return func(input *inputDataContainer) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						emitLeaveParty(l, span)(input.Data.Attributes.WorldId, input.Data.Attributes.ChannelId, memberId)
						w.WriteHeader(http.StatusAccepted)
					}
				}
			}
		}
	}
}

type MemberIdHandler func(memberId uint32) http.HandlerFunc

func ParseMemberId(l logrus.FieldLogger, next MemberIdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memberId, err := strconv.Atoi(mux.Vars(r)["memberId"])
		if err != nil {
			l.WithError(err).Errorf("Unable to properly parse memberId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(memberId))(w, r)
	}
}

func registerJoinParty(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(JoinParty, func(span opentracing.Span) http.HandlerFunc {
		return ParseId(l, func(partyId uint32) http.HandlerFunc {
			return handleJoinParty(l)(span)(partyId)
		})
	})
}

type IdHandler func(partyId uint32) http.HandlerFunc

func ParseId(l logrus.FieldLogger, next IdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			l.WithError(err).Errorf("Unable to properly parse id from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(id))(w, r)
	}
}

func registerCreateParty(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(CreateParty, func(span opentracing.Span) http.HandlerFunc {
		return ParseCreateInput(l, func(input *inputDataContainer) http.HandlerFunc {
			return handleCreateParty(l)(span)(input)
		})
	})
}

type CreateInputHandler func(input *inputDataContainer) http.HandlerFunc

func ParseCreateInput(l logrus.FieldLogger, next CreateInputHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i := &inputDataContainer{}
		err := json.FromJSON(i, r.Body)
		if err != nil {
			l.WithError(err).Errorf("Deserializing input")
			w.WriteHeader(http.StatusBadRequest)
			err := json.ToJSON(&resource.GenericError{Message: err.Error()}, w)
			if err != nil {
				l.WithError(err).Errorf("Unable to serialize error mesage")
			}
			return
		}
		next(i)(w, r)
	}
}

func handleCreateParty(l logrus.FieldLogger) func(span opentracing.Span) func(input *inputDataContainer) http.HandlerFunc {
	return func(span opentracing.Span) func(input *inputDataContainer) http.HandlerFunc {
		return func(input *inputDataContainer) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				emitCreateParty(l, span)(input.Data.Attributes.WorldId, input.Data.Attributes.ChannelId, input.Data.Attributes.CharacterId)
				w.WriteHeader(http.StatusAccepted)
			}
		}
	}
}

func registerGetAllParties(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(GetAllParties, handleGetAllParties(l))
}

func handleGetAllParties(l logrus.FieldLogger) rest.SpanHandler {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ps := GetAll()

			result := response.NewDataContainer(false)
			for _, p := range ps {
				result.AddData(p.Id(), "parties", MakeAttribute(p), MakePartyRelationships(p))
				if strings.Contains(mux.Vars(r)["include"], "members") {
					for _, m := range p.Members() {
						result.AddIncluded(m.Id(), "members", member.MakeAttribute(m))
					}
				}
			}

			w.WriteHeader(http.StatusOK)
			err := json.ToJSON(result, w)
			if err != nil {
				l.WithError(err).Errorf("Encoding response")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func registerGetParty(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(GetParty, func(span opentracing.Span) http.HandlerFunc {
		return ParseId(l, func(partyId uint32) http.HandlerFunc {
			return handleGetParty(l)(span)(partyId)
		})
	})
}

func handleGetParty(l logrus.FieldLogger) func(span opentracing.Span) func(partyId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(partyId uint32) http.HandlerFunc {
		return func(partyId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				p, err := GetById(l, span)(partyId)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				result := response.NewDataContainer(true)
				result.AddData(p.Id(), "parties", MakeAttribute(p), MakePartyRelationships(p))
				if strings.Contains(mux.Vars(r)["include"], "members") {
					for _, m := range p.Members() {
						result.AddIncluded(m.Id(), "members", member.MakeAttribute(m))
					}
				}

				err = json.ToJSON(result, w)
				if err != nil {
					l.WithError(err).Errorf("Encoding response")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

func MakePartyRelationships(p Model) map[string]*response.Relationship {
	result := make(map[string]*response.Relationship, 0)
	result["members"] = &response.Relationship{
		ToOneType: false,
		Links: response.RelationshipLinks{
			Self:    "/ms/party/parties/" + strconv.Itoa(int(p.Id())) + "/relationships/members",
			Related: "/ms/party/parties/" + strconv.Itoa(int(p.Id())) + "/members",
		},
		Data: makeMemberRelationshipData(p.Members()),
	}
	return result
}

func makeMemberRelationshipData(members []*member.Model) []response.RelationshipData {
	result := make([]response.RelationshipData, 0)
	for _, m := range members {
		result = append(result, response.RelationshipData{
			Type: "members",
			Id:   strconv.Itoa(int(m.Id())),
		})
	}
	return result
}

func registerGetMembers(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(GetMembers, func(span opentracing.Span) http.HandlerFunc {
		return ParseId(l, func(partyId uint32) http.HandlerFunc {
			return handleGetMembers(l)(span)(partyId)
		})
	})
}

func handleGetMembers(l logrus.FieldLogger) func(span opentracing.Span) func(partyId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(partyId uint32) http.HandlerFunc {
		return func(partyId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				p, err := GetById(l, span)(partyId)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				result := response.NewDataContainer(false)
				for _, m := range p.Members() {
					result.AddData(m.Id(), "members", member.MakeAttribute(m), nil)
				}

				err = json.ToJSON(result, w)
				if err != nil {
					l.WithError(err).Errorf("Encoding response")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

func handleJoinParty(l logrus.FieldLogger) func(span opentracing.Span) func(partyId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(partyId uint32) http.HandlerFunc {
		return func(partyId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				li := &member.InputDataContainer{}
				err := json.FromJSON(li, r.Body)
				if err != nil {
					l.WithError(err).Errorf("Deserializing input.")
					w.WriteHeader(http.StatusBadRequest)
					err = json.ToJSON(&resource.GenericError{Message: err.Error()}, w)
					if err != nil {
						l.WithError(err).Fatalf("Writing error message.")
					}
					return
				}

				attr := li.Data.Attributes
				emitJoinParty(l, span)(attr.WorldId, attr.ChannelId, partyId, attr.CharacterId)
				w.WriteHeader(http.StatusAccepted)
			}
		}
	}
}
