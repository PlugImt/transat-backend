package netex

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

func DownloadAndExtractIfNeeded(url string) (string, error) {
	// download the file
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// save the file to the local filesystem to /tmp/<random id>
	fileName := fmt.Sprintf("/tmp/%s", uuid.New().String())
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	io.Copy(file, resp.Body)
	defer resp.Body.Close()
	defer os.Remove(fileName)

	// check the file's first 4 bytes to see if it's a ZIP file
	zipHeader := []byte{0x50, 0x4B, 0x03, 0x04}
	zipHeaderBytes := make([]byte, 4)
	_, err = file.ReadAt(zipHeaderBytes, 0)
	if err != nil {
		return "", err
	}
	if bytes.Equal(zipHeaderBytes, zipHeader) {
		utils.LogMessage(utils.LevelInfo, "💥 File is a ZIP file")
		dst := fmt.Sprintf("/tmp/%s", uuid.New().String())
		archive, err := zip.OpenReader(fileName)
		if err != nil {
			panic(err)
		}
		defer archive.Close()

		if len(archive.File) == 0 || len(archive.File) > 1 {
			return "", fmt.Errorf("invalid file")
		}

		zipFile := archive.File[0]

		if !strings.HasSuffix(zipFile.Name, ".xml") {
			return "", fmt.Errorf("invalid file")
		}

		filePath := filepath.Join(dst, zipFile.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return "", fmt.Errorf("invalid file path")
		}
		if zipFile.FileInfo().IsDir() {
			return "", fmt.Errorf("file is a directory")
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := zipFile.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()

		// delete the zip file
		os.Remove(fileName)

		fileName = filePath

	}

	return fileName, nil
}

func DecodeNetexData(file string) (*models.PublicationDelivery, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var netexData models.PublicationDelivery
	err = xml.Unmarshal(data, &netexData)
	if err != nil {
		return nil, err
	}

	return &netexData, nil
}

func SaveNetexToDatabase(netexData *models.PublicationDelivery, db *sql.DB) error {
	stopPlaces := netexData.DataObjects.GeneralFrame.Members.StopPlaces
	quays := netexData.DataObjects.GeneralFrame.Members.Quays

	// on crée une transaction et on supprime toutes les données existantes, avant de les insérer
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM NETEX_StopPlace")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM NETEX_Quay")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM NETEX_StopPlace_QuayRef")
	if err != nil {
		return err
	}

	for _, stopPlace := range stopPlaces {
		_, err := tx.Exec("INSERT INTO NETEX_StopPlace (id, modification, name, longitude, latitude, transport_mode, other_transport_modes, stop_place_type, weighting) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", stopPlace.ID, stopPlace.Modification, stopPlace.Name, stopPlace.Centroid.Location.Longitude, stopPlace.Centroid.Location.Latitude, stopPlace.TransportMode, stopPlace.OtherTransportModes, stopPlace.StopPlaceType, stopPlace.Weighting)
		if err != nil {
			return err
		}
	}

	for _, quay := range quays {
		_, err := tx.Exec("INSERT INTO NETEX_Quay (id, name, longitude, latitude, site_ref_stopplace_id, transport_mode) VALUES ($1, $2, $3, $4, $5, $6)", quay.ID, quay.Name, quay.Centroid.Location.Longitude, quay.Centroid.Location.Latitude, quay.SiteRef.Ref, quay.TransportMode)
		if err != nil {
			return err
		}
	}

	// now, populate the NETEX_StopPlace_QuayRef link table
	for _, stopPlace := range stopPlaces {
		for _, quayRef := range stopPlace.QuayRefs {
			_, err := tx.Exec("INSERT INTO NETEX_StopPlace_QuayRef (stop_place_id, quay_id, quay_ref_version) VALUES ($1, $2, $3)", stopPlace.ID, quayRef.Ref, quayRef.Version)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
