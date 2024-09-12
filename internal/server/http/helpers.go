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

// @Summary
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
	handleAddItem(w, r, a, "banner")
}

func addSlot(w http.ResponseWriter, r *http.Request, a app.Application) {
	handleAddItem(w, r, a, "slot")
}

func addGroup(w http.ResponseWriter, r *http.Request, a app.Application) {
	handleAddItem(w, r, a, "group")
}

func handleAddItem(w http.ResponseWriter, r *http.Request, a app.Application, itemType string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var item interface{}

	switch itemType {
	case "slot":
		item, err = addSlotToDB(w, a, body)
	case "group":
		item, err = addGroupToDB(w, a, body)
	case "banner":
		item, err = addBannerToDB(w, a, body)
	}

	if err != nil {
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func addSlotToDB(w http.ResponseWriter, a app.Application, body []byte) (storage.Slot, error) {
	var slot storage.Slot
	err := json.Unmarshal(body, &slot)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return storage.Slot{}, err
	}

	id, err := a.CreateSlot(context.Background(), slot.Descr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return storage.Slot{}, err
	}
	slot.ID = id
	return slot, nil
}

func addBannerToDB(w http.ResponseWriter, a app.Application, body []byte) (storage.Banner, error) {
	var banner storage.Banner
	err := json.Unmarshal(body, &banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return storage.Banner{}, err
	}

	id, err := a.CreateBanner(context.Background(), banner.Descr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return storage.Banner{}, err
	}
	banner.ID = id
	return banner, nil
}

func addGroupToDB(w http.ResponseWriter, a app.Application, body []byte) (storage.SosialGroup, error) {
	var group storage.SosialGroup
	err := json.Unmarshal(body, &group)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return storage.SosialGroup{}, err
	}

	id, err := a.CreateGroup(context.Background(), group.Descr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return storage.SosialGroup{}, err
	}
	group.ID = id
	return group, nil
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

func handleNotExpecterRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
}
