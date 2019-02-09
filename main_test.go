package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePonto(t *testing.T) {
	ts := []struct {
		name     string
		mPath    string
		response string
		expected interface{}
		err      error
	}{
		{
			name:     "1",
			mPath:    "matriculas",
			response: `{"result":true,"time":"0931","day":"2018-11-30","batidas_dia":["0928","0931"],"nome":"IVANHENRIQUESANDRINI","employee":{"_id":"5bd05fc65b681ba7dbd65ef3"},"only_location":false,"photo_on_punch":false,"activity_on_punch":false,"justification_permissions":{"read_write_attach":true,"add_absence":true,"add_punch":true},"face_id_on_punch":false}`,
			expected: []string{"Ponto batido com SUCESSO. matrícula: 411", "Ponto batido com SUCESSO. matrícula: 412", "Ponto batido com SUCESSO. matrícula: 413"},
			err:      nil,
		},
	}

	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.RequestURI {
					case "/batidaonline/verifyIdentification":
						w.Write([]byte(tc.response))
					case "/error":
						http.Error(w, tc.response, http.StatusUnprocessableEntity)
					}

				}),
			)
			defer ts.Close()

			rec := httptest.NewRecorder()
			req := &http.Request{}

			s, _ := NewHandler(ts.URL, tc.mPath).Hponto(rec, req)
			assert.Equal(t, tc.expected, s)
		})
	}

}
