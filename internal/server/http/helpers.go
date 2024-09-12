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

type Book struct{}

// Get Banner Rotation
// getBannerRotation             godoc
// @Summary      Get best banner for the slot and group (Bandit algorithm)
// @Description  Get Banner Rotation
// @Tags         Banner-Rotation
// @Produce      json
// @Param        slot_id    query     string  true  "slot id"
// @Param        group_id    query     string  true  "slot id"
// @Success      200
// @Router       /banner-rotation [get] .
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

// Add Banner Rotation
// addBannerRotation             godoc
// @Summary      Add Banner to a Slot
// @Description  Add Banner Rotation
// @Tags         Banner-Rotation
// @Produce      json
// @Param        banner_rotation  body    storage.Rotation  true  "JSON Baneer Rotation"
// @Success      200
// @Router       /banner-rotation [post] .
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

// Delete Banner Rotation
// deleteBannerRotation             godoc
// @Summary      Delete Banner from a Slot
// @Description  Delete Banner Rotation
// @Tags         Banner-Rotation
// @Produce      json
// @Param        banner_rotation  body    storage.Rotation  true  "JSON Baneer Rotation"
// @Success      200
// @Router       /banner-rotation [delete] .
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

// Add Banner
// addBannerToDB             godoc
// @Summary      Create New Banner
// @Description  Responds with the instance with added baneer ID
// @Tags         Banner
// @Produce      json
// @Param        banner  body  Item  true  "Banner JSON"
// @Success      200  {array}  ItemResult
// @Router       /banner [post].
func addBanner(w http.ResponseWriter, r *http.Request, a app.Application) {
	handleAddItem(w, r, a, "banner")
}

// Add Slot
// addSlotToDB             godoc
// @Summary      Create New Slot
// @Description  Responds with the instance with added Slot ID
// @Tags         Slot
// @Produce      json
// @Param        slot  body    Item  true  "SlotJSON"
// @Success      200  {array}  ItemResult
// @Router       /slot [post].
func addSlot(w http.ResponseWriter, r *http.Request, a app.Application) {
	handleAddItem(w, r, a, "slot")
}

// Add Group
// addGroupToDB           godoc
// @Summary      Create New Social Group
// @Description  Responds with the instance with added baneer ID
// @Tags         Social Group
// @Produce      json
// @Param        sosial_group body Item  true  "Social Group JSON"
// @Success      200  {array}  ItemResult
// @Router       /group [post].
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

	var inputItem Item
	err = json.Unmarshal(body, &inputItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var id int

	switch itemType {
	case "slot":
		id, err = a.CreateSlot(context.Background(), inputItem.Descr)
	case "group":
		id, err = a.CreateGroup(context.Background(), inputItem.Descr)
	case "banner":
		id, err = a.CreateBanner(context.Background(), inputItem.Descr)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data, err := json.Marshal(ItemResult{ID: id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Add Banner Rotation
// addBannerRotation             godoc
// @Summary      Update Click Banner stat
// @Description  Add Banner Rotation. Increase clicks count to one
// @Tags         Banner-Rotation
// @Produce      json
// @Param        stat body Statistic  true  "JSON Banner Stat. Clicks c"
// @Success      200
// @Router       /banner-rotation [put] .
func updateClickStat(w http.ResponseWriter, r *http.Request, a app.Application) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var stat Statistic
	err = json.Unmarshal(body, &stat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = a.UpdateClickStat(context.Background(), storage.Statistic{
		BannerID:      stat.BannerID,
		SlotID:        stat.SlotID,
		SosialGroupID: stat.SocialGroupID,
	})
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
