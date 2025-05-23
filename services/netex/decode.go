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
	"strconv"
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

func DownloadAndExtractIfNeededOffer(url string) (string, error) {
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
	if !bytes.Equal(zipHeaderBytes, zipHeader) {
		return "", fmt.Errorf("invalid file")
	}

	utils.LogMessage(utils.LevelInfo, "💥 File is a ZIP file")
	dst := fmt.Sprintf("/tmp/%s", uuid.New().String())
	archive, err := zip.OpenReader(fileName)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	if len(archive.File) == 0 {
		return "", fmt.Errorf("invalid file")
	}

	// this time, the ZIP contains a folder which contains a lot of XML files. We are only interested in the xxx_commun.xml file.
	// we need to unzip the file and return the path to the xxx_commun.xml file.

	// find the xxx_commun.xml file
	var communFileName string
	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		if !strings.HasSuffix(filePath, "_commun.xml") {
			continue
		}

		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return "", fmt.Errorf("invalid file path")
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		fmt.Println("opening file ", filePath)
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		fmt.Println("writing file ", filePath)
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()

		communFileName = filePath
	}

	// delete the zip file
	os.Remove(fileName)

	return communFileName, nil
}

func DecodeNetexStopsData(file string) (*models.PublicationDelivery, error) {
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

func DecodeNetexOfferData(fileName string) (*models.NETEXCommonFile, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var netexData models.NETEXCommonFile
	err = xml.Unmarshal(data, &netexData)
	if err != nil {
		return nil, err
	}

	return &netexData, nil
}

func SaveNetexStopsToDatabase(netexData *models.PublicationDelivery, db *sql.DB) error {
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

func SaveNetexOfferToDatabase(netexData *models.NETEXCommonFile, db *sql.DB) error {
	lines := netexData.DataObjects.GeneralFrame.Members.Lines

	// on crée une transaction et on supprime toutes les données existantes, avant de les insérer
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM NETEX_Line")
	if err != nil {
		return err
	}

	for _, line := range lines {
		routeSortOrder := 0
		if line.KeyList != nil {
			for _, keyValue := range line.KeyList {
				if keyValue.Key == "route_sort_order" {
					routeSortOrder, err = strconv.Atoi(keyValue.Value)
					if err != nil {
						return err
					}
				}
			}
		}

		_, err := tx.Exec("INSERT INTO NETEX_Line (id, version, name, short_name, transport_mode, public_code, private_code, colour, text_colour, route_sort_order) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", line.ID, line.Version, line.Name, line.ShortName, line.TransportMode, line.PublicCode, line.PrivateCode, line.Presentation.Colour, line.Presentation.TextColour, routeSortOrder)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
