package internalhttp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/otus-murashko/banners-rotation/internal/app"
	"github.com/otus-murashko/banners-rotation/internal/storage"
)

func getBannerRotation(w http.ResponseWriter, r *http.Request, a app.Application) {

	slotID, err := strconv.Atoi(r.URL.Query().Get("slot_id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	banner, err := a.GetBannerRotation(context.Background(), slotID, groupID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := json.Marshal(banner)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func addBannerRotation(w http.ResponseWriter, r *http.Request, a app.Application) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var rotation storage.Rotation
	err = json.Unmarshal(body, &rotation)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = a.AddBannerToSlot(context.Background(), rotation.BannerID, rotation.SlotID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)

}

func deleteBannerRotation(w http.ResponseWriter, r *http.Request, a app.Application) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var rotation storage.Rotation
	err = json.Unmarshal(body, &rotation)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = a.DeleteBannerFromSlot(context.Background(), rotation.BannerID, rotation.SlotID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)

}

func addBanner(w http.ResponseWriter, r *http.Request, a app.Application) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var banner storage.Banner
	err = json.Unmarshal(body, &banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := a.CreateBanner(context.Background(), banner.Descr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	banner.ID = id

	data, err := json.Marshal(banner)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func addSlot(w http.ResponseWriter, r *http.Request, a app.Application) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var slot storage.Slot
	err = json.Unmarshal(body, &slot)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := a.CreateSlot(context.Background(), slot.Descr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	slot.ID = id

	data, err := json.Marshal(slot)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func updateClickStat(w http.ResponseWriter, r *http.Request, a app.Application) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var stat storage.Statistic
	err = json.Unmarshal(body, &stat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = a.UpdateClickStat(context.Background(), stat)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func addGroup(w http.ResponseWriter, r *http.Request, a app.Application) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var group storage.SosialGroup
	err = json.Unmarshal(body, &group)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := a.CreateGroup(context.Background(), group.Descr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	group.ID = id

	data, err := json.Marshal(group)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleNotExpecterRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
}
