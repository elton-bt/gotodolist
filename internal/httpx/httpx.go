package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/elton-bt/gotodolist/internal/todo"
)

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func DecodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if decoder.More() {
		return fmt.Errorf("request body must contain a single JSON document")
	}

	return nil
}

func ErrorStatusAndMessage(err error) (int, string) {
	switch {
	case errors.Is(err, todo.ErrInvalidTitle):
		return http.StatusBadRequest, "Informe um titulo para a tarefa."
	case errors.Is(err, todo.ErrNotFound):
		return http.StatusNotFound, "Tarefa nao encontrada."
	case errors.Is(err, todo.ErrUnavailable):
		return http.StatusServiceUnavailable, "Servico temporariamente indisponivel. Verifique a conexao com o banco."
	default:
		return http.StatusInternalServerError, "Nao foi possivel concluir a operacao."
	}
}
