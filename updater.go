package client

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	update "github.com/inconshreveable/go-update"
)

// UpdateService is a service when updating the client
type UpdateService struct {
	Version    string
	Identifier string
}

// updateList contains a list of all versions from the update API
type updateList struct {
	Versions map[string]updateVersion `json:"versions"`
}

// updateVersion contains information about every available version
type updateVersion struct {
	Download      string `json:"download"`
	Checksum      string `json:"checksum"`
	Patch         bool   `json:"patch"`
	PatchChecksum string `json:"patch_checksum"`
	PatchDownload string `json:"patch_download"`
}

// Update tries to update the client
func (u *UpdateService) Update(url string) error {
	// Hämta senaste information från APIn
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url + "/download.php?software=" + u.Identifier)
	if err != nil {
		return errors.New("Can't connect to url " + url + "/download.php")
	}
	defer r.Body.Close()

	upd := updateList{}
	err = json.NewDecoder(r.Body).Decode(&upd)
	if err != nil {
		return errors.New("Can't get the latest versions: " + err.Error())
	}

	keys := make([]string, 0, len(upd.Versions))
	for k := range upd.Versions {
		keys = append(keys, k)
	}

	pos := arrayPosition(keys, u.Version) + 1
	if pos < 0 || pos >= len(keys) {
		return errors.New("no more versions")
	}

	if !u.containsPatches(upd.Versions, keys, pos) {
		// Update to latest version
		return u.internalUpdate(fmt.Sprintf("%s%s", url, upd.Versions[keys[pos]].Download), upd.Versions[keys[pos]].Checksum, update.Options{})
	}

	// Update using patches
	for i := pos; i < len(keys); i++ {
		err := u.internalUpdate(fmt.Sprintf("%s%s", url, upd.Versions[keys[i]].PatchDownload), upd.Versions[keys[i]].PatchChecksum, update.Options{
			Patcher: update.NewBSDiffPatcher(),
		})
		if err != nil {
			return errors.New("Something went wrong when updating: " + err.Error())
		}
	}

	return nil
}

// containsPatches checks if all versions contains patches
func (u *UpdateService) containsPatches(haystack map[string]updateVersion, versions []string, start int) bool {
	for i := start; i < len(versions); i++ {
		if !haystack[versions[i]].Patch {
			return false
		}
	}
	return true
}

// internalUpdate downloads a patch/update and updates the client
func (u *UpdateService) internalUpdate(url string, hexChecksum string, options update.Options) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	checksum, err := hex.DecodeString(hexChecksum)
	if err != nil {
		return err
	}

	options.Hash = crypto.SHA256
	options.Checksum = checksum

	err = update.Apply(resp.Body, options)
	if err != nil {
		fmt.Println(err)
		/*if err = update.RollbackError(err); err != nil {
			return err
		}*/
		return err
	}
	return nil
}

// arrayPosition loops through a haystack and check where the needle are
func arrayPosition(haystack []string, needle string) int {
	for k, v := range haystack {
		if v == needle {
			return k
		}
	}
	return -1
}
